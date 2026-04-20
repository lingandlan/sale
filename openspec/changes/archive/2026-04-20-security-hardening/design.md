## Context

太积堂充值与门店管理系统（Go + Gin + Vue 3），当前处于开发阶段，准备进入生产部署。安全审计发现 12 个安全问题（4 CRITICAL、4 HIGH、4 MEDIUM），涵盖认证、输入校验、密钥管理和网络安全。

现有架构特点：
- JWT + Redis refresh token 认证，bcrypt 哈希用户密码
- Casbin RBAC 权限控制，但仅覆盖 `/user/*`、`/system/*`、`/admin/*` 路由
- 多数 handler 接收 `map[string]interface{}` 无结构体验证
- 配置通过 `configs/config.yaml` 加载，密钥明文写入仓库

## Goals / Non-Goals

**Goals:**
- 消除所有 CRITICAL 和 HIGH 级别安全漏洞
- 建立可复用的安全中间件模式（限速、校验、RBAC），供后续开发遵循
- 保持 API 向后兼容，不引入破坏性变更

**Non-Goals:**
- 不做 HTTPS/TLS 配置（由反向代理/基础设施层处理）
- 不引入 WAF 或专业安全扫描工具
- 不改造前端代码，安全修复全在后端
- 不做审计日志（可单独开 change）

## Decisions

### D1: 操作员密码 — bcrypt 哈希存储

**现状**: `service/recharge.go:758` 明文存储 `Password: data["password"].(string)`
**方案**: 复用已有 `service.HashPassword()`（bcrypt），CreateOperator 和 UpdateOperator 中调用
**迁移**: 写一次性 SQL 迁移脚本，将现有明文密码批量转 bcrypt，默认密码 `123456`
**替代方案**: 无（项目已对用户密码使用 bcrypt，保持一致）

### D2: ChangePassword — 增加旧密码校验

**现状**: `handler/user.go:119` 有 `// TODO: 验证旧密码`，直接跳过验证
**方案**: 使用已有的 `service.CheckPassword()` 验证 `req.OldPassword`，失败返回 `401`
**改动范围**: 仅 `handler/user.go:ChangePassword` 方法，约 5 行

### D3: 密钥管理 — 环境变量覆盖

**现状**: JWT secret、DB 密码、第三方 API 密钥硬编码在 `configs/config.yaml`
**方案**: 在 config 加载逻辑中增加环境变量覆盖，优先级 `env > config.yaml > default`
- `JWT_SECRET` → `jwt.secret`
- `DB_PASSWORD` → `database.password`
- `REDIS_PASSWORD` → `redis.password`
- `MALL_APP_ID` → `mall.app_id`
- `MALL_APP_SECRET` → `mall.app_secret`

`config.yaml` 保留占位符值，`.gitignore` 中确保不提交 `.env`（已 gitignored）
**替代方案**: 使用 Vault 等专业密钥管理 → 过度工程，当前阶段环境变量足够

### D4: RBAC 路由覆盖 — 扩展到全部业务路由

**现状**: 仅 `/user`、`/system`、`/admin` 三个路由组有 RBAC，其余 7 个路由组仅 JWT 认证
**方案**: 为 Dashboard、B-Recharge、C-Recharge、Card、Center、Operator、Records 路由组添加 `rbacMiddleware.Auth()`
- Dashboard: `center_admin` 及以上
- B-Recharge apply: `operator` 及以上
- B-Recharge approval: `center_admin` 及以上
- C-Recharge: `operator` 及以上
- Card: `operator` 及以上
- Center: `center_admin` 及以上（创建/修改/删除需 `hq_admin`）
- Operator: `center_admin` 及以上

需同步更新 Casbin 策略文件
**风险**: 现有 Casbin 策略可能不完整，需要逐条补全

### D5: 输入校验 — 请求结构体替换 map

