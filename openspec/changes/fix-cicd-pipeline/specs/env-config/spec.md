## ADDED Requirements

### Requirement: .env 文件包含所有必需变量
`.env` 文件 SHALL 包含以下变量组，缺失任何一个 MUST 导致服务无法正常启动：

- `JWT_SECRET` — JWT 签名密钥
- `MYSQL_ROOT_PASSWORD` — MySQL root 密码
- `MYSQL_DATABASE` — 数据库名（默认 `sale_prod`）
- `MYSQL_USER` — MySQL 用户名（默认 `sale`）
- `MYSQL_PASSWORD` — MySQL 用户密码
- `DB_PASSWORD` — 后端连接密码（与 MYSQL_PASSWORD 一致）
- `MALL_APP_ID` / `MALL_APP_SECRET` / `MALL_CUSTOMER_ID` / `MALL_BASE_URL` — WSY 商城 API（可选）
- `CORS_ALLOWED_ORIGINS` — CORS 允许的域名

#### Scenario: 完整配置
- **WHEN** .env 文件包含所有必需变量
- **THEN** docker compose 读取 .env，所有服务正常启动

#### Scenario: 缺少必需变量
- **WHEN** JWT_SECRET 或 DB_PASSWORD 缺失
- **THEN** 后端容器 SHALL 因配置错误而启动失败

### Requirement: docker-compose.prod.yml 正确映射环境变量
docker-compose.prod.yml SHALL 将 .env 变量映射为容器环境变量：

- `DB_PASSWORD` → backend 容器 `APP_DATABASE_PASSWORD`
- `JWT_SECRET` → backend 容器 `JWT_SECRET`
- `MALL_*` → backend 容器对应变量
- `CORS_ALLOWED_ORIGINS` → backend 容器对应变量
- `MYSQL_*` → db 容器对应变量

#### Scenario: 后端容器接收数据库环境变量
- **WHEN** docker compose 使用 --env-file .env 启动
- **THEN** backend 容器内 `APP_DATABASE_HOST=db`、`APP_DATABASE_PASSWORD=<实际密码>` 等环境变量正确设置

#### Scenario: Go viper 读取环境变量覆盖
- **WHEN** 容器内 `APP_DATABASE_HOST=db` 环境变量存在
- **THEN** Go 后端 viper 通过 envOverrides map 将其映射到 `database.host` 配置，连接到 Docker 内的 MySQL

### Requirement: .env.example 作为配置模板
项目 SHALL 提供 `.env.example` 文件，列出所有变量及说明，不含真实值。

#### Scenario: 新环境初始化
- **WHEN** 运维在新服务器上 `cp .env.example .env`
- **THEN** 得到包含所有变量占位符的配置文件，逐项填入真实值即可
