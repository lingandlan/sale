## 概述

统一表单页、详情页、弹窗的视觉样式。

## 需求

### FD1: 详情页统一头部

BRechargeDetail, RecordDetail, CardDetail 的页面头部统一为：
- 使用 `.page-header` 样式（与列表页一致）
- 返回按钮统一：`<el-button link @click="$router.back()"><el-icon><ArrowLeft /></el-icon> 返回</el-button>`
- 标题 + 右侧操作按钮

### FD2: 信息卡片统一

详情页的信息展示区域统一为 `.info-card` 样式：
- `background: var(--color-bg-card); border-radius: var(--radius-md); border: 1px solid var(--color-border); padding: 24px`
- 信息网格：2 列 label + value 布局，label 用 `var(--color-text-muted)` 14px，value 用 `var(--color-text-primary)` 14px
- 去除硬编码 `#F9F9F9`，改用 `var(--color-bg-section)`

### FD3: 表单页布局统一

BRechargeApply, CRechargeEntry, CardIssue, CardVerify 统一：
- 居中卡片布局，`max-width: 800px; margin: 0 auto`
- 卡片样式与 `.info-card` 一致
- 表单 `label-width` 统一为 `120px`
- 操作按钮区域右对齐，统一间距

### FD4: 弹窗样式统一

所有 `el-dialog` 统一：
- 圆角 `border-radius: var(--radius-lg)`
- 头部底部分割线
- 底部按钮右对齐，主按钮使用 `.save-btn` 样式

涉及：OperatorManage, CenterManage, UserManage, CardInventory 的编辑弹窗

### FD5: 去除 inline style

- CardManage 的 stat card inline `style="background:...;border-color:..."` 改为 CSS class
- 所有 `style="width: 240px"` 等内联宽度改用 CSS class 或 `style` 绑定统一管理

## 不涉及

- 不改变表单校验规则和提交逻辑
- 不改变弹窗打开/关闭逻辑
