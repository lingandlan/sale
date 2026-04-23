-- 门店卡模块重构：对齐 PRD v1.4
-- 状态机从 3 种改为 6 种，新增发放记录表，余额改为整数
-- 开发阶段，直接 DROP 旧表重建

DROP TABLE IF EXISTS card_transactions;
DROP TABLE IF EXISTS card_issue_records;
DROP TABLE IF EXISTS store_cards;

-- 门店卡主表
CREATE TABLE store_cards (
    id VARCHAR(64) PRIMARY KEY,
    card_no VARCHAR(32) UNIQUE NOT NULL COMMENT '卡号 TJ00000001',
    card_type TINYINT DEFAULT 1 COMMENT '1=实体卡,2=虚拟卡',
    status TINYINT DEFAULT 1 COMMENT '1=已入库,2=已发放,3=已激活,4=已冻结,5=已过期,6=已作废',
    balance INT DEFAULT 1000 COMMENT '余额（元），面值固定1000',
    recharge_center_id VARCHAR(64) DEFAULT NULL COMMENT '划拨到的充值中心ID',
    user_id VARCHAR(64) DEFAULT NULL COMMENT '绑定的用户ID',
    batch_no VARCHAR(64) DEFAULT '' COMMENT '批次号',
    issue_reason VARCHAR(64) DEFAULT '' COMMENT '发放原因:购买套餐包/推荐奖励/其他',
    issued_at DATETIME DEFAULT NULL COMMENT '发放时间',
    activated_at DATETIME DEFAULT NULL COMMENT '激活时间（首次核销）',
    expired_at DATETIME DEFAULT NULL COMMENT '过期时间（激活日+1年）',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    KEY idx_status (status),
    KEY idx_recharge_center_id (recharge_center_id),
    KEY idx_user_id (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='门店卡';

-- 门店卡发放记录表
CREATE TABLE card_issue_records (
    id VARCHAR(64) PRIMARY KEY,
    card_no VARCHAR(32) NOT NULL COMMENT '卡号',
    user_id VARCHAR(64) NOT NULL COMMENT '用户ID',
    user_phone VARCHAR(32) NOT NULL COMMENT '用户手机号',
    issue_reason VARCHAR(64) NOT NULL COMMENT '发放原因:购买套餐包/推荐奖励/其他',
    issue_type TINYINT NOT NULL COMMENT '1=实体卡（运营绑定）,2=虚拟卡（用户领取）',
    recharge_center_id VARCHAR(64) NOT NULL COMMENT '充值中心ID',
    operator_id VARCHAR(64) NOT NULL COMMENT '操作员ID',
    related_user_phone VARCHAR(32) DEFAULT '' COMMENT '推荐奖励时关联购买人手机号',
    remark VARCHAR(500) DEFAULT '' COMMENT '备注',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    KEY idx_card_no (card_no),
    KEY idx_user_id (user_id),
    KEY idx_recharge_center_id (recharge_center_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='门店卡发放记录';

-- 门店卡交易记录表
CREATE TABLE card_transactions (
    id VARCHAR(64) PRIMARY KEY,
    card_no VARCHAR(32) NOT NULL COMMENT '卡号',
    type VARCHAR(32) NOT NULL COMMENT 'issue/consume/freeze/unfreeze/activate/void',
    amount INT DEFAULT 0 COMMENT '金额（元）',
    balance_before INT DEFAULT 0 COMMENT '交易前余额（元）',
    balance_after INT DEFAULT 0 COMMENT '交易后余额（元）',
    remark VARCHAR(500) DEFAULT '' COMMENT '备注',
    operator_id VARCHAR(64) DEFAULT '' COMMENT '操作员ID',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    KEY idx_card_no (card_no)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='门店卡交易记录';
