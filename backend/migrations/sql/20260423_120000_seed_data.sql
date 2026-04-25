-- 初始化基础数据：超管账号 + Casbin 权限策略
-- 使用 IF NOT EXISTS / INSERT IGNORE 确保幂等

-- ========================================
-- 1. 确保 casbin_rule 表存在（SQL 迁移先于 GORM AutoMigrate 执行）
-- ========================================
CREATE TABLE IF NOT EXISTS `casbin_rule` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `ptype` varchar(12) DEFAULT NULL,
  `v0` varchar(128) DEFAULT NULL,
  `v1` varchar(128) DEFAULT NULL,
  `v2` varchar(128) DEFAULT NULL,
  `v3` varchar(128) DEFAULT NULL,
  `v4` varchar(128) DEFAULT NULL,
  `v5` varchar(128) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_casbin_rule` (`ptype`,`v0`,`v1`,`v2`,`v3`,`v4`,`v5`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ========================================
-- 2. 确保 users 表存在
-- ========================================
CREATE TABLE IF NOT EXISTS `users` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `username` varchar(50) NOT NULL COMMENT '用户名（登录账号）',
  `phone` varchar(20) NOT NULL COMMENT '手机号',
  `password` varchar(255) NOT NULL COMMENT '密码（bcrypt加密）',
  `name` varchar(50) NOT NULL COMMENT '姓名',
  `role` varchar(20) DEFAULT 'operator' COMMENT '角色',
  `center_id` varchar(64) DEFAULT NULL COMMENT '所属充值中心ID',
  `center_name` varchar(100) DEFAULT NULL COMMENT '所属充值中心名称',
  `status` tinyint DEFAULT 1 COMMENT '状态：1=启用, 0=禁用',
  `last_login_at` datetime(3) DEFAULT NULL,
  `last_login_ip` varchar(50) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_phone` (`phone`),
  UNIQUE KEY `uk_username` (`username`),
  KEY `idx_users_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ========================================
-- 3. 插入超管账号（密码: 123456）
-- ========================================
INSERT IGNORE INTO `users` (`username`, `phone`, `password`, `name`, `role`, `status`, `created_at`, `updated_at`)
VALUES ('admin', '13800000000', '$2a$10$vXPjbfC511sMp3zdk1uFzOfxRWmtsZXnNIX7buP4C9Aq6In5YhV5S', '超级管理员', 'super_admin', 1, NOW(), NOW());

-- ========================================
-- 4. Casbin 公共权限（所有已登录用户可访问）
-- ========================================
INSERT IGNORE INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES
('p', 'super_admin', '/api/v1/auth/login', 'POST'),
('p', 'super_admin', '/api/v1/auth/refresh', 'POST'),
('p', 'hq_admin', '/api/v1/auth/login', 'POST'),
('p', 'hq_admin', '/api/v1/auth/refresh', 'POST'),
('p', 'finance', '/api/v1/auth/login', 'POST'),
('p', 'finance', '/api/v1/auth/refresh', 'POST'),
('p', 'center_admin', '/api/v1/auth/login', 'POST'),
('p', 'center_admin', '/api/v1/auth/refresh', 'POST'),
('p', 'operator', '/api/v1/auth/login', 'POST'),
('p', 'operator', '/api/v1/auth/refresh', 'POST');

-- ========================================
-- 5. hq_admin 权限（全部业务路由）
-- ========================================
INSERT IGNORE INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES
('p', 'hq_admin', '/api/v1/dashboard/*', 'GET'),
('p', 'hq_admin', '/api/v1/recharge/b-apply', 'POST'),
('p', 'hq_admin', '/api/v1/recharge/b-approval', 'GET'),
('p', 'hq_admin', '/api/v1/recharge/b-approval/*', 'GET'),
('p', 'hq_admin', '/api/v1/recharge/b-approval/action', 'POST'),
('p', 'hq_admin', '/api/v1/recharge/c-entry', 'GET'),
('p', 'hq_admin', '/api/v1/recharge/c-entry/*', 'GET'),
('p', 'hq_admin', '/api/v1/recharge/c-entry/search-member', 'GET'),
('p', 'hq_admin', '/api/v1/recharge/c-entry', 'POST'),
('p', 'hq_admin', '/api/v1/recharge/records', 'GET'),
('p', 'hq_admin', '/api/v1/recharge/records/*', 'GET'),
('p', 'hq_admin', '/api/v1/card/*', 'GET'),
('p', 'hq_admin', '/api/v1/card/*', 'POST'),
('p', 'hq_admin', '/api/v1/center', 'GET'),
('p', 'hq_admin', '/api/v1/center/*', 'GET'),
('p', 'hq_admin', '/api/v1/center', 'POST'),
('p', 'hq_admin', '/api/v1/center/*', 'PUT'),
('p', 'hq_admin', '/api/v1/center/*', 'DELETE'),
('p', 'hq_admin', '/api/v1/operator', 'GET'),
('p', 'hq_admin', '/api/v1/operator', 'POST'),
('p', 'hq_admin', '/api/v1/operator/*', 'PUT'),
('p', 'hq_admin', '/api/v1/operator/*', 'DELETE'),
('p', 'hq_admin', '/api/v1/users', 'GET'),
('p', 'hq_admin', '/api/v1/users', 'POST'),
('p', 'hq_admin', '/api/v1/users/*', 'GET'),
('p', 'hq_admin', '/api/v1/users/*', 'PUT'),
('p', 'hq_admin', '/api/v1/users/*', 'DELETE');

-- ========================================
-- 6. finance 权限（充值审批 + Dashboard + 查看）
-- ========================================
INSERT IGNORE INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES
('p', 'finance', '/api/v1/dashboard/*', 'GET'),
('p', 'finance', '/api/v1/recharge/b-approval', 'GET'),
('p', 'finance', '/api/v1/recharge/b-approval/*', 'GET'),
('p', 'finance', '/api/v1/recharge/b-approval/action', 'POST'),
('p', 'finance', '/api/v1/recharge/records', 'GET'),
('p', 'finance', '/api/v1/recharge/records/*', 'GET'),
('p', 'finance', '/api/v1/center', 'GET'),
('p', 'finance', '/api/v1/center/*', 'GET'),
('p', 'finance', '/api/v1/operator', 'GET');

-- ========================================
-- 7. center_admin 权限（B端申请/审批 + C端充值 + 卡管理 + Dashboard）
-- ========================================
INSERT IGNORE INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES
('p', 'center_admin', '/api/v1/dashboard/*', 'GET'),
('p', 'center_admin', '/api/v1/recharge/b-apply', 'POST'),
('p', 'center_admin', '/api/v1/recharge/b-approval', 'GET'),
('p', 'center_admin', '/api/v1/recharge/b-approval/*', 'GET'),
('p', 'center_admin', '/api/v1/recharge/b-approval/action', 'POST'),
('p', 'center_admin', '/api/v1/recharge/c-entry', 'GET'),
('p', 'center_admin', '/api/v1/recharge/c-entry', 'POST'),
('p', 'center_admin', '/api/v1/recharge/c-entry/search-member', 'GET'),
('p', 'center_admin', '/api/v1/recharge/c-entry/*', 'GET'),
('p', 'center_admin', '/api/v1/recharge/records', 'GET'),
('p', 'center_admin', '/api/v1/recharge/records/*', 'GET'),
('p', 'center_admin', '/api/v1/card/*', 'GET'),
('p', 'center_admin', '/api/v1/card/*', 'POST'),
('p', 'center_admin', '/api/v1/center', 'GET'),
('p', 'center_admin', '/api/v1/center/*', 'GET'),
('p', 'center_admin', '/api/v1/operator', 'GET'),
('p', 'center_admin', '/api/v1/operator', 'POST'),
('p', 'center_admin', '/api/v1/operator/*', 'PUT');

-- ========================================
-- 8. operator 权限（C端充值 + 卡查询 + Dashboard）
-- ========================================
INSERT IGNORE INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES
('p', 'operator', '/api/v1/dashboard/*', 'GET'),
('p', 'operator', '/api/v1/recharge/b-apply', 'POST'),
('p', 'operator', '/api/v1/recharge/c-entry', 'GET'),
('p', 'operator', '/api/v1/recharge/c-entry', 'POST'),
('p', 'operator', '/api/v1/recharge/c-entry/search-member', 'GET'),
('p', 'operator', '/api/v1/recharge/c-entry/*', 'GET'),
('p', 'operator', '/api/v1/recharge/records', 'GET'),
('p', 'operator', '/api/v1/recharge/records/*', 'GET'),
('p', 'operator', '/api/v1/card/*', 'GET'),
('p', 'operator', '/api/v1/card/*', 'POST');
