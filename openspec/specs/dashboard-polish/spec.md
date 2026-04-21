## 概述

Dashboard 页面及共享组件的视觉升级。

## 需求

### D1: StatCard 渐变图标升级

- 卡片左侧增加色块区域（约 60px 宽），使用语义色渐变背景
- 图标移入色块区域，改为白色渲染
- 保留现有 props 接口（icon, value, label, trend, valueColor, prefix）
- 图标映射：emoji 字符 → 对应语义色（💰→红, 🎫→绿, 👥→蓝, 🏢→橙）

### D2: QuickAction 悬浮动效

- `.quick-action` 增加 `transition: all 0.2s`
- hover 时 `transform: translateY(-2px)` + `box-shadow` 加强
- 去除 Dashboard 和 QuickAction 中的 emoji，改为 Element Plus 图标或 CSS 图标

### D3: 欢迎区域微调

- 欢迎横幅增加品牌金色底部边线（2px `var(--color-primary-gold)`）
- 背景可考虑从纯色改为品牌色 → 深红色微渐变

### D4: RechargeChart 柱状图优化

- 柱体增加圆角顶部（`border-radius` 仅上方）
- 柱体增加 hover 高亮效果
- Y 轴网格线使用虚线样式
- 增加数值标签（hover 时显示具体金额）

### D5: 待办事项区域

- 去除 emoji 图标，改用 Element Plus 图标（Warning、Clock）
- todo-item 增加 hover 过渡效果

## 不涉及

- 不改变 Dashboard 数据加载逻辑
- 不增加新的统计卡片或快捷操作
