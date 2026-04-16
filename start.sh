#!/bin/bash
# 太积堂 - Main 环境启动脚本
# 前端: 5175  后端: 8080  Redis DB: 0

set -e
ROOT="$(cd "$(dirname "$0")/.." && pwd)"

# 清理僵尸进程
echo "🧹 清理残留进程..."
lsof -ti:5175 2>/dev/null | xargs kill 2>/dev/null || true
lsof -ti:8080 2>/dev/null | xargs kill 2>/dev/null || true
sleep 1

# 启动后端
echo "🚀 启动后端 (port 8080)..."
cd "$ROOT/backend"
air &

# 启动前端
echo "🚀 启动前端 (port 5175)..."
cd "$ROOT/shop-pc"
npx vite &

echo ""
echo "✅ Main 环境已启动:"
echo "   前端: http://localhost:5175"
echo "   后端: http://localhost:8080"
echo "   按 Ctrl+C 停止所有服务"
wait
