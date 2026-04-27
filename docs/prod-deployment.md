# 太积堂充值与门店管理系统 — 生产环境部署文档

## 1. 服务器初始化

### 1.1 系统要求

- **操作系统**: Ubuntu 22.04 / CentOS 7+ / Alibaba Cloud Linux 3
- **最低配置**: 2 核 4GB 内存，40GB 系统盘
- **推荐配置**: 4 核 8GB 内存，100GB 系统盘

### 1.2 安装 Docker

```bash
# Ubuntu / Debian
curl -fsSL https://get.docker.com | sh
sudo usermod -aG docker $USER

# 退出 SSH 重新登录，使 docker 组生效
exit

# 验证
docker --version
docker compose version
```

### 1.3 创建部署用户（可选）

```bash
sudo useradd -m -s /bin/bash deploy
sudo usermod -aG docker deploy

# 配置 SSH 密钥登录
sudo mkdir -p /home/deploy/.ssh
sudo cp ~/.ssh/authorized_keys /home/deploy/.ssh/
sudo chown -R deploy:deploy /home/deploy/.ssh
sudo chmod 700 /home/deploy/.ssh
sudo chmod 600 /home/deploy/.ssh/authorized_keys
```

---

## 2. DNS 配置

在阿里云 DNS 控制台添加 A 记录：

| 记录类型 | 主机记录 | 记录值 | TTL |
|---------|---------|--------|-----|
| A | center | 服务器公网 IP | 10 分钟 |

**验证**:
```bash
dig center.xingyunxuan.cn
# 确认 ANSWER SECTION 中包含服务器公网 IP
```

---

## 3. 安全组配置

在阿里云 ECS 安全组中添加以下入方向规则：

| 端口 | 协议 | 授权对象 | 说明 |
|------|------|---------|------|
| 22 | TCP | 你的 IP/32 | SSH 登录 |
| 80 | TCP | 0.0.0.0/0 | HTTP（重定向到 HTTPS） |
| 443 | TCP | 0.0.0.0/0 | HTTPS |

> **注意**: MySQL 3306 和 Redis 6379 仅容器内部通信，**不要**对外暴露。

---

## 4. SSL 证书

### 4.1 申请证书

