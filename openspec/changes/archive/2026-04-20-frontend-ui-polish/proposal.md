## Why

前端 UI 缺乏统一的品牌形象和设计系统：logo 是占位 div（纯色圆/渐变色块），favicon 是 Vite 默认图标，颜色/字体/间距在每个组件中硬编码，菜单图标使用 emoji 导致跨平台不一致。需要先建立品牌基础（真实 logo + favicon）和设计基础设施（主题变量 + 规范图标），为后续页面美化打下基础。

## What Changes

- **替换真实 Logo**：将 Login 页面和 Sidebar 的占位 div 替换为太积堂真实 logo 图片（logo.png），更新 favicon 为品牌图标
- **建立统一主题系统**：引入 CSS 变量管理品牌色（#C00000 红金 + #FFD700 金）、字体、间距、圆角等，定制 Element Plus 全局主题覆盖默认蓝色
- **清理残留样式**：删除 style.css 中 Vite 模板的紫色主题残留，修复 Inter 字体未加载的问题（引入 Google Fonts 或替换为系统中文字体）
- **图标规范化**：将所有 emoji 菜单图标（📊💰🎫👥等）替换为 Element Plus Icons（@element-plus/icons-vue，已安装）

## Capabilities

### New Capabilities
- `theme-system`: CSS 变量主题系统 + Element Plus 全局主题定制（品牌色、字体、间距、圆角统一管理）
- `brand-assets`: 品牌 logo/favicon 资源集成（Login 页、Sidebar、favicon 替换为真实品牌资源）
- `icon-standardization`: 图标规范化（emoji 全部替换为 Element Plus Icons，统一图标风格）

### Modified Capabilities
<!-- 无现有 spec 需要修改 -->

## Impact

- **前端代码**：涉及 `style.css`、`main.ts`（Element Plus 主题注入）、`Login.vue`、`Sidebar.vue`、`Header.vue`、`Dashboard.vue` 及所有使用 emoji 图标的组件
- **静态资源**：新增 logo 图片到 `src/assets/`，替换 `public/favicon.svg`
- **依赖**：无新依赖（Element Plus Icons 已安装）
- **无 API 变更**：纯前端变更，不影响后端
