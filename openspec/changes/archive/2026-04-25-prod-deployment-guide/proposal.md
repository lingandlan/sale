## Why

项目已完成功能开发，需要编写正式环境上线文档，确保团队成员能按文档从零完成生产部署。域名 `center.xingyunxuan.cn`，需配置 HTTPS（阿里云 SSL 证书）。

## What Changes

- **Nginx HTTPS 反向代理配置**：容器内前端 Nginx 改为 80/443 双端口，SSL 证书挂载到容器内，`server_name` 改为 `center.xingyunxuan.cn`
- **docker-compose.prod.yml 更新**：前端容器暴露 443 端口，挂载 SSL 证书 volume，CORS 配置改为 `https://center.xingyunxuan.cn`
- **部署文档**：完整上线操作指南，涵盖服务器初始化、DNS、SSL、Docker 部署、验证、回滚

## Capabilities

### New Capabilities

- `https-nginx-config`: Nginx HTTPS 配置 — SSL 证书挂载、server_name、HTTP→HTTPS 重定向
- `prod-deploy-doc`: 生产环境上线文档 — 从服务器初始化到上线的完整操作指南

### Modified Capabilities

（无已有 spec 需要修改）

## Impact

- **前端 Nginx 配置**: `shop-pc/nginx.conf` 增加 SSL 配置块和 HTTP→HTTPS 301 重定向
- **docker-compose.prod.yml**: 前端服务增加 443 端口映射和证书 volume 挂载
- **.env.example**: `CORS_ALLOWED_ORIGINS` 默认值更新为 `https://center.xingyunxuan.cn`
- **文档**: 新增 `docs/prod-deployment.md` 上线操作手册
- **服务器**: 需开放 80/443 端口，配置 DNS A 记录指向服务器公网 IP
