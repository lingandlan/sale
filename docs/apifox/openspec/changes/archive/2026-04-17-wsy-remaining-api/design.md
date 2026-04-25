## Context

当前 `pkg/mall/client.go` 中的 `WSYClient` 已实现：`GetAccessToken`、`PhoneToUserID`、`GetUserIntegral`、`postForm`（通用POST）。

所有 WSY 业务接口通过同一个 URL 调用：
```
POST {baseURL}/wsy_pub/third_api/index.php?m=third_application_authorization&a=third_function
```
通过 body 中的 `act` 参数区分不同接口。

现有 `postForm(endpoint, form)` 方法封装了 HTTP 调用和响应读取，新方法可复用。

## Goals / Non-Goals

**Goals:**
- 在 WSYClient 上新增 AddUserIntegral / ReduceUserIntegral 方法
- 在 WSYClient 上新增零钱相关方法（查询/增加/扣除，act 确认后）
- 在 MemberService 层暴露积分/零钱操作，供充值流程调用
- 积分添加支持幂等（union_id 防重复）
- batchcode 生成规则：前10位时间戳 + 业务ID

**Non-Goals:**
- 不新增对外 HTTP API（WSY 为内部调用层）
- 不修改前端代码
- 不实现零钱接口的具体代码（等 WSY 文档确认 act 值后再补）
- 不做 WSY 接口的 mock/单元测试（第三方接口，集成测试为主）

## Decisions

### 1. 方法签名设计

沿用现有 WSYClient 的风格，每个方法返回业务核心数据 + error：

```go
// AddUserIntegral 添加积分，返回变动后积分
func (c *WSYClient) AddUserIntegral(userID string, integral float64, changeType, batchcode, remark, unionID string) (afterIntegral float64, err error)

// ReduceUserIntegral 扣除积分，返回变动后积分
func (c *WSYClient) ReduceUserIntegral(userID string, integral float64, changeType, batchcode, unionID string) (afterIntegral float64, err error)
```

**理由**: 返回 `afterIntegral` 可直接用于更新充值记录的 `balance_after` 字段。

### 2. batchcode 生成策略

```go
// 前10位=Unix时间戳，后接业务ID（如充值记录ID），总长度 ≤ 30
batchcode := fmt.Sprintf("%d%d", time.Now().Unix(), rechargeID)
```

**理由**: WSY 要求前10位必须为时间戳，最大30位。充值记录ID天然唯一。

### 3. 幂等控制

使用 `union_id` 参数（最长32位），值为充值记录ID的字符串形式。WSY 对重复 union_id 返回 `errcode=100100`，客户端检测此错误码视为成功。

### 4. 零钱接口预留

先定义接口签名，act 值用占位常量，等文档确认后填入：

```go
// act 常量（待确认）
const ActWalletBalance = "TBD_wallet_balance"
const ActWalletAdd    = "TBD_wallet_add"
const ActWalletReduce = "TBD_wallet_reduce"

func (c *WSYClient) GetUserWallet(userID string) (float64, error)
func (c *WSYClient) AddUserWallet(userID string, amount float64, batchcode, unionID string) (float64, error)
func (c *WSYClient) ReduceUserWallet(userID string, amount float64, batchcode, unionID string) (float64, error)
```

### 5. Service 层集成

在 `service/member.go` 新增：

```go
func (s *MemberService) AddIntegral(phone string, integral float64, batchcode, remark string) (float64, error)
```

内部流程：`PhoneToUserID` → `AddUserIntegral`，供 `service/recharge.go` 的 `CreateCRecharge` 调用。

## Risks / Trade-offs

- **[WSY 零钱 act 值未确认]** → 积分接口先行实现，零钱接口预留占位，确认后快速填充
- **[WSY 接口不稳定/超时]** → 复用现有 10s HTTP 超时；积分操作失败不阻塞充值记录创建，记录状态标记为"积分同步失败"供后续补偿
- **[batchcode 冲突]** → 时间戳+充值ID组合保证全局唯一
- **[union_id 重复检测]** → WSY 返回 100100 时视为幂等成功，需在代码中处理此特殊情况
