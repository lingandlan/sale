## Context

太积堂充值与门店管理系统，Go 后端 + Vue 前端，部署在单台服务器上。当前状态：
- 代码托管在 GitHub (`lingandlan/sale`)
- 本地 pre-push hook 跑 `go test`，但没有 CI
- 部署完全手动：SSH 到服务器，git pull，重启服务
- 服务器上直接运行（无 Docker），通过 `start.sh` 启动各环境
- 使用 worktree 隔离 alpha/beta/gamma 环境

技术栈：Go 1.22、Node 24、MySQL、Redis、Nginx

## Goals / Non-Goals

**Goals:**
- PR 提交时自动跑后端测试 + 前端构建检查，防止坏代码合入 main
- 合并到 main 后自动构建 Docker 镜像并部署到测试服务器
- 基于 git tag (v*) 触发生产发布，自动生成 changelog
- 所有流程通过 GitHub Actions 实现

**Non-Goals:**
- 不搭建 Kubernetes 或复杂编排，保持单服务器 Docker Compose 部署
- 不做多环境（staging/production）的自动区分，测试环境和生产环境在同一台机器上用不同 compose 文件
- 不做前端 E2E 测试（目前没有）
- 不改造现有的 worktree 开发模式

## Decisions

### 1. CI 流水线设计

**Decision**: 单个 `ci.yml` workflow，PR 和 push 到 main/develop 时触发，包含三个 job：`backend-test`、`frontend-build`、`lint`

**Rationale**: Go 测试和前端构建互相独立，用并行 job 加速。单独 lint job 检查代码风格。

**Alternatives considered**:
- 分成多个 workflow → 增加维护复杂度，收益不大
- 不跑前端构建 → 前端 TypeScript 编译错误容易漏过

### 2. CD 部署方式

**Decision**: SSH 到服务器执行 `docker compose pull && docker compose up -d`

**Rationale**: 项目规模小，单服务器部署，Docker Compose 足够。GitHub Actions 通过 SSH key 连接服务器。

**Alternatives considered**:
- GitHub Actions self-hosted runner → 配置复杂，安全性考量
- AWS/GCP 部署 → 当前不需要云服务
- 直接 git pull + systemctl restart → 无法回滚镜像，不够可靠

### 3. Docker 镜像策略

**Decision**: 后端构建为单个 Docker 镜像，前端用 Nginx 镜像托管静态文件。用 `docker-compose.yml` 编排后端 + 前端 + Nginx 反向代理。

**Rationale**: 前后端分离构建，Nginx 统一处理路由和 API 代理，与现有 Nginx 配置兼容。

### 4. 版本发布

**Decision**: 基于 git tag `v*` 触发 release workflow，自动构建镜像、生成 changelog（从 git log 提取）、创建 GitHub Release。

**Rationale**: 最简单直接的版本管理方式，不需要额外的版本号文件。

## Risks / Trade-offs

- [Docker 化需要迁移] → 服务器目前裸跑，Docker 化是额外工作。Mitigation: 先在测试环境验证，保留裸跑作为回退
- [SSH key 安全性] → GitHub Secrets 存储 SSH 私钥，泄露风险。Mitigation: 使用 restricted SSH key，限制可执行的命令
- [单点部署无高可用] → 服务器挂了服务就挂了。Mitigation: 当前阶段可接受，后续可扩展
- [Docker 镜像体积] → Go 后端用多阶段构建，最终镜像用 scratch/alpine 压缩体积

## Open Questions

- 服务器 IP 和 SSH 配置需要用户提供
- 是否需要数据库迁移也自动化（目前手动执行 SQL）
- 生产环境是否也 Docker 化，还是只做测试环境
