## Why

上一轮已完成主题基础设施（CSS 变量、品牌 logo/favicon、Element Plus 图标替换），但各页面仍停留在功能布局阶段，视觉层次弱、间距不统一、缺少设计细节。需要对所有页面做视觉打磨，提升整体专业感和使用体验。

## What Changes

- **Login 页面**：优化左右面板比例、背景氛围、品牌展示区视觉层次
- **Dashboard 页面**：StatCard 重设计（更现代的卡片风格）、QuickAction 网格优化、欢迎区域美化、图表区域优化
- **列表页面统一**：统一筛选栏样式（背景、间距、按钮组布局）、表格样式优化（行高、斑马纹、悬浮效果）、分页组件居中对齐
- **表单/详情页面**：统一页面头部（返回按钮 + 标题 + 操作按钮）、表单卡片布局优化、对话框样式统一
- **共享组件升级**：StatCard 增加渐变/图标背景、QuickAction 增加悬浮动效、RechargeChart 优化柱状图样式

## Capabilities

### New Capabilities
- `login-page-polish`: 登录页视觉优化 — 品牌面板氛围、表单区布局、响应式比例
- `dashboard-polish`: 仪表盘视觉升级 — StatCard 重设计、欢迎区域、快捷操作、图表样式
- `list-page-unification`: 列表页样式统一 — 筛选栏、表格、分页、页面头部的一致性标准
- `form-dialog-polish`: 表单/详情/弹窗样式统一 — 卡片布局、表单间距、详情页信息网格

### Modified Capabilities
（无已有 capability 需要修改）

## Impact

- 影响范围：`src/views/` 下所有页面组件、`src/components/` 下共享组件
- 纯 CSS/模板层改动，不涉及业务逻辑和 API 接口变更
- 不引入新依赖，使用现有 Element Plus 组件 + CSS 变量体系
