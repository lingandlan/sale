## 1. Docker 化基础

- [x] 1.1 创建后端 `backend/Dockerfile`（多阶段构建：golang:1.22 编译 → alpine 运行）
- [x] 1.2 创建前端 `shop-pc/Dockerfile`（node 构建 → nginx 托管静态文件）
- [x] 1.3 创建前端 `shop-pc/nginx.conf`（API 代理到后端 + SPA 路由 fallback）
- [x] 1.4 创建根目录 `docker-compose.yml`（backend + frontend 服务编排）
- [x] 1.5 创建 `.dockerignore` 避免无关文件进入镜像

## 2. CI 流水线

- [x] 2.1 创建 `.github/workflows/ci.yml`（触发条件：PR 到 main/design + push 到 main）
- [x] 2.2 添加 `backend-test` job：`go test ./...` + 覆盖率报告
- [x] 2.3 添加 `frontend-build` job：`npm install && npm run build`
- [x] 2.4 配置 GitHub Branch Protection：CI 通过才允许合并（需手动操作）

## 3. CD 自动部署

- [x] 3.1 创建 `.github/workflows/deploy.yml`（触发条件：push 到 main）
- [x] 3.2 添加构建 Docker 镜像步骤（后端 + 前端）
- [x] 3.3 添加 SSH 部署步骤（拉取镜像 + docker compose up -d）
- [x] 3.4 在 GitHub repo Settings 中配置 Secrets：`DEPLOY_HOST`、`DEPLOY_USER`、`DEPLOY_SSH_KEY`、`DEPLOY_PATH`（需手动操作）

## 4. 版本发布管理

- [x] 4.1 创建 `.github/workflows/release.yml`（触发条件：推送 `v*` tag）
- [x] 4.2 添加自动生成 changelog 步骤（从 git log 提取 commit）
- [x] 4.3 添加创建 GitHub Release 步骤（标题 + changelog body）

## 5. 验证

- [x] 5.1 本地验证 `docker compose up` 能正常启动所有服务
- [x] 5.2 创建测试 PR 验证 CI 流水线运行
- [x] 5.3 合并测试 PR 验证自动部署
- [x] 5.4 推送测试 tag 验证 release 流程
