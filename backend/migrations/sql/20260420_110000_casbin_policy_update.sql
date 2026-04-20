-- 补全 Casbin 策略，匹配当前路由结构
-- 角色映射: super_admin, hq_admin, finance, center_admin, operator

-- 先清理旧的不匹配策略（保留 login/refresh 公共权限）
DELETE FROM casbin_rule WHERE ptype = 'p' AND v1 NOT LIKE '/api/v1/auth%';

-- super_admin: 已在 RBAC 中间件硬编码跳过，不需要策略

-- hq_admin: 全部业务路由
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES
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
('p', 'hq_admin', '/api/v1/operator/*', 'DELETE');

-- finance: 充值审批 + Dashboard + 查看
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES
('p', 'finance', '/api/v1/dashboard/*', 'GET'),
('p', 'finance', '/api/v1/recharge/b-approval', 'GET'),
('p', 'finance', '/api/v1/recharge/b-approval/*', 'GET'),
('p', 'finance', '/api/v1/recharge/b-approval/action', 'POST'),
('p', 'finance', '/api/v1/recharge/records', 'GET'),
('p', 'finance', '/api/v1/recharge/records/*', 'GET'),
('p', 'finance', '/api/v1/center', 'GET'),
('p', 'finance', '/api/v1/center/*', 'GET'),
('p', 'finance', '/api/v1/operator', 'GET');

-- center_admin: B端申请/审批 + C端充值 + 卡管理 + Dashboard
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES
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

-- operator: C端充值 + 卡查询 + Dashboard
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES
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
