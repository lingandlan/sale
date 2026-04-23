-- 初始化超管用户（手机号 13800000000，密码 123456）
-- 使用 INSERT IGNORE 避免重复插入
INSERT IGNORE INTO users (username, phone, password, name, role, status, created_at, updated_at)
VALUES ('admin', '13800000000', '$2a$10$h9fcB.C9nI2F7Z67Fjof1ux1eM4G42ITjix0EgHcXp3P.MS1oZ4QS', '超级管理员', 'super_admin', 1, NOW(), NOW());
