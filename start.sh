#!/bin/bash
# 太积堂 - Alpha 环境启动脚本
# 从 shop-pc/.env.local 读取端口配置
# 后端端口通过环境变量 APP_SERVER_PORT / APP_REDIS_DB 覆盖 config.yaml

set -e

# ROOT 指向 alpha 目录（脚本所在目录）
ROOT="$(cd "$(dirname "$0")" && pwd)"

# 从 shop-pc/.env.local 读取端口配置
ENV_FILE="$ROOT/shop-pc/.env.local"
FRONTEND_PORT="5178"
BACKEND_PORT="8081"
REDIS_DB="1"
if [ -f "$ENV_FILE" ]; then
  while IFS='=' read -r key value; do
    case "$key" in
      VITE_PORT) FRONTEND_PORT="$value" ;;
      VITE_API_PORT) BACKEND_PORT="$value" ;;
      APP_REDIS_DB) REDIS_DB="$value" ;;
    esac
  done < "$ENV_FILE"
fi

# 清理僵尸进程
echo "🧹 清理残留进程..."
lsof -ti:$FRONTEND_PORT 2>/dev/null | xargs kill 2>/dev/null || true
lsof -ti:$BACKEND_PORT 2>/dev/null | xargs kill 2>/dev/null || true
sleep 1

# 启动后端（通过环境变量覆盖 config.yaml）
echo "🚀 启动后端 (port $BACKEND_PORT, Redis DB $REDIS_DB)..."
cd "$ROOT/backend"
APP_SERVER_PORT=$BACKEND_PORT APP_REDIS_DB=$REDIS_DB air &

# 启动前端（vite 自动读取 .env.local 中的 VITE_PORT）
echo "🚀 启动前端 (port $FRONTEND_PORT)..."
cd "$ROOT/shop-pc"
npx vite &

echo ""
echo "✅ Alpha 环境已启动:"
echo "   前端: http://localhost:$FRONTEND_PORT"
echo "   后端: http://localhost:$BACKEND_PORT"
echo "   按 Ctrl+C 停止所有服务"
wait
