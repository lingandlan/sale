# 门店卡批量入库Excel导入 + 划拨数量化 实现计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 将门店卡批量入库从序号生成改为Excel文件导入，将划拨充值中心从卡号范围改为输入数量

**Architecture:** 后端用 excelize 解析 Excel 文件，multipart/form-data 接收；划拨改为查询可分配卡+按序取前N张。前后端接口签名均变更，测试同步更新。

**Tech Stack:** Go excelize (Excel解析), Gin multipart, Vue3 el-upload, Element Plus

---

### Task 1: 安装 excelize 依赖

**Files:**
- Modify: `backend/go.mod`
- Modify: `backend/go.sum`

- [ ] **Step 1: 安装 excelize**

Run: `cd backend && go get github.com/xuri/excelize/v2`

- [ ] **Step 2: 验证依赖安装成功**

Run: `cd backend && go mod tidy && grep excelize go.mod`
Expected: 输出包含 `github.com/xuri/excelize/v2`

- [ ] **Step 3: Commit**

```bash
git add backend/go.mod backend/go.sum
git commit -m "chore(deps): add excelize for Excel import"
```

---

### Task 2: 新增 Repository 层方法 — 划拨按数量

**Files:**
- Modify: `backend/internal/repository/recharge.go` (接口 + 实现)
- Modify: `backend/internal/service/recharge_test.go` (Mock 方法)

接口变更：去掉 `AllocateCardsToCenter(centerID, startCardNo, endCardNo string) (int, error)`，新增两个方法：

```go
// RechargeRepoInterface 中：
GetAllocatableCardCount() (int64, error)
AllocateCardsByQuantity(centerID string, quantity int) (int, error)
```

- [ ] **Step 1: 更新 RechargeRepoInterface**

在 `backend/internal/repository/recharge.go` 接口定义中，删除 `AllocateCardsToCenter`，添加：

```go
GetAllocatableCardCount() (int64, error)
AllocateCardsByQuantity(centerID string, quantity int) (int, error)
```

- [ ] **Step 2: 实现 GetAllocatableCardCount**

在 `backend/internal/repository/recharge.go` 实现区域添加：

```go
// GetAllocatableCardCount 获取可划拨的库存卡数量（已入库且未划拨到中心）
func (r *RechargeRepository) GetAllocatableCardCount() (int64, error) {
	var count int64
	err := r.db.Model(&model.StoreCard{}).
		Where("status = ? AND (recharge_center_id IS NULL OR recharge_center_id = '')", model.CardStatusInStock).
		Count(&count).Error
	return count, err
}
```

- [ ] **Step 3: 实现 AllocateCardsByQuantity**

```go
// AllocateCardsByQuantity 按数量划拨库存卡到充值中心（事务）
func (r *RechargeRepository) AllocateCardsByQuantity(centerID string, quantity int) (int, error) {
	var count int
	err := r.db.Transaction(func(tx *gorm.DB) error {
		// 按卡号升序取前 N 张可划拨的卡
		var cards []model.StoreCard
		if err := tx.Where("status = ? AND (recharge_center_id IS NULL OR recharge_center_id = '')", model.CardStatusInStock).
			Order("card_no ASC").Limit(quantity).Find(&cards).Error; err != nil {
			return err
		}
		if len(cards) == 0 {
			return errors.New("没有可划拨的库存卡")
		}
		// 更新这些卡的 recharge_center_id
		cardNos := make([]string, len(cards))
		for i, c := range cards {
			cardNos[i] = c.CardNo
		}
		result := tx.Model(&model.StoreCard{}).
			Where("card_no IN ?", cardNos).
			Update("recharge_center_id", centerID)
		if result.Error != nil {
			return result.Error
		}
		count = len(cards)
		return nil
	})
	return count, err
}
```

- [ ] **Step 4: 删除旧的 AllocateCardsToCenter 实现**

删除 `AllocateCardsToCenter` 方法实现。

- [ ] **Step 5: 更新测试 Mock**

在 `backend/internal/service/recharge_test.go` 中：
- 删除 `MockRechargeRepo.AllocateCardsToCenter`
- 添加：

