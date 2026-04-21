-- C端充值记录增加 status 字段，用于补偿流程
-- pending: 已创建记录+扣中心余额，等待WSY加积分
-- success: 全部完成
-- failed: WSY加积分失败

ALTER TABLE c_recharges ADD COLUMN status VARCHAR(32) NOT NULL DEFAULT 'success' COMMENT 'pending/success/failed';
