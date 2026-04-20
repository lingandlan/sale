## ADDED Requirements

### Requirement: 菜单图标替换为 Element Plus Icons
Sidebar 中所有 emoji 菜单图标 MUST 替换为 Element Plus Icons 组件（`<el-icon><Xxx /></el-icon>`），使用 `main.ts` 中已全局注册的图标组件。

映射关系：
| 菜单组 | 原 emoji | 替换图标组件 |
|--------|----------|-------------|
| 数据概览 | 📊 | DataAnalysis |
| 充值管理 | 💰 | Wallet |
| 充值中心 | 🏦 | OfficeBuilding |
| 门店卡 | 🎫 | Ticket |
| 用户管理 | 👥 | User |
| 系统设置 | ⚙️ | Setting |

#### Scenario: 菜单显示 SVG 图标
- **WHEN** 用户查看 Sidebar 菜单
- **THEN** 每个菜单组标题前显示对应 Element Plus SVG 图标，而非 emoji 字符

#### Scenario: 图标跨平台一致
- **WHEN** 在不同操作系统（macOS/Windows/Linux）上查看 Sidebar
- **THEN** 图标渲染效果一致，无平台差异

#### Scenario: 图标在收缩状态下显示
- **WHEN** Sidebar 收缩到 64px
- **THEN** 图标仍然正常显示在菜单项中

### Requirement: 退出登录图标替换
Sidebar 底部的退出登录 emoji（🚪）MUST 替换为 Element Plus 的 `SwitchButton` 图标组件。

#### Scenario: 退出登录使用规范图标
- **WHEN** 用户查看 Sidebar 底部
- **THEN** 退出登录按钮显示 SwitchButton SVG 图标，而非 🚪 emoji
