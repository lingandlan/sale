## Why

agent-browser 验证前端页面时，每次 `open` 新 URL 会触发完整页面加载，导致 localStorage 中的 JWT token 丢失，每次都要重新登录。同时，Pinia user store 中的 `userInfo` 仅存在内存中，页面刷新后丢失，即使 token 仍有效也需重新获取用户信息。

## What Changes

- **Pinia user store 持久化** — 将 `userInfo` 同步到 localStorage，页面加载时自动恢复
- **Token 注入方式优化** — agent-browser 可通过单次 eval 注入 token + userInfo，避免登录流程

## Capabilities

### New Capabilities

- `auth-state-persistence`: 前端认证状态持久化 — Pinia store 的 userInfo 持久化到 localStorage，页面刷新后自动恢复登录态

### Modified Capabilities

（无）

## Impact

- **代码**: `shop-pc/src/stores/user.ts` — 新增 localStorage 读写 userInfo
- **无后端变更**，无 API 变更
