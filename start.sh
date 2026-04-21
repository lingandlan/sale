#!/bin/bash
# 太积堂 - 通用启动脚本
# 端口从 .env.worktree 读取（gitignored，每个环境独立维护）

set -e

ROOT="$(cd "$(dirname "$0")" && pwd)"

# 读取本地环境配置（不存在则用 Main 默认值）
if [ -f "$ROOT/.env.worktree" ]; then
  source "$ROOT/.env.worktree"
fi

FRONTEND_PORT="${FRONTEND_PORT:-5175}"
BACKEND_PORT="${BACKEND_PORT:-8080}"
REDIS_DB="${REDIS_DB:-0}"

# 写回 .env.local 确保前端端口正确
echo "VITE_PORT=$FRONTEND_PORT
VITE_API_PORT=$BACKEND_PORT" > "$ROOT/shop-pc/.env.local"

# 清理僵尸进程
echo "🧹 清理残留进程..."
lsof -ti:$FRONTEND_PORT 2>/dev/null | xargs kill 2>/dev/null || true
lsof -ti:$BACKEND_PORT 2>/dev/null | xargs kill 2>/dev/null || true
sleep 1

# 启动后端（通过环境变量覆盖 config.yaml）
echo "🚀 启动后端 (port $BACKEND_PORT, Redis DB $REDIS_DB)..."
cd "$ROOT/backend"
APP_SERVER_PORT=$BACKEND_PORT APP_REDIS_DB=$REDIS_DB air &

# 启动前端
echo "🚀 启动前端 (port $FRONTEND_PORT)..."
cd "$ROOT/shop-pc"
VITE_PORT=$FRONTEND_PORT VITE_API_PORT=$BACKEND_PORT npx vite &

echo ""
echo "✅ 环境已启动:"
echo "   前端: http://localhost:$FRONTEND_PORT"
echo "   后端: http://localhost:$BACKEND_PORT"
echo "   Redis DB: $REDIS_DB"
echo "   按 Ctrl+C 停止所有服务"
wait
