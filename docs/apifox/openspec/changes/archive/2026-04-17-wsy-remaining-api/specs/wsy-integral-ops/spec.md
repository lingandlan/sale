## ADDED Requirements

### Requirement: Add user integral
WSYClient SHALL provide `AddUserIntegral(userID, integral, changeType, batchcode, remark, unionID)` 方法，调用 WSY `10000_integral_add` 接口，返回变动后积分 `afterIntegral`。

#### Scenario: Successful integral add
- **WHEN** 调用 AddUserIntegral(userID="965298296", integral=100, changeType="recharge", batchcode="171330000012345", remark="充值中心充值", unionID="recharge_001")
- **THEN** WSY 返回 errcode=0，方法返回 afterIntegral=50392.47 且 err=nil

#### Scenario: Idempotent duplicate request
- **WHEN** 调用 AddUserIntegral 且 WSY 返回 errcode=100100（union_id 重复）
- **THEN** 方法不返回错误，返回 afterIntegral=0 且 err=nil（视为幂等成功）

#### Scenario: WSY returns error
- **WHEN** 调用 AddUserIntegral 且 WSY 返回 errcode 非 0 且非 100100
- **THEN** 方法返回 err 包含 errcode 和 errmsg 信息

### Requirement: Reduce user integral
WSYClient SHALL provide `ReduceUserIntegral(userID, integral, changeType, batchcode, unionID)` 方法，调用 WSY `10000_integral_reduce` 接口，返回变动后积分。

#### Scenario: Successful integral reduce
- **WHEN** 调用 ReduceUserIntegral(userID="965298296", integral=50, changeType="refund", batchcode="171330000012345", unionID="refund_001")
- **THEN** WSY 返回 errcode=0，方法返回 afterIntegral=50292.47 且 err=nil

#### Scenario: Insufficient integral
- **WHEN** 调用 ReduceUserIntegral 且用户积分不足
- **THEN** WSY 返回错误，方法返回 err 包含错误信息

### Requirement: Batchcode generation rule
积分操作的 batchcode 参数 SHALL 遵循规则：前10位为 Unix 时间戳，后接业务 ID，总长度不超过30位。

#### Scenario: Batchcode from recharge ID
- **WHEN** 充值记录 ID 为 12345，当前时间戳为 1713300000
- **THEN** batchcode 为 "171330000012345"

#### Scenario: Batchcode length limit
- **WHEN** 业务 ID 导致 batchcode 总长度超过30位
- **THEN** 截断业务 ID 部分，确保总长度 ≤ 30

### Requirement: Service layer integral operation
MemberService SHALL 提供 `AddIntegral(phone, integral, batchcode, remark)` 方法，内部调用 PhoneToUserID → AddUserIntegral，返回变动后积分。

#### Scenario: Full flow via phone number
- **WHEN** 调用 AddIntegral(phone="17615860006", integral=100, batchcode="171330000012345", remark="充值")
- **THEN** 方法先通过手机号获取 userID，再调用 AddUserIntegral，返回变动后积分
