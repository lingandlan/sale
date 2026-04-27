## Why

前端 `vue-tsc -b` 存在约 100 个类型错误（21 个文件），导致 `npm run build` 失败，Docker 构建前端镜像时报错。这些错误长期被 Docker 层缓存掩盖，修改 Dockerfile 后缓存失效才暴露。需要一次性修复，恢复 `vue-tsc` 类型检查使构建流程完整。

## What Changes

- 修复 API 函数返回类型声明，统一与 axios interceptor unwrap 行为一致（去掉多余的 `{ data: ... }` 包装层）
- 修复所有 `.vue` 文件中因类型不匹配导致的属性访问错误（`res.data.xxx` → `res.data.data.xxx` 或修正类型声明）
- 清理未使用的变量/导入（TS6133）
- 修复 tsconfig `baseUrl` 废弃警告
- 恢复 `package.json` 中 `build` 脚本为 `vue-tsc -b && vite build`

## Capabilities

### New Capabilities

- `api-type-consistency`: API 函数返回类型声明规范 — 统一与 interceptor unwrap 行为对齐，确保 `res.data` 直接访问业务数据

### Modified Capabilities

（无已有 spec 需要修改）

## Impact

- **API 类型声明**: `src/api/*.ts` 中所有函数的返回类型泛型参数
- **Vue 页面**: 21 个文件（views、layouts、stores、tests）
- **构建流程**: `package.json` build 脚本恢复类型检查
- **无运行时行为变更**: 纯类型修复，不影响实际功能
