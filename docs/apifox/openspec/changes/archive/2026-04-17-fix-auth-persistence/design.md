## Context

当前 token 存储在 localStorage（`access_token` / `refresh_token`），但 `userInfo` 只在 Pinia 内存中。router guard 检查 `access_token` 存在后调 `fetchUserInfo()`，若失败则跳转登录。

问题：agent-browser open 新页面时 localStorage 会被清空，导致 token + userInfo 同时丢失。

## Goals / Non-Goals

**Goals:**
- userInfo 持久化到 localStorage，与 token 一起恢复
- 页面刷新/agent-browser 重新 open 后，若 token 有效则自动恢复登录态，无需重新登录

**Non-Goals:**
- 不改后端
- 不改 token 刷新逻辑
- 不引入 pinia-plugin-persistedstate 等第三方库

## Decisions

### 1. 手动持久化 userInfo

在 user store 的 `setUserInfo` 和 `clear` 方法中同步写/删 localStorage：

```ts
const USER_INFO_KEY = 'user_info'

function setUserInfo(info: UserInfo) {
  userInfo.value = info
  localStorage.setItem(USER_INFO_KEY, JSON.stringify(info))
}

function initFromStorage() {
  const raw = localStorage.getItem(USER_INFO_KEY)
  if (raw) {
    try { userInfo.value = JSON.parse(raw) } catch {}
  }
}
```

在 store 创建时调用 `initFromStorage()`。

**理由**: 项目只用了一个 store，不值得引入持久化插件。手动实现 3 行代码搞定。

### 2. router guard 优化

当前逻辑：有 token 但无 userInfo → 调 fetchUserInfo。持久化后，localStorage 恢复 userInfo 后直接可用，减少一次 API 调用。若 fetchUserInfo 失败（token 过期），仍走 401 刷新流程。

## Risks / Trade-offs

- **[userInfo 与 token 不同步]** → 退出登录时同时清除 localStorage 中的 userInfo；fetchUserInfo 失败时也清除
- **[localStorage 被清空]** → 这是 agent-browser 的限制，持久化后至少在同一 session 内 token + userInfo 同时存在或同时丢失，状态一致
