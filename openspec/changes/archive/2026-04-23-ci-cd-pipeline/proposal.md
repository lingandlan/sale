## Why

项目目前没有任何 CI/CD 配置，只有本地 pre-push hook 跑测试。代码合并到 main 后需要手动部署，没有自动化测试保障、没有版本管理、没有自动发布流程。随着项目功能增多，手动流程容易出错且效率低。

## What Changes

- 新增 GitHub Actions CI 流水线：PR 提交时自动运行 Go 后端测试 + 前端 lint/build 检查
- 新增 GitHub Actions CD 流水线：PR 合并到 main 后自动构建 Docker 镜像并部署到服务器
- 新增版本发布管理：基于 git tag 触发自动构建，生成 changelog
- 新增 Dockerfile 和 docker-compose.yml 支持容器化部署
- 新增 GitHub Actions 所需的 Secrets 配置文档

## Capabilities

### New Capabilities

- `ci-pipeline`: PR 和 push 时的自动化测试、lint、构建检查
- `cd-deploy`: 合并到 main 后自动构建并部署到服务器
- `release-management`: 基于 tag 的版本发布和 changelog 生成

### Modified Capabilities

（无）

## Impact

- **新增文件**: `.github/workflows/ci.yml`、`.github/workflows/deploy.yml`、`.github/workflows/release.yml`、`Dockerfile`、`docker-compose.yml`
- **依赖**: 需要在 GitHub repo Settings 中配置 Secrets（SSH key、服务器地址等）
- **服务器**: 需要安装 Docker + Docker Compose
- **现有流程**: pre-push hook 保持不变，CI 作为额外保障
