## Why

SQL migration 文件存在于 `backend/migrations/sql/`，且有独立的 `cmd/migrate` 工具可以按顺序执行未跑过的 migration。但 Docker 容器启动时只运行 `./server`，从未执行 migration，导致部署后新表/新字段不会自动创建。

## What Changes

- 后端容器启动前自动执行 `cmd/migrate`，运行所有未执行的 SQL migration
- 在 deploy.yml 中配置 migration 所需的数据库连接环境变量
- 复用已有的 `schema_migrations` 表跟踪执行状态，幂等安全

## Capabilities

### New Capabilities
- `docker-migration`: Docker 容器启动时自动执行 SQL migration 的机制

### Modified Capabilities
（无已有 spec 需要修改）

## Impact

- `backend/Dockerfile`: 编译 `cmd/migrate` 二进制，添加 entrypoint 脚本先跑 migration 再启动 server
- `.github/workflows/deploy.yml`: 传递数据库连接环境变量给 docker-compose
- `docker-compose.yml` / `docker-compose.prod.yml`: 添加 DB 连接环境变量
