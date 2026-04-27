## Context

项目已完成开发，部署方案采用 Docker Compose（`docker-compose.prod.yml`），后端 Go + 前端 Vue 均容器化。当前 Nginx 配置只监听 80 端口，`server_name` 为通配符 `_`。生产环境需要绑定域名 `center.xingyunxuan.cn` 并启用 HTTPS。

阿里云 ECS 服务器，操作系统 Linux，已安装 Docker + Docker Compose。SSL 证书使用阿里云免费/付费证书，下载 Nginx 格式（`.pem` + `.key`）。

## Goals / Non-Goals

**Goals:**
- 编写完整的上线操作文档，团队成员可按文档从零完成生产部署
- Nginx 配置 HTTPS（443）+ HTTP→HTTPS 自动重定向（80）
- SSL 证书正确挂载到前端容器
- docker-compose.prod.yml 适配 HTTPS 部署
- 包含 DNS、安全组、验证、回滚、备份等运维操作

**Non-Goals:**
- 不做 CI/CD 流水线改造（已有 flow.yaml）
- 不做自动证书续期（阿里云证书需手动续期，文档中说明）
- 不做多机集群/负载均衡（当前单机部署足够）
- 不做监控告警系统搭建

## Decisions

### D1: HTTPS 终止在容器内 Nginx

**方案**: 在前端容器内的 Nginx 直接处理 SSL，监听 80 + 443。

**替代方案**: 宿主机装 Nginx 做 SSL 终止，反向代理到容器 — 放弃，因为多一层代理增加复杂度，且容器内 Nginx 已有完整配置（SPA fallback、API 代理、静态资源缓存）。

**实现**: 修改 `shop-pc/nginx.conf`，增加 443 server block + SSL 配置，80 端口做 301 重定向。证书文件通过 docker volume 挂载到容器内 `/etc/nginx/ssl/`。

### D2: 证书文件管理

**方案**: 证书存放在服务器 `/opt/sale/ssl/` 目录，通过 `docker-compose.prod.yml` 的 volumes 挂载到容器。

**续期**: 阿里云证书到期前在控制台续期，下载新证书替换文件，`docker compose restart frontend` 即可。

### D3: 前端 Dockerfile 不变

当前 `shop-pc/Dockerfile` 是多阶段构建，最终阶段 `FROM nginx:alpine`，`COPY nginx.conf`。SSL 相关配置全在 `nginx.conf` 中，证书通过 volume 挂载，Dockerfile 不需要改动。

### D4: 文档结构

按操作顺序组织：
1. 服务器初始化（Docker 安装、用户创建）
2. DNS 配置
3. 安全组/防火墙
4. SSL 证书下载与放置
5. 代码部署（git clone + .env 配置 + docker compose up）
6. 验证清单
7. 日常运维（更新部署、查看日志、数据库备份）
8. 回滚方案

### D5: 数据库备份

使用 `mysqldump` 定时任务，备份文件保留 7 天，存放在 `/opt/sale/backups/`。通过 crontab 配置每日凌晨 3 点自动备份。

## Risks / Trade-offs

- **[证书过期]** → 文档中注明到期时间，建议设置日历提醒。阿里云免费证书有效期 3 个月，付费证书 1 年。
- **[容器内挂载证书路径写死]** → 在 nginx.conf 中用 `/etc/nginx/ssl/` 固定路径，docker-compose volume 映射到宿主机 `/opt/sale/ssl/`，路径变更只需改 compose 文件。
- **[单机部署无高可用]** → 当前阶段可接受。如后续需要高可用，可迁移到阿里云 SLB + 多 ECS。
- **[MySQL 数据持久化]** → 已通过 `mysql-data` volume 持久化，但需定期备份到其他位置（如 OSS），防止单机故障丢数据。
