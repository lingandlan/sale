-- =============================================================================
-- 太积堂用户表 Migration
-- =============================================================================

-- 用户表（操作员、管理员）
CREATE TABLE IF NOT EXISTS `users` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `phone` VARCHAR(20) NOT NULL COMMENT '手机号（登录账号）',
    `password` VARCHAR(255) NOT NULL COMMENT '密码（bcrypt加密）',
    `name` VARCHAR(50) NOT NULL COMMENT '姓名',
    `role` VARCHAR(20) NOT NULL DEFAULT 'operator' COMMENT '角色：super_admin=超管, admin=管理员, operator=操作员',
    `center_id` BIGINT UNSIGNED COMMENT '所属充值中心ID',
    `center_name` VARCHAR(100) COMMENT '所属充值中心名称（冗余字段）',
    `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1=启用, 0=禁用',
    `last_login_at` DATETIME COMMENT '最后登录时间',
    `last_login_ip` VARCHAR(50) COMMENT '最后登录IP',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` DATETIME COMMENT '删除时间（软删除）',

    UNIQUE KEY `uk_phone` (`phone`),
    KEY `idx_center_id` (`center_id`),
    KEY `idx_role` (`role`),
    KEY `idx_status` (`status`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

-- 初始数据：超级管理员
-- 密码：admin123（需要通过bcrypt加密后替换）
INSERT INTO `users` (`phone`, `password`, `name`, `role`, `status`) VALUES
('13800000000', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iAt6Z5EH', '超级管理员', 'super_admin', 1);
