## 设计决策

### 后端：已有的角色过滤机制

后端 3 个 Dashboard API 已通过 `getOperatorInfo(c)` 从 JWT 提取 role 和 centerID，center_admin/operator 自动按 centerID 过滤数据。**后端无需改动过滤逻辑**。

需改动的：
- `GetDashboardStatistics`：移除 `memberCount` / `memberTrend` 字段（当前硬编码为 0 和 "—"）
- 返回结构保持 `todayRecharge`, `rechargeTrend`, `todayConsumption`, `consumptionTrend`, `activeCenters`, `centerTrend` 六个字段

### 前端：条件渲染策略

利用 userStore 已有的 `canSelectAllCenters` 计算属性（super_admin / hq_admin / finance 为 true），无需新增角色判断逻辑。

#### 统计卡片
- `isHeadquarters`（canSelectAllCenters）: 今日充值金额 + 今日核销金额 + 活跃中心数
- `!isHeadquarters`: 今日充值金额 + 今日核销金额

#### 快捷操作
- `isHeadquarters`: 4 项全展示
- `!isHeadquarters`: v-if 过滤，只展示 C端充值录入、门店卡核销、绑定卡号

#### 待办事项
- `isHeadquarters`: 展示（待审批申请、即将过期卡）
- `!isHeadquarters`: 隐藏整个 todo 区域

#### 充值趋势图
- 后端已按角色过滤，前端无需额外处理

### 类型变更

`Statistics` 接口移除 `memberCount` 和 `memberTrend` 字段。

## 涉及文件

| 文件 | 改动 |
|------|------|
| `shop-pc/src/types/dashboard.ts` | 移除 memberCount/memberTrend |
| `shop-pc/src/views/Dashboard.vue` | 移除总会员数卡片、条件渲染快捷操作/待办/活跃中心数 |
| `shop-pc/src/api/dashboard.ts` | 无改动 |
| `backend/internal/service/recharge.go` | GetDashboardStatistics 移除 memberCount/memberTrend |
