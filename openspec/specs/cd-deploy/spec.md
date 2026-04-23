## ADDED Requirements

### Requirement: 合并 main 自动部署
系统 SHALL 在 PR 合并到 main 分支后自动触发部署流程，构建 Docker 镜像并部署到服务器。

#### Scenario: 合并触发部署
- **WHEN** PR 合并到 main 分支
- **THEN** 自动构建后端和前端 Docker 镜像

#### Scenario: 部署成功
- **WHEN** Docker 镜像构建完成
- **THEN** 通过 SSH 连接服务器，拉取新镜像并重启服务

#### Scenario: 部署失败回滚
- **WHEN** 部署过程中出错
- **THEN** 保留上一个镜像版本，可手动回滚

### Requirement: Docker 镜像构建
系统 SHALL 为后端（Go）和前端（Nginx）分别构建 Docker 镜像。

#### Scenario: 后端镜像构建
- **WHEN** 触发构建流程
- **THEN** 使用多阶段构建，最终镜像基于 alpine，包含编译后的 Go 二进制

#### Scenario: 前端镜像构建
- **WHEN** 触发构建流程
- **THEN** 构建前端静态文件，用 Nginx 镜像托管

### Requirement: Docker Compose 编排
系统 SHALL 通过 `docker-compose.yml` 编排后端、前端、Nginx 反向代理。

#### Scenario: 一键启动
- **WHEN** 执行 `docker compose up -d`
- **THEN** 后端、前端、Nginx 容器全部启动并互联

### Requirement: 部署所需 Secrets
部署 SHALL 使用以下 GitHub Secrets：`DEPLOY_HOST`、`DEPLOY_USER`、`DEPLOY_SSH_KEY`、`DEPLOY_PATH`。

#### Scenario: 缺少 Secret
- **WHEN** 必要的 Secret 未配置
- **THEN** 部署 workflow 报错并跳过
