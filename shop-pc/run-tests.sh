#!/bin/bash

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}========================================${NC}"
echo -e "${YELLOW}  太积堂系统 - 单元测试${NC}"
echo -e "${YELLOW}========================================${NC}"

# 检查是否在前端项目目录
if [ ! -f "package.json" ]; then
    echo -e "${RED}错误: 请在前端项目根目录运行此脚本${NC}"
    exit 1
fi

# 检查测试依赖
echo -e "\n${YELLOW}检查测试依赖...${NC}"
if ! npm list vitest > /dev/null 2>&1; then
    echo -e "${RED}vitest 未安装，正在安装...${NC}"
    npm install -D vitest @vitest/ui jsdom @vue/test-utils
fi

echo -e "${GREEN}✓ 测试依赖已就绪${NC}"

# 运行测试
echo -e "\n${YELLOW}运行前端单元测试...${NC}"
npm run test

# 检查测试结果
if [ $? -eq 0 ]; then
    echo -e "\n${GREEN}========================================${NC}"
    echo -e "${GREEN}✓ 所有测试通过！${NC}"
    echo -e "${GREEN}========================================${NC}"
    exit 0
else
    echo -e "\n${RED}========================================${NC}"
    echo -e "${RED}✗ 测试失败，请检查错误信息${NC}"
    echo -e "${RED}========================================${NC}"
    exit 1
fi
