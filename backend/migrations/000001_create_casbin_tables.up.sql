-- =============================================================================
-- Casbin Policy Tables Migration
-- =============================================================================

-- casbin_rule 表用于存储策略和角色关联
CREATE TABLE IF NOT EXISTS `casbin_rule` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `ptype` VARCHAR(128) NOT NULL DEFAULT '' COMMENT '策略类型: p=策略, g=角色关联',
    `v0` VARCHAR(128) NOT NULL DEFAULT '' COMMENT '角色/用户',
    `v1` VARCHAR(128) NOT NULL DEFAULT '' COMMENT '路径/角色',
    `v2` VARCHAR(128) NOT NULL DEFAULT '' COMMENT '方法',
    `v3` VARCHAR(128) NOT NULL DEFAULT '' COMMENT '额外字段',
    `v4` VARCHAR(128) NOT NULL DEFAULT '' COMMENT '额外字段',
    `v5` VARCHAR(128) NOT NULL DEFAULT '' COMMENT '额外字段',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_casbin_ptype` (`ptype`, `v0`, `v1`, `v2`),
    KEY `idx_casbin_v0` (`v0`),
    KEY `idx_casbin_v1` (`v1`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Casbin 权限策略表';

-- =============================================================================
-- 初始数据：系统默认角色权限
-- =============================================================================

-- 管理员角色：拥有所有权限
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES
('p', 'admin', '/api/v1/admin/*', '*');

-- 商家角色：商家相关权限
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES
('p', 'merchant', '/api/v1/merchant/*', '*'),
('p', 'merchant', '/api/v1/products/*', '*'),
('p', 'merchant', '/api/v1/orders/*', '*');

-- 普通用户角色：基础权限
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES
('p', 'user', '/api/v1/user/info', 'GET'),
('p', 'user', '/api/v1/user/info', 'PUT'),
('p', 'user', '/api/v1/products', 'GET'),
('p', 'user', '/api/v1/products/:id', 'GET'),
('p', 'user', '/api/v1/cart/*', '*'),
('p', 'user', '/api/v1/orders/*', '*');

-- =============================================================================
-- 回滚语句
-- =============================================================================
-- DROP TABLE IF EXISTS `casbin_rule`;
