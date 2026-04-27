## 1. 确认 interceptor 行为

- [ ] 1.1 读取 `shop-pc/src/utils/request.ts`，确认 interceptor 返回值结构（是 `response.data` 还是二次 unwrap）

## 2. 修复 API 类型声明（核心）

- [ ] 2.1 `src/api/admin.ts` — 去掉返回类型中的 `{ data: ... }` 包装
- [ ] 2.2 `src/api/recharge.ts` — 去掉返回类型中的 `{ data: ... }` 包装
- [ ] 2.3 `src/api/card.ts` — 去掉返回类型中的 `{ data: ... }` 包装（如有）
- [ ] 2.4 `src/api/center.ts` — 去掉返回类型中的 `{ data: ... }` 包装（如有）
- [ ] 2.5 `src/api/dashboard.ts` — 去掉返回类型中的 `{ data: ... }` 包装（如有）
- [ ] 2.6 `src/api/operator.ts` — 去掉返回类型中的 `{ data: ... }` 包装（如有）
- [ ] 2.7 `src/stores/user.ts` — 修复类型赋值错误

## 3. 修复未使用变量/导入（TS6133/TS6196）

- [ ] 3.1 `src/layouts/Header.vue` — 删除未使用的 `ref`、`Expand`、`Fold`
- [ ] 3.2 `src/layouts/Sidebar.vue` — 删除未使用的 `watch`、`defaultOpeneds`
- [ ] 3.3 `src/router/index.ts` — 删除未使用的 `from` 参数
- [ ] 3.4 `src/views/Dashboard.vue` — 删除未使用的 `Todos`、`RechargeTrend`
- [ ] 3.5 `src/views/operator/OperatorManage.vue` — 删除未使用的 `ElMessageBox`、`deleteOperator`
- [ ] 3.6 `src/views/recharge/BRechargeDetail.vue` — 删除未使用的 `ElMessageBox`
- [ ] 3.7 `src/views/recharge/BRechargeList.vue` — 删除未使用的 `extractErrorMessage`、`totalPages`
- [ ] 3.8 `src/views/recharge/CRechargeEntry.vue` — 删除未使用的 `router`
- [ ] 3.9 `src/views/card/CardInventory.vue` — 清理未使用导入

## 4. 修复其他类型错误

- [ ] 4.1 `src/main.ts` — 添加 `element-plus/dist/locale/zh-cn.mjs` 类型声明
- [ ] 4.2 `src/views/card/CardIssue.vue` — 修复类型不匹配
- [ ] 4.3 `src/views/card/CardManage.vue` — 修复类型不匹配
- [ ] 4.4 `src/views/center/CenterManage.vue` — 修复类型不匹配（如有）

## 5. 修复测试文件

- [ ] 5.1 `src/api/__tests__/admin.spec.ts` — 修复 status 类型 string vs number
- [ ] 5.2 `src/api/__tests__/card.spec.ts` — 修复不存在导出 `issueCard`/`toggleCardStatus`
- [ ] 5.3 `src/views/__tests__/CardManage.spec.ts` — 修复不存在的方法和导出
- [ ] 5.4 `src/views/__tests__/Login.spec.ts` — 修复类型错误
- [ ] 5.5 `src/views/__tests__/UserManage.spec.ts` — 修复类型错误

## 6. 恢复构建脚本并验证

- [ ] 6.1 `package.json` build 改回 `vue-tsc -b && vite build`，删除 `build:check`
- [ ] 6.2 运行 `npx vue-tsc -b --noEmit` 确认零错误
- [ ] 6.3 运行 `npm run build` 确认构建成功
