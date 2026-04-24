## Tasks

- [ ] SQL: 新建 center_monthly_consumption 表（id, center_id, month, consumption, created_at, updated_at），唯一约束 (center_id, month)
- [ ] SQL: recharge_applications 表新增 last_month_consumption 字段
- [ ] 后端: model 新增 CenterMonthlyConsumption GORM 模型
- [ ] 后端: repository 新增 CRUD 方法（UpsertMonthlyConsumption、GetMonthlyConsumption、ListMonthlyConsumption）
- [ ] 后端: service 新增 GetCenterLastMonthConsumption（查上月记录，算返还比例）
- [ ] 后端: service 新增 UpsertMonthlyConsumption（录入/更新）和 ImportMonthlyConsumption（Excel 批量导入）
- [ ] 后端: handler 新增 GetCenterLastMonthConsumption、UpsertMonthlyConsumption、ImportMonthlyConsumption、ListMonthlyConsumption
- [ ] 后端: router 注册新路由，补充 Casbin 权限规则
- [ ] 后端: service CreateBRechargeApplication 存档 lastMonthConsumption 到申请记录
- [ ] 前端: api/recharge.ts 新增上月消费查询和消费录入接口
- [ ] 前端: BRechargeApply.vue 选中充值中心后调用接口获取上月消费和返还比例，移除硬编码判断
- [ ] 前端: 新增消费数据管理页面（录入 + Excel 导入），hq_admin/finance 角色可见
- [ ] 验证: 录入消费数据 → B端申请查看返还比例 → 提交确认存档
