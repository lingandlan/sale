# 太积堂系统 - 开发运行指南

## 📋 前置条件

### 必需软件

1. **Node.js**: v18+ （推荐 v24）
   ```bash
   # 检查版本
   node --version

   # 如需安装 nvm 管理多版本
   curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash
   nvm install 20
   nvm use 20
   ```

2. **Go**: 1.25.6+
3. **Docker**: 用于 MySQL 和 Redis
4. **Make**: 用于后端自动化命令

### Node.js 版本说明

- ✅ **Node.js v24**：推荐，完全兼容
- ✅ **Node.js v20**：兼容版本
- ⚠️ **Node.js < v18**：不支持

---

## 🚀 快速启动

### 一键启动（推荐）

```bash
# 1. 启动数据库（首次需要）
cd /Users/zhangdaodong/code/sale/backend
docker-compose up -d

# 2. 启动后端
make dev

# 3. 新终端窗口，启动前端
cd /Users/zhangdaodong/code/sale/shop-h5
source ~/.nvm/nvm.sh && nvm use 20
npm run dev
```

### 详细启动步骤

#### 1. 数据库服务

```bash
cd /Users/zhangdaodong/code/sale/backend
docker-compose up -d

# 验证服务状态
docker-compose ps

# 应该看到以下服务运行中：
# - MySQL: localhost:3306
# - Redis: localhost:6379
```

**首次运行需要初始化数据库**：
```bash
make migrate-up
```

#### 2. 后端服务

```bash
cd /Users/zhangdaodong/code/sale/backend

# 方式1：使用 Make（推荐）
make dev

# 方式2：直接运行
go run ./cmd/server/main.go
```

**验证后端**：
- 地址：http://localhost:8080
- 健康检查：`curl http://localhost:8080/health`

#### 3. 前端服务

```bash
cd /Users/zhangdaodong/code/sale/shop-h5

# 启动开发服务器（Node.js v18+ 均可）
npm run dev
```

**验证前端**：
- 地址：http://localhost:5174
- 浏览器打开应该看到登录页面

---

## 🔧 开发工具

### Playwright（浏览器自动化）

**版本要求**：agent-browser >= 0.25.0（推荐 0.25.3+）

**安装**：
```bash
nvm use 20
npx playwright install chromium
```

**使用方式**：
```bash
# 方式1：使用 npx（推荐，自动使用最新版）
npx agent-browser@latest open http://localhost:5174/
npx agent-browser@latest screenshot /tmp/screenshot.png
npx agent-browser@latest console

# 方式2：更新全局版本后直接使用
sudo npm update -g agent-browser@latest
agent-browser open http://localhost:5174/

# 方式3：使用脚本
node /path/to/test-script.js
```

**⚠️ 注意**：
- 全局安装的 0.9.2 版本过旧，不兼容新版 Playwright
- 推荐使用 `npx agent-browser@latest` 确保使用最新稳定版（0.25.3+）

### 后端调试工具

**密码重置**：
```bash
cd /Users/zhangdaodong/code/sale/backend
go run ./cmd/resetpwd/main.go -phone 13800000000 -password 新密码
```

**诊断工具**：
```bash
go run ./cmd/diagnose/main.go
```

---

## 📦 服务端口

| 服务 | 地址 | 说明 |
|------|------|------|
| 前端 | http://localhost:5174 | UniApp H5 |
| 后端 | http://localhost:8080 | Go API |
| MySQL | localhost:3306 | 数据库 |
| Redis | localhost:6379 | 缓存 |
| Apifox | https://app.apifox.com | API 文档 |

---

## 🔑 测试账号

```
手机号: 13800000000
密码: Test123456
角色: 超级管理员 (super_admin)
```

**重置密码**：
```bash
go run ./cmd/resetpwd/main.go -phone 13800000000 -password Test123456
```

---

## 🛠️ 常用命令

### 后端命令

```bash
cd /Users/zhangdaodong/code/sale/backend

make dev           # 启动开发服务器（Air热重载）
make build         # 构建生产二进制
make test          # 运行测试
make lint          # 代码检查
make migrate-up    # 执行数据库迁移
make migrate-down  # 回滚数据库迁移
make fmt           # 格式化代码
```

### 前端命令

```bash
cd /Users/zhangdaodong/code/sale/shop-h5

npm run dev        # 启动开发服务器
npm run build      # 构建生产版本
npm run lint       # 代码检查
npm run type-check # TypeScript 类型检查
```

### Docker 命令

```bash
cd /Users/zhangdaodong/code/sale/backend

docker-compose up -d          # 启动服务
docker-compose down           # 停止服务
docker-compose ps             # 查看状态
docker-compose logs -f mysql  # 查看 MySQL 日志
docker-compose logs -f redis  # 查看 Redis 日志
```

---

## 🐛 常见问题排查

### 问题1：前端页面空白

**症状**：浏览器打开 http://localhost:5174 显示空白页面

**原因**：Vue 应用未正确挂载

**解决方案**：
1. 检查浏览器控制台错误
2. 确认 Node.js 版本 >= v18：`node --version`
3. 清除缓存重启：
   ```bash
   rm -rf node_modules/.vite
   npm run dev
   ```

### 问题2：接口返回 401

**症状**：登录后访问接口返回 `401 Unauthorized`

**原因**：Token 过期或无效