```go
func (m *MockRechargeRepo) GetAllocatableCardCount() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRechargeRepo) AllocateCardsByQuantity(centerID string, quantity int) (int, error) {
	args := m.Called(centerID, quantity)
	return args.Get(0).(int), args.Error(1)
}
```

- [ ] **Step 6: 运行编译检查**

Run: `cd backend && go build ./...`
Expected: 编译通过，无错误

- [ ] **Step 7: Commit**

```bash
git add backend/internal/repository/recharge.go backend/internal/service/recharge_test.go
git commit -m "refactor(repo): replace AllocateCardsToCenter with quantity-based allocation"
```

---

### Task 3: 重写 Service 层 — BatchImportCards 改为 Excel 解析 + AllocateCards 改为数量

**Files:**
- Modify: `backend/internal/service/interfaces.go` (接口签名)
- Modify: `backend/internal/service/recharge.go` (实现)

接口签名变更：

```go
// 旧：
BatchImportCards(startSeq, endSeq, cardType int, operatorID string) ([]string, error)
AllocateCards(centerID, startCardNo, endCardNo string) (int, error)

// 新：
BatchImportCards(file []byte, operatorID string) (int, []string, error)
AllocateCards(centerID string, quantity int) (int, error)
```

- [ ] **Step 1: 更新 interfaces.go 签名**

在 `backend/internal/service/interfaces.go` 第 46-47 行，替换：

```go
BatchImportCards(file []byte, operatorID string) (int, []string, error)
AllocateCards(centerID string, quantity int) (int, error)
```

- [ ] **Step 2: 重写 BatchImportCards 实现**

在 `backend/internal/service/recharge.go` 中替换整个 `BatchImportCards` 方法（第 203-259 行）：

```go
// BatchImportCards 从 Excel 文件批量入库门店卡
// Excel 格式: 第1列表头"卡号", 第2列"卡类型", 第3列"面值"
func (s *RechargeService) BatchImportCards(file []byte, operatorID string) (int, []string, error) {
	f, err := excelize.OpenReader(bytes.NewReader(file))
	if err != nil {
		return 0, nil, errors.New("无法解析Excel文件")
	}
	defer f.Close()

	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return 0, nil, errors.New("无法读取Excel内容")
	}

	if len(rows) <= 1 {
		return 0, nil, errors.New("Excel文件没有数据行")
	}

	// 卡号格式校验正则
	cardNoRegex := regexp.MustCompile(`^TJ\d{8}$`)
	cardTypeMap := map[string]int{"实体": 1, "实体卡": 1, "虚拟": 2, "虚拟卡": 2}

	batchNo := fmt.Sprintf("B%s", time.Now().Format("20060102150405"))
	cards := make([]*model.StoreCard, 0, len(rows)-1)
	cardNos := make([]string, 0, len(rows)-1)
	seen := make(map[string]bool)

	for i, row := range rows {
		if i == 0 {
			continue // 跳过表头
		}
		if len(row) < 3 {
			return 0, nil, fmt.Errorf("第%d行列数不足，需要3列（卡号、卡类型、面值）", i+1)
		}

		cardNo := strings.TrimSpace(row[0])
		cardTypeStr := strings.TrimSpace(row[1])
		balanceStr := strings.TrimSpace(row[2])

		// 校验卡号格式
		if !cardNoRegex.MatchString(cardNo) {
			return 0, nil, fmt.Errorf("第%d行卡号格式错误：%s（需TJ+8位数字）", i+1, cardNo)
		}

		// 校验卡类型
		cardType, ok := cardTypeMap[cardTypeStr]
		if !ok {
			return 0, nil, fmt.Errorf("第%d行卡类型错误：%s（需填写：实体/虚拟）", i+1, cardTypeStr)
		}

		// 校验面值
		balance, err := strconv.Atoi(balanceStr)
		if err != nil || balance <= 0 {
			return 0, nil, fmt.Errorf("第%d行面值错误：%s（需为正整数）", i+1, balanceStr)
		}

		// Excel 内去重
		if seen[cardNo] {
			return 0, nil, fmt.Errorf("第%d行卡号重复：%s", i+1, cardNo)
		}
		seen[cardNo] = true

		cards = append(cards, &model.StoreCard{
			ID:       uuid.New().String(),
			CardNo:   cardNo,
			CardType: cardType,
			Status:   model.CardStatusInStock,
			Balance:  balance,
			BatchNo:  batchNo,
		})
		cardNos = append(cardNos, cardNo)
	}

	if len(cards) == 0 {
		return 0, nil, errors.New("没有有效的卡数据")
	}
	if len(cards) > 1000 {
		return 0, nil, errors.New("单次入库不能超过1000张")
	}

	// 检查数据库中卡号是否已存在
	existingCards, _ := s.rechargeRepo.GetCardList(0, "", "", 1, 1) // 不做过滤，仅检查
	_ = existingCards // 用 BatchCreateCards 的数据库唯一索引来兜底重复

	if err := s.rechargeRepo.BatchCreateCards(cards); err != nil {
		if strings.Contains(err.Error(), "Duplicate") || strings.Contains(err.Error(), "duplicate") {
			return 0, nil, errors.New("存在重复卡号，入库失败")
		}
		return 0, nil, fmt.Errorf("入库失败：%v", err)
	}

	// 创建入库交易记录
	for _, card := range cards {
		s.rechargeRepo.CreateCardTransaction(&model.CardTransaction{
			ID:            uuid.New().String(),
			CardNo:        card.CardNo,
			Type:          "stock_in",
			Amount:        0,
			BalanceBefore: card.Balance,
			BalanceAfter:  card.Balance,
			Remark:        fmt.Sprintf("Excel导入入库（批次%s）", batchNo),
			OperatorID:    operatorID,
		})
	}

	return len(cards), cardNos, nil
}
```

