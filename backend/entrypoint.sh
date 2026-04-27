#!/bin/sh
set -e

# 等待 MySQL 就绪（最多 60 秒）
echo "Waiting for MySQL at ${APP_DATABASE_HOST:-db}:${APP_DATABASE_PORT:-3306}..."
for i in $(seq 1 30); do
  if wget --no-verbose --tries=1 --spider "${APP_DATABASE_HOST:-db}:${APP_DATABASE_PORT:-3306}/" 2>/dev/null; then
    echo "MySQL is ready"
    break
  fi
  if [ "$i" -eq 30 ]; then
    echo "ERROR: MySQL not ready after 60 seconds"
    exit 1
  fi
  echo "  Waiting... ($i/30)"
  sleep 2
done

echo "Running database migrations..."
./migrate

echo "Starting server..."
exec ./server
