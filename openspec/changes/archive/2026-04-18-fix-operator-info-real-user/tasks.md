## 1. 后端 handler

- [x] 1.1 新增 `getOperatorInfo()` 函数，替代 `getOperatorCenter()`，返回 `userID, role, centerID, name, error`
- [x] 1.2 修改 `CreateCRecharge` 调用 `getOperatorInfo()`，将返回值写入 `req["operatorId"]` 和 `req["operatorName"]`
- [x] 1.3 适配所有调用方 `getOperatorCenter` → `getOperatorInfo` 的返回值变化（共7处）

## 2. 验证

- [x] 2.1 API 测试：`POST /recharge/c-entry` 创建充值，验证 `operatorId` 和 `operatorName` 为真实用户信息