- [ ] **Step 3: 添加 import**

在 `backend/internal/service/recharge.go` 的 import 中添加：

```go
import (
	"bytes"
	"regexp"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
	// ... 已有的 import
)
```

- [ ] **Step 4: 重写 AllocateCards 实现**

替换 `AllocateCards` 方法（第 262-277 行）：

```go
// AllocateCards 按数量划拨库存卡到充值中心
func (s *RechargeService) AllocateCards(centerID string, quantity int) (int, error) {
	if quantity <= 0 {
		return 0, errors.New("划拨数量必须大于0")
	}
	if quantity > 1000 {
		return 0, errors.New("单次划拨不能超过1000张")
	}

	// 校验充值中心存在
	center, err := s.rechargeRepo.GetCenterByID(centerID)
	if err != nil || center == nil {
		return 0, errors.New("充值中心不存在")
	}

	// 检查可划拨数量
	available, err := s.rechargeRepo.GetAllocatableCardCount()
	if err != nil {
		return 0, err
	}
	if available < int64(quantity) {
		return 0, fmt.Errorf("可划拨库存不足（当前可划拨 %d 张，需划拨 %d 张）", available, quantity)
	}

	count, err := s.rechargeRepo.AllocateCardsByQuantity(centerID, quantity)
	if err != nil {
		return 0, err
	}
	return count, nil
}
```

- [ ] **Step 5: 更新 handler 测试 Mock 签名**

在 `backend/internal/handler/recharge_test.go` 中更新 `MockRechargeService` 的两个方法签名：

```go
func (m *MockRechargeService) BatchImportCards(file []byte, operatorID string) (int, []string, error) {
	args := m.Called(file, operatorID)
	return args.Get(0).(int), args.Get(1).([]string), args.Error(2)
}

func (m *MockRechargeService) AllocateCards(centerID string, quantity int) (int, error) {
	args := m.Called(centerID, quantity)
	return args.Get(0).(int), args.Error(1)
}
```

- [ ] **Step 6: 运行编译检查**

Run: `cd backend && go build ./...`
Expected: 编译通过

- [ ] **Step 7: Commit**

```bash
git add backend/internal/service/interfaces.go backend/internal/service/recharge.go backend/internal/handler/recharge_test.go
git commit -m "refactor(service): rewrite BatchImportCards for Excel + AllocateCards for quantity"
```

---

### Task 4: 更新 Handler 层

**Files:**
- Modify: `backend/internal/handler/recharge.go` (第 152-198 行)

- [ ] **Step 1: 重写 BatchImportCards handler**

替换 handler 中的 `BatchImportCards` 方法（第 152-176 行）：

