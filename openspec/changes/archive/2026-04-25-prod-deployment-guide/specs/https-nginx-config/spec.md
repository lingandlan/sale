## ADDED Requirements

### Requirement: Nginx 监听 HTTPS 443 端口
Nginx SHALL 监听 443 端口，使用 SSL 证书启用 HTTPS，`server_name` 为 `center.xingyunxuan.cn`。

#### Scenario: HTTPS 正常访问
- **WHEN** 用户通过浏览器访问 `https://center.xingyunxuan.cn`
- **THEN** 页面正常加载，浏览器显示安全锁标识，证书有效

#### Scenario: SSL 证书无效
- **WHEN** SSL 证书过期或未正确挂载
- **THEN** 浏览器提示证书错误，Nginx 不提供服务

### Requirement: HTTP 自动重定向到 HTTPS
Nginx SHALL 在 80 端口监听 HTTP 请求，并返回 301 永久重定向到 HTTPS 对应地址。

#### Scenario: HTTP 请求自动跳转
- **WHEN** 用户访问 `http://center.xingyunxuan.cn` 或 `http://center.xingyunxuan.cn/api/...`
- **THEN** Nginx 返回 301 重定向到 `https://center.xingyunxuan.cn` 对应路径

### Requirement: SSL 证书文件通过 Docker Volume 挂载
`docker-compose.prod.yml` 中前端服务 SHALL 将宿主机 `/opt/sale/ssl/` 目录挂载到容器内 `/etc/nginx/ssl/`，包含 `cert.pem` 和 `key.pem` 两个文件。

#### Scenario: 证书文件正确挂载
- **WHEN** 前端容器启动
- **THEN** 容器内 `/etc/nginx/ssl/cert.pem` 和 `/etc/nginx/ssl/key.pem` 文件存在且可读

#### Scenario: 证书文件缺失
- **WHEN** 宿主机 `/opt/sale/ssl/` 目录下证书文件不存在
- **THEN** Nginx 启动失败，容器日志输出 SSL 证书加载错误

### Requirement: 前端容器暴露 443 端口
`docker-compose.prod.yml` 前端服务 SHALL 同时暴露 80 和 443 端口到宿主机。

#### Scenario: 端口映射配置
- **WHEN** 查看 docker-compose.prod.yml 前端服务配置
- **THEN** ports 包含 `"80:80"` 和 `"443:443"` 两条映射

### Requirement: CORS 配置更新为 HTTPS 域名
`docker-compose.prod.yml` 和 `.env.example` 中 `CORS_ALLOWED_ORIGINS` SHALL 包含 `https://center.xingyunxuan.cn`。

#### Scenario: CORS 允许 HTTPS 源
- **WHEN** 前端从 `https://center.xingyunxuan.cn` 发起 API 请求
- **THEN** 后端返回正确的 CORS 响应头，请求正常处理
