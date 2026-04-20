## ADDED Requirements

### Requirement: 登录接口限速

登录接口 `POST /api/v1/auth/login` SHALL 实施基于 IP 地址的滑动窗口限速，使用 Redis 存储计数。

限速规则：
- 每个 IP 每分钟最多 10 次登录请求
- key 格式: `ratelimit:login:{ip}`
- 超出限速返回 HTTP 429，响应体包含 `retry_after` 秒数

#### Scenario: 正常登录不受影响

- **WHEN** 同一 IP 在 1 分钟内发起 ≤10 次登录请求
- **THEN** 所有请求正常处理，不限速

#### Scenario: 超出登录限速

- **WHEN** 同一 IP 在 1 分钟内发起第 11 次登录请求
- **THEN** 返回 HTTP 429，响应 `{"code": 429, "message": "请求过于频繁，请稍后再试", "retry_after": <秒>}`

#### Scenario: 限速窗口自动重置

- **WHEN** 限速触发后等待超过窗口时间（1 分钟）
- **THEN** 该 IP 的登录请求恢复正常处理

### Requirement: 通用 API 限速

所有需要认证的 API 接口 SHALL 实施基于用户 ID 的滑动窗口限速。

限速规则：
- 每个用户每分钟最多 60 次请求
- key 格式: `ratelimit:api:{user_id}`
- 超出限速返回 HTTP 429

#### Scenario: 正常 API 调用不受影响

- **WHEN** 用户在 1 分钟内发起 ≤60 次 API 请求
- **THEN** 所有请求正常处理

#### Scenario: 超出 API 限速

- **WHEN** 用户在 1 分钟内发起第 61 次请求
- **THEN** 返回 HTTP 429，响应 `{"code": 429, "message": "请求过于频繁"}`

#### Scenario: 未认证请求不触发 API 限速

- **WHEN** 请求未携带有效 JWT token
- **THEN** 不应用用户级限速（由认证中间件先行拦截）
