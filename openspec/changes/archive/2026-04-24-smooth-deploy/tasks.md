## 1. 配置外部化

- [x] 1.1 后端 config.go 扩展 envOverrides：添加 MALL_CUSTOMER_ID、MALL_BASE_URL、CORS_ALLOWED_ORIGINS 映射
- [x] 1.2 后端 config.go 添加 CORS_ALLOWED_ORIGINS 解析逻辑（逗号分隔字符串 → []string）
- [x] 1.3 docker-compose.yml 中所有明文密码改为 `${VAR:-default}` 格式
- [x] 1.4 创建 `.env.example` 模板文件（列出所有需要配置的环境变量及说明）
- [x] 1.5 `.gitignore` 中确认 `.env` 已排除

## 2. 生产环境 Docker Compose

- [x] 2.1 新建 `docker-compose.prod.yml`，backend/frontend 使用 `image:` 而非 `build:`，其余与 docker-compose.yml 一致
- [x] 2.2 docker-compose.prod.yml 中密码全部引用 `${VAR}`（无默认值，强制从 .env 读取）

## 3. 服务器准备

- [ ] 3.1 在服务器创建 `/opt/sale/.env` 文件，填入生产环境密钥（JWT_SECRET、数据库密码、WSY 配置、CORS 域名）— 需手动执行
- [ ] 3.2 服务器上传 `docker-compose.prod.yml` — 推送后 git pull 获取

## 4. Codeup Flow 构建阶段

- [ ] 4.1 Flow YAML 增加构建阶段：clone 代码 → docker build backend → docker build frontend
- [ ] 4.2 构建产物 `docker save` 打包为 tar.gz
- [ ] 4.3 验证 Flow 公共构建机是否支持 Docker daemon（如不支持，回退方案见任务 6）

## 5. Codeup Flow 部署阶段

- [ ] 5.1 部署脚本改为：接收镜像 tar → `docker load` → `docker compose -f docker-compose.prod.yml up -d --no-deps --force-recreate backend frontend`
- [ ] 5.2 db/redis 只在首次或配置变更时重建

## 6. 回退方案（如 Flow 不支持 Docker）

- [ ] 6.1 如果 Flow 构建机无 Docker，保留当前服务器构建方案，只做配置外部化和平滑重建
- [ ] 6.2 部署脚本改为：`cd /opt/sale && git pull && docker compose up -d --no-deps --force-recreate backend`（按需扩展到 frontend）

## 7. 验证

- [ ] 7.1 本地用 `.env` + `docker compose up` 验证配置读取正常
- [ ] 7.2 触发一次完整部署，确认服务器正常运行
- [ ] 7.3 确认 `.env` 不在 git 仓库中
