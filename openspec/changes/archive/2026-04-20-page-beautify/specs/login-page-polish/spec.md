## 概述

Login 页面视觉优化，提升品牌氛围感和表单区精致度。

## 需求

### L1: 品牌面板装饰

- `.brand-section` 增加 CSS 伪元素装饰：1-2 个半透明圆环或几何线条，绝对定位，不遮挡 logo 和文字
- 装饰颜色使用 `rgba(255,255,255,0.08)` 或 `rgba(255,215,0,0.1)` 等品牌色低透明度变体

### L2: 表单卡片阴影

- `.login-card` 添加 `box-shadow: var(--shadow-md)` 替代当前纯边框
- 可选保留 `border` 但减弱为 `1px solid rgba(0,0,0,0.06)`

### L3: 输入框聚焦态

- 输入框聚焦时边框颜色从默认蓝色改为品牌红色 `var(--color-primary)`
- 通过全局 CSS 覆盖：`.el-input__wrapper.is-focus { box-shadow: 0 0 0 1px var(--color-primary) inset; }`

## 不涉及

- 不改变登录逻辑、表单校验、API 调用
- 不改变左右面板宽度比例
