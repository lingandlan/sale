#!/bin/bash

# =============================================================================
# Integration Test Helper Script
# =============================================================================

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo_step() {
    echo -e "${GREEN}==>${NC} $1"
}

echo_warn() {
    echo -e "${YELLOW}WARNING:${NC} $1"
}

echo_error() {
    echo -e "${RED}ERROR:${NC} $1"
}

# 测试数据库连接
TEST_DSN="${TEST_DSN:-root:root123@tcp(localhost:3307)/test_db?parseTime=true&charset=utf8mb4}"

# 启动测试数据库
start_test_db() {
    echo_step "启动测试数据库..."
    cd "$(dirname "$0")"
    docker-compose -f docker-compose.yml up -d mysql-test
    
    echo_step "等待数据库就绪..."
    for i in {1..30}; do
        if docker exec sale-mysql-test mysqladmin ping -h localhost -uroot -proot123 &>/dev/null; then
            echo_step "数据库已就绪"
            return 0
        fi
        echo "等待中... ($i/30)"
        sleep 2
    done
    
    echo_error "数据库启动超时"
    return 1
}

# 停止测试数据库
stop_test_db() {
    echo_step "停止测试数据库..."
    cd "$(dirname "$0")"
    docker-compose -f docker-compose.yml down
}

# 运行集成测试
run_tests() {
    export TEST_DSN="$TEST_DSN"
    
    echo_step "运行集成测试..."
    go test -v -tags=integration ./tests/integration/...
}

# 显示帮助
show_help() {
    echo "集成测试脚本"
    echo ""
    echo "用法: $0 <命令>"
    echo ""
    echo "命令:"
    echo "  start     启动测试数据库"
    echo "  stop      停止测试数据库"
    echo "  test      运行集成测试 (自动启动数据库)"
    echo "  clean     清理测试数据"
    echo "  help      显示帮助"
    echo ""
    echo "环境变量:"
    echo "  TEST_DSN  测试数据库连接字符串"
    echo "           默认: root:root123@tcp(localhost:3307)/test_db"
}

# 主逻辑
case "${1:-help}" in
    start)
        start_test_db
        ;;
    stop)
        stop_test_db
        ;;
    test)
        start_test_db
        trap 'stop_test_db' EXIT
        run_tests
        ;;
    clean)
        echo_step "清理测试数据库..."
        docker exec sale-mysql-test mysql -uroot -proot123 -e "DROP DATABASE IF EXISTS test_db; CREATE DATABASE test_db;"
        ;;
    help|*)
        show_help
        ;;
esac
