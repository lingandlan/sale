## ADDED Requirements

### Requirement: UpdateCenter 接受 status 参数
后端 `PUT /center/:id` 接口 SHALL 接受 `status` 字段，合法值为 `active` 和 `frozen`。`UpdateCenterRequest` 模型 SHALL 包含 `Status string` 字段，带 `binding:"omitempty,oneof=active frozen"` 验证。Handler SHALL 将 status 传入 service 层的 data map。

#### Scenario: 冻结充值中心
- **WHEN** 前端调用 `PUT /center/:id` 传入 `{ "status": "frozen" }`
- **THEN** 后端更新该中心 status 为 `frozen`，返回更新后的中心数据

#### Scenario: 解冻充值中心
- **WHEN** 前端调用 `PUT /center/:id` 传入 `{ "status": "active" }`
- **THEN** 后端更新该中心 status 为 `active`，返回更新后的中心数据

#### Scenario: 非法 status 值
- **WHEN** 前端调用 `PUT /center/:id` 传入 `{ "status": "invalid" }`
- **THEN** 后端返回参数错误 400

### Requirement: UpdateCenter 接受完整字段
`UpdateCenterRequest` SHALL 支持以下可选字段：`name`, `address`, `phone`, `status`, `province`, `city`, `district`, `managerId`。Handler SHALL 将所有非空字段传入 service 层。

#### Scenario: 编辑充值中心信息
- **WHEN** 前端调用 `PUT /center/:id` 传入 `{ "name": "新名称", "province": "山东省", "city": "济南市", "district": "历下区", "managerId": "op-123" }`
- **THEN** 后端更新所有传入字段，返回更新后的中心数据

### Requirement: GetCenters 正确关联管理员
`GetCenters` 返回数据中的 `managerName` 和 `managerPhone` SHALL 通过 `managerId`（操作员 ID）匹配操作员，而非通过 `centerId` 匹配。

#### Scenario: 中心有管理员
- **WHEN** 充值中心的 `managerId` 对应的操作员存在
- **THEN** 返回数据包含 `managerName`（操作员姓名）和 `managerPhone`（操作员手机号）

#### Scenario: 中心无管理员
- **WHEN** 充值中心的 `managerId` 为空或对应操作员不存在
- **THEN** 返回数据不包含 `managerName` 和 `managerPhone` 字段

### Requirement: 前端冻结/解冻操作
前端 CenterManage 页面 SHALL 在用户点击冻结/解冻按钮后，调用 `updateCenter(id, { status: "frozen"/"active" })` 并刷新列表。

#### Scenario: 点击冻结按钮
- **WHEN** 用户点击状态为 `normal` 的中心的"冻结"按钮并确认
- **THEN** 前端调用 `PUT /center/:id` 传入 `{ status: "frozen" }`，列表刷新后该中心显示为"冻结"状态

#### Scenario: 点击解冻按钮
- **WHEN** 用户点击状态为 `frozen` 的中心的"解冻"按钮并确认
- **THEN** 前端调用 `PUT /center/:id` 传入 `{ status: "active" }`，列表刷新后该中心显示为"正常"状态
