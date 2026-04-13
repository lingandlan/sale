#!/bin/bash
set -e

echo "=========================================="
echo "完整测试套件运行器"
echo "=========================================="

cd "$(dirname "$0")"

UNIT_PASSED=0
INTEGRATION_PASSED=0
E2E_PASSED=0

echo ""
echo ">>> 1. 单元测试..."
echo ""

if go test ./internal/... -v 2>&1; then
    echo "✅ 单元测试通过"
    UNIT_PASSED=1
else
    echo "❌ 单元测试失败"
fi

echo ""
echo ">>> 2. 集成测试..."
echo ""

if go run ./integration/integration_test.go 2>&1; then
    echo "✅ 集成测试通过"
    INTEGRATION_PASSED=1
else
    echo "❌ 集成测试失败"
fi

echo ""
echo ">>> 3. E2E 测试..."
echo ""

if command -v npx &> /dev/null; then
    cd e2e
    if npm test 2>&1; then
        echo "✅ E2E 测试通过"
        E2E_PASSED=1
    else
        echo "❌ E2E 测试失败"
    fi
    cd ..
else
    echo "⚠️ 跳过 E2E 测试 (npx 未安装)"
fi

echo ""
echo "=========================================="
echo "测试结果汇总"
echo "=========================================="
echo "单元测试:    $([ $UNIT_PASSED -eq 1 ] && echo '✅ 通过' || echo '❌ 失败')"
echo "集成测试:    $([ $INTEGRATION_PASSED -eq 1 ] && echo '✅ 通过' || echo '❌ 失败')"
echo "E2E 测试:    $([ $E2E_PASSED -eq 1 ] && echo '✅ 通过' || echo '❌ 失败')"
echo "=========================================="

if [ $UNIT_PASSED -eq 1 ] && [ $INTEGRATION_PASSED -eq 1 ]; then
    echo "核心测试通过 ✅"
    exit 0
else
    echo "有测试失败 ❌"
    exit 1
fi