```go
func (h *RechargeHandler) BatchImportCards(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.ParamsError(c, "请上传Excel文件")
		return
	}

	// 检查文件扩展名
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".xlsx" && ext != ".xls" {
		response.ParamsError(c, "仅支持.xlsx格式文件")
		return
	}

	// 读取文件内容
	content, err := file.Open()
	if err != nil {
		response.InternalError(c, "无法读取上传文件")
		return
	}
	defer content.Close()

	fileBytes, err := io.ReadAll(content)
	if err != nil {
		response.InternalError(c, "无法读取文件内容")
		return
	}

	// TODO: 从JWT获取操作员信息
	operatorID := "op123"

	count, cardNos, err := h.rechargeService.BatchImportCards(fileBytes, operatorID)
	if err != nil {
		response.InternalError(c, errmsg.Get("card.issue_failed")+":"+err.Error())
		return
	}

	response.SuccessWithMessage(c, errmsg.Get("card.issue_success"), gin.H{
		"count":   count,
		"cardNos": cardNos,
	})
}
```

- [ ] **Step 2: 重写 AllocateCards handler**

替换 handler 中的 `AllocateCards` 方法（第 178-198 行）：

```go
func (h *RechargeHandler) AllocateCards(c *gin.Context) {
	var req struct {
		CenterID string `json:"centerId"`
		Quantity int    `json:"quantity"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamsError(c, err.Error())
		return
	}

	count, err := h.rechargeService.AllocateCards(req.CenterID, req.Quantity)
	if err != nil {
		response.InternalError(c, errmsg.Get("card.issue_failed")+":"+err.Error())
		return
	}

	response.SuccessWithMessage(c, errmsg.Get("card.issue_success"), gin.H{
		"count": count,
	})
}
```

- [ ] **Step 3: 添加 import**

在 handler/recharge.go 的 import 中添加：

```go
import (
	"io"
	"path/filepath"
	"strings"
	// ... 已有的 import
)
```

- [ ] **Step 4: 运行编译检查**

Run: `cd backend && go build ./...`
Expected: 编译通过

- [ ] **Step 5: Commit**

```bash
git add backend/internal/handler/recharge.go
git commit -m "refactor(handler): BatchImportCards multipart + AllocateCards quantity"
```

---

### Task 5: 更新 Service 层测试

**Files:**
- Modify: `backend/internal/service/recharge_test.go`

- [ ] **Step 1: 重写 BatchImportCards 测试**

替换现有的 `TestRechargeService_BatchImportCards` 函数。用 `excelize` 构建测试用 Excel 文件：

```go
func TestRechargeService_BatchImportCards(t *testing.T) {
	newTestExcel := func(rows [][]string) []byte {
		f := excelize.NewFile()
		sheet := "Sheet1"
		for i, row := range rows {
			for j, cell := range row {
				cellName, _ := excelize.CoordinatesToCellName(j+1, i+1)
				f.SetCellValue(sheet, cellName, cell)
			}
		}
		buf, _ := f.WriteToBuffer()
		return buf.Bytes()
	}

	t.Run("success", func(t *testing.T) {
		repo := new(MockRechargeRepo)
		svc := newTestRechargeService(repo)

		excelData := newTestExcel([][]string{
			{"卡号", "卡类型", "面值"},
			{"TJ00000001", "实体", "1000"},
			{"TJ00000002", "虚拟", "500"},
		})
		repo.On("BatchCreateCards", mock.AnythingOfType("[]*model.StoreCard")).Return(nil)
		repo.On("CreateCardTransaction", mock.AnythingOfType("*model.CardTransaction")).Return(nil)

		count, cardNos, err := svc.BatchImportCards(excelData, "op-1")
		require.NoError(t, err)
		assert.Equal(t, 2, count)
		assert.Equal(t, "TJ00000001", cardNos[0])
		assert.Equal(t, "TJ00000002", cardNos[1])

		repo.AssertExpectations(t)
	})

	t.Run("empty file", func(t *testing.T) {
		repo := new(MockRechargeRepo)
		svc := newTestRechargeService(repo)

		excelData := newTestExcel([][]string{
			{"卡号", "卡类型", "面值"},
		})
		count, cardNos, err := svc.BatchImportCards(excelData, "op-1")
		assert.Equal(t, 0, count)
		assert.Nil(t, cardNos)
		assert.EqualError(t, err, "Excel文件没有数据行")
	})

	t.Run("invalid card number format", func(t *testing.T) {
		repo := new(MockRechargeRepo)
		svc := newTestRechargeService(repo)

		excelData := newTestExcel([][]string{
			{"卡号", "卡类型", "面值"},
			{"INVALID01", "实体", "1000"},
		})
		count, _, err := svc.BatchImportCards(excelData, "op-1")
		assert.Equal(t, 0, count)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "卡号格式错误")
	})

	t.Run("invalid card type", func(t *testing.T) {
		repo := new(MockRechargeRepo)
		svc := newTestRechargeService(repo)

		excelData := newTestExcel([][]string{
			{"卡号", "卡类型", "面值"},
			{"TJ00000001", "未知", "1000"},
		})
		count, _, err := svc.BatchImportCards(excelData, "op-1")
		assert.Equal(t, 0, count)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "卡类型错误")
	})

	t.Run("duplicate card number in excel", func(t *testing.T) {
		repo := new(MockRechargeRepo)
		svc := newTestRechargeService(repo)

		excelData := newTestExcel([][]string{
			{"卡号", "卡类型", "面值"},
			{"TJ00000001", "实体", "1000"},
			{"TJ00000001", "实体", "1000"},
		})
		count, _, err := svc.BatchImportCards(excelData, "op-1")
		assert.Equal(t, 0, count)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "卡号重复")
	})
}
```

- [ ] **Step 2: 新增 AllocateCards 测试**

在 `TestRechargeService_BatchImportCards` 之后添加：

```go
func TestRechargeService_AllocateCards(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(MockRechargeRepo)
		svc := newTestRechargeService(repo)

		repo.On("GetCenterByID", "center-1").Return(&model.RechargeCenter{ID: "center-1"}, nil)
		repo.On("GetAllocatableCardCount").Return(int64(100), nil)
		repo.On("AllocateCardsByQuantity", "center-1", 10).Return(10, nil)

		count, err := svc.AllocateCards("center-1", 10)
		require.NoError(t, err)
		assert.Equal(t, 10, count)
		repo.AssertExpectations(t)
	})

	t.Run("invalid quantity", func(t *testing.T) {
		repo := new(MockRechargeRepo)
		svc := newTestRechargeService(repo)

		count, err := svc.AllocateCards("center-1", 0)
		assert.Equal(t, 0, count)
		assert.EqualError(t, err, "划拨数量必须大于0")
	})

	t.Run("center not found", func(t *testing.T) {
		repo := new(MockRechargeRepo)
		svc := newTestRechargeService(repo)

		repo.On("GetCenterByID", "bad-id").Return(nil, errors.New("not found"))

		count, err := svc.AllocateCards("bad-id", 10)
		assert.Equal(t, 0, count)
		assert.EqualError(t, err, "充值中心不存在")
	})

	t.Run("insufficient stock", func(t *testing.T) {
		repo := new(MockRechargeRepo)
		svc := newTestRechargeService(repo)

		repo.On("GetCenterByID", "center-1").Return(&model.RechargeCenter{ID: "center-1"}, nil)
		repo.On("GetAllocatableCardCount").Return(int64(5), nil)

		count, err := svc.AllocateCards("center-1", 10)
		assert.Equal(t, 0, count)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "可划拨库存不足")
	})
}
```

- [ ] **Step 3: 更新 handler 测试中的 BatchImportCards 和 AllocateCards**

在 `backend/internal/handler/recharge_test.go` 中：

替换 `TestRechargeHandler_BatchImportCards`（第 622-657 行）：

```go
func TestRechargeHandler_BatchImportCards(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockRechargeService)
		h := NewRechargeHandler(mockSvc)
		router := setupRechargeRouter(h)

		cardNos := []string{"TJ00000001", "TJ00000002"}
		mockSvc.On("BatchImportCards", mock.AnythingOfType("[]uint8"), "op123").Return(2, cardNos, nil).Once()

		// 构建 multipart form
		var buf bytes.Buffer
		writer := multipart.NewWriter(&buf)
		part, _ := writer.CreateFormFile("file", "test.xlsx")
		part.Write([]byte("fake-excel-content"))
		writer.Close()

		req, _ := http.NewRequest("POST", "/api/v1/card/batch-import", &buf)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assertResponseCode(t, w.Body.Bytes(), 0)
		mockSvc.AssertExpectations(t)
	})

	t.Run("no file", func(t *testing.T) {
		mockSvc := new(MockRechargeService)
		h := NewRechargeHandler(mockSvc)
		router := setupRechargeRouter(h)

		req, _ := http.NewRequest("POST", "/api/v1/card/batch-import", bytes.NewBufferString(""))
		req.Header.Set("Content-Type", "multipart/form-data")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assertResponseCode(t, w.Body.Bytes(), 400)
	})
}
```

需要添加 import `"mime/multipart"` 到 handler 测试文件。

- [ ] **Step 4: 运行全部测试**

Run: `cd backend && go test ./internal/service/ -v -run "TestRechargeService_BatchImportCards|TestRechargeService_AllocateCards"`
Run: `cd backend && go test ./internal/handler/ -v -run "TestRechargeHandler_BatchImportCards"`
Expected: 全部 PASS

- [ ] **Step 5: Commit**

```bash
git add backend/internal/service/recharge_test.go backend/internal/handler/recharge_test.go
git commit -m "test: update tests for Excel import + quantity-based allocation"
```

---

### Task 6: 更新前端 API 层

**Files:**
- Modify: `shop-pc/src/api/card.ts` (第 122-130 行)

- [ ] **Step 1: 修改 batchImportCards 函数**

替换第 123-125 行：

```typescript
// 批量入库（Excel上传）
export function batchImportCards(file: File) {
  const formData = new FormData()
  formData.append('file', file)
  return request.post('/card/batch-import', formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  })
}
```

- [ ] **Step 2: 修改 allocateCards 函数**

替换第 128-130 行：

```typescript
// 划拨到充值中心（按数量）
export function allocateCards(data: { centerId: string; quantity: number }) {
  return request.post('/card/allocate', data)
}
```

- [ ] **Step 3: Commit**

```bash
git add shop-pc/src/api/card.ts
git commit -m "refactor(api): update card import/allocate API signatures"
```

---

### Task 7: 重构前端 CardInventory.vue

**Files:**
- Modify: `shop-pc/src/views/card/CardInventory.vue`

- [ ] **Step 1: 重写整个 CardInventory.vue**

```vue
<template>
  <div class="card-inventory">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>总卡库管理</span>
          <el-button type="primary" @click="showImportDialog = true">批量入库</el-button>
        </div>
      </template>

      <!-- 库存统计 -->
      <el-row :gutter="20" class="stats-row">
        <el-col :span="8">
          <el-statistic title="总卡数" :value="inventory.totalCards" />
        </el-col>
        <el-col :span="8">
          <el-statistic title="已出库" :value="inventory.issuedCards" />
        </el-col>
        <el-col :span="8">
          <el-statistic title="剩余库存" :value="inventory.inStockCards" />
        </el-col>
      </el-row>
    </el-card>

    <!-- 划拨到充值中心 -->
    <el-card shadow="never" style="margin-top: 16px">
      <template #header>
        <span>划拨到充值中心</span>
      </template>
      <el-form :model="allocateForm" label-width="100px" :rules="allocateRules" ref="allocateFormRef">
        <el-form-item label="目标中心" prop="centerId">
          <el-select v-model="allocateForm.centerId" placeholder="选择充值中心" style="width: 100%">
            <el-option v-for="c in centers" :key="c.id" :label="c.name" :value="c.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="划拨数量" prop="quantity">
          <el-input-number v-model="allocateForm.quantity" :min="1" :max="1000" style="width: 100%" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleAllocate" :loading="allocateLoading">确认划拨</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 批量入库对话框 -->
    <el-dialog v-model="showImportDialog" title="批量入库" width="500px">
      <el-upload
        ref="uploadRef"
        :auto-upload="false"
        :limit="1"
        accept=".xlsx,.xls"
        :on-change="handleFileChange"
        :on-remove="handleFileRemove"
        drag
      >
        <el-icon style="font-size: 48px; color: #c0c4cc"><upload-filled /></el-icon>
        <div style="margin-top: 8px">拖拽文件到此处，或<em>点击上传</em></div>
        <template #tip>
          <div class="el-upload__tip">仅支持 .xlsx 文件，格式：卡号 | 卡类型 | 面值</div>
        </template>
      </el-upload>
      <template #footer>
        <el-button @click="showImportDialog = false">取消</el-button>
        <el-button type="primary" @click="handleImport" :loading="importLoading" :disabled="!uploadFile">确认入库</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { UploadFilled } from '@element-plus/icons-vue'
