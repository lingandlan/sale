## Context

当前部署流程：CI 在 GitHub Actions 中构建镜像 → push 到 ghcr.io → SSH 到服务器 pull 镜像 → docker compose up。国内服务器从 ghcr.io pull 镜像耗时 10+ 分钟。

## Goals / Non-Goals

**Goals:**
- 镜像推送和拉取改用阿里云 ACR，国内服务器秒级拉取
- 最小改动，只换镜像仓库地址和认证方式
- 保留 ghcr.io 作为可选备用（不删除 GHCR 登录步骤，只是不再使用）

**Non-Goals:**
- 不改 ACR 的仓库结构（命名空间、仓库名在阿里云控制台手动创建）
- 不做镜像同步（GHCR ↔ ACR 双写）
- 不改 CI/CD 的触发条件和部署脚本逻辑

## Decisions

### 1. 使用个人版 ACR

阿里云容器镜像服务个人版免费，足够当前项目使用。

**替代方案：** 企业版支持跨区域同步 — 当前不需要。

### 2. 镜像地址格式

`registry.cn-<region>.aliyuncs.com/<namespace>/<repo>:latest`

用户需在阿里云控制台创建命名空间和仓库后，提供具体地址。

### 3. 认证方式

CI 和服务器都通过 `docker login` 使用 ACR 的固定密码（在阿里云控制台设置镜像仓库的固定密码）。

存为 GitHub Secrets：`ACR_REGISTRY`、`ACR_USERNAME`、`ACR_PASSWORD`。

### 4. 同时保留 GHCR permissions

deploy.yml 中保留 `permissions: packages: write`，但实际不再推送到 GHCR。后续可清理。

## Risks / Trade-offs

- **[ACR 凭据泄露]** → 使用 GitHub Secrets 存储，不硬编码
- **[ACR 服务不可用]** → 概率极低，阿里云核心服务
- **[需要手动在阿里云创建仓库]** → 一次性操作，不可自动化（个人版）
