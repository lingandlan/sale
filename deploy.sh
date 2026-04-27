#!/bin/bash
# 太积堂部署脚本
# 用法: ./deploy.sh [full|app]
#   full - 首次部署，启动所有服务（db/redis/backend/frontend）
#   app  - 日常更新，只重建 backend 和 frontend（默认）

set -e

cd /opt/sale

MODE=${1:-app}

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

  # 构建应用镜像
  echo "构建应用镜像..."
  docker compose build backend frontend

  # 启动应用
  docker compose -f docker-compose.prod.yml --env-file .env up -d --no-deps backend frontend
  echo "=== 首次部署完成 ==="

elif [ "$MODE" = "app" ]; then
  echo "=== 日常更新：重建应用服务 ==="

  # 检查基础设施是否运行
  if ! docker ps --format '{{.Names}}' | grep -q 'sale-mysql'; then
    echo "错误：MySQL 未运行，请先执行 ./deploy.sh full"
    exit 1
  fi

  # 加载新镜像（如果有）
  if [ -d artifacts ]; then
    echo "加载构建产物..."
    # 后端
    if [ -f artifacts/backend.tar.gz ]; then
      tar xzf artifacts/backend.tar.gz -C artifacts/ 2>/dev/null || true
      docker build -f Dockerfile.prod-backend -t sale-backend:latest artifacts/backend/ 2>/dev/null || true
    fi
    # 前端
    if [ -f artifacts/frontend.tar.gz ]; then
      tar xzf artifacts/frontend.tar.gz -C artifacts/ 2>/dev/null || true
      docker build -f Dockerfile.prod-frontend -t sale-frontend:latest artifacts/frontend/ 2>/dev/null || true
    fi
  fi

  # 按服务重建
  docker compose -f docker-compose.prod.yml --env-file .env up -d --no-deps --force-recreate backend frontend
  echo "=== 日常更新完成 ==="

else
  echo "用法: $0 [full|app]"
  echo "  full - 首次部署，启动所有服务"
  echo "  app  - 日常更新，只重建 backend/frontend"
  exit 1
fi

echo ""
docker compose -f docker-compose.prod.yml ps
