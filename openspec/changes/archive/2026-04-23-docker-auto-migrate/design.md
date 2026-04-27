## Context

已有 `cmd/migrate/main.go` 工具，通过 `schema_migrations` 表跟踪已执行的 SQL 文件，按文件名排序执行未跑过的 migration。当前 Docker 容器启动只运行 `./server`，不会调用 migrate 工具。

后端 Dockerfile 已是多阶段构建，最终镜像为 Alpine，包含 `./server`、`configs/`、`migrations/`。

## Goals / Non-Goals

**Goals:**
- 容器启动时自动执行 SQL migration，无需手动干预
- 幂等安全：已执行的 migration 不重复执行
- 失败时容器不启动，避免应用跑在错误的 schema 上

**Non-Goals:**
- 不做 migration 回滚（只进不退）
- 不替换现有 GORM AutoMigrate（两者共存）
- 不改动 `cmd/migrate/main.go` 本身

## Decisions

### 1. 编译 migrate 二进制加入镜像

在 Dockerfile 构建阶段额外编译 `cmd/migrate`，最终镜像包含 `./server` + `./migrate`。

**替代方案：** 在 entrypoint 中用 shell 调用 mysql 客户端 — 放弃，因为需要额外安装 mysql client 且不如复用已有 Go 工具。

### 2. 用 entrypoint.sh 先 migrate 再启动 server

创建 `entrypoint.sh`：
```bash
#!/bin/sh
./migrate
exec ./server
```

migrate 失败时 `exit 1`，容器不启动。`exec` 确保 server 进程接收信号。

**替代方案：** 在 docker-compose.yml 中用 `command: sh -c "./migrate && ./server"` — 可行但 entrypoint.sh 更清晰可维护。

### 3. DB 连接环境变量

migrate 工具通过环境变量读取 DB 连接信息（`DB_HOST`、`DB_USER` 等）。Docker 环境中 host 为 `db`（docker-compose 服务名），需在 docker-compose.yml 中配置。

复用已有的 `APP_DATABASE_HOST` 等变量，migrate 工具的环境变量名保持不变。

### 4. 启动顺序：等 MySQL 就绪

MySQL 容器启动后需要几秒才能接受连接。用 `depends_on` 只保证容器启动，不保证服务就绪。

migrate 工具中已有 `db.Ping()` 检测，失败会退出。配合 `restart: unless-stopped`，容器会自动重试直到 MySQL 就绪。

## Risks / Trade-offs

- **[migration SQL 有语法错误]** → 容器启动失败，不会影响正在运行的旧版本。修复 SQL 后重新部署即可。
- **[大型 migration 执行时间长]** → 当前 migration 都很小（ALTER TABLE 级别），暂无风险。后续如需大数据量 migration，可在 entrypoint.sh 中加超时。
- **[MySQL 未就绪导致 migrate 失败]** → `restart: unless-stopped` 会自动重启容器，最多重试几次后 MySQL 就会就绪。
