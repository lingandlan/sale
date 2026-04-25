## 1. User Store 持久化

- [ ] 1.1 在 `shop-pc/src/stores/user.ts` 新增 `USER_INFO_KEY` 常量和 `initFromStorage()` 方法，store 初始化时调用
- [ ] 1.2 在 `fetchUserInfo()` 成功后同步写入 localStorage
- [ ] 1.3 在 `clear()` 方法中同时删除 localStorage `user_info`

## 2. Header 退出登录

- [ ] 2.1 确认 `shop-pc/src/layouts/Header.vue` 退出登录时调用了 store.clear()（应已包含）

## 3. 验证

- [ ] 3.1 编译通过（`cd shop-pc && npx vue-tsc --noEmit`）
- [ ] 3.2 浏览器登录后刷新页面，确认仍保持登录态
- [ ] 3.3 退出登录后刷新，确认跳转到登录页
