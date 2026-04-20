## ADDED Requirements

### Requirement: 敏感配置通过环境变量覆盖

配置加载逻辑 SHALL 支持环境变量覆盖 config.yaml 中的敏感字段，优先级为：环境变量 > config.yaml > 默认值。

支持的环境变量映射：

| 环境变量 | 覆盖配置路径 | 说明 |
|----------|-------------|------|
| `JWT_SECRET` | `jwt.secret` | JWT 签名密钥 |
| `DB_PASSWORD` | `database.password` | 数据库密码 |
| `REDIS_PASSWORD` | `redis.password` | Redis 密码 |
| `MALL_APP_ID` | `mall.app_id` | 第三方 API app_id |
| `MALL_APP_SECRET` | `mall.app_secret` | 第三方 API app_secret |

#### Scenario: 环境变量覆盖 config.yaml 值

- **WHEN** 设置了环境变量 `JWT_SECRET=my-production-secret`
- **THEN** 应用使用 `my-production-secret` 作为 JWT 签名密钥，忽略 config.yaml 中的值

#### Scenario: 环境变量未设置时使用 config.yaml

- **WHEN** 环境变量 `JWT_SECRET` 未设置
- **THEN** 应用使用 config.yaml 中 `jwt.secret` 的值

#### Scenario: config.yaml 保留占位符

- **WHEN** config.yaml 中 `jwt.secret` 为 `your-super-secret-key-change-in-production`
- **THEN** 应用仍可正常运行（本地开发场景），但启动时输出 warning 日志提示需更换

### Requirement: config.yaml 不含真实密钥

`configs/config.yaml` 中的敏感字段 SHALL 使用占位符值，真实密钥通过环境变量注入。`.env` 文件已在 `.gitignore` 中，确保不提交到仓库。

#### Scenario: 新部署时设置密钥

- **WHEN** 首次部署到新环境
- **THEN** 运维人员通过环境变量设置所有敏感配置，无需修改代码仓库中的文件
