# 单元测试说明

## 前端单元测试

### 测试文件位置

```
src/
├── views/__tests__/
│   └── Login.spec.ts          # 登录组件测试
├── api/__tests__/
│   └── auth.spec.ts           # 登录API测试
├── utils/__tests__/
│   └── request.spec.ts        # HTTP请求工具测试
└── test/
    └── setup.ts               # 测试环境配置
```

### 运行测试

```bash
# 运行所有测试
npm run test

# 运行测试并打开UI界面
npm run test:ui

# 运行测试并生成覆盖率报告
npm run test:coverage
```

### 测试覆盖范围

#### Login.spec.ts
- ✓ 组件渲染测试
- ✓ 表单验证测试
- ✓ 登录功能测试
- ✓ UI元素测试

#### auth.spec.ts
- ✓ 登录API调用测试
- ✓ 请求参数验证
- ✓ 错误处理测试

#### request.spec.ts
- ✓ 请求拦截器测试
- ✓ 响应拦截器测试
- ✓ Token处理测试

## 后端单元测试

### 测试文件位置

```
internal/handler/
└── auth_test.go               # 认证处理器测试
```

### 运行测试

```bash
cd backend

# 运行所有测试
go test ./...

# 运行测试并查看覆盖率
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# 运行特定测试
go test ./internal/handler -v
```

### 测试覆盖范围

#### auth_test.go
- ✓ 登录接口测试
- ✓ Token刷新测试
- ✓ 用户注册测试
- ✓ 参数验证测试
- ✓ 错误处理测试

## 测试最佳实践

1. **编写测试**：在开发功能的同时编写测试
2. **运行测试**：每次提交代码前运行测试
3. **覆盖率目标**：核心业务逻辑测试覆盖率≥60%
4. **Mock使用**：合理使用Mock隔离外部依赖
5. **测试命名**：使用清晰的测试名称描述测试场景

## CI/CD集成

测试会在以下情况自动运行：
- Pull Request创建时
- 代码提交到main分支时
- 每日定时构建

测试失败将阻止代码合并。
