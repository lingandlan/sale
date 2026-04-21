-- 充值中心月度消费记录表（手动录入商城消费数据）
CREATE TABLE IF NOT EXISTS center_monthly_consumption (
    id VARCHAR(64) NOT NULL PRIMARY KEY,
    center_id VARCHAR(64) NOT NULL COMMENT '充值中心ID',
    month VARCHAR(7) NOT NULL COMMENT '月份 YYYY-MM',
    consumption DOUBLE NOT NULL DEFAULT 0 COMMENT '消费金额',
    created_at DATETIME(3) DEFAULT NULL,
    updated_at DATETIME(3) DEFAULT NULL,
    UNIQUE KEY uk_center_month (center_id, month),
    INDEX idx_center_id (center_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- B端充值申请存档上月消费金额
ALTER TABLE recharge_applications ADD COLUMN last_month_consumption DOUBLE NOT NULL DEFAULT 0;
