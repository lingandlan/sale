#!/bin/bash

# 太积堂后端数据库迁移脚本

echo "🚀 开始数据库迁移..."

# 检查MySQL是否运行
if ! command -v mysql &> /dev/null; then
    echo "❌ 错误: 未找到MySQL命令"
    echo "请确保MySQL已安装并运行"
    exit 1
fi

# 检查配置文件
if [ ! -f "configs/config.yaml" ]; then
    echo "❌ 错误: 未找到配置文件 configs/config.yaml"
    exit 1
fi

# 读取数据库配置
DB_USER=$(grep "user:" configs/config.yaml | awk '{print $2}')
DB_PASS=$(grep "password:" configs/config.yaml | awk '{print $2}')
DB_NAME=$(grep "name:" configs/config.yaml | awk '{print $2}')

echo "📋 数据库配置:"
echo "   用户: $DB_USER"
echo "   数据库: $DB_NAME"
echo ""

# 检查数据库连接
echo "🔗 检查数据库连接..."
mysql -u"$DB_USER" -p"$DB_PASS" -e "SELECT 1;" &> /dev/null
if [ $? -ne 0 ]; then
    echo "❌ 错误: 无法连接到MySQL数据库"
    echo "请检查用户名和密码是否正确"
    exit 1
fi
echo "✅ 数据库连接成功"

# 创建数据库（如果不存在）
echo "📦 创建数据库（如果不存在）..."
mysql -u"$DB_USER" -p"$DB_PASS" -e "CREATE DATABASE IF NOT EXISTS $DB_NAME CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
if [ $? -eq 0 ]; then
    echo "✅ 数据库 $DB_NAME 已就绪"
else
    echo "❌ 创建数据库失败"
    exit 1
fi

# 运行迁移程序
echo ""
echo "🔧 运行数据库迁移..."
go run migrations/migrate.go
if [ $? -eq 0 ]; then
    echo ""
    echo "✅ 数据库迁移完成!"
    echo ""
    echo "📝 默认账号信息:"
    echo "   管理员账号: admin"
    echo "   管理员密码: admin123"
    echo ""
    echo "🎯 下一步:"
    echo "   1. 启动后端服务: ./bin/server"
    echo "   2. 启动前端服务: cd ../shop-pc && npm run dev"
    echo "   3. 访问系统: http://localhost:5177"
else
    echo "❌ 数据库迁移失败"
    exit 1
fi
