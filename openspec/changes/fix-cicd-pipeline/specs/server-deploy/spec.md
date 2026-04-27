## ADDED Requirements

### Requirement: 加载预构建镜像
部署脚本 SHALL 在服务器上执行 `docker load < frontend-image.tar.gz` 和 `docker load < backend-image.tar.gz`，将 CI 构建的镜像加载到本地 Docker。

#### Scenario: 加载镜像成功
- **WHEN** 制品文件存在于服务器工作目录
- **THEN** 镜像加载到 Docker，`docker images` 可看到 `sale-backend:latest` 和 `sale-frontend:latest`

#### Scenario: 制品文件不存在
- **WHEN** tar.gz 文件不在工作目录
- **THEN** 跳过加载，继续执行（不中断部署，支持手动部署场景）

### Requirement: 使用 docker-compose.prod.yml 启动服务
部署脚本 SHALL 执行 `docker compose -f docker-compose.prod.yml --env-file .env up -d` 启动所有服务。

#### Scenario: 全量启动
- **WHEN** 执行部署命令
- **THEN** backend、frontend、db、redis 四个容器全部启动，网络互通

#### Scenario: 增量重启
- **WHEN** 只更新应用镜像（db 和 redis 已运行）
- **THEN** 只有 backend 和 frontend 容器被重建重启，db 和 redis 不受影响

### Requirement: 后端容器启动前等待 MySQL 就绪
后端 `entrypoint.sh` SHALL 在执行 migrate 之前循环检测 `db:3306` 是否可达，最多等待 60 秒，超时则退出。

#### Scenario: MySQL 在等待时间内就绪
- **WHEN** MySQL 容器启动完成，3306 端口可达
- **THEN** entrypoint.sh 继续执行 migrate

#### Scenario: MySQL 等待超时
- **WHEN** 60 秒内 MySQL 3306 端口不可达
- **THEN** entrypoint.sh 输出错误信息并以非零退出码退出，容器状态为非 running

### Requirement: 后端容器自动执行数据库迁移
后端 `entrypoint.sh` SHALL 在 MySQL 就绪后、启动 server 前执行 `./migrate`。

#### Scenario: migrate 成功
- **WHEN** migrate 执行成功（退出码 0）
- **THEN** 继续启动 server

#### Scenario: migrate 失败
- **WHEN** migrate 执行失败（非零退出码）
- **THEN** entrypoint.sh 以相同退出码退出（`set -e`），不启动 server

### Requirement: 前端容器挂载 SSL 证书
docker-compose.prod.yml SHALL 将 `/opt/sale/ssl` 挂载到前端容器的 `/etc/nginx/ssl:ro`，nginx.conf 引用 `/etc/nginx/ssl/cert.pem` 和 `/etc/nginx/ssl/cert.key`。

#### Scenario: SSL 证书正常加载
- **WHEN** 服务器上 `/opt/sale/ssl/cert.pem` 和 `/opt/sale/ssl/cert.key` 存在
- **THEN** nginx 启动成功，HTTPS 443 端口正常监听

#### Scenario: SSL 证书文件缺失
- **WHEN** 证书文件不存在
- **THEN** nginx 启动失败，前端容器 Restarting

### Requirement: 部署后健康检查验证
部署脚本 SHALL 在容器启动后等待 5 秒，执行 `docker logs sale-backend --tail 10` 输出后端日志，确认服务启动正常。

#### Scenario: 验证输出
- **WHEN** 部署脚本执行完成
- **THEN** 控制台显示后端最近 10 行日志，用于确认服务状态
