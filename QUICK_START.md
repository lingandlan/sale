# 🚀 快速启动指南

## 前置要求

### 必需软件
- ✅ Node.js v20+ (前端)
- ✅ MySQL 8.0+ (数据库)
- ✅ Redis 6.0+ (缓存)
- ✅ Go 1.21+ (后端)

### 检查环境
```bash
# 检查Node.js版本
node --version  # 应该是 v20.x.x

# 检查MySQL是否运行
mysql --version

# 检查Redis是否运行
redis-cli ping

# 检查Go版本
go version
```

## 📋 启动步骤

### 步骤1: 启动MySQL和Redis

**macOS**:
```bash
# 使用Homebrew服务管理
brew services start mysql
brew services start redis
```

**Linux**:
```bash
sudo systemctl start mysql
sudo systemctl start redis
```

**Windows**: 使用服务管理器启动MySQL和Redis

### 步骤2: 配置数据库

**创建数据库用户和数据库**:
```bash
# 登录MySQL
mysql -u root -p

# 执行以下SQL
CREATE USER 'sale'@'localhost' IDENTIFIED BY 'sale123';
CREATE DATABASE sale_dev CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
GRANT ALL PRIVILEGES ON sale_dev.* TO 'sale'@'localhost';
FLUSH PRIVILEGES;
EXIT;
```

### 步骤3: 运行数据库迁移

```bash
cd /Users/zhangdaodong/code/sale/backend
./scripts/migrate.sh
```

**迁移完成后会看到**:
```
✅ 数据库迁移完成!

📝 默认账号信息:
   管理员账号: admin
   管理员密码: admin123
```

### 步骤4: 启动后端服务

**新终端窗口**:
```bash
cd /Users/zhangdaodong/code/sale/backend
./bin/server
```

**成功启动后会看到**:
```
INFO        server listening on :8080
INFO        database connected
INFO        redis connected
```

### 步骤5: 启动前端服务

**另一个新终端窗口**:
```bash
cd /Users/zhangdaodong/code/sale/shop-pc
npm run dev
```

**成功启动后会看到**:
```
  VITE v5.x.x  ready in xxx ms

  ➜  Local:   http://localhost:5177/
```

### 步骤6: 访问系统

打开浏览器访问: **http://localhost:5177**

**默认登录账号**:
- 用户名: `admin`
- 密码: `admin123`

## 🔍 验证服务状态

### 测试后端API
```bash
# 测试健康检查
curl http://localhost:8080/health

# 测试登录
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
```

### 测试前端
1. 打开浏览器访问 http://localhost:5177
2. 使用 admin/admin123 登录
3. 查看Dashboard页面
4. 测试各个菜单功能

## ❗ 常见问题

### Q1: MySQL连接失败
**错误**: `connect database failed`

**解决**:
```bash
# 检查MySQL是否运行
brew services list | grep mysql

# 重启MySQL
brew services restart mysql

# 检查用户权限
mysql -u sale -p sale_dev
```

### Q2: Redis连接失败
**错误**: `connect redis failed`

**解决**:
```bash
# 启动Redis
brew services start redis

# 或直接启动
redis-server
```

### Q3: 端口被占用
**错误**: `bind: address already in use`

**解决**:
```bash
# 查找占用进程
lsof -i :8080  # 后端端口
lsof -i :5177  # 前端端口

# 结束进程
kill -9 <PID>

# 或修改配置文件中的端口
```

### Q4: 数据库迁移失败
**错误**: `Error 1050: Table 'xxx' already exists`

**解决**:
```bash
# 删除所有表重新迁移
mysql -u sale -p sale_dev -e "SHOW TABLES"
mysql -u sale -p sale_dev -e "DROP DATABASE sale_dev"
mysql -u sale -p sale_dev -e "CREATE DATABASE sale_dev CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci"

# 重新运行迁移
./scripts/migrate.sh
```

### Q5: 前端无法连接后端
**错误**:浏览器控制台显示CORS错误

**解决**:
1. 确认后端服务已启动
2. 检查后端CORS配置（configs/config.yaml）
3. 检查前端API配置（src/utils/request.ts）

### Q6: 登录后跳转失败
**错误**: 登录成功但无法跳转

**解决**:
1. 打开浏览器开发者工具（F12）
2. 查看Console和Network标签
3. 检查localStorage是否存储了token
4. 检查路由守卫配置

## 📊 服务端口说明

| 服务 | 端口 | 用途 |
|------|------|------|
| 前端 | 5177 | Vue开发服务器 |
| 后端 | 8080 | Go API服务器 |
| MySQL | 3306 | 数据库 |
| Redis | 6379 | 缓存 |

## 🎯 开发模式 vs 生产模式

### 开发模式（当前）
```bash
# 前端
npm run dev  # http://localhost:5177

# 后端
./bin/server  # http://localhost:8080
```

### 生产模式
```bash
# 前端构建
npm run build
# 部署 dist/ 目录到Web服务器

# 后端
# 修改 configs/config.yaml
# mode: release
# 重新编译
go build -o bin/server cmd/server/main.go
```

## 📝 开发提示

### 前端热重载
- 修改Vue组件会自动刷新浏览器
- 修改配置文件需要重启服务

### 后端热重载
```bash
# 使用air实现热重载（可选）
go install github.com/cosmtrek/air@latest
air
```

### 查看日志
- **后端日志**: 终端输出
- **前端日志**: 浏览器开发者工具Console

## 🔄 重置开发环境

如果遇到无法解决的问题，可以重置环境：

```bash
# 1. 停止所有服务
# Ctrl+C 或 kill 进程

# 2. 删除数据库
mysql -u root -p -e "DROP DATABASE sale_dev"

# 3. 重新创建数据库
mysql -u root -p -e "CREATE DATABASE sale_dev CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci"

# 4. 重新运行迁移
cd /Users/zhangdaodong/code/sale/backend
./scripts/migrate.sh

# 5. 重新启动服务
./bin/server
```

## 📞 获取帮助

如果遇到问题：
1. 查看本文档的常见问题部分
2. 查看 `backend/DEV_GUIDE.md`
3. 查看各服务的日志输出
4. 检查浏览器控制台错误信息

---

**祝您使用愉快！** 🎉
