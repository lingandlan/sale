## 1. entrypoint.sh 添加 MySQL 就绪等待

- [x] 1.1 在 `backend/entrypoint.sh` 中 migrate 执行前添加 MySQL `db:3306` TCP 检测循环，最多等待 60 秒，超时 exit 1
- [x] 1.2 在 `Dockerfile.prod-backend` 中安装 `nc` 或确保 `wget` 可用（alpine 已有 wget），用于端口探测

## 2. 确认 Dockerfile COPY 路径

- [x] 2.1 验证 `Dockerfile.prod-backend` 的 COPY 路径与 CI 构建上下文（仓库根目录）一致 ✓
- [x] 2.2 修复 `Dockerfile.prod-frontend` build context：flow.yaml 统一用仓库根目录作为 context

## 3. 确认环境变量覆盖链路

- [x] 3.1 添加 `APP_DATABASE_PASSWORD` 到 config.go envOverrides（修复密码无法传递的 bug）
- [x] 3.2 添加 `APP_DATABASE_PASSWORD` 到 migrate.go 显式读取

## 4. 确认 docker-compose.prod.yml 配置

- [x] 4.1 确认 backend environment 完整 ✓
- [x] 4.2 确认 frontend SSL 卷挂载 ✓
- [x] 4.3 确认 db MYSQL_DATABASE=sale_prod ✓
- [x] 4.4 确认无 build 指令 ✓

## 5. 重写 flow.yaml

- [x] 5.1 编写 build stage：build_frontend job ✓
- [x] 5.2 编写 build stage：build_backend job ✓
- [x] 5.3 编写 deploy stage：VMDeploy job ✓（前端 build context 修复为仓库根目录）

## 6. 本地验证

- [x] 6.1 后端 Go 编译验证通过 ✓
- [ ] 6.2 前端 Docker 镜像构建（需要 CI 环境验证）
- [x] 6.3 docker-compose.prod.yml config 验证通过 ✓
