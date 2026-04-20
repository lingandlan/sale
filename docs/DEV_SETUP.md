# 太积堂充值与门店管理系统 - 开发环境部署指南

## 1. 系统要求

| 软件 | 最低版本 | 推荐安装方式 |
|------|---------|-------------|
| Go | 1.22+ | https://go.dev/dl/ |
| Node.js | 18+ | nvm / https://nodejs.org/ |
| MySQL | 5.6+ | Homebrew / 系统包管理器 |
| Redis | 6.0+ | Homebrew / 系统包管理器 |
| Git | 2.x | 系统包管理器 |

### macOS 推荐安装

```bash
# Homebrew（如果还没装）
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# Go
brew install go

# Node.js（推荐用 nvm 管理版本）
brew install nvm
nvm install 20
nvm use 20

# MySQL
brew install mysql
brew services start mysql

# Redis
brew install redis
brew services start redis

# air（Go 热重载工具）
go install github.com/cosmtrek/air@latest
```

### Linux (Ubuntu/Debian)

```bash
# Go
wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# Node.js
curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
sudo apt-get install -y nodejs

# MySQL
sudo apt install mysql-server
sudo systemctl start mysql

# Redis
sudo apt install redis-server
sudo systemctl start redis

# air
go install github.com/cosmtrek/air@latest
```

---

## 2. 获取代码

```bash
git clone <仓库地址> sale
cd sale
```

---

## 3. 数据库配置

### 3.1 创建 MySQL 用户和数据库

```bash
mysql -u root -p
```

```sql
-- 创建数据库
CREATE DATABASE sale_dev CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 创建用户（如果不存在）
CREATE USER 'sale'@'localhost' IDENTIFIED BY 'sale123';

-- 授权
GRANT ALL PRIVILEGES ON sale_dev.* TO 'sale'@'localhost';
FLUSH PRIVILEGES;
```

> **注意**：MySQL 8.x 使用 `caching_sha2_password` 认证插件。如果遇到认证错误：
> ```sql
> ALTER USER 'sale'@'localhost' IDENTIFIED WITH caching_sha2_password BY 'sale123';
> ```

### 3.2 运行数据库迁移

迁移通过 GORM AutoMigrate 执行，会自动创建所有表和种子数据。

```bash
cd backend
go run migrations/migrate.go
```

迁移完成后会自动创建：
- **表**：users, recharge_centers, recharge_operators, recharge_applications, c_recharges, store_cards, card_issue_records, card_transactions, casbin_rule
- **种子数据**：
  - 管理员账号：`13800000000` / `123456`（super_admin）
  - 3 个测试充值中心（北京朝阳、北京海淀、上海浦东）
  - 2 个测试操作员

### 3.3 增量 SQL 迁移

如果已有数据库需要同步到最新 schema，执行增量 SQL：

```bash
cd backend
# 按 sql/ 目录下的文件名顺序执行
mysql -u sale -psale123 sale_dev < migrations/sql/000003_add_username_to_users.up.sql
mysql -u sale -psale123 sale_dev < migrations/sql/20260415_100000_store_card_redesign.sql
mysql -u sale -psale123 sale_dev < migrations/sql/20260418_100000_user_center_id_to_string.sql
```

> 如果是全新环境，GORM AutoMigrate 已包含最新 schema，无需手动执行这些 SQL。

### 3.4 RBAC 权限数据

```bash
mysql -u sale -psale123 sale_dev < migrations/000001_create_casbin_tables.up.sql
mysql -u sale -psale123 sale_dev < migrations/000002_create_users_table.up.sql
```

> 同样，全新环境下 GORM AutoMigrate 已处理，无需额外执行。

---

## 4. Redis

确保 Redis 已启动：

```bash
# macOS
brew services start redis

# Linux
sudo systemctl start redis

# 验证
redis-cli ping
# 应返回 PONG
```

本项目 Redis 无需密码，使用默认 6379 端口，DB 0。

---

## 5. 后端配置

### 5.1 配置文件

配置文件位于 `backend/configs/config.yaml`：

```yaml
server:
  port: 8080
  mode: debug

database:
  host: localhost
  port: 3306
  user: sale
  password: sale123
  name: sale_dev

redis:
  host: localhost
  port: 6379
  password: ""
  db: 0

jwt:
  secret: your-super-secret-key-change-in-production
  expire_hours: 24
```

