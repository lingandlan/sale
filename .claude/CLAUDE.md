# 太积堂充值与门店管理系统

## 端口分配表

| Session | Backend | 前端 PC | 数据库 | Redis DB |
|---------|---------|---------|--------|----------|
| Main | 8080 | 5175 | sale_dev | 0 |
| Alpha | 8081 | 5178 | sale_alpha | 1 |
| Beta | 8082 | 5179 | sale_beta | 2 |
| Gamma | 8083 | 5177 | sale_gamma | 3 |

## 技术栈

- **后端：** Go 1.22 + Gin + GORM + MySQL 5.6 + Redis + Casbin
- **前端 PC：** Vue 3 + TypeScript + Element Plus + Pinia + Axios
- **前端 H5：** UniApp

## 全栈开发规范

### GORM
- 以 GORM model 为唯一 schema 来源
- GORM model 变更必须同步写增量 SQL：`backend/migrations/sql/YYYYMMDD_HHMMSS_描述.sql`
- **禁止用 `Save()` 做部分更新**，必须用 `Updates(map[string]interface{})`，`Save()` 会用零值覆盖未传字段
- MySQL DSN 必须指定 `charset=utf8mb4`，不依赖运行时 SET NAMES

### 前端
- Axios interceptor 已 unwrap `response.data`，代码里用 `res.data` 不是 `res.data.data`
- Vue SFC 回调中不能用 await，先存变量
- Element Plus el-dropdown trigger 用数组 `["hover"]`
- 提 PR 前搜索 TODO / Mock / 硬编码数据，确保已清理

### 通用
- 前后端接口字段对齐，遵循已有 API 响应格式
- 路由变更需补充 harness 测试
- 遵循项目已有代码风格，不引入新的模式
- 合并分支时必须检查 `configs/config.yaml` 和 `vite.config.ts` 端口配置是否冲突
- 启动前端/后端前先清理僵尸进程：`lsof -ti:PORT | xargs kill`
- 后端启动命令：`cd backend && air`（热加载）

## 业务设计决策

- **充值中心余额**：绑定到充值中心 ID（center_id），不绑定到个人用户
- **用户登录**：手机号(phone)作为登录账号，username 作为显示名称
- **用户角色**：super_admin > hq_admin > finance > center_admin > operator

## 测试账号

- 主分支：13800000000 / 123456
- Alpha：13900000001 / 123456
