## ADDED Requirements

### Requirement: 请求结构体绑定与校验

所有业务接口 SHALL 使用带 validator tag 的结构体接收请求参数，通过 Gin 的 `ShouldBindJSON` 自动解析和校验，禁止使用 `map[string]interface{}` 直接接收用户输入。

涉及接口及校验规则：

| 接口 | 结构体字段 | 校验规则 |
|------|-----------|----------|
| CreateBRechargeApplication | memberId, centerId, amount, paymentMethod | required, amount>0 |
| ApprovalRechargeApplication | id, action(approve/reject), reason | required |
| CreateCRecharge | memberId, centerId, amount | required, amount>0 |
| CreateCenter | name, address, phone | required, min=1 |
| UpdateCenter | name, address, phone | omitempty |
| CreateOperator | name, phone, password, centerId, role | required, password min=6 |
| UpdateOperator | name, phone, password, role, status, centerId | omitempty |

#### Scenario: 合法请求通过校验

- **WHEN** 客户端发送符合结构体校验规则的 JSON 请求
- **THEN** 请求正常解析为结构体，handler 调用 service 处理并返回 200

#### Scenario: 缺少必填字段

- **WHEN** 客户端发送的 JSON 缺少 `binding:"required"` 标注的字段
- **THEN** 返回 400，错误信息明确指出缺失字段名

#### Scenario: 字段值不满足约束

- **WHEN** 客户端发送的 amount 为负数或 password 少于 6 位
- **THEN** 返回 400，错误信息说明具体约束

#### Scenario: 类型不匹配

- **WHEN** 客户端发送字符串到期望数字的字段
- **THEN** 返回 400，不会触发 panic