1. 登录 [阿里云 SSL 证书控制台](https://yundun.console.aliyun.com/?p=cas)
2. 购买/申请免费证书，域名填 `center.xingyunxuan.cn`
3. 按提示完成 DNS 验证
4. 证书签发后，点击"下载"，选择 **Nginx** 格式

### 4.2 放置证书文件

下载后解压会得到两个文件（文件名可能含域名前缀），重命名并放到服务器：

```bash
# 在服务器上操作
sudo mkdir -p /opt/sale/ssl

# 上传证书（在本地执行）
scp <下载的证书文件>.pem root@<服务器IP>:/opt/sale/ssl/cert.pem
scp <下载的密钥文件>.key root@<服务器IP>:/opt/sale/ssl/cert.key

# 在服务器上设置权限
sudo chmod 644 /opt/sale/ssl/cert.pem
sudo chmod 600 /opt/sale/ssl/cert.key
```

> **注意**：证书文件名必须为 `cert.pem` 和 `cert.key`，与 nginx.conf 中的配置对应。

---

## 5. 首次部署

### 5.1 克隆代码

添加主机组 https://flow.aliyun.com/setting/deploy 
在服务上手动执行runner，
复制系统给我的脚本直接到服务器执行便可。

```bash
cd /opt
sudo git clone git@codeup.aliyun.com:69e9953b7b6e0a0129630f7f/sale.git sale
cd /opt/sale
```

### 5.2 配置环境变量

```bash
cp .env.example .env
```

编辑 `.env`，填写真实值：

```env
# 必须修改
JWT_SECRET=<openssl rand -hex 32 生成的随机值>
MYSQL_ROOT_PASSWORD=<强密码>
MYSQL_PASSWORD=<与 DB_PASSWORD 保持一致>
DB_PASSWORD=<MySQL 用户密码>

# WSY 商城 API（如需对接）
MALL_APP_ID=<从 WSY 后台获取>
MALL_APP_SECRET=<从 WSY 后台获取>
MALL_CUSTOMER_ID=<从 WSY 后台获取>

# CORS（已预填，无需修改）
CORS_ALLOWED_ORIGINS=https://center.xingyunxuan.cn
```

生成 JWT 密钥：
```bash
openssl rand -hex 32
```

### 5.3 启动服务

```bash
./deploy.sh full
```

脚本会自动：
1. 启动 MySQL 和 Redis 容器
2. 等待 MySQL 就绪（最多 60 秒）
3. 构建并启动 backend 和 frontend 容器（**使用 docker-compose.prod.yml**，确保 SSL 卷挂载生效）
4. backend 容器启动时自动执行数据库 migration 和超管初始化

### 5.4 确认服务状态

```bash
docker compose -f docker-compose.prod.yml ps
```

预期输出 4 个服务均为 `running` / `Up` 状态。

---

## 6. 部署验证清单

部署完成后逐项验证：

- [ ] **HTTPS 访问**: 浏览器打开 `https://center.xingyunxuan.cn`，页面正常加载，浏览器显示安全锁
- [ ] **HTTP 重定向**: 访问 `http://center.xingyunxuan.cn`，自动跳转到 HTTPS
- [ ] **API Health**: `curl https://center.xingyunxuan.cn/api/v1/health` 返回 `{"status":"ok"}`
- [ ] **登录功能**: 使用超管账号登录（手机号 + 密码）
- [ ] **数据库连接**: 登录后能看到 Dashboard 数据（非空页面）
- [ ] **上传功能**: 测试图片上传，确认 `/uploads/` 路径可访问

---

## 7. 日常运维

### 7.1 更新部署

代码推送到 Codeup master 分支后，CI/CD 流水线会自动部署。

手动更新：
```bash
cd /opt/sale
git pull

# 强制重新构建（避免缓存导致配置不更新）
docker compose -f docker-compose.prod.yml build --no-cache frontend backend

# 重启前后端容器
docker compose -f docker-compose.prod.yml up -d --no-deps --force-recreate backend frontend
```

> **重要**：生产环境**必须使用 `-f docker-compose.prod.yml`**，不要使用默认的 `docker-compose.yml`（后者没有 SSL 证书卷挂载，会导致前端容器无法启动）。

`deploy.sh app` 只重建 backend 和 frontend 容器，**不影响数据库和 Redis**。

### 7.2 查看日志

```bash
# 查看所有服务日志
docker compose -f docker-compose.prod.yml logs

# 实时跟踪后端日志
docker compose -f docker-compose.prod.yml logs -f backend

# 查看最近 100 行前端日志
docker compose -f docker-compose.prod.yml logs --tail 100 frontend

# 进入后端容器排查
docker exec -it sale-backend sh
```

### 7.3 重启单个服务

> **重要**：所有 `docker compose` 命令必须加 `-f docker-compose.prod.yml`，否则不会挂载 SSL 证书，前端无法启动。

```bash
# 只重启后端
docker compose -f docker-compose.prod.yml restart backend

# 只重启前端（如更新了 nginx 配置）
docker compose -f docker-compose.prod.yml restart frontend
```

> 如果前端仍然无法启动，先确认卷挂载是否生效：`docker inspect sale-frontend --format '{{json .Mounts}}'`

---

## 8. 数据库备份

### 8.1 自动备份脚本

创建备份脚本 `/opt/sale/scripts/backup-db.sh`：

```bash
#!/bin/bash
BACKUP_DIR="/opt/sale/backups"
DB_NAME="sale_dev"
DB_USER="sale"
DB_PASS="<你的数据库密码>"
RETAIN_DAYS=7

mkdir -p $BACKUP_DIR
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="$BACKUP_DIR/${DB_NAME}_${TIMESTAMP}.sql.gz"

docker exec sale-mysql mysqldump -u$DB_USER -p$DB_PASS $DB_NAME | gzip > $BACKUP_FILE

# 清理过期备份
find $BACKUP_DIR -name "*.sql.gz" -mtime +$RETAIN_DAYS -delete

echo "[$(date)] 备份完成: $BACKUP_FILE"
```

```bash
chmod +x /opt/sale/scripts/backup-db.sh
```

### 8.2 配置定时任务

```bash
crontab -e
```

添加以下行（每日凌晨 3 点执行备份）：

```
0 3 * * * /opt/sale/scripts/backup-db.sh >> /opt/sale/backups/backup.log 2>&1
```

### 8.3 手动恢复备份

```bash
# 解压备份文件
gunzip /opt/sale/backups/sale_dev_20260424_030000.sql.gz

# 恢复到数据库
docker exec -i sale-mysql mysql -u sale -p<密码> sale_dev < /opt/sale/backups/sale_dev_20260424_030000.sql
```

---

## 9. 回滚方案

### 9.1 查看部署历史

```bash
cd /opt/sale
git log --oneline -10
```

### 9.2 回滚到指定版本

```bash
# 回滚到上一个版本
git checkout <commit-hash>

# 重新构建并部署
docker compose -f docker-compose.prod.yml build --no-cache backend frontend
docker compose -f docker-compose.prod.yml up -d --no-deps --force-recreate backend frontend

# 验证回滚成功后，将 HEAD 指回最新（可选）
git checkout master
```

> **注意**: 回滚只影响应用代码，不会回退数据库。如果 migration 已执行，数据库 schema 不会自动回退。构建时建议加 `--no-cache` 避免缓存导致旧配置不生效。

---

## 10. SSL 证书续期

### 10.1 续期步骤

1. 登录 [阿里云 SSL 证书控制台](https://yundun.console.aliyun.com/?p=cas)
2. 找到即将过期的证书，点击"续费/续期"
3. 完成域名验证，等待证书重新签发
4. 下载新的 Nginx 格式证书

### 10.2 替换证书

```bash
# 上传新证书
scp <新证书>.pem root@<服务器IP>:/opt/sale/ssl/cert.pem
scp <新密钥>.key root@<服务器IP>:/opt/sale/ssl/cert.key

# 重启前端容器加载新证书
cd /opt/sale
docker compose -f docker-compose.prod.yml restart frontend

# 验证新证书生效
curl -vI https://center.xingyunxuan.cn 2>&1 | grep "expire"
```

> **提醒**: 阿里云免费证书有效期 3 个月，建议在到期前 1 周续期。可在日历中设置提醒。
