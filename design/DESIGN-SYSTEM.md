# 太积堂 Design System

## 概述

基于太积堂品牌 VI 系统创建的设计系统，支持浅色和暗色两种主题。

## 文件结构

```
design/
├── taijido-design-system.pen    # Design System 主文件
├── logo/
│   └── logo_1.png               # 太积堂品牌 Logo
└── pages/
    └── login-page.pen           # 登录页面原型
```

## 设计令牌

### 颜色系统

#### 浅色主题
- **主色**: `#C00000` (太积堂品牌红 - 深红)
- **辅助色**: `#FFD700` (金色)
- **成功色**: `#52C41A` (绿色)
- **警告色**: `#FAAD14` (橙色)
- **错误色**: `#FF4D4F` (红色)
- **信息色**: `#1677FF` (蓝色)

#### 中性色 (灰度)
- 灰50: `#F0F0F0`
- 灰100: `#D9D9D9`
- 灰200: `#BFBFBF`
- 灰300: `#8C8C8C`
- 灰400: `#595959`
- 灰500: `#262626` (文本主色)

#### 暗色主题
- **背景色**: `#1F1F1F`
- **卡片背景**: `#2C2C2C`
- **文本主色**: `#E8E8E8`
- **文本辅色**: `#A0A0A0`
- **边框色**: `#3C3C3C`
- **主色**: `#D61919` (调整后的深红，在暗色背景上更明显)

### 字体系统

#### 字体家族
- **中文字体**: PingFang SC, Microsoft YaHei, sans-serif
- **英文字体**: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif

#### 字号
- H1 (标题1): 24px, fontWeight: 600
- H2 (标题2): 20px, fontWeight: 600
- H3 (标题3): 18px, fontWeight: 600
- H4 (标题4): 16px, fontWeight: 600
- Body 1 (正文1): 14px, fontWeight: 400
- Body 2 (正文2): 12px, fontWeight: 400
- Caption (辅助文字): 10px, fontWeight: 400

### 间距系统

- **xs**: 8px
- **sm**: 12px
- **md**: 16px
- **lg**: 24px
- **xl**: 32px

### 圆角系统

- **小圆角**: 2px
- **中圆角**: 4px
- **大圆角**: 8px

## 组件示例

### 按钮

#### 主要按钮
- 背景: `#C00000`
- 文字: `#FFFFFF`
- 高度: 48px
- Padding: [10, 16]
- 圆角: 4px

#### 次要按钮
- 背景: `#FFFFFF`
- 边框: `#C00000` (1px)
- 文字: `#C00000`
- 高度: 48px
- Padding: [10, 16]
- 圆角: 4px

### 输入框

- 高度: 40px (表单) / 48px (登录)
- Padding: [8, 16] / [12, 16]
- 边框: `#D9D9D9` (1px)
- 圆角: 4px
- 占位符颜色: `#BFBFBF`

### 卡片

- 背景: `#FFFFFF` (浅色) / `#2C2C2C` (暗色)
- 边框: `#E5E5E5` (1px)
- 圆角: 8px
- Padding: 24px

## 登录页面

### 设计特点

1. **分栏布局**: 左侧品牌区域 (560px) + 右侧表单区域 (自适应)
2. **品牌展示**:
   - 深红色背景 (#C00000)
   - 金色 Logo 圆形徽章
   - "太积堂" 品牌名称 (48px)
   - 系统标题 (18px)

3. **登录表单**:
   - 宽度: 480px
   - Tab 切换: 密码登录 / 验证码登录
   - 手机号输入
   - 密码/验证码输入
   - 记住密码 + 忘记密码
   - 登录按钮
   - 用户协议提示

### 交互元素

- **Tab 切换**: 密码登录 (选中，红色下划线) / 验证码登录 (灰色)
- **记住密码**: 复选框 + 文字
- **忘记密码**: 红色链接
- **登录按钮**: 主色按钮，全宽

## 使用 uView Plus

### 主题覆盖

在项目中使用 uView Plus 3.2 组件库时，通过以下方式覆盖主题变量：

```scss
// uni.scss
$u-primary: #C00000;
$u-success: #52C41A;
$u-warning: #FAAD14;
$u-error: #FF4D4F;
$u-info: #1677FF;

// 文字颜色
$u-main-color: #262626;
$u-content-color: #595959;
$u-tips-color: #8C8C8C;
$u-light-color: #BFBFBF;

// 背景颜色
$u-bg-color: #F5F5F5;

// 边框颜色
$u-border-color: #E5E5E5;
```

### 暗色主题配置

```scss
// dark.scss
$u-main-color: #E8E8E8;
$u-content-color: #A0A0A0;
$u-bg-color: #1F1F1F;
$u-border-color: #3C3C3C;
```

## 设计原则

1. **清晰性**: 信息层级分明，重点突出
2. **一致性**: 统一的视觉语言和交互模式
3. **高效性**: 减少操作步骤，提高完成效率
4. **友好性**: 明确的提示和错误处理
5. **品牌性**: 融入太积堂传统中式美学

## 下一步

- [ ] 创建其他核心页面原型 (B端充值、C端充值、门店卡管理)
- [ ] 设计响应式布局适配
- [ ] 创建 uView Plus 组件主题配置文件
- [ ] 导出设计规范文档
