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
-- 初始数据：太积堂系统默认角色权限
-- =============================================================================

-- 超级管理员角色：拥有所有权限
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES
('p', 'super_admin', '/api/v1/*', '*');

-- 总部管理员角色：充值中心、门店、操作员管理
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES
('p', 'admin', '/api/v1/recharge-centers/*', '*'),
('p', 'admin', '/api/v1/stores/*', '*'),
('p', 'admin', '/api/v1/operators/*', '*'),
('p', 'admin', '/api/v1/recharge-applications/*', '*'),
('p', 'admin', '/api/v1/dashboard/*', 'GET');

-- 充值中心财务/运营角色：充值申请、审批
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES
('p', 'center_finance', '/api/v1/recharge-applications', 'POST'),
('p', 'center_finance', '/api/v1/recharge-applications/*', 'GET'),
('p', 'center_finance', '/api/v1/recharge-centers/*', 'GET'),
('p', 'center_finance', '/api/v1/dashboard/*', 'GET');

-- 充值中心操作员角色：C端充值录入
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES
('p', 'operator', '/api/v1/recharge-records', 'POST'),
('p', 'operator', '/api/v1/recharge-records/*', 'GET'),
('p', 'operator', '/api/v1/member-cards/*', 'GET'),
('p', 'operator', '/api/v1/member-cards/consume', 'POST'),
('p', 'operator', '/api/v1/mall/members/*', 'GET'),
('p', 'operator', '/api/v1/recharge-centers/*', 'GET');

-- 门店操作员角色：会员卡核销
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES
('p', 'store_operator', '/api/v1/member-cards/*', 'GET'),
('p', 'store_operator', '/api/v1/member-cards/consume', 'POST'),
('p', 'store_operator', '/api/v1/stores/*', 'GET');

-- 通用权限：所有角色都有
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES
('p', 'super_admin', '/api/v1/auth/login', 'POST'),
('p', 'admin', '/api/v1/auth/login', 'POST'),
('p', 'center_finance', '/api/v1/auth/login', 'POST'),
('p', 'operator', '/api/v1/auth/login', 'POST'),
('p', 'store_operator', '/api/v1/auth/login', 'POST'),
('p', 'super_admin', '/api/v1/auth/refresh', 'POST'),
('p', 'admin', '/api/v1/auth/refresh', 'POST'),
('p', 'center_finance', '/api/v1/auth/refresh', 'POST'),
('p', 'operator', '/api/v1/auth/refresh', 'POST'),
('p', 'store_operator', '/api/v1/auth/refresh', 'POST');

-- =============================================================================
-- 回滚语句
-- =============================================================================
-- DROP TABLE IF EXISTS `casbin_rule`;
