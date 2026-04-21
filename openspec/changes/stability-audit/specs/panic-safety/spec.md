## ADDED Requirements

### Requirement: Map parameter safe parsing helpers
The system SHALL provide `getFloat64`, `getString`, `getInt` helper functions that safely extract values from `map[string]interface{}` without triggering panic. When a key is missing or the value type does not match, the helper SHALL return an `errno` business error with a descriptive message.

#### Scenario: Missing required field
- **WHEN** a service function calls `getString(data, "memberPhone")` and `data` does not contain key "memberPhone"
- **THEN** the function SHALL return `("", errno.New(errno.CodeInvalidParam, "参数 memberPhone 不能为空"))` without panicking

#### Scenario: Wrong value type
- **WHEN** a service function calls `getFloat64(data, "amount")` and `data["amount"]` is a string "100" instead of float64
- **THEN** the function SHALL return `(0, errno.New(errno.CodeInvalidParam, "参数 amount 类型错误"))` without panicking

#### Scenario: Correct value type
- **WHEN** a service function calls `getFloat64(data, "amount")` and `data["amount"]` is float64(100.0)
- **THEN** the function SHALL return `(100.0, nil)`

### Requirement: Replace all bare type assertions in service layer
All `data["key"].(type)` bare assertions in `internal/service/recharge.go` SHALL be replaced with the safe helper functions. The affected functions include `CreateBRechargeApplication`, `CreateCRecharge`, `CreateCenter`, `CreateOperator`, and any other function using `map[string]interface{}` input.

#### Scenario: C recharge with missing memberPhone
- **WHEN** a C recharge request is submitted without "memberPhone" field
- **THEN** the service SHALL return a business error (not panic), and the handler SHALL return HTTP 400 with the error message

#### Scenario: B recharge application with all valid fields
- **WHEN** a B recharge application is submitted with all required fields of correct types
- **THEN** the service SHALL process normally, identical to current behavior

### Requirement: Recovery middleware logs panic stack traces
The `Recovery()` middleware in `internal/middleware/recovery.go` SHALL log the panic error and full stack trace using `logger.Error()` when a panic is recovered. The `_ = stack` suppression SHALL be removed.

#### Scenario: Handler triggers panic
- **WHEN** a handler function panics with "runtime error: index out of range"
- **THEN** the recovery middleware SHALL:
  1. Log the panic message and full stack trace via `logger.Error`
  2. Return HTTP 500 with error response
  3. Abort the request

#### Scenario: No panic occurs
- **WHEN** a request completes without panic
- **THEN** the recovery middleware SHALL add no overhead (defer function executes but recover() returns nil)

### Requirement: Safe context helper functions
The `GetUserID`, `GetPhone`, `GetRole` functions in `internal/middleware/auth.go` SHALL use comma-ok type assertion pattern. When the context value is missing or wrong type, they SHALL return the zero value instead of panicking.

#### Scenario: GetUserID called without auth middleware
- **WHEN** `GetUserID(c)` is called and `c.Get("user_id")` returns nil
- **THEN** the function SHALL return `int64(0)` instead of panicking

### Requirement: Safe type assertions in RBAC middleware
All bare `role.(string)` assertions in `internal/middleware/rbac.go` SHALL use comma-ok pattern with proper fallback handling.

#### Scenario: Role not set in context
- **WHEN** RBAC middleware reads `c.Get("role")` and role is nil or not a string
- **THEN** the middleware SHALL return HTTP 401 Unauthorized instead of panicking
