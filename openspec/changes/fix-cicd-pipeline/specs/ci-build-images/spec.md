## ADDED Requirements

### Requirement: CI 编译 Go 后端二进制
流水线 SHALL 在 CI runner 上执行 `cd backend && go mod download && go build -o ../server ./cmd/server && go build -o ../migrate ./migrations`，生成 `server` 和 `migrate` 两个 Linux/amd64 二进制文件。

#### Scenario: Go 编译成功
- **WHEN** CI runner 检出代码并执行 Go 编译命令
- **THEN** 生成可执行的 `server` 和 `migrate` 二进制文件，编译无错误

#### Scenario: Go 编译失败
- **WHEN** 代码存在编译错误
- **THEN** CI 流水线 SHALL 失败并显示编译错误信息，不继续后续步骤

### Requirement: CI 构建前端产物
流水线 SHALL 在 CI runner 上执行 `cd shop-pc && npm install && npm run build`，生成 `shop-pc/dist/` 目录。

#### Scenario: 前端构建成功
- **WHEN** CI runner 执行 npm install 和 npm run build
- **THEN** `shop-pc/dist/` 目录生成，包含 index.html 和 JS/CSS 资产

#### Scenario: TypeScript 类型检查失败
- **WHEN** 前端代码存在 TypeScript 类型错误
- **THEN** CI 流水线 SHALL 失败，不生成 Docker 镜像

### Requirement: CI 打包后端 Docker 镜像
流水线 SHALL 使用仓库根目录作为 Docker build context，执行 `docker build -t sale-backend:latest -f Dockerfile.prod-backend .`。

#### Scenario: 后端镜像构建成功
- **WHEN** `server` 和 `migrate` 二进制已生成，Dockerfile.prod-backend 存在
- **THEN** 生成 `sale-backend:latest` 镜像，包含二进制、entrypoint.sh、configs、migrations

### Requirement: CI 打包前端 Docker 镜像
流水线 SHALL 使用 `shop-pc/` 作为 Docker build context，执行 `docker build -t sale-frontend:latest -f Dockerfile.prod-frontend shop-pc/`。

#### Scenario: 前端镜像构建成功
- **WHEN** `shop-pc/dist/` 和 `shop-pc/nginx.conf` 已生成
- **THEN** 生成 `sale-frontend:latest` 镜像，包含 nginx 配置和静态文件

### Requirement: CI 导出镜像制品
流水线 SHALL 将两个镜像导出为 gzip 压缩的 tar 文件：`docker save sale-backend:latest | gzip > backend-image.tar.gz` 和 `docker save sale-frontend:latest | gzip > frontend-image.tar.gz`，并上传为流水线制品。

#### Scenario: 镜像导出并上传
- **WHEN** 两个 Docker 镜像构建完成
- **THEN** 生成 `backend-image.tar.gz` 和 `frontend-image.tar.gz`，上传为 CI 制品（retention 1 天）

### Requirement: 构建阶段并行执行
前端和后端的构建任务 SHALL 可并行执行（在 build stage 内作为独立 jobs）。

#### Scenario: 并行构建
- **WHEN** build stage 触发
- **THEN** `build_frontend` 和 `build_backend` 两个 jobs 同时开始执行
