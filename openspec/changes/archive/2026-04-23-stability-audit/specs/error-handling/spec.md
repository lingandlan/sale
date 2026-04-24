## ADDED Requirements

### Requirement: Logger nil protection
The `pkg/logger` package SHALL protect against nil `log` variable access. `GetLogger()` SHALL return `zap.NewNop()` when `log` is nil. All package-level convenience functions (Info, Error, Warn, Debug, Fatal) SHALL check for nil before calling methods on `log`.

#### Scenario: GetLogger called before Init
- **WHEN** `logger.GetLogger()` is called before `logger.Init()`
- **THEN** it SHALL return a no-op logger (`zap.NewNop()`) instead of causing a nil pointer panic

#### Scenario: Package-level function called before Init
- **WHEN** `logger.Info("message")` is called before `logger.Init()`
- **THEN** the call SHALL be silently ignored, not panic

### Requirement: GORM error propagation in repository
All GORM operations in `internal/repository/recharge.go` that currently discard errors SHALL be fixed:
- `.Count(&total)` return values SHALL be checked
- `.Scan(&result)` return values SHALL be checked
- Functions like `GetCenterTotalRecharge`, `GetCenterTotalConsumed`, `GetCenterCardStats` SHALL return `error` in their signature
- `gorm.go` `SET NAMES utf8mb4` Exec error SHALL be checked and logged

#### Scenario: Count query fails in GetRechargeApplications
- **WHEN** the GORM `.Count(&total)` query fails due to a database connection error
- **THEN** the function SHALL return `(nil, 0, err)` instead of silently returning zero total

#### Scenario: Scan fails in GetCardStats
- **WHEN** the GORM `.Scan(&result)` query fails
- **THEN** the function SHALL return the error to the caller

#### Scenario: GetCenterTotalRecharge signature change
- **WHEN** service code calls `GetCenterTotalRecharge(centerID)`
- **THEN** the function SHALL return `(int64, error)` instead of just `int64`

### Requirement: Casbin role sync error handling
The admin handler SHALL check the return value of Casbin role operations (`AddRoleForUser`, `UpdateUserRole`). If the Casbin operation fails, the handler SHALL return an error response to the client instead of silently ignoring the failure.

#### Scenario: Create user with Casbin failure
- **WHEN** a new user is created in the database but `casbinSvc.AddRoleForUser()` fails
- **THEN** the handler SHALL return an error response indicating the role assignment failed

#### Scenario: Update user role with Casbin failure
- **WHEN** `casbinSvc.UpdateUserRole()` fails during a user role update
- **THEN** the handler SHALL return an error response instead of returning success

### Requirement: Server startup and shutdown robustness
The HTTP server in `cmd/server/main.go` SHALL:
- Replace `log.Fatal` in the ListenAndServe goroutine with error logging and signal-based shutdown
- Set `ReadHeaderTimeout: 5 * time.Second` on the HTTP server
- Check the return value of `gormDB.DB()` for nil before calling `.Close()`

#### Scenario: ListenAndServe fails on startup
- **WHEN** `srv.ListenAndServe()` fails (e.g., port already in use)
- **THEN** the error SHALL be logged (not `log.Fatal`) and the shutdown signal SHALL be sent to trigger graceful cleanup

#### Scenario: Graceful shutdown with GORM DB
- **WHEN** the server is shutting down and `gormDB.DB()` returns nil sqlDB
- **THEN** the code SHALL skip the `.Close()` call instead of panicking

### Requirement: Database and Redis connection pool configuration
The database pool in `repository/gorm.go` SHALL set `SetConnMaxIdleLifetime(30 * time.Minute)`. The Redis client in `repository/redis.go` SHALL configure `MinIdleConns: 5`, `DialTimeout: 5s`, `ReadTimeout: 3s`, `WriteTimeout: 3s` instead of using all-default values.

#### Scenario: Database idle connection timeout
- **WHEN** a database connection has been idle for more than 30 minutes
- **THEN** the connection pool SHALL recycle it via `SetConnMaxIdleLifetime`

#### Scenario: Redis connection under load
- **WHEN** the application starts and no Redis connections exist yet
- **THEN** the pool SHALL pre-create `MinIdleConns: 5` connections

### Requirement: Debug logging cleanup
The `fmt.Printf("登录失败: phone=%s, error=%v\n", ...)` in `internal/handler/auth.go` SHALL be replaced with `logger.Warn("login failed", zap.String("phone", req.Phone), zap.Error(err))`.

#### Scenario: Login failure logging
- **WHEN** a login attempt fails
- **THEN** the system SHALL log the failure using structured zap logging (not fmt.Printf)
- **THEN** the log SHALL use Warn level (not stdout print)
