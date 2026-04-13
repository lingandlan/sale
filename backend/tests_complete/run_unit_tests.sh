#!/bin/bash
set -e

echo "=========================================="
echo "运行完整测试套件"
echo "=========================================="

cd "$(dirname "$0")/.."

echo ""
echo ">>> 1. 运行单元测试..."
echo ""

go test ./internal/... -coverprofile=coverage.out -covermode=atomic -v 2>&1 | tee unit_test.log

if [ ${PIPESTATUS[0]} -ne 0 ]; then
    echo ""
    echo "❌ 单元测试失败!"
    exit 1
fi

echo ""
echo ">>> 2. 生成覆盖率报告..."
echo ""

go tool cover -html=coverage.out -o coverage.html 2>/dev/null || true

echo ""
echo ">>> 3. 单元测试完成 <<<"
echo ""
echo "覆盖率报告: coverage.html"
echo ""

exit 0