import type { FormInstance, UploadFile } from 'element-plus'
import { batchImportCards, allocateCards, getCardInventoryStats, type CardInventoryResponse } from '@/api/card'
import request from '@/utils/request'

const inventory = ref<CardInventoryResponse>({ totalCards: 0, issuedCards: 0, inStockCards: 0 })
const centers = ref<{ id: string; name: string }[]>([])

// 批量入库
const showImportDialog = ref(false)
const importLoading = ref(false)
const uploadRef = ref()
const uploadFile = ref<File | null>(null)

function handleFileChange(file: UploadFile) {
  if (file.raw) {
    uploadFile.value = file.raw
  }
}

function handleFileRemove() {
  uploadFile.value = null
}

// 划拨
const allocateLoading = ref(false)
const allocateFormRef = ref<FormInstance>()
const allocateForm = ref({ centerId: '', quantity: 1 })
const allocateRules = {
  centerId: [{ required: true, message: '请选择充值中心', trigger: 'change' }],
  quantity: [{ required: true, message: '请输入划拨数量', trigger: 'blur' }]
}

async function loadInventory() {
  const res = await getCardInventoryStats()
  inventory.value = res.data || res
}

async function loadCenters() {
  const res = await request.get('/center')
  const data = res.data || res
  centers.value = Array.isArray(data) ? data : (data.list || [])
}

