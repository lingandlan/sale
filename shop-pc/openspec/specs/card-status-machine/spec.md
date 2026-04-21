## MODIFIED Requirements

### Requirement: 卡状态转换规则
系统 SHALL 按以下规则执行卡状态转换：
- 任何非冻结状态（已入库1、已发放2、已激活3、已过期5）的卡 SHALL 可被冻结（转为4）
- 冻结(4)状态的卡 SHALL 只能解冻（转为3），不能执行核销、发放等其他操作
- 过期状态由系统自动标记（VerifyCard 检查 expiredAt），不通过手动状态转换触发

#### Scenario: 已入库卡冻结
- **WHEN** 对 status=1（已入库）的卡执行冻结
- **THEN** 状态变更为 4（已冻结），写入冻结交易记录

#### Scenario: 已发放卡冻结
- **WHEN** 对 status=2（已发放）的卡执行冻结
- **THEN** 状态变更为 4（已冻结），写入冻结交易记录

#### Scenario: 已激活卡冻结
- **WHEN** 对 status=3（已激活）的卡执行冻结
- **THEN** 状态变更为 4（已冻结），写入冻结交易记录

#### Scenario: 已过期卡冻结
- **WHEN** 对 status=5（已过期）的卡执行冻结
- **THEN** 状态变更为 4（已冻结），写入冻结交易记录

#### Scenario: 冻结卡解冻
- **WHEN** 对 status=4（已冻结）的卡执行解冻
- **THEN** 状态变更为 3（已激活），写入解冻交易记录

#### Scenario: 冻结卡不能核销
- **WHEN** 对 status=4（已冻结）的卡执行核销
- **THEN** 返回错误"卡已冻结"

## REMOVED Requirements

### Requirement: 作废卡功能
**Reason**: 业务上不需要作废操作，简化状态机
**Migration**: `/card/:cardNo/void` API 已移除，已有的 status=6 历史数据保留不变

#### Scenario: 作废 API 不可用
- **WHEN** 调用 POST /api/v1/card/:cardNo/void
- **THEN** 返回 404

## ADDED Requirements

### Requirement: 前端操作按钮一致性
前端 CardManage 页面操作列 SHALL 按以下规则显示按钮：
- 非冻结状态（1/2/3/5）：显示"详情"和"冻结"按钮
- 冻结状态(4)：显示"详情"和"解冻"按钮
- 不显示"作废"按钮

#### Scenario: 已激活卡的操作按钮
- **WHEN** 卡状态为 3（已激活）
- **THEN** 显示"详情"和"冻结"按钮，不显示"作废"按钮

#### Scenario: 冻结卡的操作按钮
- **WHEN** 卡状态为 4（已冻结）
- **THEN** 显示"详情"和"解冻"按钮，不显示"冻结"和"作废"按钮

### Requirement: 总卡库统计不含作废
CardInventory 页面统计区 SHALL 展示 6 种计数：总卡数、已入库、已发放、已激活、已冻结、已过期。不展示已作废计数。

#### Scenario: 统计区展示
- **WHEN** 用户访问总卡库管理页面
- **THEN** 统计区显示：总卡数、已入库、已发放、已激活、已冻结、已过期，共 6 项
