## Why

当前 CI/CD 流水线（flow.yaml）存在多个问题导致部署持续失败：Go 编译在服务器上缺少 Go 环境、Docker 镜像在 CI 端构建但 flow.yaml 语法可能不兼容云效 Flow、后端 migrate 连不上 Docker 内的 MySQL、前端容器 SSL 证书挂载路径不匹配。需要一个完整可用的流水线方案，让代码推送后自动构建并部署到生产服务器。

## What Changes

- 重写 `flow.yaml`，在 CI 环境完成 Go 编译、前端构建、Docker 镜像打包，制品传输到服务器后只做 `docker load` + `docker compose up`
- 确保 `docker-compose.prod.yml` 环境变量正确，后端 migrate 和 server 都能通过 `APP_DATABASE_*` 环境变量连接 Docker 内的 MySQL 和 Redis
- 确保 `entrypoint.sh` 中 migrate 先于 server 执行，且数据库已就绪
- 确保 `Dockerfile.prod-backend` 和 `Dockerfile.prod-frontend` 的 COPY 路径与 CI 构建上下文匹配
- 验证云效 Flow YAML 语法兼容性，确保 stages/jobs/steps 写法正确

## Capabilities

### New Capabilities

- `ci-build-images`: CI 环境中编译后端 Go 二进制、构建前端产物、打包 Docker 镜像并导出为 tar.gz 制品
- `server-deploy`: 服务器端加载预构建镜像、使用 docker-compose.prod.yml 启动所有服务（含健康检查验证）
- `env-config`: 生产环境变量配置（数据库、Redis、JWT、CORS、WSY 商城）通过 .env 文件注入 docker-compose.prod.yml

### Modified Capabilities

## Impact

- **flow.yaml**: 完全重写，分为 build 和 deploy 两个 stage
- **docker-compose.prod.yml**: 确认无 build 指令，仅用 image + env
- **Dockerfile.prod-backend / Dockerfile.prod-frontend**: COPY 路径与 CI 构建上下文对齐
- **backend/entrypoint.sh**: 确保 migrate → server 启动顺序
- **backend/internal/config/config.go**: 确保 APP_DATABASE_* / APP_REDIS_* 环境变量覆盖生效
- **backend/migrations/migrate.go**: 确保读取环境变量连接正确的数据库
- **服务器 /opt/sale/.env**: 需要正确配置所有生产环境变量
