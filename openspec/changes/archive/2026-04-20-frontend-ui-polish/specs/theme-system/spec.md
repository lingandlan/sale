## ADDED Requirements

### Requirement: CSS 变量主题 Token 系统
系统 SHALL 在 `style.css` 的 `:root` 中定义完整的 CSS 变量 token 集，覆盖颜色、字体、间距、圆角、阴影等设计属性。所有组件 MUST 引用 CSS 变量而非硬编码色值。

变量分类：
- **颜色**：`--color-primary`（#C00000）、`--color-primary-gold`（#FFD700）、`--color-bg`（#F5F5F5）、`--color-bg-card`（#FFFFFF）、`--color-border`（#E5E5E5）、`--color-text-primary`（#262626）、`--color-text-secondary`（#595959）、`--color-text-muted`（#8C8C8C）、语义色（success/warning/danger/info）
- **字体**：`--font-family`、`--font-size-*`（xs/sm/base/lg/xl/2xl/3xl）
- **间距**：`--spacing-*`（xs:4px / sm:8px / md:12px / base:16px / lg:24px / xl:32px）
- **圆角**：`--radius-sm`（4px）、`--radius-md`（8px）、`--radius-lg`（12px）
- **阴影**：`--shadow-sm`、`--shadow-md`

#### Scenario: 组件引用主题变量
- **WHEN** 开发者在任何 Vue SFC 中需要使用品牌色
- **THEN** MUST 使用 `var(--color-primary)` 而非硬编码 `#C00000`

#### Scenario: 全局主题修改
- **WHEN** 需要更改品牌主色
- **THEN** 只需修改 `style.css` 中 `:root` 的 `--color-primary` 值，所有组件自动生效

### Requirement: Element Plus 全局主题覆盖
系统 SHALL 通过 CSS 变量覆盖 Element Plus 默认主题，使所有 Element Plus 组件（按钮、输入框、复选框、菜单、对话框等）默认使用品牌色 `#C00000` 作为 primary 色。

覆盖的关键 Element Plus CSS 变量：
- `--el-color-primary` → `#C00000`
- `--el-color-primary-light-3` 至 `--el-color-primary-dark-2` → 基于 #C00000 的梯度色
- `--el-font-family` → 系统字体栈
- `--el-border-radius-base` → 4px

#### Scenario: Element Plus 按钮使用品牌色
- **WHEN** 页面使用 `<el-button type="primary">`
- **THEN** 按钮背景色为 #C00000，hover 状态自动变深，无需 inline style 覆盖

#### Scenario: Element Plus 输入框聚焦使用品牌色
- **WHEN** 用户聚焦 el-input
- **THEN** 输入框边框色为 #C00000 系列色，而非默认蓝色

### Requirement: 清理 Vite 模板残留样式
系统 MUST 删除 `style.css` 中所有 Vite 模板残留内容（紫色主题变量、dark mode 媒体查询、.hero/.counter/#next-steps 等 Vite 示例样式），仅保留项目级全局基础样式和 CSS 变量定义。

#### Scenario: style.css 无残留代码
- **WHEN** 检查 style.css 内容
- **THEN** 不包含 `--accent: #aa3bff`、`color-scheme: light dark`、`.hero`、`#next-steps`、`.counter`、`.ticks` 等 Vite 模板样式

### Requirement: 系统中文字体栈
系统 SHALL 使用以 PingFang SC / Microsoft YaHei 为主的系统字体栈，不引入外部字体文件或 CDN。所有 `font-family: 'Inter'` 声明 MUST 被替换为 `var(--font-family)` 或直接使用系统字体栈。

#### Scenario: 中文字体正确渲染
- **WHEN** 页面在 macOS 上加载
- **THEN** 中文文字使用 PingFang SC 渲染

#### Scenario: 中文字体正确渲染（Windows）
- **WHEN** 页面在 Windows 上加载
- **THEN** 中文文字使用 Microsoft YaHei 渲染

#### Scenario: 无外部字体请求
- **WHEN** 页面加载
- **THEN** 不向 Google Fonts 或任何外部 CDN 发起字体请求
