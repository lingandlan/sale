## ADDED Requirements

### Requirement: Logo 资源部署
系统 SHALL 将真实品牌 logo 文件（logo.png）部署到前端静态资源目录 `src/assets/` 中，供 Login 页和 Sidebar 引用。

#### Scenario: Logo 文件存在
- **WHEN** 检查 `src/assets/logo.png`
- **THEN** 文件存在且为太积堂品牌 logo 图片

### Requirement: Login 页 Logo 替换
Login 页面左侧品牌区域的金色圆形占位 div（`.logo-circle`）MUST 替换为真实 logo 图片（`src/assets/logo.png`），并保持居中展示。

#### Scenario: Login 页显示真实 Logo
- **WHEN** 用户访问登录页
- **THEN** 左侧品牌区域显示太积堂真实 logo 图片，而非纯色圆形占位

### Requirement: Sidebar Logo 替换
Sidebar 顶部的渐变色 div（`.logo-icon` 含 "太" 字）MUST 替换为真实 logo 图片。展开状态下显示 logo + 系统名称文字，收缩状态下仅显示 logo 缩略图。

#### Scenario: Sidebar 展开状态显示 Logo
- **WHEN** Sidebar 处于展开状态（240px）
- **THEN** 顶部显示真实 logo 图片 + "太积堂管理系统" 文字

#### Scenario: Sidebar 收缩状态显示 Logo
- **WHEN** Sidebar 处于收缩状态（64px）
- **THEN** 顶部仅显示真实 logo 缩略图，无文字

### Requirement: Favicon 替换
浏览器标签页的 favicon MUST 从 Vite 默认闪电图标替换为太积堂品牌图标（使用 logo.png 作为 favicon）。

#### Scenario: 浏览器标签显示品牌图标
- **WHEN** 用户在浏览器中打开系统
- **THEN** 浏览器标签页显示太积堂品牌 logo，而非 Vite 闪电图标
