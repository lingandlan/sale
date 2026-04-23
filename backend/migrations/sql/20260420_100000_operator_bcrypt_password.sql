-- 操作员密码从明文迁移到 bcrypt 哈希
-- 此脚本需要配合 Go 程序执行，纯 SQL 无法调用 bcrypt

-- 使用方式：
-- 1. dry-run: go run backend/cmd/migrate_operator_password/main.go --dry-run
-- 2. 实际执行: go run backend/cmd/migrate_operator_password/main.go
-- 3. 执行前请备份数据库

-- 验证：迁移后所有密码应以 $2a$ 开头（bcrypt 哈希特征）
-- SELECT id, LEFT(password, 4) as prefix FROM recharge_operators;
