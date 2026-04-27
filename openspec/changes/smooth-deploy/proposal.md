## Why

当前部署有三个问题：
1. `docker compose down` 会中断所有服务，前后端更新都有约 10 秒不可用
2. 敏感配置（JWT 密钥、WSY API 密钥、数据库密码）硬编码在 config.yaml 和 docker-compose.yml 中，无法按环境区分，且存在泄露风险
3. 服务器承担构建工作（go mod download、npm install），在国内网络下不稳定

## What Changes

### 1. 构建迁移到流水线
- Codeup Flow 流水线中执行 `docker build`，不在服务器上构建
- 构建产物通过 `docker save` 打包为 tar，SSH 传到服务器
- 服务器只执行 `docker load` + `docker compose up`，不再需要编译环境
- 好处：服务器不依赖国内网络下载 Go/Node 依赖，部署更快更稳定

### 2. 平滑部署
- Flow 部署脚本改为按服务单独重建（`--no-deps`），避免全量 down
- 前端变更只重建 frontend，后端变更只重建 backend
- 中断时间从 30+ 秒降到 5 秒内

### 3. 配置外部化
- 后端：config.yaml 中的敏感值全部通过环境变量注入（JWT_SECRET、DB_PASSWORD、MALL_*）
- Docker：docker-compose.yml 中的密码使用服务器上的 .env 文件
- 生产环境：服务器上维护 `/opt/sale/.env` 文件存储所有密钥，不进代码仓库
- CORS allowed_origins 通过环境变量配置，支持生产域名

## Capabilities

### New Capabilities
- `pipeline-build`: Codeup Flow 中构建 Docker 镜像并传到服务器
- `smooth-deploy`: 按服务单独重建的部署方案，最小化中断时间

### Modified Capabilities
- `config-externalization`: 敏感配置从代码中提取到环境变量 + .env 文件

## Impact

- Codeup Flow YAML 需要增加构建步骤
- 后端 config.go 需要扩展 envOverrides 映射
- docker-compose.yml 密码改为 `${VAR}` 引用
- 服务器需要创建 .env 文件
- 服务器不再需要源码中的 Go/Node 依赖，但仍需要代码用于 docker compose 编排
- 不影响本地开发（本地继续用 config.yaml 默认值）
