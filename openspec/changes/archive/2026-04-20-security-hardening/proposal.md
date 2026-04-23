## Why

安全审计发现多个高危和严重漏洞，包括：明文存储的操作员密码、未验证旧密码的修改密码接口、硬编码在代码仓库中的 API 密钥和 JWT secret、路由缺少 RBAC 中间件、无登录限速、CORS 配置过于宽松等。这些问题在生产环境中可能导致数据泄露、未授权访问和暴力破解攻击。

## What Changes

- **CRITICAL**: 修复 ChangePassword 接口，增加旧密码验证
- **CRITICAL**: 操作员密码改为 bcrypt 哈希存储，不再明文保存
- **CRITICAL**: 将 JWT secret、数据库密码、第三方 API 密钥移至环境变量，config.yaml 不再包含真实密钥
- **HIGH**: 为 Card/Center/Operator/Recharge 等路由组添加 RBAC 中间件保护
- **HIGH**: 为所有接收 `map[string]interface{}` 的接口引入请求结构体验证（bind + validate）
- **HIGH**: 添加登录和 API 限速中间件
- **HIGH**: 收紧 CORS 策略，移除 `Access-Control-Allow-Origin: *`
- **MEDIUM**: Logout 接口移入需要认证的路由组
- **MEDIUM**: 文件上传增加大小限制和 MIME 类型校验
- **MEDIUM**: Panic recovery 中间件不再向客户端返回内部错误详情
- **MEDIUM**: 移除登录失败日志中的手机号输出

## Capabilities

### New Capabilities
- `input-validation`: 统一的请求结构体定义与 bind/validate 校验，替代 map[string]interface{} 入参
- `rate-limiting`: 基于 Redis 的登录限速和通用 API 限速中间件
- `secret-management`: 密钥管理方案，将敏感配置迁移到环境变量

### Modified Capabilities
- `recharge-c-entry`: 操作员密码改为 bcrypt 哈希存储
- `recharge-record-list`: 路由组添加 RBAC 中间件

## Impact

- **后端代码**: handler、service、middleware、router、config 多处修改
- **API 行为**: 部分接口请求/响应格式可能微调（结构体验证替换 map），但总体保持兼容
- **部署方式**: 需要设置环境变量（JWT_SECRET、DB_PASSWORD、MALL_APP_ID、MALL_APP_SECRET），部署文档需更新
- **数据库**: 需要一次性迁移脚本将现有操作员明文密码转为 bcrypt 哈希
- **无前端改动**: 安全修复主要在后端，前端无需变更
