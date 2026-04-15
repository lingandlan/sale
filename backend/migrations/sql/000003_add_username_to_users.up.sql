-- 用户表增加 username 字段
ALTER TABLE users ADD COLUMN username VARCHAR(50) NOT NULL DEFAULT '' COMMENT '用户名（登录账号）' AFTER id;
ALTER TABLE users ADD UNIQUE KEY uk_username (username);

-- 已有用户用 phone 作为默认 username
UPDATE users SET username = phone WHERE username = '';
