-- 补充 c_recharges 表缺少的 status 列
ALTER TABLE c_recharges ADD COLUMN status varchar(32) NOT NULL DEFAULT 'success' AFTER balance_after;
