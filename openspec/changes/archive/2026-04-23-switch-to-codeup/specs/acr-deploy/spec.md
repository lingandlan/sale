## ADDED Requirements

### Requirement: CI 推送镜像到阿里云 ACR
deploy.yml SHALL 登录阿里云 ACR 并将 backend 和 frontend 镜像推送到指定的 ACR 仓库。

#### Scenario: 镜像构建后推送到 ACR
- **WHEN** push 到 main 分支触发 deploy workflow
- **THEN** CI 使用 `ACR_REGISTRY`、`ACR_USERNAME`、`ACR_PASSWORD` secrets 登录 ACR，构建并推送 `sale-backend:latest` 和 `sale-frontend:latest`

### Requirement: 服务器从 ACR 拉取镜像
SSH 部署脚本 SHALL 登录 ACR 并拉取最新镜像。

#### Scenario: 服务器 pull 镜像
- **WHEN** SSH 部署步骤执行 `docker compose pull`
- **THEN** 服务器先 `docker login` 到 ACR，然后从 ACR 拉取镜像

## MODIFIED Requirements

### Requirement: docker-compose 镜像地址使用 ACR
docker-compose.prod.yml 和 deploy.yml 内嵌的 docker-compose.yml 中，backend 和 frontend 镜像地址 SHALL 使用 ACR 地址替代 ghcr.io 地址。

#### Scenario: 镜像地址格式
- **WHEN** 查看 docker-compose 配置
- **THEN** backend 镜像地址为 `<ACR_REGISTRY>/sale-backend:latest`，frontend 为 `<ACR_REGISTRY>/sale-frontend:latest`
