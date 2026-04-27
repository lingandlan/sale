## Context

太积堂生产系统部署在阿里云 ECS，使用云效 Flow CI/CD。当前架构：

- **CI（云效 Flow）**：推送代码触发流水线，在 CI runner（Ubuntu）上构建
- **部署**：通过 VMDeploy 组件将制品传输到 ECS，在服务器上 `docker load` + `docker compose up`
- **运行时**：4 个容器（MySQL 8.0 + Redis 7 + Go backend + nginx frontend）

当前问题：
1. `flow.yaml` 使用 GitHub Actions 风格语法，与云效 Flow 实际语法可能有差异
2. 后端 `migrate` 二进制在容器内执行，需要通过 `APP_DATABASE_*` 环境变量连接 Docker 内的 MySQL（`db` 主机名）
3. `entrypoint.sh` 缺少 MySQL 就绪等待，migrate 可能在 MySQL 未启动时执行失败
4. `Dockerfile.prod-frontend` 的 `shop-pc/dist/` COPY 路径需要在 CI 构建上下文中验证

## Goals / Non-Goals

**Goals:**
- 代码推送到 master 后自动触发完整构建和部署
- CI 环境完成全部编译（Go + npm）和 Docker 镜像打包
- 服务器零编译，只做 `docker load` + `docker compose up`
- 后端容器启动时自动等待 MySQL 就绪、执行 migration、启动 server
- 前端容器正确挂载 SSL 证书，HTTPS 正常工作

**Non-Goals:**
- 不搭建私有镜像仓库（直接 tar.gz 传输）
- 不做蓝绿部署或金丝雀发布（直接滚动重启）
- 不改后端业务逻辑或 API

## Decisions

### 1. CI 构建镜像 vs 服务器构建镜像

**选择：CI 构建**

理由：服务器无 Go/Node 环境，安装编译工具链增加维护成本。CI runner 是 Ubuntu 标准环境，安装 Go/Node 很方便。

构建流程：
1. CI runner 检出代码
2. `cd backend && go build` 生成 `server` 和 `migrate` 二进制
3. `cd shop-pc && npm install && npm run build` 生成 `dist/`
4. `docker build -f Dockerfile.prod-backend -t sale-backend:latest .` 打包后端镜像
5. `docker build -f Dockerfile.prod-frontend -t sale-frontend:latest shop-pc/` 打包前端镜像
6. `docker save | gzip` 导出为 tar.gz 制品

注意 Docker build context：
- 后端 Dockerfile 使用仓库根目录作为 context（`COPY server .` + `COPY backend/entrypoint.sh .`）
- 前端 Dockerfile 使用 `shop-pc/` 作为 context（`COPY dist/` + `COPY nginx.conf`）

### 2. entrypoint.sh 添加 MySQL 就绪等待

**选择：在 entrypoint.sh 中循环检测 MySQL**

`db:3306` TCP 端口可达即视为就绪，使用 `nc` 或 `wget` 检测。等待最多 60 秒后超时退出。

替代方案：docker-compose `depends_on` + `healthcheck` — 更优雅但需要在 docker-compose.prod.yml 中为 db 添加 healthcheck，增加复杂度。entrypoint.sh 方案更简单可靠。

### 3. 环境变量传递机制

**选择：.env 文件 + docker-compose.prod.yml env 注入**

`.env` 文件由运维手动配置在服务器 `/opt/sale/.env`，docker-compose 通过 `--env-file .env` 读取。

关键变量传递链：
```
.env → docker-compose.prod.yml environment → 容器环境变量 → Go viper envOverrides
```

两套变量名对应关系：
- `DB_PASSWORD` → `APP_DATABASE_PASSWORD`（docker-compose 映射）
- `MYSQL_PASSWORD` / `MYSQL_ROOT_PASSWORD` → MySQL 容器初始化
- `JWT_SECRET` → `JWT_SECRET`（直接传递）
- `MALL_*` → `MALL_*`（直接传递）

### 4. 云效 Flow YAML 语法

**选择：严格遵循云效 Flow 文档语法**

关键差异（vs GitHub Actions）：
- 使用 `stages` → `jobs` → `steps` 层级
- 步骤用 `- name: ...` + `run: |` 格式
- 制品用 `uses: actions/upload-artifact@v4` 和 `uses: actions/download-artifact@v4`
- VMDeploy 用 `component: "VMDeploy"` + `with.downloadArtifact: true`

如果云效不支持 `upload-artifact`，备选方案是将镜像 tar.gz 直接 scp 到服务器。

## Risks / Trade-offs

- **云效 Flow YAML 兼容性** → 需要在实际流水线中测试。如果语法不兼容，根据云效报错调整
- **制品体积** → 后端镜像 tar.gz 约 50-100MB，前端约 10-30MB，传输时间取决于 CI runner 到服务器带宽
- **migrate 失败** → entrypoint.sh 中 migrate 失败应导致容器退出（`set -e`），不会静默跳过
- **数据库连接延迟** → entrypoint.sh 等待最多 60 秒，极端情况 MySQL 首次启动可能更慢（首次执行初始化脚本）
