## Why

国内服务器从 GHCR 拉取 Docker 镜像极慢（10+ 分钟），严重影响部署效率。将整个 CI/CD 迁移到阿里云 Codeup Pipeline + ACR，可实现：
- 国内构建和拉取镜像秒级完成
- 源码和流水线在同一平台，减少外部依赖
- GitHub 只做代码托管，不再承担构建部署

## What Changes

- 在 Codeup 创建流水线：push 到 main 自动构建镜像 → 推送到 ACR → SSH 部署到服务器
- 保留 GitHub 仓库用于代码协作，通过镜像同步或手动 push 到 Codeup
- 移除 GitHub Actions 的 deploy.yml（CI 流水线可保留用于 PR 检查）
- ACR 存储镜像，服务器从 ACR pull

## Capabilities

### New Capabilities
- `codeup-pipeline`: 阿里云 Codeup 流水线 + ACR 镜像仓库的完整 CI/CD 方案

### Modified Capabilities
- `cd-deploy`: 部署方式从 GitHub Actions SSH 改为 Codeup Pipeline SSH
- `ci-pipeline`: GitHub Actions CI 保留用于 PR 检查，但不再负责部署

## Impact

- 需在 Codeup 配置流水线（YAML 或控制台配置）
- 需创建 ACR 命名空间和仓库
- 服务器部署脚本从 ACR 拉取镜像
- GitHub Actions 的 deploy.yml 和 release.yml 可移除或禁用
