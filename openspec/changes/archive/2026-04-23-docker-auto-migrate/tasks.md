## 1. Dockerfile 改造

- [x] 1.1 `backend/Dockerfile`: 构建阶段额外编译 `cmd/migrate`，输出 `migrate` 二进制
- [x] 1.2 `backend/Dockerfile`: 将 `migrate` 二进制 COPY 到最终镜像

## 2. Entrypoint 脚本

- [x] 2.1 创建 `backend/entrypoint.sh`：先执行 `./migrate`，成功后 `exec ./server`，失败则 `exit 1`
- [x] 2.2 `backend/Dockerfile`: COPY entrypoint.sh 并设为 ENTRYPOINT

## 3. 环境变量配置

- [x] 3.1 `docker-compose.yml`: backend 服务添加 `DB_HOST=db`、`DB_PORT=3306`、`DB_USER=sale`、`DB_PASSWORD=sale123`、`DB_NAME=sale_dev` 环境变量
- [x] 3.2 `docker-compose.prod.yml`: 同步添加相同环境变量

## 4. 部署流程更新

- [x] 4.1 `.github/workflows/deploy.yml`: SSH 脚本中的 docker-compose.yml 内同步添加 DB 环境变量

## 5. 验证

- [x] 5.1 本地 `docker compose up` 验证：空库启动后自动创建 `schema_migrations` 表并执行所有 migration
- [x] 5.2 重复启动验证：第二次启动跳过已执行的 migration
