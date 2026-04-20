## Context

太积堂管理系统前端（Vue 3 + Element Plus）目前缺少统一的设计基础设施：

- **Logo/品牌资源**：Login 页用 120px 金色纯 div 圆作占位，Sidebar 用渐变色 div + "太" 字作图标，favicon 是 Vite 默认闪电 SVG
- **主题系统**：15+ 组件各自硬编码 `#C00000`、`#FFD700`、`#F5F5F5` 等颜色值，无 CSS 变量。Element Plus 使用默认蓝色主题，仅通过 inline style 或 `:deep()` 逐个覆盖
- **字体**：全局声明 `font-family: 'Inter'` 但从未导入字体文件，实际回退到系统字体。中文场景下 Inter 并非最佳选择
- **图标**：菜单用 emoji（📊💰🎫👥🏦⚙️🚪），跨平台渲染不一致
- **残留代码**：`style.css` 包含整个 Vite 模板样式（紫色 #aa3bff 主题、dark mode、1126px max-width），与实际应用冲突

用户已提供真实 logo 文件：`logo.png`（带透明背景）和 `logo_1.jpg`（带背景版本）。

## Goals / Non-Goals

**Goals:**
- 替换所有占位 logo/favicon 为真实品牌资源
- 建立基于 CSS 变量的主题 token 系统，统一管理颜色、字体、间距
- 定制 Element Plus 全局主题，使默认组件样式与品牌一致（红金配色）
- 将所有 emoji 图标替换为 Element Plus Icons（已安装已全局注册）
- 清理 Vite 模板残留样式

**Non-Goals:**
- 页面布局/组件结构的美化（留给后续 change）
- 响应式设计改造
- ECharts 图表替换（手写 div chart → ECharts）
- 新增功能或修改业务逻辑

## Decisions

### 1. 主题方案：CSS 变量 + Element Plus SCSS 变量覆盖

**选择**：在 `style.css` 中定义 `:root` CSS 变量作为 design tokens，同时在 `main.ts` 中通过 Element Plus 的 CSS 变量覆盖（`--el-color-primary` 等）实现全局主题定制。

**理由**：
- Element Plus 2.x 支持通过 CSS 变量覆盖主题（不需要 SCSS 编译），与我们的 CSS 变量方案天然兼容
- 不引入新依赖（如 SCSS preprocessor），保持项目简洁
- 所有组件直接引用 `var(--color-primary)` 而非硬编码色值

**替代方案**：
- SCSS 变量覆盖：需要安装 sass、配置 vite preprocess，增加构建复杂度，且项目目前无 SCSS
- CSS-in-JS：与 Vue SFC scoped style 生态不符，过度工程

### 2. 字体方案：系统字体栈（不引入外部字体）

**选择**：使用以 PingFang SC / Microsoft YaHei 为主的中文系统字体栈，去掉 Inter 引用。

```css
font-family: -apple-system, BlinkMacSystemFont, 'PingFang SC', 'Hiragino Sans GB',
  'Microsoft YaHei', 'Helvetica Neue', sans-serif;
```

**理由**：
- 管理系统用户以中文为主，PingFang SC（macOS）和 Microsoft YaHei（Windows）覆盖主要平台
- 系统字体零加载延迟，无需引入 Google Fonts CDN（国内访问不稳定）
- 当前 Inter 从未加载，实际已是系统字体回退，只是声明不一致

### 3. Logo 使用方案

**选择**：
- **Sidebar**：使用 `logo.png`（透明背景），收缩时显示小 logo（32x32），展开时显示 logo + 系统名
- **Login 页**：使用 `logo.png` 替换金色圆形占位 div
- **Favicon**：从 `logo.png` 生成 SVG favicon 或直接使用 PNG favicon

**理由**：logo.png 带透明背景，适配深色（sidebar）和红色（login）两种背景

### 4. 图标替换映射

**选择**：使用已全局注册的 Element Plus Icons Vue 组件替换 emoji：

| 原 emoji | 新图标组件 | 用途 |
|----------|-----------|------|
| 📊 | `DataAnalysis` | 数据概览 |
| 💰 | `Wallet` | 充值管理 |
| 🏦 | `OfficeBuilding` | 充值中心 |
| 🎫 | `Ticket` | 门店卡 |
| 👥 | `User` | 用户管理 |
| ⚙️ | `Setting` | 系统设置 |
| 🚪 | `SwitchButton` | 退出登录 |

**理由**：Element Plus Icons 已在 `main.ts` 中全局注册，直接用 `<el-icon><Xxx /></el-icon>` 即可，与 Element Plus 菜单组件原生配合

### 5. style.css 处理方案

**选择**：完全重写 `style.css`，删除所有 Vite 模板内容，替换为项目级 CSS 变量定义和全局基础样式。

**理由**：现有 style.css 的内容 100% 是 Vite 模板残留，与实际应用无任何关系。部分规则（如 `#app { width: 1126px }`）与 App.vue 中的 `#app { width: 100%; height: 100vh }` 冲突。

## Risks / Trade-offs

- **[CSS 变量兼容性]** → 现代浏览器均支持，Element Plus 2.x 本身也依赖 CSS 变量，无实际风险
- **[Logo 尺寸适配]** → logo.png 需要测试在 sidebar 小尺寸（32px）下是否清晰，可能需要单独制作小尺寸版本 → 可先尝试 CSS 缩放，效果不佳再补充资源
- **[Element Plus 主题覆盖范围]** → CSS 变量覆盖无法覆盖所有组件的所有状态色（如部分组件用 hardcoded 色值），后续可能需要补充 `:deep()` 覆盖 → 第一轮覆盖主要 token（primary/success/warning/danger），观察效果
- **[图标语义匹配]** → Element Plus Icons 不一定有完美匹配的图标（如"门店卡"对应的 Ticket 图标） → 选最接近的，后续可考虑自定义 SVG
