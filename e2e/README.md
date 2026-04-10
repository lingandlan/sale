# E2E 测试报告

## 运行方式

```bash
# 安装依赖
npm install

# 运行所有测试
npm run test:e2e

# 运行登录测试
npm run test:e2e:login

# UI 模式 (可视化)
npm run test:e2e:ui
```

## 测试用例

| ID | 用例名称 | 描述 | 优先级 |
|----|----------|------|--------|
| TC-001 | 空表单提交 | 验证必填提示 | 高 |
| TC-002 | 只填手机号 | 验证密码必填提示 | 高 |
| TC-003 | 只填密码 | 验证手机号必填提示 | 高 |
| TC-004 | 错误密码 | 验证错误提示 | 高 |
| TC-005 | 正确登录 | 验证登录成功跳转 | 高 |
| TC-006 | 记住密码 | 验证复选框功能 | 中 |
| TC-007 | 忘记密码 | 验证链接跳转 | 中 |

## CI/CD 集成

添加到 GitHub Actions:

```yaml
- name: E2E Tests
  run: |
    cd e2e
    npm install
    npm run test:e2e
```

## 本地开发

确保前端服务运行中:
```bash
cd shop-h5 && npm run dev
```
