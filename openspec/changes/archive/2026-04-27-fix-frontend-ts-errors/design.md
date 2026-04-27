## Context

前端 `vue-tsc -b` 有约 100 个类型错误（21 个文件），错误分布在 4 类：

1. **TS2339 (68个)** — 属性不存在。主因是 API 函数返回类型声明为 `request.get<{ data: Foo }>(...)`，但 axios interceptor 已 unwrap `response.data`，实际返回的是 `{ data: Foo }`（即整个响应体），导致 `res.data.xxx` 应为 `res.data.data.xxx`，或者类型声明应去掉外层 `{ data: ... }`。
2. **TS6133/TS6196 (21个)** — 未使用的变量/导入，多个 agent 开发时遗留。
3. **TS2322/TS2345 (6个)** — 类型赋值不匹配（string vs number 等）。
4. **其他 (7个)** — TS2724/TS2305 导出成员不存在、TS7016 缺少类型声明、TS2740 缺少属性。

## Goals / Non-Goals

**Goals:**
- 修复全部 TS 错误，`vue-tsc -b` 零错误通过
- 恢复 `build` 脚本为 `vue-tsc -b && vite build`
- API 类型声明与 interceptor 行为统一，后续开发不再踩坑

**Non-Goals:**
- 不重构 API 层架构
- 不改 axios interceptor 逻辑
- 不改后端接口

## Decisions

### D1: 修复策略 — 改调用方而非类型声明

**方案**: 保留 API 函数的类型声明 `request.get<{ data: Foo }>(...)` 不变，修改调用方从 `res.data.xxx` 改为 `res.data.data.xxx`。

**理由**: interceptor 返回 `response.data`（即 `{ code, data, message }`），`res` = `{ code, data, message }`，`res.data` = `{ data: Foo }` 的 `data` 字段。但 CLAUDE.md 说"代码用 `res.data`"——说明 interceptor 做了二次 unwrap。

**实际确认**: interceptor `return response.data`，response.data 是 `{ code:0, data:{...}, message:"" }`。所以 `res = { code, data, message }`，`res.data` 就是业务数据。类型声明 `{ data: Foo }` 是对的，但代码写法 `res.data.items` 就变成 `{ data: Foo }.items` — 不存在。

**结论**: 类型声明是对的，调用方需要改为 `res.data.data.xxx`。但这与 CLAUDE.md 矛盾。需要确认 interceptor 实际行为。

**替代方案（更简单）**: 去掉 API 类型声明中多余的 `{ data: ... }` 包装，改为直接 `request.get<Foo>(...)`。这样 `res.data` 类型就是 `Foo`，代码 `res.data.items` 就能对上。且与 CLAUDE.md "代码用 `res.data`" 一致。

**最终选择**: 改类型声明（去掉 `{ data: ... }` 包装），因为：
- 调用方代码量大（21 个文件），改动多
- 类型声明集中在 `api/*.ts`，改动少
- 与 CLAUDE.md 规范一致

### D2: 未使用变量 — 删除或加下划线前缀

**方案**: 直接删除未使用的导入和变量。如果是解构/回调参数，加 `_` 前缀。

### D3: 测试文件中的 TS 错误

**方案**: 同步修复。测试文件引用的 API 函数如果改名或类型变了，一并更新。

### D4: 恢复 build 脚本

修复完成后将 `package.json` 的 `build` 改回 `vue-tsc -b && vite build`，删除临时 `build:check`。

## Risks / Trade-offs

- **[改动类型声明可能遗漏]** → 全量跑 `vue-tsc` 验证零错误
- **[与 interceptor 实际行为不符]** → 先读 `request.ts` 确认 interceptor 逻辑再动手
