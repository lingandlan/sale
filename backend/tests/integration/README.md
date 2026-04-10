# 集成测试

## 概述

集成测试验证多个组件的真实交互，包括 Service、Repository 与真实数据库的组合。

## 目录结构

```
tests/integration/
├── db_test.go              # 数据库测试工具
├── user_service_test.go    # 用户服务集成测试
├── docker-compose.yml      # 测试数据库配置
├── run.sh                  # 测试运行脚本
└── README.md               # 本文档
```

## 运行方式

### 方式 1: 使用脚本（推荐）

```bash
cd backend/tests/integration

# 启动数据库并运行测试
./run.sh test

# 仅启动数据库
./run.sh start

# 仅停止数据库
./run.sh stop

# 清理测试数据
./run.sh clean
```

### 方式 2: 手动运行

```bash
# 1. 启动测试数据库
docker-compose -f docker-compose.yml up -d

# 2. 运行测试
TEST_DSN="root:root123@tcp(localhost:3307)/test_db?parseTime=true" \
  go test -v -tags=integration ./tests/integration/...

# 3. 停止数据库
docker-compose -f docker-compose.yml down
```

### 方式 3: 跳过集成测试

```bash
# 只运行单元测试（跳过集成测试）
go test ./...

# 运行所有测试
go test -short=false ./...
```

## 环境变量

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `TEST_DSN` | 测试数据库连接字符串 | `root:root123@tcp(localhost:3307)/test_db` |

## 测试用例

| 测试文件 | 用例 | 描述 |
|----------|------|------|
| `user_service_test.go` | TestUserService_Integration | 用户服务 CRUD 完整流程 |
| `user_service_test.go` | TestAuthService_Integration | 认证服务集成测试 |
| `user_service_test.go` | TestRepository_Integration | Repository 层 CRUD 测试 |

## CI/CD 集成

### GitHub Actions

```yaml
- name: Integration Tests
  env:
    TEST_DSN: ${{ secrets.TEST_DSN }}
  run: |
    docker-compose -f tests/integration/docker-compose.yml up -d
    sleep 10
    go test -v -tags=integration ./tests/integration/...
    docker-compose -f tests/integration/docker-compose.yml down
```

## 注意事项

1. **数据隔离**: 每个测试前清理数据，确保独立性
2. **测试顺序**: 使用 `t.Run()` 组织测试顺序
3. **资源清理**: 使用 `defer` 确保资源释放
4. **跳过机制**: 数据库不可用时自动跳过

## 测试金字塔

```
       /E2E\
      /______\          E2E 测试 (agent-browser + Playwright)
     /Integration\
    /______________\    集成测试 (真实数据库)
   /                \
  /    Unit Tests    \  单元测试 (Mock + SQLite)
 /____________________\
```

## 故障排除

### 数据库连接失败

```bash
# 检查数据库状态
docker ps | grep mysql-test

# 查看日志
docker logs sale-mysql-test
```

### 测试数据残留

```bash
# 清理并重建
./run.sh clean
```
