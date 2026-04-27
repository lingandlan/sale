## 1. 阿里云准备（手动）

- [ ] 1.1 在阿里云容器镜像服务创建命名空间（如 `lingandlan`）
- [ ] 1.2 创建仓库 `sale-backend` 和 `sale-frontend`
- [ ] 1.3 设置仓库固定密码
- [ ] 1.4 在 GitHub repo Settings → Secrets 添加：`ACR_REGISTRY`、`ACR_USERNAME`、`ACR_PASSWORD`

## 2. 修改 deploy.yml

- [ ] 2.1 GHCR 登录步骤改为 ACR 登录（使用 ACR secrets）
- [ ] 2.2 镜像 tag 从 `ghcr.io/...` 改为 `${{ secrets.ACR_REGISTRY }}/...`
- [ ] 2.3 SSH 脚本中 docker login 改为 ACR
- [ ] 2.4 SSH 脚本中内嵌 docker-compose.yml 的镜像地址改为 ACR

## 3. 修改 docker-compose.prod.yml

- [ ] 3.1 backend 和 frontend 镜像地址改为 ACR

## 4. 验证

- [ ] 4.1 合并后 deploy workflow 成功推送镜像到 ACR
- [ ] 4.2 服务器从 ACR 拉取镜像并正常启动
