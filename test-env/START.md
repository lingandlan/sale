# 测试环境启动指南

## 📋 目录结构

```
sale/                              # 项目根目录
├── shop-pc/                       # ⬅️ 前端开发版 (Claude Code 在用)
├── backend/                       # ⬅️ 后端开发版 (Claude Code 在用)
├── test-env/                      # ⬅️ 测试专用环境 (本目录)
│   ├── backend/                   # 后端测试版 (Git 拉取)
│   └── frontend/                  # 前端测试版 (Git 拉取)
```

---

## 🎯 分离策略

- **开发目录** (`shop-pc/`, `backend/`): 留给 Claude Code 开发
- **测试目录** (`test-env/`): 独立测试，不受开发影响

---

## 🚀 快速启动

### 1. 初始化测试环境 (首次)

```bash
cd test-env

# 克隆后端代码
git clone ../backend.git backend
# 或者 git clone <remote-url> backend

# 克隆前端代码
git clone ../shop-pc.git frontend
# 或者 git clone <remote-url> frontend
```

### 2. 启动后端测试服务 (端口 8081)

```bash
cd test-env/backend
go run cmd/server/main.go -port 8081
```

### 3. 启动前端测试服务 (端口 5176)

```bash
cd test-env/frontend

# 如需修改 API 指向测试后端:
# 编辑 vite.config.ts 或 .env 文件
# VITE_API_BASE=http://localhost:8081/api/v1

npm run dev -- --port 5176
```

### 4. 运行测试

```bash
agent-browser open http://localhost:5176/login
```

---

## 🔄 同步最新代码

每次测试前，拉取最新代码：

```bash
# 后端更新
cd test-env/backend
git pull

# 前端更新
cd test-env/frontend
git pull
```

---

## 📝 测试账号

| 账号 | 密码 | 角色 |
|-----|------|-----|
| 13800000000 | 123456 | 超级管理员 |

---

## 🔗 正确页面 URL (前端 → 后端)

| 页面 | 前端路径 | 后端API |
|-----|---------|---------|
| 首页仪表盘 | /dashboard | /api/v1/dashboard/statistics |
| B端充值申请 | /recharge/b-apply | /api/v1/recharge/b-apply |
| B端充值审批 | /recharge/b-approval | /api/v1/recharge/b-approval |
| C端充值录入 | /recharge/c-entry | /api/v1/recharge/c-entry |
| 充值记录 | /recharge/records | /api/v1/recharge/records |
| 充值中心管理 | /center/manage | /api/v1/center |
| 门店卡管理 | /card/manage | /api/v1/card |
| 用户管理 | /user/manage | /api/v1/admin/users |

---

## 🐛 待修复 Bug

### BUG-001: 退出登录功能失效

- **位置**: `frontend/src/layouts/Sidebar.vue:147-159`
- **现象**: 点击「退出登录」后不跳转
- **状态**: ❌ 待修复

---

## 📊 测试历史

### 2026-04-11 测试结果

- **正常页面**: 8/8 (除退出登录外全部正常)
- **失败页面**: 1 (退出登录)

---

## 🔧 常用命令

```bash
# 运行 Go 单元测试
cd test-env/backend
go test ./internal/... -v

# 运行 Playwright E2E 测试
cd test-env/backend/tests_complete/e2e
npm install
npx playwright test

# API 一致性检查
cd test-env/backend
go run ./tests_complete/api_consistency_check.go
```

---

## ⚠️ 注意事项

1. **Git 克隆**: 测试环境是独立的 Git 仓库副本
2. **端口区分**: 
   - 开发后端: 8080
   - 测试后端: 8081
   - 开发前端: 5175
   - 测试前端: 5176
3. **独立环境**: Claude Code 在 `backend/` 和 `shop-pc/` 开发，不影响测试环境
4. **同步代码**: 测试前用 `git pull` 更新测试环境的代码