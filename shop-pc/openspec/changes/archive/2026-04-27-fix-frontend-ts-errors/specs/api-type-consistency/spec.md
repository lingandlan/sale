## ADDED Requirements

### Requirement: API 函数返回类型不包含额外 data 包装
所有 `src/api/*.ts` 中的 API 函数返回类型 SHALL 直接使用业务数据类型，不包装 `{ data: ... }`。因为 axios interceptor 已 unwrap `response.data`，`res.data` 直接就是业务数据。

#### Scenario: 列表接口类型正确
- **WHEN** 调用 `getAdminUsers(params)` 返回结果 `res`
- **THEN** `res.data` 类型为 `AdminUserListResponse`（含 items/total），可直接 `res.data.items` 访问

#### Scenario: 详情接口类型正确
- **WHEN** 调用 `searchMember(phone)` 返回结果 `res`
- **THEN** `res.data` 类型为 `MemberInfo`（含 userId/phone/balance），可直接 `res.data.userId` 访问

### Requirement: 所有 vue-tsc 类型检查零错误
运行 `vue-tsc -b` SHALL 零错误通过，包括 TS2339/TS6133/TS2322 等所有类型。

#### Scenario: 完整类型检查通过
- **WHEN** 在 shop-pc 目录运行 `npx vue-tsc -b --noEmit`
- **THEN** 退出码为 0，无任何错误输出

### Requirement: build 脚本包含类型检查
`package.json` 的 `build` 脚本 SHALL 为 `vue-tsc -b && vite build`，构建时执行类型检查。

#### Scenario: npm run build 包含类型检查
- **WHEN** 运行 `npm run build`
- **THEN** 先执行 vue-tsc 类型检查，通过后再 vite build

### Requirement: 无未使用变量或导入
所有 `.vue` 和 `.ts` 文件中 SHALL 无未使用的导入（TS6133）和未使用的声明（TS6196）。

#### Scenario: 清理未使用导入
- **WHEN** vue-tsc 检查完成
- **THEN** 无 TS6133 或 TS6196 错误
