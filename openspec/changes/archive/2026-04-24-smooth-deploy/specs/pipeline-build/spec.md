## ADDED Requirements

### Requirement: Codeup Flow 构建 Docker 镜像
Codeup Flow 流水线 SHALL 在云端构建 backend 和 frontend Docker 镜像，不在服务器上执行构建。

#### Scenario: 后端镜像构建
- **WHEN** push 到 master 触发流水线
- **THEN** Flow 构建步骤执行 `docker build -f backend/Dockerfile -t sale-backend:latest ./backend`，使用 `GOPROXY=https://goproxy.cn,direct` 加速依赖下载

#### Scenario: 前端镜像构建
- **WHEN** push 到 master 触发流水线
- **THEN** Flow 构建步骤执行 `docker build -f shop-pc/Dockerfile -t sale-frontend:latest ./shop-pc`

### Requirement: 镜像传输到服务器
构建完成后 SHALL 将镜像通过 docker save/load 方式传输到部署服务器。

#### Scenario: 镜像打包与传输
- **WHEN** 后端或前端镜像构建成功
- **THEN** 流水线执行 `docker save sale-backend:latest | gzip > backend.tar.gz`，通过 SSH 传输到服务器 `/opt/sale/images/` 目录

#### Scenario: 服务器加载镜像
- **WHEN** 镜像 tar 文件传输到服务器
- **THEN** 服务器执行 `docker load < /opt/sale/images/backend.tar.gz`，加载镜像到本地 Docker

### Requirement: 服务器无需编译环境
服务器 SHALL 只执行 docker load + docker compose up，不执行任何构建命令。

#### Scenario: 部署过程无 build 步骤
- **WHEN** Flow 触发服务器部署
- **THEN** 服务器执行的脚本中不包含 `docker compose build` 或 `go build` 或 `npm install`
