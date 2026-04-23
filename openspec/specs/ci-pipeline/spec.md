## ADDED Requirements

### Requirement: PR 和 push 自动触发 CI
系统 SHALL 在 PR 提交到 main/design 分支以及直接 push 到 main 时自动触发 CI 流水线。

#### Scenario: PR 提交触发 CI
- **WHEN** 开发者创建或更新 PR 目标为 main 或 design 分支
- **THEN** GitHub Actions 自动触发 CI 流水线

#### Scenario: Push 到 main 触发 CI
- **WHEN** 代码直接 push 到 main 分支
- **THEN** GitHub Actions 自动触发 CI 流水线

### Requirement: 后端测试 Job
CI 流水线 SHALL 包含 `backend-test` job，运行 Go 后端全部单元测试。

#### Scenario: 后端测试通过
- **WHEN** CI 流水线运行后端测试 job
- **THEN** 执行 `go test ./...` 并报告覆盖率

#### Scenario: 后端测试失败
- **WHEN** 后端测试有 case 失败
- **THEN** CI 标记为失败，PR 不允许合并

### Requirement: 前端构建检查 Job
CI 流水线 SHALL 包含 `frontend-build` job，运行前端 TypeScript 编译和构建。

#### Scenario: 前端构建通过
- **WHEN** CI 流水线运行前端构建 job
- **THEN** 执行 `npm install && npm run build` 成功

#### Scenario: 前端构建失败
- **WHEN** 前端代码有 TypeScript 编译错误
- **THEN** CI 标记为失败，PR 不允许合并

### Requirement: CI 结果反馈
CI 结果 SHALL 显示在 PR 页面上，作为合并的必要条件。

#### Scenario: PR 合并保护
- **WHEN** CI 流水线未全部通过
- **THEN** GitHub 阻止 PR 合并
