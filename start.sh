#!/bin/bash
# 太积堂 - Beta 环境启动脚本
# 前端: 5179  后端: 8082  Redis DB: 2

set -e
ROOT="$(cd "$(dirname "$0")/.." && pwd)"

# 清理僵尸进程
echo "🧹 清理残留进程..."
lsof -ti:5179 2>/dev/null | xargs kill 2>/dev/null || true
lsof -ti:8082 2>/dev/null | xargs kill 2>/dev/null || true
sleep 1

# 启动后端（环境变量覆盖端口和 Redis DB）
echo "🚀 启动后端 (port 8082)..."
cd "$ROOT/backend"
APP_SERVER_PORT=8082 APP_REDIS_DB=2 air &

# 启动前端
echo "🚀 启动前端 (port 5179)..."
cd "$ROOT/shop-pc"
npx vite &

echo ""
echo "✅ Beta 环境已启动:"
echo "   前端: http://localhost:5179"
echo "   后端: http://localhost:8082"
echo "   按 Ctrl+C 停止所有服务"
wait
