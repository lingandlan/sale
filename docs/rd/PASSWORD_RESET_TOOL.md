# 密码重置工具使用说明

## 工具位置
`/Users/zhangdaodong/code/sale/backend/cmd/resetpwd/main.go`

## 功能
- 重置用户登录密码
- 支持指定新密码或自动生成
- 自动进行bcrypt加密
- 验证密码强度

## 使用方法

### 1. 重置为指定密码
```bash
cd /Users/zhangdaodong/code/sale/backend
go run ./cmd/resetpwd/main.go -phone 13800000000 -password 新密码123
```

### 2. 自动生成随机密码
```bash
cd /Users/zhangdaodong/code/sale/backend
go run ./cmd/resetpwd/main.go -phone 13800000000
```

## 参数说明
- `-phone`: 用户手机号（必需）
- `-password`: 新密码（可选，不指定则自动生成8位随机密码）
- `-dsn`: 数据库连接字符串（可选，默认为本地开发配置）

## 输出示例
```
🔗 连接数据库...
🔍 查找用户: 13800000000
✅ 找到用户:
   ID: 1
   手机号: 13800000000
   姓名: 超级管理员
   角色: super_admin
   状态: 1

📝 新密码: Test123456
🔐 生成密码hash...
💾 更新数据库...

✅ 密码重置成功！

═══════════════════════════════════════
  手机号: 13800000000
  新密码: Test123456
═══════════════════════════════════════

💡 提示: 请妥善保管新密码
```

## 密码要求
- 最小长度: 6位
- 最大长度: 32位
- 自动生成的密码包含: 大小写字母、数字、特殊字符

## 故障排除

### 错误: 用户不存在
```bash
❌ 用户不存在: 13900000000
```
解决: 检查手机号是否正确，或用户是否被软删除

### 错误: 密码长度不符合要求
```bash
❌ 密码长度不能少于6位
```
解决: 使用6-32位的密码

### 错误: 连接数据库失败
```bash
❌ 连接数据库失败: dial tcp: connection refused
```
解决:
1. 检查Docker服务是否运行
2. 启动数据库: `docker compose up -d`

## 相关工具
- 诊断工具: `go run ./cmd/diag/main.go` - 检查用户信息和密码验证
- 数据库测试: `go run ./cmd/testdb/main.go` - 测试数据库连接

## 当前测试账号
```
手机号: 13800000000
密码: Test123456
角色: 超级管理员
```
