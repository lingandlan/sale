#!/bin/bash
# 太积堂部署脚本
# 用法: ./deploy.sh [full|app]
#   full - 首次部署，启动所有服务（db/redis/backend/frontend）
#   app  - 日常更新，拉代码 + 编译 + 重建 backend/frontend（默认）

set -e

cd /opt/sale

MODE=${1:-app}

# 编译 Go 二进制
build_backend() {
  echo "=== 编译后端 ==="
  cd backend
  GOPROXY=https://goproxy.cn,direct CGO_ENABLED=0 GOOS=linux go build -o ../server ./cmd/server
  GOPROXY=https://goproxy.cn,direct CGO_ENABLED=0 GOOS=linux go build -o ../migrate ./cmd/migrate
  cd ..
  echo "编译完成: server, migrate"
}

# 构建前端
build_frontend() {
  echo "=== 构建前端 ==="
  cd shop-pc
  npm install --prefer-offline --registry=https://registry.npmmirror.com
  npm run build
  cd ..
  echo "前端构建完成"
}

# 构建 Docker 镜像并启动服务
deploy_app() {
  # 构建后端镜像
  echo "构建后端镜像..."
  docker build -f Dockerfile.prod-backend -t sale-backend:latest .

  # 构建前端镜像
  echo "构建前端镜像..."
  docker build -f Dockerfile.prod-frontend -t sale-frontend:latest .

  # 重启服务
  docker compose -f docker-compose.prod.yml --env-file .env up -d --no-deps --force-recreate backend frontend
}

if [ "$MODE" = "full" ]; then
  echo "=== 首次部署：启动所有服务 ==="

  # 启动基础设施
  docker compose -f docker-compose.prod.yml --env-file .env up -d db redis

  # 等待 MySQL 就绪
  echo "等待 MySQL 就绪..."
  for i in $(seq 1 30); do
    if docker exec sale-mysql mysqladmin ping -h localhost --silent 2>/dev/null; then
      echo "MySQL 就绪"
      break
    fi
    echo "  等待中... ($i/30)"
    sleep 2
  done

  # 拉代码 + 编译 + 部署
  git pull
  build_backend
  build_frontend
  deploy_app

  echo "=== 首次部署完成 ==="

elif [ "$MODE" = "app" ]; then
  echo "=== 日常更新 ==="

  # 检查基础设施是否运行
  if ! docker ps --format '{{.Names}}' | grep -q 'sale-mysql'; then
    echo "错误：MySQL 未运行，请先执行 ./deploy.sh full"
    exit 1
  fi

  # 拉代码 + 编译 + 部署
  git pull
  build_backend
  deploy_app

  echo "=== 日常更新完成 ==="

else
  echo "用法: $0 [full|app]"
  echo "  full - 首次部署，启动所有服务（含前端构建）"
  echo "  app  - 日常更新，拉代码 + 编译后端 + 重启"
  exit 1
fi

echo ""
docker compose -f docker-compose.prod.yml ps
