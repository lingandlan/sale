## ADDED Requirements

### Requirement: 前端详情弹窗
CenterManage 页面 SHALL 在用户点击"详情"按钮时打开详情弹窗，调用 `GET /center/:id` 获取完整中心信息并展示。

#### Scenario: 查看充值中心详情
- **WHEN** 用户点击某中心行的"详情"按钮
- **THEN** 前端调用 `GET /center/:id`，在弹窗中展示中心名称、级别、省市区、地址、管理员、余额、累计充值、已消耗、状态等信息

#### Scenario: 详情数据加载失败
- **WHEN** 调用 `GET /center/:id` 返回错误
- **THEN** 前端显示错误提示信息，不打开弹窗

### Requirement: 前端详情 API 调用
前端 `center.ts` API 模块 SHALL 导出 `getCenterDetail(id)` 函数，调用 `GET /center/:id`。

#### Scenario: 调用详情 API
- **WHEN** 前端调用 `getCenterDetail("center-123")`
- **THEN** 发送 `GET /center/center-123` 请求，返回中心详情数据
