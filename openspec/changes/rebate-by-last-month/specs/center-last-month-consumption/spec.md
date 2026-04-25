## 概述

手动录入每个充值中心每月的消费金额，B端申请时查询上月消费计算返还比例。

## 数据库

### 新建表 center_monthly_consumption

| 字段 | 类型 | 说明 |
|------|------|------|
| id | VARCHAR(64) PK | UUID |
| center_id | VARCHAR(64) INDEX | 充值中心ID |
| month | VARCHAR(7) | 月份 YYYY-MM |
| consumption | DOUBLE | 消费金额 |
| created_at | DATETIME | 创建时间 |
| updated_at | DATETIME | 更新时间 |

唯一约束: UNIQUE(center_id, month)

### recharge_applications 变更

新增字段 `last_month_consumption` DOUBLE DEFAULT 0

## 接口

### GET /api/v1/center/:id/last-month-consumption

查询指定充值中心上月消费金额及返还比例。

**权限**: operator, center_admin, hq_admin, super_admin

**响应**:
```json
{
  "code": 0,
  "data": {
    "consumption": 120000.00,
    "rebateRate": 2,
    "month": "2026-03"
  }
}
```

### POST /api/v1/center-monthly-consumption

录入/更新某月消费数据（单条）。

**权限**: hq_admin, super_admin, finance

**请求**:
```json
{
  "centerId": "center-bj-cy",
  "month": "2026-03",
  "consumption": 120000.00
}
```

### POST /api/v1/center-monthly-consumption/import

Excel 批量导入消费数据。

**权限**: hq_admin, super_admin, finance

**请求**: multipart/form-data，Excel 包含列：充值中心ID、月份(YYYY-MM)、消费金额

### GET /api/v1/center-monthly-consumption

查询消费记录列表，支持按月份筛选。

**权限**: hq_admin, super_admin, finance

## 前端交互

1. B端申请页选中充值中心 → 调用 last-month-consumption → 显示上月消费和返还比例
2. 新增管理页面：录入/导入每月消费数据（hq_admin/finance 角色）
3. 移除前端硬编码的比例判断逻辑