> 默认配置已适配本地开发环境，通常无需修改。

### 5.2 环境变量（可选）

`backend/.env` 提供了环境变量覆盖方式，默认已配置好。如需修改端口等参数可直接编辑。

---

## 6. 前端配置

### 6.1 安装依赖

```bash
cd shop-pc
npm install
```

### 6.2 环境配置

前端默认配置已在 `.env.development` 中：

```
VITE_API_BASE_URL=http://localhost:8080/api/v1
```

创建 `.env.local`（已被 gitignore）自定义端口：

```
VITE_PORT=5175
VITE_API_PORT=8080
```

> 前端 Vite 代理已配置 `/api` -> `http://localhost:8080`，无需额外设置。

---

## 7. 启动服务

### 方式一：一键启动（推荐）

```bash
# 项目根目录
./start.sh
```

脚本会自动：
1. 读取 `shop-pc/.env.local` 的端口配置
2. 清理残留进程
3. 启动后端（air 热重载）
4. 启动前端（vite dev server）

### 方式二：手动分步启动

**终端 1 - 后端：**

```bash
cd backend
air
# 或者不带热重载：go run ./cmd/server
```

**终端 2 - 前端：**

```bash
cd shop-pc
npm run dev
```

### 启动后验证

```bash
# 后端健康检查
curl http://localhost:8080/health

# 登录测试
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"phone":"13800000000","password":"123456"}'
```

---

## 8. 访问系统

| 服务 | 地址 |
|------|------|
| 前端（B端管理后台） | http://localhost:5175 |
| 后端 API | http://localhost:8080 |

**测试账号：** `13800000000` / `123456`（超级管理员）

---

## 9. 导入已有数据（可选）

如果从其他环境迁移，已有数据 dump 文件：

```bash
mysql -u sale -psale123 sale_dev < docs/sale_dev.sql
```

> 全新环境无需此步，`go run migrations/migrate.go` 会创建表结构和种子数据。

---

## 10. 端口分配表

| 环境 | 后端 | 前端 PC | Redis DB |
|------|------|---------|----------|
| Main（默认） | 8080 | 5175 | 0 |
| Alpha | 8081 | 5178 | 1 |
| Beta | 8082 | 5179 | 2 |
| Gamma | 8083 | 5177 | 3 |

> 多环境并行开发时，通过 `shop-pc/.env.local` 和环境变量 `APP_SERVER_PORT` / `APP_REDIS_DB` 隔离。

---

## 11. 常见问题

### MySQL 连接失败

```bash
# 检查 MySQL 状态
brew services list | grep mysql   # macOS
sudo systemctl status mysql       # Linux

# 测试连接
mysql -u sale -psale123 sale_dev -e "SELECT 1"
```

### Redis 连接失败

```bash
redis-cli ping  # 应返回 PONG
```

### 端口被占用

```bash
# 查找并清理
lsof -ti:8080 | xargs kill   # 后端
lsof -ti:5175 | xargs kill   # 前端
```

### 重新初始化数据库

```bash
mysql -u root -p -e "DROP DATABASE sale_dev; CREATE DATABASE sale_dev CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
cd backend
go run migrations/migrate.go
```

### 前端启动报错 `node_modules` 问题

```bash
cd shop-pc
rm -rf node_modules package-lock.json
npm install
```

---

## 12. 项目结构

```
sale/
├── backend/           # Go 后端
│   ├── cmd/server/    # 入口 main.go
│   ├── configs/       # 配置文件 (config.yaml, rbac_model.conf)
│   ├── internal/      # 业务代码 (handler/service/repository/model/router/middleware)
│   ├── pkg/           # 公共包 (errno/logger/mall/response)
│   ├── migrations/    # 数据库迁移
│   ├── .air.toml      # air 热重载配置
│   └── Makefile       # 后端构建命令
├── shop-pc/           # B端管理后台 (Vue 3 + Element Plus)
│   ├── src/           # 前端源码
│   └── vite.config.ts # Vite 配置
├── shop-h5/           # C端 H5 (UniApp，早期阶段)
├── start.sh           # 一键启动脚本
└── docs/              # PRD / 设计文档
```
