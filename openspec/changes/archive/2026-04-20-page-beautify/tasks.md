## 1. CSS 基础层

- [x] 1.1 `style.css` 新增变量：`--color-primary-hover: #A00000`、`--color-bg-header: #FAFAFA`、`--color-bg-section: #F9F9F9`
- [x] 1.2 `style.css` 新增全局表格表头样式规则，替代 inline `:header-cell-style`
- [x] 1.3 `style.css` 新增全局 el-input 聚焦态品牌红色边框覆盖
- [x] 1.4 `style.css` 新增全局 el-dialog 圆角和按钮样式覆盖

## 2. 硬编码色值清理

- [x] 2.1 全局搜索替换 `#A00000` → `var(--color-primary-hover)`（约 5 个文件）
- [x] 2.2 全局搜索替换 `#FAFAFA` → `var(--color-bg-header)`（约 4 个文件，不含 style.css 新增规则）
- [x] 2.3 全局搜索替换 `#F9F9F9` → `var(--color-bg-section)`（CRechargeEntry, BRechargeDetail 等）
- [x] 2.4 全局搜索替换 `#D9D9D9` → `var(--color-border)`（CRechargeEntry, CardManage 等）
- [x] 2.5 移除所有 `<el-table :header-cell-style="...">` inline 覆盖（6 个文件）

## 3. Login 页面

- [x] 3.1 `.brand-section` 增加 CSS 伪元素装饰（半透明圆环/几何线条）
- [x] 3.2 `.login-card` 增加 `box-shadow` 替代纯边框

## 4. Dashboard 组件升级

- [x] 4.1 `StatCard.vue` 重设计：左侧色块 + 白色图标，emoji → 语义色映射
- [x] 4.2 `QuickAction.vue` 增加 hover 浮起动效（translateY + shadow）
- [x] 4.3 `RechargeChart.vue` 优化：柱体圆角顶部 + hover 高亮 + Y轴虚线网格
- [x] 4.4 `Dashboard.vue` emoji 清理：快捷操作/待办标题/待办图标改用 Element Plus 图标或纯文字

## 5. 列表页统一

- [x] 5.1 `BRechargeList.vue` 统一：`.filter-bar` → `.filter-card`、`.table-card` → `.list-card`、自定义分页 → `<el-pagination>`
- [x] 5.2 `RecordList.vue` 统一筛选栏：添加 label、统一 `.filter-card` 样式
- [x] 5.3 `CenterManage.vue` / `UserManage.vue` 统一筛选栏标签和按钮样式
- [x] 5.4 `CardManage.vue` / `CardStats.vue`：补充表头样式（已被全局规则覆盖，确认移除 inline）、stat card inline style 改 CSS class
- [x] 5.5 `OperatorManage.vue`：确认 `.list-card` / `.filter-card` / `.pagination-row` 模式一致

## 6. 详情/表单页统一

- [x] 6.1 `BRechargeDetail.vue` / `RecordDetail.vue` / `CardDetail.vue`：统一 `.page-header` + 返回按钮模式
- [x] 6.2 `BRechargeApply.vue` / `CRechargeEntry.vue`：统一居中卡片布局，去除硬编码背景色
- [x] 6.3 `CardVerify.vue` / `CardIssue.vue`：统一表单卡片样式
- [x] 6.4 所有弹窗（OperatorManage, CenterManage, UserManage, CardInventory）：确认 `el-dialog` 使用全局圆角样式

## 7. 验证

- [x] 7.1 启动前端，逐页截图验证：Login、Dashboard、BRechargeList、BRechargeApply、BRechargeDetail、RecordList、RecordDetail、CardManage、CardStats、CardVerify、CenterManage、UserManage、OperatorManage、SystemConfig
- [x] 7.2 确认无硬编码色值残留（grep `#A00000`、`#FAFAFA`、`#F9F9F9`、`#D9D9D9`）
- [x] 7.3 确认全局输入框聚焦态为品牌红色（非默认蓝色）
