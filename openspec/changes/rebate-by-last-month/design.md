## Context

B端充值申请的返还积分比例规则：上月净消费≥10万元返2%，否则返1%。PRD 定义数据来源为"商城净消费金额"，需调用商城接口获取。因商城对接有滞后性，目前采用手动录入方式记录每个充值中心每月的消费金额。

## Goals / Non-Goals

**Goals:**
- 新建 `center_monthly_consumption` 表存储每个充值中心每月的消费金额
- 提供手动录入/导入每月消费数据的接口
- B端申请页选中充值中心后查询上月消费，计算返还比例
- 申请记录存档上月消费金额，保证审批时可追溯

**Non-Goals:**
- 不修改返还比例规则本身（≥10万2%，否则1%）
- 不做商城接口自动对接（后续对接后可扩展）
- 不做返还比例的动态配置化

## Decisions

1. **新建表 `center_monthly_consumption`**：核心字段 center_id、month（YYYY-MM）、consumption（消费金额）。唯一约束 (center_id, month)，同一中心同月只能有一条记录
2. **手动录入方式**: 通过管理页面录入，支持单条和 Excel 批量导入
3. **查询上月消费**: `GET /api/v1/center/:id/last-month-consumption`，根据当前月份查上月的 consumption 记录，无记录默认 0（按1%算）
4. **前端交互**: B端申请页选中充值中心后调用接口获取上月消费，填充 `lastMonthConsumption`，实时更新积分预览
5. **存档**: recharge_applications 表新增 `last_month_consumption` 字段，提交时写入

## Risks / Trade-offs

- 手动录入有出错风险，后续对接商城接口后可自动化校验
- 上月消费为 0 或无记录的新中心，返还比例固定 1%，符合 PRD "首月处理"规则
