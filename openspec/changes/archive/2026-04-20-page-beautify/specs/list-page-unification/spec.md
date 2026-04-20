## 概述

统一所有列表页（含筛选、表格、分页）的视觉样式和布局模式。

## 需求

### LU1: 新增 CSS 变量

在 `style.css` 中补充：
- `--color-primary-hover: #A00000`
- `--color-bg-header: #FAFAFA`
- `--color-bg-section: #F9F9F9`

所有文件的硬编码色值 `#A00000` → `var(--color-primary-hover)`、`#FAFAFA` → `var(--color-bg-header)`。

### LU2: 全局表格样式

在 `style.css` 中添加全局表头样式规则，移除所有 `<el-table :header-cell-style="...">` 内联覆盖。

### LU3: 统一页面头部

所有列表页的 `.page-header` 统一为：
- `height: 64px; background: var(--color-bg-card); border-bottom: 1px solid var(--color-border); padding: 16px 24px; display: flex; justify-content: space-between; align-items: center`
- `.page-title` 统一：`font-size: 20px; font-weight: 600; color: var(--color-text-primary)`

涉及文件：BRechargeList, RecordList, CenterManage, UserManage, CardManage, OperatorManage

### LU4: 统一筛选栏

所有列表页筛选栏统一为 `.filter-card` 样式：
- `background: var(--color-bg-card); border-radius: var(--radius-md); border: 1px solid var(--color-border); padding: 16px`
- 使用 `<el-form :inline="true">` 布局，每个筛选项配 `<el-form-item label="xxx">`
- 搜索/重置按钮右对齐

BRechargeList 的 `.filter-bar` → `.filter-card`，补充 border。

### LU5: 统一分页

所有列表页分页统一为 `<el-pagination>`：
- `layout="total, sizes, prev, pager, next"`
- `:page-sizes="[10, 20, 50]"`
- 外层 `.pagination-row`：`display: flex; justify-content: center; margin-top: 16px`

BRechargeList 的自定义分页替换为 `<el-pagination>`。

### LU6: 表格卡片

统一为 `.list-card`：
- `background: var(--color-bg-card); border-radius: var(--radius-md); border: 1px solid var(--color-border); padding: 24px`

BRechargeList 的 `.table-card` → `.list-card`。

### LU7: 按钮 hover 色

所有 `.save-btn:hover` 和 `.search-btn:hover` 统一使用 `var(--color-primary-hover)` 替代硬编码 `#A00000`。

## 涉及文件

BRechargeList, RecordList, RecordDetail, CenterManage, UserManage, CardManage, CardStats, CardInventory, OperatorManage, style.css

## 不涉及

- 不改变筛选逻辑、表格列定义、分页数据加载
