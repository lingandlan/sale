## Context

项目已完成主题基础设施（CSS 变量体系、品牌 logo/favicon、Element Plus 图标替换），但各页面仍存在明显的样式不一致和视觉粗糙问题：

- **列表页面**：9 个页面中有 3 种不同的筛选栏样式、2 种分页实现、表头样式通过 inline `:header-cell-style` 硬编码
- **硬编码色值**：`#A00000`（hover 色）、`#FAFAFA`（表头背景）、`#F9F9F9`、`#D9D9D9` 等散落在多个文件中
- **组件视觉**：StatCard 是平铺白卡片、QuickAction 无悬浮反馈、RechargeChart 是基础 CSS 柱状图
- **Login 页面**：品牌面板单调（纯红背景），缺少视觉层次

## Goals / Non-Goals

**Goals:**
- 统一所有列表页的筛选栏、表格、分页组件样式
- 统一详情/表单页的页面头部和卡片布局模式
- 升级 Dashboard 共享组件（StatCard、QuickAction、RechargeChart）的视觉表现
- 提升 Login 页面的品牌氛围感
- 将所有硬编码色值收归 CSS 变量体系

**Non-Goals:**
- 不引入新的 UI 框架或第三方组件库
- 不改变页面布局结构（不增加/删除页面区域）
- 不涉及业务逻辑变更
- 不做响应式适配（仅 PC 端 1280px+ 宽度）
- 不做暗色模式

## Decisions

### D1: 新增 CSS 变量补充 token 缺口

在 `style.css` 的 `:root` 中补充以下变量：

```
--color-primary-hover: #A00000     /* 按钮/链接 hover 态 */
--color-bg-header: #FAFAFA         /* 表头背景 */
--color-bg-section: #F9F9F9        /* 信息区块背景 */
```

**理由**：当前 `#A00000` 在 5 个文件中硬编码为 hover 色、`#FAFAFA` 在 4 个文件中硬编码为表头背景。收归变量后一处修改全局生效。

### D2: 全局表格样式替代 inline header-cell-style

在 `style.css` 中添加全局 Element Plus 表头样式规则，移除所有 `<el-table :header-cell-style="...">` 的 inline 覆盖：

```css
.el-table th.el-table__cell {
  background-color: var(--color-bg-header) !important;
  color: var(--color-text-primary);
  font-weight: 600;
}
```

**理由**：4 个文件重复传入相同的 inline style 对象，1 个文件遗漏。全局规则更 DRY 且不会遗漏。

### D3: 列表页三件套标准化

所有列表页统一采用以下结构：

```
.page-header（标题 + 操作按钮）
.filter-card（筛选条件卡片，使用 el-form 行内布局）
.list-card（表格 + 分页）
```

- 筛选栏：统一 `.filter-card` 类名、使用 `<el-form :inline="true">` 布局
- 分页：统一 `<el-pagination>` 组件，`layout="total, sizes, prev, pager, next"`
- BRechargeList 的自定义分页替换为 `<el-pagination>`

### D4: StatCard 升级为渐变图标风格

每个 StatCard 添加图标背景色块（左上角带圆角的色块区域），图标使用白色 SVG 渲染在色块上。渐变方向从语义色到稍深色调。保留当前 props 接口不变。

### D5: Login 品牌面板增加装饰元素

在 `.brand-section` 增加微妙的装饰图案（CSS 伪元素绘制的几何线条/圆环），保持品牌红色主调，增加视觉层次。表单卡片增加 `box-shadow`。

### D6: 详情页统一页面头部组件

所有带返回按钮的详情页（BRechargeDetail、RecordDetail、CardDetail）统一使用 `.page-header` + `.back-btn` 模式，标题左对齐、操作按钮右对齐。

## Risks / Trade-offs

- **全局表格样式可能覆盖特定页面需求** → 使用 `:root` 级别的样式，如需覆盖可在页面 scoped style 中用 `:deep()` 覆盖
- **StatCard 改版可能影响 Dashboard 布局** → 保持卡片尺寸（280x120）不变，只修改内部渲染
- **BRechargeList 分页重写可能影响功能** → 分页逻辑保持 JS 层不变，只替换模板和样式
