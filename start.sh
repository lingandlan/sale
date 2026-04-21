#!/bin/bash
# 太积堂 - Beta 环境启动脚本
# 端口写死，不依赖 .env.local

set -e

ROOT="$(cd "$(dirname "$0")" && pwd)"

FRONTEND_PORT="5179"
BACKEND_PORT="8082"
REDIS_DB="2"

# 清理僵尸进程
echo "🧹 清理残留进程..."
lsof -ti:$FRONTEND_PORT 2>/dev/null | xargs kill 2>/dev/null || true
lsof -ti:$BACKEND_PORT 2>/dev/null | xargs kill 2>/dev/null || true
sleep 1

# 启动后端（通过环境变量覆盖 config.yaml）
echo "🚀 启动后端 (port $BACKEND_PORT, Redis DB $REDIS_DB)..."
cd "$ROOT/backend"
APP_SERVER_PORT=$BACKEND_PORT APP_REDIS_DB=$REDIS_DB air &

# 启动前端，通过环境变量指定端口和代理目标（不依赖 .env.local）
echo "🚀 启动前端 (port $FRONTEND_PORT)..."
cd "$ROOT/shop-pc"
VITE_PORT=$FRONTEND_PORT VITE_API_PORT=$BACKEND_PORT npx vite &

echo ""
echo "✅ Beta 环境已启动:"
echo "   前端: http://localhost:$FRONTEND_PORT"
echo "   后端: http://localhost:$BACKEND_PORT"
echo "   按 Ctrl+C 停止所有服务"
wait