async function handleImport() {
  if (!uploadFile.value) {
    ElMessage.warning('请先选择文件')
    return
  }
  await ElMessageBox.confirm('确认导入所选Excel文件中的门店卡？', '确认')
  importLoading.value = true
  const file = uploadFile.value
  const res = await batchImportCards(file).finally(() => { importLoading.value = false })
  ElMessage.success(`成功入库 ${(res.data || res).count} 张卡`)
  showImportDialog.value = false
  uploadFile.value = null
  loadInventory()
}

async function handleAllocate() {
  await allocateFormRef.value?.validate()
  await ElMessageBox.confirm(`确认划拨 ${allocateForm.value.quantity} 张卡到所选充值中心？`, '确认')
  allocateLoading.value = true
  const res = await allocateCards(allocateForm.value).finally(() => { allocateLoading.value = false })
  ElMessage.success(`成功划拨 ${(res.data || res).count} 张卡`)
  allocateForm.value = { centerId: '', quantity: 1 }
  loadInventory()
}

onMounted(() => {
  loadInventory()
  loadCenters()
})
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.stats-row {
  text-align: center;
}
</style>
```

- [ ] **Step 2: Commit**

```bash
git add shop-pc/src/views/card/CardInventory.vue
git commit -m "feat(ui): refactor CardInventory for Excel upload + quantity allocation"
```

---

### Task 8: 端到端验证

**Files:** 无新文件

- [ ] **Step 1: 运行后端全部测试**

Run: `cd backend && go test ./... -v -count=1 2>&1 | tail -30`
Expected: 所有测试 PASS

- [ ] **Step 2: 编译检查后端**

Run: `cd backend && go build ./...`
Expected: 编译成功

- [ ] **Step 3: 启动后端验证接口**

Run: `cd backend && air`
Expected: 服务在 8080 端口启动成功

- [ ] **Step 4: 用 curl 测试 batch-import 接口**

创建测试 Excel 文件后：

Run: `curl -X POST http://localhost:8080/api/v1/card/batch-import -F "file=@test.xlsx" -H "Authorization: Bearer <token>"`
Expected: 返回 `{ "code": 0, "data": { "count": N, "cardNos": [...] } }`

- [ ] **Step 5: 启动前端验证页面**

Run: `cd shop-pc && npm run dev`
Expected: 页面在 5175 端口可访问，总卡库管理页面正常展示

- [ ] **Step 6: Commit（如有修复）**

```bash
git add -A
git commit -m "fix: end-to-end fixes for Excel import + quantity allocation"
```