**解决方案**：
1. 检查后端服务是否运行：`curl http://localhost:8080/health`
2. 清除浏览器 localStorage，重新登录
3. 检查 JWT 配置：`configs/config.dev.yaml`

### 问题3：数据库连接失败

**症状**：后端日志显示 `database connection error`

**原因**：Docker 服务未启动

**解决方案**：
```bash
cd /Users/zhangdaodong/code/sale/backend
docker-compose up -d
docker-compose ps  # 确认服务运行中
```

### 问题4：端口被占用

**症状**：启动时提示 `port 5174/8080 already in use`

**解决方案**：
```bash
# 查找占用进程
lsof -ti:5174  # 前端
lsof -ti:8080  # 后端

# 杀死进程
lsof -ti:5174 | xargs kill -9
lsof -ti:8080 | xargs kill -9
```

### 问题5：agent-browser 版本过低

**症状**：`Executable doesn't exist at chromium_headless_shell-1208`

**原因**：
- 全局安装的 agent-browser 0.9.2 版本过旧
- 与新版 Playwright v1217 不兼容

**解决方案**：
```bash
# 方案1：使用 npx（推荐）
npx agent-browser@latest open http://localhost:5174/

# 方案2：更新全局版本
sudo npm update -g agent-browser@latest

# 验证版本
agent-browser --version  # 应显示 0.25.3 或更高
```

**版本对照**：
| agent-browser | Playwright | 状态 |
|---------------|------------|------|
| 0.9.2 | 1208 | ❌ 过旧 |
| 0.25.3+ | 1217 | ✅ 兼容 |

---

## 📂 关键配置文件

### 后端配置

**configs/config.dev.yaml** - 开发环境配置
```yaml
server:
  port: 8080
  mode: debug

database:
  host: localhost
  port: 3306
  name: sale_dev
  username: root
  password: root
  tls: false

jwt:
  secret: your-secret-key
  expire: 24h
```

### 前端配置

**vite.config.js** - Vite 配置
```javascript
export default defineConfig({
  plugins: [uni()],
  resolve: {
    alias: {
      '@': '/src'
    }
  },
  server: {
    port: 5174,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true
      }
    }
  }
})
```

---

## 🧪 测试流程

### 1. 手动测试登录

```bash
# 1. 启动所有服务
# 见"快速启动"部分

# 2. 使用 agent-browser 测试
npx agent-browser@latest open http://localhost:5174/

# 3. 填写表单
# 手机号: 13800000000
# 密码: Test123456

# 4. 提交登录
```

### 2. API 测试

```bash
# 使用 curl 测试
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"phone":"13800000000","password":"Test123456"}'

# 应该返回：
# {"code":0,"message":"success","data":{"access_token":"...","refresh_token":"..."}}
```

### 3. 前端集成测试

```bash
# 1. 使用 Playwright 自动化测试
npx agent-browser@latest eval "document.querySelector('input[type=\"tel\"]').value = '13800000000'"

# 2. 截图验证
npx agent-browser@latest screenshot /tmp/test-login.png
```

---

## 📝 开发工作流

### 新功能开发流程

1. **API 设计**（Apifox）
   - 在 Apifox 中定义接口
   - 配置 Mock 数据
   - 前后端确认

2. **后端开发**
   ```bash
   # Model → Repository → Service → Handler
   make dev  # 启动后端
   make test  # 运行测试
   ```

3. **前端开发**
   ```bash
   npm run dev  # 启动前端
   npm run lint  # 代码检查
   ```

4. **联调测试**
   - 接口联调
   - 功能测试
   - Bug 修复

---

## 🔄 开发环境重置

如果遇到无法解决的依赖或缓存问题：

```bash
# 1. 停止所有服务
cd /Users/zhangdaodong/code/sale/backend
docker-compose down
lsof -ti:8080 | xargs kill -9

cd /Users/zhangdaodong/code/sale/shop-h5
lsof -ti:5174 | xargs kill -9

# 2. 清理前端缓存
rm -rf node_modules/.vite
rm -rf dist
rm -rf unpackage

# 3. 重新启动
# 见"快速启动"部分
```

---

## 📚 相关文档

- **后端架构**: `/Users/zhangdaodong/code/sale/backend/README.md`
- **前端设计**: `/Users/zhangdaodong/code/sale/design/DESIGN-SYSTEM.md`
- **API 文档**: https://app.apifox.com/project/8082766
- **登录问题解决**: `/Users/zhangdaodong/code/sale/backend/LOGIN_401_SOLUTION.md`
- **密码重置工具**: `/Users/zhangdaodong/code/sale/backend/PASSWORD_RESET_TOOL.md`

---

## 🆘 获取帮助

### 检查服务状态

```bash
# 一键检查所有服务
cd /Users/zhangdaodong/code/sale/backend
docker-compose ps
lsof -ti:8080 && echo "✅ 后端运行中" || echo "❌ 后端未运行"

cd /Users/zhangdaodong/code/sale/shop-h5
lsof -ti:5174 && echo "✅ 前端运行中" || echo "❌ 前端未运行"
```

### 日志查看

```bash
# 后端日志
tail -f logs/app.log

# Docker 日志
docker-compose logs -f

# 前端控制台
# 浏览器 F12 → Console
```

---

**最后更新**: 2026-04-10
**uni-app 版本**: 3.0.0-5000720260408001 (兼容 Node.js v18+)
**Node.js 版本要求**: >= v18 (推荐 v24)
**状态**: ✅ 可正常运行和调试
