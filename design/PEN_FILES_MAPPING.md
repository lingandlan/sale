# Pencil 设计文件详细对应关系

## 🎯 重要说明

- **`.pen` 文件** = 设计源文件（在 `pencil/` 文件夹）
- **`.png` 文件** = 展示预览图（在 `pages/` 文件夹）
- **开发时使用 `.pen` 文件作为设计参考** ✅

---

## 📋 设计文件列表

### 1. 登录页设计

| 项目 | 值 |
|------|-----|
| **.pen 文件** | `pencil/3_login-page.pen` |
| **PNG 预览** | `pages/01-登录页.png` |
| **设计ID** | `FoyHo`（最终完整版） |
| **画布尺寸** | 1200 x 800 px（横向，适合PC） |

#### 设计结构（FoyHo - 完整版）

**左侧品牌区域**（560px宽）：
```
背景色: #C00000（深红）
布局: 垂直居中，内边距 80px
元素:
  - logo圆圈: 120x120，金色 #FFD700
  - 品牌名称"太积堂": 48px，金色 #FFD700，字体粗细 600
  - 系统标题"充值与门店管理系统": 18px，白色
```

**右侧表单区域**（剩余宽度）：
```
背景色: #FFFFFF
布局: 垂直居中，内边距 80px
元素（垂直布局，间距 24px）:
  - Tab切换（密码登录/验证码登录）
  - 手机号输入框
  - 密码输入框
  - 操作行（记住密码 + 忘记密码）
  - 登录按钮（#C00000，48px高）
  - 用户协议提示（12px，灰色）
```

#### 关键设计令牌

**颜色**：
- 主色（深红）: `#C00000`
- 辅助色（金色）: `#FFD700`
- 文本主色: `#262626`
- 文本辅色: `#595959`
- 文本占位符: `#BFBFBF`
- 边框色: `#D9D9D9`（输入框）、`#E5E5E5`（卡片）
- 页面背景: `#F5F5F5`
- 卡片背景: `#FFFFFF`

**尺寸**：
- 登录卡片: 480px 宽
- 输入框高度: 48px
- 登录按钮高度: 48px
- 内边距: 40px（卡片）、80px（区域）
- 圆角: 4px（按钮）、8px（卡片）

**字体**：
- 大标题: 48px / 600
- 系统标题: 18px / 400
- 表单标题: 28px / 400
- 表单标签: 14px / 400
- 按钮文字: 16px / 400
- 协议文字: 12px / 400

---

### 2. B端充值申请页

| 项目 | 值 |
|------|-----|
| **.pen 文件** | `pencil/3_login-page.pen` (文件中包含) |
| **设计ID** | `MLq5D` |
| **PNG 预览** | `pages/02-B端充值申请.png` |
| **画布尺寸** | 1200 x 800 px |

---

## 🔍 如何使用 .pen 文件开发

### 步骤1：打开设计文件查看

```bash
# 方式A：使用 Pencil 应用（推荐）
open /Users/zhangdaodong/code/sale/design/pencil/3_login-page.pen

# 方式B：使用 MCP 工具（需要 Pencil 运行中）
# 读取设计结构
mcp__pencil__batch_get({
  filePath: "/Users/zhangdaodong/code/sale/design/pencil/3_login-page.pen",
  nodeIds: ["FoyHo"],
  readDepth: 5
})
```

### 步骤2：提取设计令牌

从设计中提取：
- 颜色值（fill、stroke）
- 尺寸值（width、height）
- 间距值（gap、padding）
- 字体大小（fontSize）
- 圆角值（cornerRadius）

### 步骤3：按照设计实现

严格遵守设计规范，**不要随意修改**：
- ✅ 使用精确的颜色值
- ✅ 使用精确的尺寸值
- ✅ 使用精确的间距值
- ✅ 使用精确的字体大小

### 步骤4：对照验收

实现完成后，与设计稿对比：
- 布局结构是否一致
- 颜色是否准确
- 间距是否正确
- 字体大小是否匹配

---

## 📊 完整设计文件清单

| 序号 | .pen 文件 | 设计ID | PNG预览 | 页面名称 | 状态 | 优先级 |
|------|-----------|--------|---------|----------|------|--------|
| 1 | `pencil/3_login-page.pen` | FoyHo | 01-登录页.png | 登录页（完整版） | ✅ 有预览 | P0 |
| 2 | `pencil/3_login-page.pen` | NEtMi | - | 登录页（简化版） | ✅ | - |
| 3 | `pencil/3_login-page.pen` | MLq5D | 02-B端充值申请.png | B端充值申请 | ✅ 有预览 | P0 |
| 4 | `pencil/1_b_recharge_.pen` | - | - | B端充值申请发起 | ✅ | P0 |
| 5 | `pencil/2_b_recharge_list.pen` | - | 03-B端充值审批列表.png | B端充值审批列表 | ✅ 有预览 | P0 |
| 6 | `pencil/4_b_recharge_detail.pen` | - | - | B端充值审批详情 | ✅ | P1 |
| 7 | `pencil/5_c_rechange.pen` | - | - | C端充值录入 | ✅ | P0 |
| 8 | `pencil/6_card_hexiao.pen` | - | - | 门店卡核销 | ✅ | P0 |
| 9 | `pencil/7_card_manage.pen` | - | - | 门店卡管理 | ✅ | P1 |
| 10 | `pencil/8_card_send.pen` | - | - | 门店卡发放 | ✅ | P1 |
| 11 | `pencil/9_card_detail.pen` | - | - | 门店卡详情 | ✅ | P1 |
| 12 | `pencil/10_card_tongji.pen` | - | - | 门店卡统计 | ✅ | P2 |
| 13 | `pencil/11_recharge_list.pen` | - | - | 充值记录列表 | ✅ | P1 |
| 14 | `pencil/12_system.pen` | - | - | 系统设置 | ✅ | P2 |
| 15 | `pencil/13_user_mange.pen` | - | - | 用户管理 | ✅ | P1 |
| 16 | `pencil/14_recharge_center_management.pen` | - | - | 充值中心管理 | ✅ | P1 |
| 17 | `pencil/15_dashboard.pen` | - | - | 首页仪表盘 | ✅ | P0 |
| 18 | `pencil/16_recharge_detail.pen` | - | - | 充值记录详情 | ✅ | P1 |
| 19 | `pencil/17_recharge_operater.pen` | - | - | 充值操作员管理 | ✅ | P1 |

---

## ⚠️ 开发规则（强制）

1. **所有页面必须先查看 .pen 设计文件，再开始开发** ✅
2. **严格遵守设计令牌，不要随意修改** ✅
3. **颜色、尺寸、间距、字体必须精确匹配** ✅
4. **实现完成后与设计稿对照验收** ✅

---

## 🎯 下一步：PC端登录页开发

### 设计源文件
- **.pen 文件**: `pencil/3_login-page.pen`
- **设计ID**: `FoyHo`（完整版）
- **布局**: 横向（1200x800），左侧品牌区 + 右侧表单区

### 技术栈
- **框架**: Vue 3 + Vite
- **UI库**: Element Plus
- **路由**: Vue Router
- **状态管理**: Pinia
- **HTTP**: Axios

### 开发前确认
1. **是否开始实现 PC端登录页？**
2. **是否需要我先查看更多设计细节？**

记住：**不会直接起手做，一定先确认设计** ✅
