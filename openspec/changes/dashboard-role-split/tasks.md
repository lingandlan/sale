## 1. 后端 — 移除 memberCount

- [x] 1.1 `service/recharge.go` GetDashboardStatistics 返回值移除 `memberCount` 和 `memberTrend` 字段
- [x] 1.2 验证：`go build ./...` 编译通过

## 2. 前端 — 类型 + 移除总会员数卡片

- [x] 2.1 `types/dashboard.ts` Statistics 接口移除 `memberCount` 和 `memberTrend`
- [x] 2.2 `Dashboard.vue` 移除"总会员数"统计卡片（User icon 那张）
- [x] 2.3 验证：总部角色访问 Dashboard，确认只剩 3 张卡片（充值、核销、活跃中心）

## 3. 前端 — 条件渲染快捷操作

- [x] 3.1 `Dashboard.vue` 引入 userStore，用 `canSelectAllCenters` 判断是否总部人员
- [x] 3.2 "B端充值申请"快捷操作加 `v-if="isHeadquarters"` 条件
- [x] 3.3 验证：center_admin 角色只看到 3 个快捷操作

## 4. 前端 — 条件渲染统计卡片和待办

- [x] 4.1 "活跃中心数"卡片加 `v-if="isHeadquarters"` 条件
- [x] 4.2 整个待办事项区域加 `v-if="isHeadquarters"` 条件
- [x] 4.3 验证：center_admin 角色只看到 2 张卡片（充值、核销），无待办区域

## 5. 集成验证

- [ ] 5.1 用总部账号（13900000001）访问 Dashboard，确认显示 3 张卡片 + 4 个快捷操作 + 待办 + 趋势图
- [ ] 5.2 用充值中心账号访问 Dashboard，确认显示 2 张卡片 + 3 个快捷操作 + 无待办 + 趋势图仅本中心数据
