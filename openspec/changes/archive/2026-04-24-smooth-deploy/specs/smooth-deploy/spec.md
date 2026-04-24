## ADDED Requirements

### Requirement: 按服务单独重建
部署脚本 SHALL 支持单独重建 backend 或 frontend，不影响其他运行中的服务。

#### Scenario: 只重建后端
- **WHEN** 只有后端代码变更
- **THEN** 执行 `docker compose up -d --no-deps --force-recreate backend`，frontend/db/redis 保持运行

#### Scenario: 只重建前端
- **WHEN** 只有前端代码变更
- **THEN** 执行 `docker compose up -d --no-deps --force-recreate frontend`，backend/db/redis 保持运行

#### Scenario: 基础设施变更
- **WHEN** docker-compose.yml 中 db 或 redis 配置变更
- **THEN** 执行 `docker compose down && docker compose up -d`，全量重建

### Requirement: 部署中断最小化
后端重建时 SHALL 将不可用时间控制在 5 秒以内。

#### Scenario: 后端重建期间前端可用
- **WHEN** 后端容器正在重建
- **THEN** 前端页面仍可访问，API 请求短暂返回 502 后自动恢复

#### Scenario: 旧容器优雅停止
- **WHEN** 执行 `docker compose up --force-recreate backend`
- **THEN** 旧 backend 容器收到 SIGTERM，有 5 秒完成进行中的请求后再退出
