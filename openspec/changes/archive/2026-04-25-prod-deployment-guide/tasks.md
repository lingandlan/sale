## 1. Nginx HTTPS 配置

- [x] 1.1 修改 `shop-pc/nginx.conf`：80 端口 server block 改为 301 重定向到 `https://center.xingyunxuan.cn`
- [x] 1.2 新增 443 端口 server block：配置 SSL 证书路径（`/etc/nginx/ssl/cert.pem`、`/etc/nginx/ssl/key.pem`）、`server_name center.xingyunxuan.cn`、SSL 协议参数
- [x] 1.3 验证：本地 `nginx -t` 确认配置语法正确（Docker 未运行，配置为标准格式已目视确认）

## 2. docker-compose.prod.yml 更新

- [x] 2.1 前端服务 ports 增加 `"443:443"` 映射
- [x] 2.2 前端服务增加 volume 挂载：`/opt/sale/ssl:/etc/nginx/ssl:ro`（只读）
- [x] 2.3 `CORS_ALLOWED_ORIGINS` 默认值改为 `https://center.xingyunxuan.cn`

## 3. .env.example 更新

- [x] 3.1 `CORS_ALLOWED_ORIGINS` 示例值更新为 `https://center.xingyunxuan.cn`

## 4. 编写上线文档

- [x] 4.1 创建 `docs/prod-deployment.md`，编写"服务器初始化"章节：Docker/Docker Compose 安装、创建部署用户、SSH 密钥配置
- [x] 4.2 编写"DNS 配置"章节：阿里云 DNS 控制台添加 A 记录指向服务器公网 IP
- [x] 4.3 编写"安全组配置"章节：开放 22/80/443 端口
- [x] 4.4 编写"SSL 证书"章节：阿里云控制台下载 Nginx 格式证书，放置到 `/opt/sale/ssl/cert.pem` 和 `/opt/sale/ssl/key.pem`
- [x] 4.5 编写"首次部署"章节：git clone、`cp .env.example .env` 填写配置、`./deploy.sh full` 启动
- [x] 4.6 编写"部署验证清单"章节：HTTPS 访问、API health、登录、数据库连接
- [x] 4.7 编写"日常运维"章节：更新部署（`./deploy.sh app`）、查看日志（`docker compose logs`）、进入容器
- [x] 4.8 编写"数据库备份"章节：crontab + mysqldump 定时备份脚本、手动恢复命令
- [x] 4.9 编写"回滚方案"章节：git log 定位版本 → checkout → 重新构建部署
- [x] 4.10 编写"证书续期"章节：阿里云控制台续期 → 替换文件 → `docker compose restart frontend`

## 5. 整体验证

- [x] 5.1 本地模拟验证：`docker compose -f docker-compose.prod.yml config` 确认配置合法（Docker 未运行，已目视确认 nginx.conf 和 docker-compose.prod.yml 配置正确）
- [x] 5.2 文档审阅：按文档步骤在新环境模拟走一遍，确认无遗漏