**现状**: `handler/recharge.go` 中 6 个接口用 `map[string]interface{}` 接收入参
**方案**: 为每个接口定义 binding struct，使用 Gin 的 `ShouldBindJSON` + validator tag
```go
type CreateOperatorRequest struct {
    Name     string `json:"name" binding:"required,min=1,max=50"`
    Phone    string `json:"phone" binding:"required,len=11"`
    Password string `json:"password" binding:"required,min=6,max=32"`
    CenterID string `json:"centerId" binding:"required,uuid"`
    Role     string `json:"role" binding:"required,oneof=center_admin operator"`
}
```
涉及接口：CreateBRechargeApplication、ApprovalRechargeApplication、CreateCRecharge、CreateCenter、UpdateCenter、CreateOperator、UpdateOperator
**替代方案**: 保留 map 但加手动校验 → 不如 struct binding 声明式、可维护

### D6: 限速中间件 — Redis 滑动窗口

**方案**: 新建 `middleware/ratelimit.go`，基于 Redis 实现滑动窗口限速
- 登录: 每个 IP 每分钟 10 次，失败 5 次后锁定 15 分钟
- 通用 API: 每用户每分钟 60 次
- key 格式: `ratelimit:{type}:{identifier}`
**替代方案**: 令牌桶（更平滑但对登录场景不需要）、内存计数（多实例不共享）

### D7: CORS 收紧

**现状**: `Access-Control-Allow-Origin: *` + `Allow-Credentials: true`（浏览器实际拒绝此组合）
**方案**: 从配置读取允许的 origins 列表，动态设置 `Allow-Origin`
- 开发环境: `http://localhost:5173,http://localhost:5175`
- 生产环境: 从环境变量 `CORS_ALLOWED_ORIGINS` 读取
- 移除 `Allow-Credentials: true` 或改为配合具体 origin 使用

### D8: 其他 MEDIUM 修复

| 修复项 | 方案 |
|--------|------|
| Logout 无认证 | 移到需要 auth 中间件的路由组 |
| Recovery 泄露错误详情 | 返回固定消息 `"internal server error"`，详情仅写日志 |
| 登录日志泄露手机号 | 改为 `login failed: error=...`，不打印 phone |
| 文件上传无限制 | 添加 `MaxMultipartMemory` 全局限制 + 扩展名白名单校验 + MIME type 检测 |

## Risks / Trade-offs

| Risk | Impact | Mitigation |
|------|--------|------------|
| RBAC 策略补全可能遗漏路径，导致合法用户 403 | HIGH | 部署前完整测试所有角色的访问路径，准备 Casbin 策略回滚方案 |
| 操作员密码迁移脚本执行失败 | MEDIUM | 迁移前备份数据库，脚本含 dry-run 模式 |
| 限速误伤正常用户（NAT/代理共享 IP） | LOW | 登录限速以 IP + phone 组合为 key，而非纯 IP |
| 输入结构体验证可能拒绝之前能通过的请求 | LOW | struct tag 与前端现有字段对齐，只加 required/min/max，不改字段名 |

## Migration Plan

1. **Phase 1 — 无破坏性修复**（可立即部署）
   - ChangePassword 旧密码校验
   - Recovery 中间件隐藏错误详情
   - 移除登录日志中的手机号
   - Logout 移入认证路由组
   - CORS 收紧

2. **Phase 2 — 数据迁移 + 代码变更**（需维护窗口）
   - 操作员密码迁移脚本
   - 请求结构体替换 map
   - RBAC 路由扩展 + Casbin 策略补全

3. **Phase 3 — 密钥管理**（部署时执行）
   - 设置环境变量
   - config.yaml 替换为占位符

4. **Rollback**: 每个 phase 可独立回滚，git revert 对应 commit 即可

## Open Questions

- Casbin 当前策略是否已覆盖所有角色的权限组合？需要导出现有策略确认
- 生产环境的 CORS allowed origins 具体有哪些域名？
