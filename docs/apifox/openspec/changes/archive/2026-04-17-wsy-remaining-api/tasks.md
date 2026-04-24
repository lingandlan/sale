## 1. WSYClient 积分操作方法

- [x] 1.1 在 `backend/pkg/mall/client.go` 新增 `AddUserIntegral` 方法，调用 `10000_integral_add`，返回 afterIntegral + error
- [x] 1.2 在 `AddUserIntegral` 中处理 errcode=100100 幂等场景，视为成功不返回错误
- [x] 1.3 在 `backend/pkg/mall/client.go` 新增 `ReduceUserIntegral` 方法，调用 `10000_integral_reduce`，返回 afterIntegral + error

## 2. WSYClient 零钱操作方法（预留）

- [x] 2.1 在 `backend/pkg/mall/client.go` 新增零钱 act 常量（TBD 占位）
- [x] 2.2 新增 `GetUserWallet(userID)` 方法骨架，act 使用占位常量
- [x] 2.3 新增 `AddUserWallet(userID, amount, batchcode, unionID)` 方法骨架
- [x] 2.4 新增 `ReduceUserWallet(userID, amount, batchcode, unionID)` 方法骨架

## 3. Batchcode 生成工具

- [x] 3.1 在 `backend/pkg/mall/client.go` 新增 `generateBatchcode(businessID string) string` 辅助函数，规则：前10位时间戳 + businessID，总长度 ≤ 30

## 4. Service 层集成

- [x] 4.1 在 `backend/internal/service/member.go` 新增 `AddIntegral(phone, integral, batchcode, remark)` 方法（PhoneToUserID → AddUserIntegral）
- [x] 4.2 更新 `MemberServiceInterface` 接口定义（如有）

## 5. 验证

- [x] 5.1 后端编译通过（`cd backend && go build ./...`）
- [x] 5.2 启动后端，通过 search-member 接口验证现有功能不受影响
- [x] 5.3 手动测试 AddUserIntegral（通过充值流程触发，验证积分变动）
