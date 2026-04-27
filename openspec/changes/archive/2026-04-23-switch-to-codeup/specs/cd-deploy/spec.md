## MODIFIED Requirements

### Requirement: 镜像仓库使用阿里云 ACR
CI 构建的镜像推送目标从 ghcr.io 改为阿里云容器镜像服务（ACR）。服务器从 ACR 拉取镜像。

#### Scenario: CI 登录并推送镜像
- **WHEN** deploy workflow 执行
- **THEN** 使用 ACR 凭据登录 `registry.cn-<region>.aliyuncs.com`，推送 backend 和 frontend 镜像

#### Scenario: 服务器拉取镜像
- **WHEN** SSH 部署执行 docker compose pull
- **THEN** 服务器使用 ACR 凭据登录，从 ACR 拉取镜像
