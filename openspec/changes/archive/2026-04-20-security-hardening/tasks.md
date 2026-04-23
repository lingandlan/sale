## 1. Phase 1 — 无破坏性修复

- [x] 1.1 修复 ChangePassword 旧密码校验：在 `handler/user.go:ChangePassword` 中调用 `service.CheckPassword()` 验证 `req.OldPassword`，失败返回 401
- [x] 1.2 Recovery 中间件隐藏错误详情：`middleware/recovery.go` 中移除 `Detail: fmt.Sprintf("%v", err)`，改为固定消息，详情仅写 stdout
- [x] 1.3 移除登录日志中的手机号：`handler/auth.go` 中 `fmt.Printf` 改为不输出 phone 字段
- [x] 1.4 Logout 移入认证路由组：`router/router.go` 中将 `/auth/logout` 移到需要 auth 中间件的组内
- [x] 1.5 CORS 收紧：`middleware/cors.go` 从配置读取 allowed origins，移除 `*` 通配，支持环境变量 `CORS_ALLOWED_ORIGINS`
- [x] 1.6 文件上传安全加固：为 BatchImportCards 添加文件大小限制（10MB）和 MIME type 校验

## 2. Phase 2 — 密码哈希 + 输入校验

- [x] 2.1 操作员密码改 bcrypt：`service/recharge.go` 中 CreateOperator 和 UpdateOperator 调用 `HashPassword()` 哈希密码
- [x] 2.2 操作员登录验证适配（操作员使用统一用户认证，无需单独适配）
- [x] 2.3 写操作员密码迁移脚本：`backend/migrations/sql/` 下创建 SQL + Go 脚本，将现有明文密码批量转 bcrypt，支持 dry-run
- [x] 2.4 为 CreateBRechargeApplication 定义请求结构体（含 binding tag），替换 map[string]interface{}
- [x] 2.5 为 ApprovalRechargeApplication 定义请求结构体，替换 map
- [x] 2.6 为 CreateCRecharge 定义请求结构体，替换 map
- [x] 2.7 为 CreateCenter / UpdateCenter 定义请求结构体，替换 map
- [x] 2.8 为 CreateOperator / UpdateOperator 定义请求结构体，替换 map
- [x] 2.9 更新 service 层方法签名：将 `map[string]interface{}` 参数改为对应结构体指针，内部逻辑适配

## 3. Phase 3 — RBAC 路由扩展

- [x] 3.1 为全部业务路由组添加 `rbacMiddleware.Auth()`
- [x] 3.2 补全 Casbin 策略：按 design 中角色-路由映射表，为各角色配置允许访问的路径和方法
- [x] 3.3 验证各角色访问权限（单元测试全部通过，集成测试待部署后验证）

## 4. Phase 4 — 限速中间件

- [x] 4.1 新建 `middleware/ratelimit.go`：实现基于 Redis 滑动窗口的通用限速中间件
- [x] 4.2 登录限速：在 `/auth/login` 路由上应用，每 IP 每分钟 10 次，超限返回 429 + retry_after
- [x] 4.3 通用 API 限速：在认证中间件后应用，每用户每分钟 60 次，超限返回 429
- [x] 4.4 在 router 中注册限速中间件到对应路由组

## 5. Phase 5 — 密钥管理

- [x] 5.1 配置加载支持环境变量覆盖
- [x] 5.2 config.yaml 占位符化（已自带占位符，添加 warning 日志）
- [x] 5.3 更新 `.env.example` 和部署文档

## 6. 验证

- [x] 6.1 启动后端确认所有接口正常响应
- [x] 6.2 验证限速：连续请求登录接口触发 429
- [x] 6.3 验证 RBAC：用 operator 账号访问 center 管理接口确认 403
- [x] 6.4 验证密码：确认操作员密码为 bcrypt 哈希、ChangePassword 需旧密码
