## 1. 接口层

- [x] 1.1 `api/recharge.ts` — `RechargeRecordItem` 接口移除 `paymentMethod`，补充 `operatorId` + `operatorName`

## 2. 列表页

- [x] 2.1 `RecordList.vue` — 删除表格"支付方式"列（el-table-column）
- [x] 2.2 `RecordList.vue` — 新增"操作员"列，展示 `operatorName`，宽度 120px

## 3. 详情页

- [x] 3.1 `RecordDetail.vue` — 删除详情卡中"支付方式"行
- [x] 3.2 `RecordDetail.vue` — 保留"操作员"行（已展示 operatorName，无需修改）

## 4. 验证

- [x] 4.1 启动前后端服务，打开 /recharge/records，确认无支付方式列，有操作员列
- [x] 4.2 点击某条记录进入详情页，确认无支付方式行，有操作员姓名
