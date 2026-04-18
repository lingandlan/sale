-- 将 users.center_id 从 bigint unsigned 改为 varchar(64)，与 recharge_centers.id 类型一致
ALTER TABLE users MODIFY COLUMN center_id varchar(64) DEFAULT NULL COMMENT '所属充值中心ID';
