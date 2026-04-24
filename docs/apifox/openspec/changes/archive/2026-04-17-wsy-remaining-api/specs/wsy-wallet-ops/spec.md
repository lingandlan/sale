## ADDED Requirements

### Requirement: Query user wallet balance
WSYClient SHALL provide `GetUserWallet(userID)` 方法，调用 WSY 零钱余额查询接口，返回零钱余额。

#### Scenario: Successful wallet query
- **WHEN** 调用 GetUserWallet(userID="965298296")
- **THEN** WSY 返回 errcode=0，方法返回零钱余额（float64）且 err=nil

#### Scenario: User not found
- **WHEN** 调用 GetUserWallet 且 userID 不存在
- **THEN** WSY 返回错误，方法返回 err 包含错误信息

### Requirement: Add user wallet
WSYClient SHALL provide `AddUserWallet(userID, amount, batchcode, unionID)` 方法，调用 WSY 零钱增加接口，返回变动后余额。

#### Scenario: Successful wallet add
- **WHEN** 调用 AddUserWallet(userID="965298296", amount=50, batchcode="171330000012345", unionID="wallet_add_001")
- **THEN** WSY 返回 errcode=0，方法返回变动后余额且 err=nil

### Requirement: Reduce user wallet
WSYClient SHALL provide `ReduceUserWallet(userID, amount, batchcode, unionID)` 方法，调用 WSY 零钱扣除接口，返回变动后余额。

#### Scenario: Successful wallet reduce
- **WHEN** 调用 ReduceUserWallet(userID="965298296", amount=20, batchcode="171330000012345", unionID="wallet_reduce_001")
- **THEN** WSY 返回 errcode=0，方法返回变动后余额且 err=nil

#### Scenario: Insufficient wallet balance
- **WHEN** 调用 ReduceUserWallet 且用户零钱不足
- **THEN** WSY 返回错误，方法返回 err 包含错误信息

### Requirement: Wallet act values configurable
零钱接口的 act 值 SHALL 使用常量定义，待 WSY 文档确认后填入实际值。

#### Scenario: Placeholder act values
- **WHEN** 零钱接口 act 值尚未确认
- **THEN** 代码中使用 TBD 占位常量，编译通过但不实际调用

#### Scenario: Confirmed act values
- **WHEN** WSY 文档确认零钱接口 act 值
- **THEN** 替换占位常量为实际值，接口可正常调用
