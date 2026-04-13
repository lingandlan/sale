# 后端开发指南

## 快速启动后端开发

### 1. 数据库初始化

#### 运行迁移
```bash
cd /Users/zhangdaodong/code/sale/backend

# 创建数据表
go run migrations/migrate.go up
```

#### 手动创建表（临时方案）
```sql
-- 充值申请表
CREATE TABLE recharge_applications (
    id VARCHAR(64) PRIMARY KEY,
    center_id VARCHAR(64),
    center_name VARCHAR(100),
    amount DECIMAL(10,2),
    points INT,
    base_points INT,
    rebate_points INT,
    rebate_rate INT,
    applicant_id VARCHAR(64),
    applicant_name VARCHAR(100),
    transaction_no VARCHAR(100),
    screenshot TEXT,
    remark TEXT,
    status VARCHAR(20) DEFAULT 'pending',
    approved_by VARCHAR(64),
    approved_at TIMESTAMP,
    approval_remark TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_status (status),
    INDEX idx_center (center_id),
    INDEX idx_applicant (applicant_id)
);

-- 门店卡表
CREATE TABLE store_cards (
    id VARCHAR(64) PRIMARY KEY,
    card_no VARCHAR(20) UNIQUE NOT NULL,
    holder_id VARCHAR(64),
    holder_name VARCHAR(100),
    holder_phone VARCHAR(20),
    balance DECIMAL(10,2) DEFAULT 0,
    status VARCHAR(20) DEFAULT 'active',
    issue_center_id VARCHAR(64),
    issue_center_name VARCHAR(100),
    issue_date DATE,
    expiry_date DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_card_no (card_no),
    INDEX idx_holder_phone (holder_phone),
    INDEX idx_status (status)
);

-- C端充值记录表
CREATE TABLE c_recharges (
    id VARCHAR(64) PRIMARY KEY,
    member_id VARCHAR(64),
    member_name VARCHAR(100),
    member_phone VARCHAR(20),
    center_id VARCHAR(64),
    center_name VARCHAR(100),
    amount DECIMAL(10,2),
    points INT,
    payment_method VARCHAR(20),
    operator_id VARCHAR(64),
    operator_name VARCHAR(100),
    remark TEXT,
    balance_before INT,
    balance_after INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_member_phone (member_phone),
    INDEX idx_center (center_id)
);
```

### 2. 配置API基础URL

#### 前端配置 (src/utils/request.ts)
```typescript
// 根据环境配置API地址
const BASE_URL = import.meta.env.DEV
  ? 'http://localhost:8080'
  : 'https://api.taijitang.com'
```

#### 后端CORS配置 (.env)
```bash
# 允许前端域名
CORS_ALLOWED_ORIGINS=http://localhost:5173,http://localhost:5174,http://localhost:5175
```

### 3. 测试后端API

#### 启动后端服务器
```bash
cd /Users/zhangdaodong/code/sale/backend
go run cmd/server/main.go
```

#### 测试认证接口
```bash
# 登录
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'

# 获取用户信息（需要token）
curl -X GET http://localhost:8080/api/v1/user/info \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 4. 关键业务实现

#### 积分计算逻辑
```go
// 在 internal/service/recharge.go 中完善
func (s *RechargeService) CalculatePoints(amount, lastMonthConsumption float64) (int, int, int) {
    basePoints := int(amount)
    var rebateRate int
    if lastMonthConsumption >= 100000 {
        rebateRate = 2
    } else {
        rebateRate = 1
    }
    rebatePoints := int(float64(basePoints) * float64(rebateRate) / 100)
    totalPoints := basePoints + rebatePoints
    return basePoints, rebatePoints, totalPoints
}
```

#### 门店卡核销逻辑
```go
func (s *RechargeService) ConsumeCard(cardNo string, amount float64, remark, operatorID string) error {
    // 1. 查询卡信息
    card, err := s.rechargeRepo.GetCardByCardNo(cardNo)
    if err != nil {
        return err
    }

    // 2. 检查余额
    balance := card["balance"].(float64)
    if amount > balance {
        return errors.New("余额不足")
    }

    // 3. 扣除余额
    newBalance := balance - amount
    // TODO: 更新数据库余额字段

    // 4. 记录交易
    transaction := map[string]interface{}{
        "cardNo":       cardNo,
        "type":         "consume",
        "amount":       amount,
        "balanceAfter": newBalance,
        "remark":       remark,
        "operatorId":   operatorID,
    }
    return s.rechargeRepo.CreateCardTransaction(transaction)
}
```

### 5. API接口清单

#### Dashboard API
- `GET /api/v1/dashboard/statistics` - 获取统计数据
- `GET /api/v1/dashboard/todos` - 获取待办事项
- `GET /api/v1/dashboard/recharge-trends` - 获取充值趋势

#### B端充值API
- `POST /api/v1/recharge/b-apply` - 创建充值申请
- `GET /api/v1/recharge/b-approval` - 获取审批列表
- `GET /api/v1/recharge/b-approval/:id` - 获取审批详情
- `POST /api/v1/recharge/b-approval/action` - 审批操作

#### C端充值API
- `POST /api/v1/recharge/c-entry` - C端充值

#### 门店卡API
- `GET /api/v1/card/verify/:cardNo` - 验证卡号
- `POST /api/v1/card/consume` - 核销
- `GET /api/v1/card/list` - 卡列表
- `GET /api/v1/card/detail/:cardNo` - 卡详情
- `GET /api/v1/card/stats` - 卡统计
- `POST /api/v1/card/issue` - 发放卡

#### 充值记录API
- `GET /api/v1/recharge/records` - 充值记录列表
- `GET /api/v1/recharge/records/:id` - 充值记录详情

## 🐛 常见问题

### Q1: 后端启动失败
**检查项**:
1. 端口8080是否被占用
2. MySQL是否运行
3. Redis是否运行
4. configs/config.yaml配置是否正确

### Q2: 前端无法连接后端
**检查项**:
1. 后端是否正常启动
2. CORS配置是否正确
3. API URL配置是否正确
4. 浏览器控制台是否有CORS错误

### Q3: 登录后跳转失败
**解决**:
1. 检查JWT token是否正确返回
2. 检查localStorage是否存储了token
3. 检查路由守卫配置

## 📚 参考文档

- Go后端项目结构: `/Users/zhangdaodong/code/sale/backend/README.md`
- 数据库模型: `internal/model/recharge.go`
- API路由: `internal/router/router.go`
- 业务逻辑: `internal/service/recharge.go`
