## 1. 主题基础设施

- [x] 1.1 重写 `style.css`：删除所有 Vite 模板残留内容，定义 `:root` CSS 变量（颜色、字体、间距、圆角、阴影 token）
- [x] 1.2 在 `style.css` 中覆盖 Element Plus CSS 变量（`--el-color-primary: #C00000` 等梯度色、`--el-font-family`、`--el-border-radius-base`）
- [x] 1.3 更新 `App.vue` body 字体声明：移除 `'Inter'` 引用，改为 `var(--font-family)` 或系统字体栈

## 2. 组件样式迁移

- [x] 2.1 `Login.vue`：所有硬编码色值（#C00000、#FFD700、#F5F5F5、#E5E5E5、#262626 等）替换为 CSS 变量引用；`font-family: 'Inter'` 替换为 `var(--font-family)`；移除登录按钮的 inline `backgroundColor: '#C00000'` 覆盖（Element Plus 主题已覆盖）
- [x] 2.2 `Sidebar.vue`：硬编码色值替换为 CSS 变量；`font-family: 'Inter'` 替换
- [x] 2.3 `Header.vue`：硬编码色值替换为 CSS 变量（如有）
- [x] 2.4 `Dashboard.vue` 及子组件（StatCard、QuickAction 等）：硬编码色值替换为 CSS 变量；emoji 图标（如有）替换
- [x] 2.5 其余 views 目录下的页面组件：批量替换 `font-family: 'Inter'` 和硬编码色值为 CSS 变量

## 3. 品牌资源集成

- [x] 3.1 将 `logo.png` 复制到 `src/assets/logo.png`
- [x] 3.2 `Login.vue`：将 `.logo-circle` 占位 div 替换为 `<img src="@/assets/logo.png">` 真实 logo 图片
- [x] 3.3 `Sidebar.vue`：将 `.logo-icon` 渐变 div 替换为 logo 图片；展开状态显示 logo + 文字，收缩状态仅显示 logo
- [x] 3.4 替换 `public/favicon.svg`（或新增 `public/favicon.png`）为品牌 logo，更新 `index.html` 的 favicon 引用

## 4. 图标规范化

- [x] 4.1 `Sidebar.vue` 菜单图标映射：📊→`<DataAnalysis>`、💰→`<Wallet>`、🏦→`<OfficeBuilding>`、🎫→`<Ticket>`、👥→`<User>`、⚙️→`<Setting>`，使用 `<el-icon>` 包裹
- [x] 4.2 `Sidebar.vue` 退出登录图标：🚪→`<SwitchButton>`
- [x] 4.3 调整 Sidebar 菜单模板结构，适配 `el-icon` 组件替换 emoji 文本（`<span class="menu-icon">{{ group.icon }}</span>` → `<el-icon><component :is="group.icon" /></el-icon>`，icon 字段从 emoji 字符串改为组件名）

## 5. 验证

- [x] 5.1 启动前端，验证 Login 页 logo 显示、Sidebar logo 显示、favicon 显示
- [x] 5.2 验证所有页面 Element Plus 组件（按钮、输入框、菜单等）使用品牌红色主题，无蓝色泄漏
- [x] 5.3 验证 Sidebar 菜单图标全部为 SVG 图标，无 emoji 残留
- [x] 5.4 验证字体渲染：macOS 上 PingFang SC 正确加载，无外部字体请求
