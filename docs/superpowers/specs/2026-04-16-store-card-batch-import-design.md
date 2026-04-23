# 门店卡功能调整设计

> 日期: 2026-04-16
> 状态: 待审核

## 背景

门店卡总卡库管理页面 (CardInventory) 当前有两个功能需要调整：
1. 批量入库 — 目前通过输入起始/结束序号自动生成卡号，改为上传 Excel
2. 划拨充值中心 — 目前通过输入起止卡号范围，改为输入数量

## 变更 1：批量入库改为 Excel 导入

### 现有实现

- 前端：输入 `startSeq/endSeq/cardType`，系统自动计算数量
- 后端：接收 JSON，校验序号不冲突，生成 `TJ%08d` 格式卡号，`CreateInBatches` 入库
- 卡号格式固定，面值硬编码 1000

### 新设计

**Excel 格式：**

| 卡号 | 卡类型 | 面值 |
|------|--------|------|
| TJ00000001 | 实体 | 1000 |
| TJ00000002 | 虚拟 | 500 |

**前端变更：**
- `CardInventory.vue` 批量入库弹窗改为 `el-upload` 上传 `.xlsx` 文件
- 提供模板下载按钮（可选）
- 上传后后端返回 `{ count, cardNos }` 展示结果

**后端变更：**
- 引入 `excelize` 库解析 Excel
- 路由不变 `POST /card/batch-import`，改为 `multipart/form-data` 接收文件
- Handler 层：从 `c.Request.FormFile("file")` 读取文件
- Service 层逻辑：
  1. 解析 Excel 逐行读取卡号/卡类型/面值
  2. 跳过表头行
  3. 校验：卡号格式 `TJ\d{8}`、卡类型有效值(1=实体/2=虚拟)、面值 > 0
  4. 去重：检查卡号是否已存在于数据库
  5. 生成批次号 `BATCH-{timestamp}`
  6. `CreateInBatches` 入库，创建 `stock_in` 交易记录
- 去掉 `startSeq/endSeq` 相关逻辑

**接口变更：**
```
POST /card/batch-import
Content-Type: multipart/form-data

请求：file (xlsx 文件)
响应：{ "code": 0, "data": { "count": 50, "cardNos": ["TJ00000001", ...] } }
```

## 变更 2：划拨充值中心改为输入数量

### 现有实现

- 前端：输入 `centerId + startCardNo + endCardNo`
- 后端：`WHERE card_no >= ? AND card_no <= ? AND status = 1` 批量更新
- 无事务包裹，无交易记录

### 新设计

**前端变更：**
- `CardInventory.vue` 划拨表单改为：选择充值中心（下拉） + 输入数量
- 去掉 startCardNo/endCardNo 字段

**后端变更：**
- 接口改为 `POST /card/allocate`，请求体 `{ centerId, quantity }`
- Service 层逻辑（事务内）：
  1. 验证充值中心存在
  2. 查询可划拨卡数量：`status=1 AND recharge_center_id IS NULL` 的数量
  3. 如果可划拨数量 < 请求数量，返回错误提示
  4. 按卡号升序取前 N 张，更新 `recharge_center_id`
  5. 为每张卡创建 `allocate` 类型的交易记录

**接口变更：**
```
POST /card/allocate
请求：{ "centerId": "xxx", "quantity": 10 }
响应：{ "code": 0, "data": { "count": 10 } }
```

## 影响范围

| 文件 | 变更类型 |
|------|----------|
| `backend/internal/handler/recharge.go` | 修改 BatchImportCards/AllocateCards |
| `backend/internal/service/recharge.go` | 重写 BatchImportCards/AllocateCards |
| `backend/internal/service/interfaces.go` | 更新接口签名 |
| `backend/internal/repository/recharge.go` | 修改 AllocateCardsToCenter，新增查询可划拨卡方法 |
| `shop-pc/src/api/card.ts` | 更新 API 调用方式 |
| `shop-pc/src/views/card/CardInventory.vue` | 重构批量入库/划拨 UI |
| `backend/go.mod` | 新增 excelize 依赖 |

## 不在范围内

- Excel 模板下载功能（后续可加）
- 卡号导出为 Excel 功能
- 批量入库失败行的详细错误报告（仅返回整体成功/失败）
