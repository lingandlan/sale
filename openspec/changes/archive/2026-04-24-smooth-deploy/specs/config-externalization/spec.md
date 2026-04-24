## ADDED Requirements

### Requirement: 敏感配置通过环境变量注入
后端 config.go SHALL 支持通过环境变量覆盖所有敏感配置项。

#### Scenario: JWT 密钥环境变量
- **WHEN** 设置环境变量 `JWT_SECRET`
- **THEN** 后端使用该值作为 JWT 签名密钥，覆盖 config.yaml 中的默认值

#### Scenario: 数据库密码环境变量
- **WHEN** 设置环境变量 `DB_PASSWORD` 或 `APP_DATABASE_PASSWORD`
- **THEN** 后端使用该值连接数据库

#### Scenario: WSY 商城配置环境变量
- **WHEN** 设置环境变量 `MALL_APP_ID`、`MALL_APP_SECRET`、`MALL_CUSTOMER_ID`、`MALL_BASE_URL`
- **THEN** 后端使用这些值访问 WSY 商城 API

### Requirement: CORS 域名可配置
CORS allowed_origins SHALL 通过环境变量配置，支持生产和开发环境使用不同域名。

#### Scenario: 配置 CORS 域名
- **WHEN** 设置环境变量 `CORS_ALLOWED_ORIGINS=http://example.com,http://www.example.com`
- **THEN** 后端 CORS 中间件允许这些域名的跨域请求

### Requirement: Docker Compose 密码使用 .env 文件
docker-compose.yml SHALL 使用 `${VAR}` 语法引用服务器上的 .env 文件，不在 YAML 中硬编码密码。

#### Scenario: 数据库密码从 .env 读取
- **WHEN** 服务器 `/opt/sale/.env` 文件包含 `MYSQL_ROOT_PASSWORD=xxx`
- **THEN** docker-compose.yml 中 `MYSQL_ROOT_PASSWORD: ${MYSQL_ROOT_PASSWORD}` 使用该值

#### Scenario: .env 文件不进代码仓库
- **WHEN** 查看 .gitignore
- **THEN** `.env` 文件被排除在版本控制之外

### Requirement: 服务器维护独立 .env 文件
服务器上 SHALL 维护 `/opt/sale/.env` 文件存储所有生产环境密钥和配置。

#### Scenario: .env 包含所有必需变量
- **WHEN** 查看服务器 .env 文件
- **THEN** 包含：JWT_SECRET、MYSQL_ROOT_PASSWORD、MYSQL_PASSWORD、DB_PASSWORD、MALL_APP_ID、MALL_APP_SECRET、MALL_CUSTOMER_ID、CORS_ALLOWED_ORIGINS

#### Scenario: 本地开发不受影响
- **WHEN** 本地开发未设置环境变量
- **THEN** 后端继续使用 config.yaml 中的默认值正常工作
