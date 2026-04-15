#!/bin/bash
# Worktree Alpha 启动脚本
# 设置环境变量覆盖默认端口，避免合并冲突
export APP_SERVER_PORT=8081
exec air "$@"
