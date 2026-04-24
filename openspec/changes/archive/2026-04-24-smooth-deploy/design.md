## Context

当前部署流程：push 到 Codeup master → Flow SSH 到服务器 → `git pull + docker compose down + up --build`。服务器承担编译构建（Go + Node），在国内网络下 `go mod download` 经常卡住。所有敏感配置（JWT 密钥、WSY API 密钥、数据库密码）硬编码在 config.yaml 和 docker-compose.yml 中。

## Goals / Non-Goals

**Goals:**
- 构建从服务器迁移到 Codeup Flow 云端，服务器只运行容器
- 按服务单独重建，前端/后端/基础设施变更互不影响
- 敏感配置通过环境变量注入，生产环境使用服务器上的 .env 文件

**Non-Goals:**
- 不做蓝绿部署或零中断（当前单服务器场景不必要）
- 不搭建镜像仓库（ACR 企业版收费，个人版不可用）
- 不改本地开发流程（本地继续用 config.yaml）
- 不引入 K8s 或 Docker Swarm

## Decisions

### 1. 镜像传输方案：docker save/load via SSH

**选择**：Pipeline 中 `docker save | gzip` → SSH `docker load`

**替代方案**：
- ACR 个人版：用户账号下只有企业版入口，不可用
- docker registry 自建：增加运维复杂度，单服务器不值得
- 直接 SSH pipe `docker save | ssh server "docker load"`：可行但大镜像会占带宽

**方案**：Pipeline 构建后 `docker save` 打包，作为构建产物上传，SSH 部署时 `docker load` 加载。镜像通过 Flow 的产物机制传递。

### 2. Codeup Flow YAML 结构

```
stages:
  build:
    - 构建 backend 镜像 (docker build)
    - 构建 frontend 镜像 (docker build)
  deploy:
    - SSH 到服务器执行部署脚本
```

构建在 Flow 的公共构建机上完成（有 Docker 环境），部署通过主机部署 SSH 执行。

### 3. 按服务重建策略

检测变更范围，选择重建方式：
- 只有 `backend/` 变更 → `docker compose up -d --no-deps --force-recreate backend`
- 只有 `shop-pc/` 变更 → `docker compose up -d --no-deps --force-recreate frontend`
- `docker-compose.yml` 或基础设施变更 → `docker compose down && docker compose up -d`

简化处理：始终用 `--no-deps --force-recreate` 重建变更的服务，不做变更检测（复杂度不值得）。

### 4. 配置外部化方案

**后端 config.go 扩展 envOverrides**：

| 环境变量 | 覆盖配置项 | 优先级 |
|----------|-----------|--------|
| JWT_SECRET | jwt.secret | 最高 |
| DB_PASSWORD / APP_DATABASE_PASSWORD | database.password | 最高 |
| MALL_APP_ID | mall.app_id | 最高 |
| MALL_APP_SECRET | mall.app_secret | 最高 |
| MALL_CUSTOMER_ID | mall.customer_id | 新增 |
| MALL_BASE_URL | mall.base_url | 新增 |
| CORS_ALLOWED_ORIGINS | cors.allowed_origins | 新增 |

**CORS 配置格式**：环境变量用逗号分隔（`http://a.com,http://b.com`），config.go 解析为 `[]string`。

**docker-compose.yml**：密码改为 `${MYSQL_ROOT_PASSWORD:-sale123}` 格式，开发环境有默认值。

**服务器 .env**：`/opt/sale/.env` 存储生产密钥，.gitignore 排除 `.env`。

### 5. docker-compose.yml 调整

backend 服务去掉 `build` 字段，改为 `image: sale-backend:latest`。前端同理。服务器上 `docker load` 先加载镜像，然后 `docker compose up` 直接使用已加载的镜像。

但保留 build 字段作为 fallback（本地开发用），通过 `docker-compose.override.yml` 或 profile 区分。实际方案：生产用单独的 `docker-compose.prod.yml`，只有 `image` 没有 `build`。

## Risks / Trade-offs

- **[镜像体积大]** → backend 镜像约 50MB（Alpine + 静态编译），frontend 约 30MB（Nginx + 静态文件），gzip 后更小，传输可接受
- **[Flow 构建环境 Docker 可用性]** → Codeup Flow 公共构建机支持 Docker，但需验证；若不支持，回退到服务器构建方案
- **[Flow 构建超时]** → 默认超时可能不够（Go 编译需要时间），需设置合理超时（15 分钟）
- **[镜像未版本化]** → 始终用 latest tag，回滚需要重新构建；可接受因为 git 本身有版本记录

## Migration Plan

1. 先改代码：config.go 扩展 envOverrides + docker-compose.yml 用 `${VAR}` + 创建 .env 模板
2. 验证本地：用 .env 启动 docker compose 确认配置读取正常
3. 改 Flow YAML：增加构建阶段 + 调整部署脚本
4. 服务器创建 .env 文件
5. 触发一次完整部署验证
6. 回滚策略：Flow YAML 可随时改回服务器构建方案，代码变更兼容两种模式

## Open Questions

- Codeup Flow 构建环境是否支持 Docker daemon？（需实际验证）
- 如果不支持，是否需要自建 Docker 构建机或回退到服务器构建 + 只做配置外部化？
