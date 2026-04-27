## ADDED Requirements

### Requirement: 服务器初始化步骤
文档 SHALL 包含服务器初始化步骤：安装 Docker、Docker Compose、创建部署用户、配置 SSH 密钥登录。

#### Scenario: 按文档初始化新服务器
- **WHEN** 运维人员在全新 ECS 上按文档操作
- **THEN** Docker 和 Docker Compose 可用，部署用户可通过 SSH 密钥登录

### Requirement: DNS 配置说明
文档 SHALL 说明在阿里云 DNS 控制台添加 A 记录，将 `center.xingyunxuan.cn` 指向服务器公网 IP。

#### Scenario: DNS 解析生效
- **WHEN** DNS A 记录已添加
- **THEN** `dig center.xingyunxuan.cn` 返回正确的服务器 IP 地址

### Requirement: 安全组端口配置说明
文档 SHALL 说明阿里云安全组需开放的端口：22（SSH）、80（HTTP）、443（HTTPS）。

#### Scenario: 安全组配置完成后端口可达
- **WHEN** 安全组规则已添加
- **THEN** 外部可通过 80/443 访问服务，通过 22 SSH 登录

### Requirement: SSL 证书下载与放置步骤
文档 SHALL 说明从阿里云 SSL 控制台下载 Nginx 格式证书（`.pem` + `.key`），放置到服务器 `/opt/sale/ssl/` 目录，并命名为 `cert.pem` 和 `key.pem`。

#### Scenario: 证书正确放置
- **WHEN** 运维人员按文档操作完成证书下载和放置
- **THEN** `/opt/sale/ssl/cert.pem` 和 `/opt/sale/ssl/key.pem` 文件存在

### Requirement: 代码部署步骤
文档 SHALL 包含完整部署步骤：git clone 代码、复制并填写 `.env` 文件、执行 `./deploy.sh full` 启动所有服务。

#### Scenario: 首次部署成功
- **WHEN** 运维人员按文档完成所有步骤
- **THEN** `docker compose ps` 显示 backend、frontend、db、redis 四个服务均为 running/up 状态，`curl -k https://localhost/health` 返回 `{"status":"ok"}`

### Requirement: 部署验证清单
文档 SHALL 包含部署后的验证清单：HTTPS 可访问、API 正常响应、登录功能正常、数据库连接正常。

#### Scenario: 按清单逐项验证
- **WHEN** 部署完成后按清单验证
- **THEN** 所有验证项通过

### Requirement: 日常更新部署步骤
文档 SHALL 说明日常代码更新时的部署操作：`git pull` + `./deploy.sh app`，仅重建应用容器不影响数据库。

#### Scenario: 日常更新部署
- **WHEN** 运维人员执行日常更新步骤
- **THEN** 应用更新完成，数据库数据未丢失

### Requirement: 数据库备份策略
文档 SHALL 包含 MySQL 自动备份方案：crontab 定时任务每日凌晨 3 点执行 mysqldump，备份文件保留 7 天，存放在 `/opt/sale/backups/`。

#### Scenario: 自动备份正常执行
- **WHEN** crontab 配置完成且经过一个凌晨 3 点
- **THEN** `/opt/sale/backups/` 目录下生成对应日期的 SQL 备份文件

#### Scenario: 手动恢复备份
- **WHEN** 运维人员按文档执行恢复命令
- **THEN** 数据库恢复到备份时的状态

### Requirement: 查看日志方法
文档 SHALL 说明查看容器日志的命令：`docker compose logs`、`docker compose logs -f backend` 等。

#### Scenario: 排查后端错误
- **WHEN** 运维人员需要排查后端问题
- **THEN** 按文档命令可查看实时后端日志

### Requirement: 回滚方案
文档 SHALL 说明回滚步骤：`git log` 找到上一个版本 → `git checkout <commit>` → 重新构建部署。

#### Scenario: 部署失败后回滚
- **WHEN** 新版本部署后发现问题，需要回滚到上一版本
- **THEN** 按文档步骤回滚成功，服务恢复到上一版本正常运行
