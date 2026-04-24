package repository

import (
	"errors"
	"fmt"
	"marketplace/backend/internal/model"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RechargeRepoInterface 充值仓库接口
//go:generate mockgen -destination=mock_recharge_repo.go -package=repository marketplace/backend/internal/repository RechargeRepoInterface
type RechargeRepoInterface interface {
	CreateRechargeApplication(app *model.RechargeApplication) error
	GetRechargeApplications(status string, centerID string, page, pageSize int) ([]model.RechargeApplication, int64, error)
	GetRechargeApplicationByID(id string) (*model.RechargeApplication, error)
	UpdateRechargeApplicationStatus(id, status, approvedBy, remark string) error
	CreateCRecharge(recharge *model.CRecharge) error
	GetCRechargeList(memberPhone, centerID, startDate, endDate string, page, pageSize int) ([]model.CRecharge, int64, error)
	GetCRechargeByID(id string) (*model.CRecharge, error)
	UpdateCRecharge(id string, updates map[string]interface{}) error
	// 门店卡
	CreateCard(card *model.StoreCard) error
	BatchCreateCards(cards []*model.StoreCard) error
	GetCardList(status int, cardNo, centerID string, page, pageSize int) ([]model.StoreCard, int64, error)
	GetCardByCardNo(cardNo string) (*model.StoreCard, error)
	GetMaxCardSequence() (int, error)
	UpdateCardByMap(cardNo string, updates map[string]interface{}) error
	GetAllocatableCardCount() (int64, error)
	AllocateCardsByQuantity(centerID string, quantity int) (int, error)
	BindCardToUser(cardNo string, updates map[string]interface{}, record *model.CardIssueRecord, txn *model.CardTransaction) error
	ConsumeCardInTx(cardNo string, amount int, operatorID, remark string) error
	CreateCardTransaction(transaction *model.CardTransaction) error
		TransitionCardStatusTX(cardNo string, updates map[string]interface{}, txn *model.CardTransaction) error
		BatchCreateCardTransactions(transactions []*model.CardTransaction) error
	GetCardTransactions(cardNo string, page, pageSize int) ([]model.CardTransaction, int64, error)
	GetCardStats(centerID string) (map[string]int64, error)
	GetCardInventoryStats() (map[string]int64, error)
	GetMonthlyTrend(centerID string) ([]MonthlyTrendItem, error)
	GetCenterCardStats(centerID string) ([]CenterCardStatsItem, error)
		GetAvailableCardNos(centerID string, keyword string) ([]string, error)
		GetAvailableCardCount(centerID string) (int64, error)
	// 充值中心
	GetCenterByID(id string) (*model.RechargeCenter, error)
	AddCenterBalance(id string, amount float64) error
	DeductCenterBalance(id string, amount float64) (float64, error)
	GetCenterTotalRecharge(centerID string) (int64, error)
	GetCenterTotalConsumed(centerID string) (float64, error)
	GetCenters() ([]model.RechargeCenter, error)
	CreateCenter(center *model.RechargeCenter) error
	UpdateCenter(id string, updates map[string]interface{}) error
	DeleteCenter(id string) error
	// 操作员
	GetOperators() ([]model.RechargeOperator, error)
	CreateOperator(operator *model.RechargeOperator) error
	UpdateOperator(id string, updates map[string]interface{}) error
	DeleteOperator(id string) error
	// Dashboard
	GetTodayRechargeTotal(centerID string) (float64, error)
	GetTodayConsumptionTotal(centerID string) (float64, error)
	GetActiveCenterCount(centerID string) (int64, error)
	GetYesterdayRechargeTotal(centerID string) (float64, error)
	GetYesterdayConsumptionTotal(centerID string) (float64, error)
	CountPendingApprovals(centerID string) (int64, error)
	CountExpiringCards(centerID string) (int64, error)
	GetRechargeTrends(days int, centerID string) ([]string, []float64, error)
	// 月度消费
	UpsertMonthlyConsumption(record *model.CenterMonthlyConsumption) error
	GetMonthlyConsumption(centerID, month string) (*model.CenterMonthlyConsumption, error)
	ListMonthlyConsumption(month string) ([]model.CenterMonthlyConsumption, error)
}

var _ RechargeRepoInterface = (*RechargeRepository)(nil)

type RechargeRepository struct {
	db *gorm.DB
}

func NewRechargeRepository(db *gorm.DB) *RechargeRepository {
	return &RechargeRepository{db: db}
}

// ========== B端充值申请 ==========

// CreateRechargeApplication 创建充值申请
func (r *RechargeRepository) CreateRechargeApplication(app *model.RechargeApplication) error {
	return r.db.Create(app).Error
}

// GetRechargeApplications 获取充值申请列表
func (r *RechargeRepository) GetRechargeApplications(status string, centerID string, page, pageSize int) ([]model.RechargeApplication, int64, error) {
	var list []model.RechargeApplication
	var total int64

	query := r.db.Model(&model.RechargeApplication{})
	if status != "" {
		statuses := strings.Split(status, ",")
		query = query.Where("status IN ?", statuses)
	}
	if centerID != "" {
		query = query.Where("center_id = ?", centerID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

// GetRechargeApplicationByID 根据ID获取充值申请
func (r *RechargeRepository) GetRechargeApplicationByID(id string) (*model.RechargeApplication, error) {
	var app model.RechargeApplication
	err := r.db.Where("id = ?", id).First(&app).Error
	return &app, err
}

// UpdateRechargeApplicationStatus 更新充值申请状态
func (r *RechargeRepository) UpdateRechargeApplicationStatus(id, status, approvedBy, remark string) error {
	updates := map[string]interface{}{
		"status": status,
	}
	if approvedBy != "" {
		updates["approved_by"] = approvedBy
	}
	if remark != "" {
		updates["approval_remark"] = remark
	}

	now := time.Now()
	updates["approved_at"] = &now

	return r.db.Model(&model.RechargeApplication{}).Where("id = ?", id).Updates(updates).Error
}

// ========== C端充值 ==========

// CreateCRecharge 创建C端充值记录
func (r *RechargeRepository) CreateCRecharge(recharge *model.CRecharge) error {
	return r.db.Create(recharge).Error
}

// GetCRechargeList 获取C端充值列表
func (r *RechargeRepository) GetCRechargeList(memberPhone, centerID, startDate, endDate string, page, pageSize int) ([]model.CRecharge, int64, error) {
	var list []model.CRecharge
	var total int64

	query := r.db.Model(&model.CRecharge{})
	if memberPhone != "" {
		query = query.Where("member_phone = ?", memberPhone)
	}
	if centerID != "" {
		query = query.Where("center_id = ?", centerID)
	}
	if startDate != "" {
		query = query.Where("created_at >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("created_at <= ?", endDate+" 23:59:59")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

// GetCRechargeByID 根据ID获取C端充值记录
func (r *RechargeRepository) GetCRechargeByID(id string) (*model.CRecharge, error) {
	var recharge model.CRecharge
	err := r.db.Where("id = ?", id).First(&recharge).Error
	return &recharge, err
}

// UpdateCRecharge 更新C端充值记录
func (r *RechargeRepository) UpdateCRecharge(id string, updates map[string]interface{}) error {
	return r.db.Model(&model.CRecharge{}).Where("id = ?", id).Updates(updates).Error
}

// ========== 门店卡 ==========

// CreateCard 创建门店卡
func (r *RechargeRepository) CreateCard(card *model.StoreCard) error {
	return r.db.Create(card).Error
}

// BatchCreateCards 批量创建卡
func (r *RechargeRepository) BatchCreateCards(cards []*model.StoreCard) error {
	return r.db.CreateInBatches(cards, 100).Error
}

// GetCardList 获取门店卡列表
func (r *RechargeRepository) GetCardList(status int, cardNo, centerID string, page, pageSize int) ([]model.StoreCard, int64, error) {
	var list []model.StoreCard
	var total int64

	query := r.db.Model(&model.StoreCard{})
	if status > 0 {
		query = query.Where("status = ?", status)
	}
	if cardNo != "" {
		query = query.Where("card_no LIKE ?", "%"+cardNo+"%")
	}
	if centerID != "" {
		query = query.Where("recharge_center_id = ?", centerID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

// GetCardByCardNo 根据卡号获取门店卡
func (r *RechargeRepository) GetCardByCardNo(cardNo string) (*model.StoreCard, error) {
	var card model.StoreCard
	err := r.db.Where("card_no = ?", cardNo).First(&card).Error
	return &card, err
}

// GetMaxCardSequence 从数据库获取最大卡号序号
func (r *RechargeRepository) GetMaxCardSequence() (int, error) {
	var maxCardNo string
	err := r.db.Model(&model.StoreCard{}).Select("MAX(card_no)").Scan(&maxCardNo).Error
	if err != nil || maxCardNo == "" {
		return 0, nil
	}
	// 从 "TJ00000001" 提取序号部分
	var seq int
	fmt.Sscanf(maxCardNo, "TJ%d", &seq)
	return seq, nil
}

// UpdateCardByMap 通用更新方法
func (r *RechargeRepository) UpdateCardByMap(cardNo string, updates map[string]interface{}) error {
	return r.db.Model(&model.StoreCard{}).Where("card_no = ?", cardNo).Updates(updates).Error
}

// GetAllocatableCardCount 获取可划拨库存卡数量
func (r *RechargeRepository) GetAllocatableCardCount() (int64, error) {
	var count int64
	err := r.db.Model(&model.StoreCard{}).
		Where("status = ? AND (recharge_center_id IS NULL OR recharge_center_id = '')", model.CardStatusInStock).
		Count(&count).Error
	return count, err
}

// AllocateCardsByQuantity 按数量划拨卡到充值中心
func (r *RechargeRepository) AllocateCardsByQuantity(centerID string, quantity int) (int, error) {
	var allocated int
	err := r.db.Transaction(func(tx *gorm.DB) error {
		// 按card_no升序取前N张可划拨的卡
		var cards []model.StoreCard
		if err := tx.Where("status = ? AND (recharge_center_id IS NULL OR recharge_center_id = '')", model.CardStatusInStock).
			Order("card_no ASC").
			Limit(quantity).
			Find(&cards).Error; err != nil {
			return err
		}

		if len(cards) == 0 {
			return nil
		}

		// 收集卡号
		cardNos := make([]string, len(cards))
		for i, c := range cards {
			cardNos[i] = c.CardNo
		}

		// 批量更新
		result := tx.Model(&model.StoreCard{}).
			Where("card_no IN ?", cardNos).
			Update("recharge_center_id", centerID)
		if result.Error != nil {
			return result.Error
		}
		allocated = int(result.RowsAffected)
		return nil
	})
	return allocated, err
}

// BindCardToUser 绑定卡号到用户，同时创建发放记录（在一个事务中）
func (r *RechargeRepository) BindCardToUser(cardNo string, updates map[string]interface{}, record *model.CardIssueRecord, txn *model.CardTransaction) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 更新卡状态和绑定信息
		if err := tx.Model(&model.StoreCard{}).Where("card_no = ?", cardNo).Updates(updates).Error; err != nil {
			return err
		}
		// 创建发放记录
		if err := tx.Create(record).Error; err != nil {
			return err
		}
		// 创建发放交易记录
		if txn != nil {
			if err := tx.Create(txn).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// ConsumeCardInTx 事务核销（行锁 + 首次激活）
func (r *RechargeRepository) ConsumeCardInTx(cardNo string, amount int, operatorID, remark string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. 行锁查卡
		var card model.StoreCard
		if err := tx.Set("gorm:query_option", "FOR UPDATE").
			Where("card_no = ?", cardNo).First(&card).Error; err != nil {
			return err
		}

		// 2. 状态校验
		if card.Status != model.CardStatusIssued && card.Status != model.CardStatusActive {
			return errors.New("卡状态异常，无法核销")
		}
		if card.Status == model.CardStatusActive && card.ExpiredAt != nil && time.Now().After(*card.ExpiredAt) {
			return errors.New("卡已过期，无法核销")
		}

		// 3. 余额校验
		if amount < 100 {
			return errors.New("最低消费100元")
		}
		if amount > card.Balance {
			return errors.New("余额不足")
		}

		// 4. 扣减余额
		newBalance := card.Balance - amount
		updates := map[string]interface{}{
			"balance": newBalance,
		}

		// 5. 首次核销激活
		if card.ActivatedAt == nil {
			now := time.Now()
			expiredAt := now.AddDate(1, 0, 0)
			updates["activated_at"] = &now
			updates["expired_at"] = &expiredAt
			updates["status"] = model.CardStatusActive
		}

		if err := tx.Model(&model.StoreCard{}).Where("card_no = ?", cardNo).Updates(updates).Error; err != nil {
			return err
		}

		// 6. 创建交易记录
		txn := &model.CardTransaction{
			ID:            uuid.New().String(),
			CardNo:        cardNo,
			Type:          "consume",
			Amount:        amount,
			BalanceBefore: card.Balance,
			BalanceAfter:  newBalance,
			Remark:        remark,
			OperatorID:    operatorID,
		}
		return tx.Create(txn).Error
	})
}

// CreateCardTransaction 创建卡交易记录
func (r *RechargeRepository) CreateCardTransaction(transaction *model.CardTransaction) error {
	return r.db.Create(transaction).Error
}

// TransitionCardStatusTX 在事务中更新卡状态并创建交易记录
func (r *RechargeRepository) TransitionCardStatusTX(cardNo string, updates map[string]interface{}, txn *model.CardTransaction) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.StoreCard{}).Where("card_no = ?", cardNo).Updates(updates).Error; err != nil {
			return err
		}
		if txn != nil {
			if err := tx.Create(txn).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// BatchCreateCardTransactions 批量创建卡交易记录（事务）
func (r *RechargeRepository) BatchCreateCardTransactions(transactions []*model.CardTransaction) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, t := range transactions {
			if err := tx.Create(t).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// GetCardTransactions 获取卡交易记录（分页）
func (r *RechargeRepository) GetCardTransactions(cardNo string, page, pageSize int) ([]model.CardTransaction, int64, error) {
	var list []model.CardTransaction
	var total int64

	query := r.db.Model(&model.CardTransaction{}).Where("card_no = ?", cardNo)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

// GetCardStats 获取门店卡统计
func (r *RechargeRepository) GetCardStats(centerID string) (map[string]int64, error) {
	stats := make(map[string]int64)

	// 基础 query，按中心过滤
	baseQuery := r.db.Model(&model.StoreCard{})
	if centerID != "" {
		baseQuery = baseQuery.Where("recharge_center_id = ?", centerID)
	}

	// 总卡数
	var total int64
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, err
	}
	stats["totalCards"] = total

	// 按6种状态统计
	statusFields := map[string]int{
		"inStockCards": model.CardStatusInStock,
		"issuedCards":  model.CardStatusIssued,
		"activeCards":  model.CardStatusActive,
		"frozenCards":  model.CardStatusFrozen,
		"expiredCards": model.CardStatusExpired,
		"voidedCards":  model.CardStatusVoided,
	}
	for field, status := range statusFields {
		var count int64
		q := r.db.Model(&model.StoreCard{}).Where("status = ?", status)
		if centerID != "" {
			q = q.Where("recharge_center_id = ?", centerID)
		}
		if err := q.Count(&count).Error; err != nil {
			return nil, err
		}
		stats[field] = count
	}

	// 总余额（活跃+已冻结+已发放的卡）
	var totalBalance struct{ Total int }
	q := r.db.Model(&model.StoreCard{}).
		Where("status IN ?", []int{model.CardStatusActive, model.CardStatusFrozen, model.CardStatusIssued})
	if centerID != "" {
		q = q.Where("recharge_center_id = ?", centerID)
	}
	if err := q.Select("COALESCE(SUM(balance), 0) as total").Scan(&totalBalance).Error; err != nil {
		return nil, err
	}
	stats["totalBalance"] = int64(totalBalance.Total)

	// 今日消费（按中心过滤需 JOIN store_cards）
	today := time.Now().Format("2006-01-02")
	var todayConsume struct{ Total int }
	tq := r.db.Model(&model.CardTransaction{}).
		Where("type = ? AND DATE(card_transactions.created_at) = ?", "consume", today)
	if centerID != "" {
		tq = tq.Joins("JOIN store_cards ON store_cards.card_no = card_transactions.card_no").
			Where("store_cards.recharge_center_id = ?", centerID)
	}
	if err := tq.Select("COALESCE(SUM(amount), 0) as total").Scan(&todayConsume).Error; err != nil {
		return nil, err
	}
	stats["todayConsume"] = int64(todayConsume.Total)

	// 7天内过期
	sevenDaysLater := time.Now().AddDate(0, 0, 7)
	var expireIn7Days int64
	eq := r.db.Model(&model.StoreCard{}).
		Where("status = ? AND expired_at IS NOT NULL AND expired_at <= ?", model.CardStatusActive, sevenDaysLater)
	if centerID != "" {
		eq = eq.Where("recharge_center_id = ?", centerID)
	}
	if err := eq.Count(&expireIn7Days).Error; err != nil {
		return nil, err
	}
	stats["expireIn7Days"] = expireIn7Days

	return stats, nil
}

// GetCardInventoryStats 总卡库统计（返回所有状态计数）
func (r *RechargeRepository) GetCardInventoryStats() (map[string]int64, error) {
	stats := make(map[string]int64)

	// 总卡数
	var total int64
	if err := r.db.Model(&model.StoreCard{}).Count(&total).Error; err != nil {
		return nil, err
	}
	stats["totalCards"] = total

	// 按6种状态统计
	statusFields := map[string]int{
		"inStockCards": model.CardStatusInStock,
		"issuedCards":  model.CardStatusIssued,
		"activeCards":  model.CardStatusActive,
		"frozenCards":  model.CardStatusFrozen,
		"expiredCards": model.CardStatusExpired,
		"voidedCards":  model.CardStatusVoided,
	}
	for field, status := range statusFields {
		var count int64
		if err := r.db.Model(&model.StoreCard{}).Where("status = ?", status).Count(&count).Error; err != nil {
			return nil, err
		}
		stats[field] = count
	}

	return stats, nil
}

// MonthlyTrendItem 月度趋势数据
type MonthlyTrendItem struct {
	Month   string `json:"month"`
	Issue   int64  `json:"issue"`
	Consume int64  `json:"consume"`
}

// GetMonthlyTrend 获取月度发放/核销趋势
func (r *RechargeRepository) GetMonthlyTrend(centerID string) ([]MonthlyTrendItem, error) {
	// 生成最近6个月的月份列表
	months := make([]string, 6)
	now := time.Now()
	for i := 5; i >= 0; i-- {
		t := now.AddDate(0, -i, 0)
		months[5-i] = t.Format("2006-01")
	}

	// 查询发放数据
	type monthCount struct {
		Month string
		Cnt   int64
	}
	var issueCounts []monthCount
	iq := r.db.Model(&model.CardTransaction{}).
		Where("type = ?", "issue").
		Where("DATE_FORMAT(card_transactions.created_at, '%Y-%m') IN ?", months)
	if centerID != "" {
		iq = iq.Joins("JOIN store_cards ON store_cards.card_no = card_transactions.card_no").
			Where("store_cards.recharge_center_id = ?", centerID)
	}
	if err := iq.Select("DATE_FORMAT(card_transactions.created_at, '%Y-%m') as month, COUNT(*) as cnt").
		Group("DATE_FORMAT(card_transactions.created_at, '%Y-%m')").Scan(&issueCounts).Error; err != nil {
		return nil, err
	}

	var consumeCounts []monthCount
	cq := r.db.Model(&model.CardTransaction{}).
		Where("type = ?", "consume").
		Where("DATE_FORMAT(card_transactions.created_at, '%Y-%m') IN ?", months)
	if centerID != "" {
		cq = cq.Joins("JOIN store_cards ON store_cards.card_no = card_transactions.card_no").
			Where("store_cards.recharge_center_id = ?", centerID)
	}
	if err := cq.Select("DATE_FORMAT(card_transactions.created_at, '%Y-%m') as month, COUNT(*) as cnt").
		Group("DATE_FORMAT(card_transactions.created_at, '%Y-%m')").Scan(&consumeCounts).Error; err != nil {
		return nil, err
	}

	issueMap := make(map[string]int64)
	for _, ic := range issueCounts {
		issueMap[ic.Month] = ic.Cnt
	}
	consumeMap := make(map[string]int64)
	for _, cc := range consumeCounts {
		consumeMap[cc.Month] = cc.Cnt
	}

	result := make([]MonthlyTrendItem, 0, 6)
	for _, m := range months {
		result = append(result, MonthlyTrendItem{
			Month:   m,
			Issue:   issueMap[m],
			Consume: consumeMap[m],
		})
	}
	return result, nil
}

// CenterCardStatsItem 充值中心卡统计
type CenterCardStatsItem struct {
	CenterName   string `json:"centerName"`
	TotalCards   int64  `json:"totalCards"`
	IssuedCards  int64  `json:"issuedCards"`
	ActiveCards  int64  `json:"activeCards"`
	FrozenCards  int64  `json:"frozenCards"`
	ExpiredCards int64  `json:"expiredCards"`
	TotalBalance int64  `json:"totalBalance"`
}

// GetCenterCardStats 按充值中心分组统计
func (r *RechargeRepository) GetCenterCardStats(centerID string) ([]CenterCardStatsItem, error) {
	results := make([]CenterCardStatsItem, 0)
	query := r.db.Model(&model.StoreCard{}).
		Select("rc.name as center_name, COUNT(*) as total_cards, "+
			"SUM(CASE WHEN sc.status = 2 THEN 1 ELSE 0 END) as issued_cards, "+
			"SUM(CASE WHEN sc.status = 3 THEN 1 ELSE 0 END) as active_cards, "+
			"SUM(CASE WHEN sc.status = 4 THEN 1 ELSE 0 END) as frozen_cards, "+
			"SUM(CASE WHEN sc.status = 5 THEN 1 ELSE 0 END) as expired_cards, "+
			"COALESCE(SUM(CASE WHEN sc.status IN (2,3,4) THEN sc.balance ELSE 0 END), 0) as total_balance").
		Joins("JOIN recharge_centers rc ON rc.id COLLATE utf8mb4_unicode_ci = sc.recharge_center_id").
		Table("store_cards sc").
		Where("sc.recharge_center_id IS NOT NULL AND sc.recharge_center_id != ''")

	if centerID != "" {
		query = query.Where("sc.recharge_center_id = ?", centerID)
	}
	if err := query.Group("rc.name").Scan(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

// GetAvailableCardNos 获取指定充值中心的可用卡号列表
func (r *RechargeRepository) GetAvailableCardNos(centerID string, keyword string) ([]string, error) {
	var cardNos []string
	query := r.db.Model(&model.StoreCard{}).
		Where("status = ? AND recharge_center_id = ?", model.CardStatusInStock, centerID)
	if keyword != "" {
		query = query.Where("card_no LIKE ?", keyword+"%")
	}
	err := query.Order("card_no ASC").Limit(20).Pluck("card_no", &cardNos).Error
	return cardNos, err
}

// GetAvailableCardCount 获取指定充值中心的可用卡数量
func (r *RechargeRepository) GetAvailableCardCount(centerID string) (int64, error) {
	var count int64
	err := r.db.Model(&model.StoreCard{}).
		Where("status = ? AND recharge_center_id = ?", model.CardStatusInStock, centerID).
		Count(&count).Error
	return count, err
}

// ========== 充值中心 ==========

// GetCenterByID 根据ID获取充值中心
func (r *RechargeRepository) GetCenterByID(id string) (*model.RechargeCenter, error) {
	var center model.RechargeCenter
	err := r.db.Where("id = ?", id).First(&center).Error
	return &center, err
}

// AddCenterBalance 增加充值中心余额（原子操作）
func (r *RechargeRepository) AddCenterBalance(id string, amount float64) error {
	result := r.db.Model(&model.RechargeCenter{}).
		Where("id = ?", id).
		Update("balance", gorm.Expr("balance + ?", amount))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// DeductCenterBalance 扣减充值中心余额（原子操作），返回扣减后余额
func (r *RechargeRepository) DeductCenterBalance(id string, amount float64) (float64, error) {
	result := r.db.Model(&model.RechargeCenter{}).
		Where("id = ? AND balance >= ?", id, amount).
		Update("balance", gorm.Expr("balance - ?", amount))
	if result.Error != nil {
		return 0, result.Error
	}
	if result.RowsAffected == 0 {
		// 区分：center 不存在 还是 余额不足
		var exists int64
		r.db.Model(&model.RechargeCenter{}).Where("id = ?", id).Count(&exists)
		if exists == 0 {
			return 0, gorm.ErrRecordNotFound
		}
		return 0, errors.New("充值中心余额不足")
	}
	// 查询扣减后余额
	var center model.RechargeCenter
	if err := r.db.Where("id = ?", id).First(&center).Error; err != nil {
		return 0, fmt.Errorf("查询扣减后余额失败: %w", err)
	}
	return center.Balance, nil
}

// GetCenters 获取充值中心列表
func (r *RechargeRepository) GetCenters() ([]model.RechargeCenter, error) {
	var list []model.RechargeCenter
	err := r.db.Find(&list).Error
	return list, err
}

// GetCenterTotalRecharge 获取中心累计充值（approved 的申请 points 总和）
func (r *RechargeRepository) GetCenterTotalRecharge(centerID string) (int64, error) {
	var total int64
	if err := r.db.Model(&model.RechargeApplication{}).
		Where("center_id = ? AND status = ?", centerID, "approved").
		Select("COALESCE(SUM(points), 0)").
		Scan(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

// GetCenterTotalConsumed 获取中心已消耗（c_recharges 的 amount 总和）
func (r *RechargeRepository) GetCenterTotalConsumed(centerID string) (float64, error) {
	var total float64
	if err := r.db.Table("c_recharges").
		Where("center_id = ?", centerID).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

// CreateCenter 创建充值中心
func (r *RechargeRepository) CreateCenter(center *model.RechargeCenter) error {
	return r.db.Create(center).Error
}

// UpdateCenter 更新充值中心
func (r *RechargeRepository) UpdateCenter(id string, updates map[string]interface{}) error {
	return r.db.Model(&model.RechargeCenter{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteCenter 删除充值中心
func (r *RechargeRepository) DeleteCenter(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.RechargeCenter{}).Error
}

// ========== 操作员 ==========

// GetOperators 获取操作员列表
func (r *RechargeRepository) GetOperators() ([]model.RechargeOperator, error) {
	var list []model.RechargeOperator
	err := r.db.Find(&list).Error
	return list, err
}

// CreateOperator 创建操作员
func (r *RechargeRepository) CreateOperator(operator *model.RechargeOperator) error {
	return r.db.Create(operator).Error
}

// UpdateOperator 更新操作员
func (r *RechargeRepository) UpdateOperator(id string, updates map[string]interface{}) error {
	return r.db.Model(&model.RechargeOperator{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteOperator 删除操作员
func (r *RechargeRepository) DeleteOperator(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.RechargeOperator{}).Error
}

// GetOperatorByUsername 根据用户名获取操作员
func (r *RechargeRepository) GetOperatorByUsername(username string) (*model.RechargeOperator, error) {
	var operator model.RechargeOperator
	err := r.db.Where("username = ?", username).First(&operator).Error
	return &operator, err
}

// ========== Dashboard ==========

func (r *RechargeRepository) GetTodayRechargeTotal(centerID string) (float64, error) {
	var total float64
	query := r.db.Model(&model.CRecharge{}).
		Where("DATE(created_at) = CURDATE()")
	if centerID != "" {
		query = query.Where("center_id = ?", centerID)
	}
	err := query.Select("COALESCE(SUM(amount), 0)").Scan(&total).Error
	return total, err
}

func (r *RechargeRepository) GetTodayConsumptionTotal(centerID string) (float64, error) {
	var total float64
	query := r.db.Model(&model.CardTransaction{}).
		Where("type = 'consume' AND DATE(card_transactions.created_at) = CURDATE()")
	if centerID != "" {
		query = query.Where("card_no IN (SELECT card_no FROM store_cards WHERE recharge_center_id = ?)", centerID)
	}
	err := query.Select("COALESCE(SUM(amount), 0)").Scan(&total).Error
	return total, err
}

func (r *RechargeRepository) GetActiveCenterCount(centerID string) (int64, error) {
	if centerID != "" {
		return 1, nil
	}
	var count int64
	err := r.db.Model(&model.CRecharge{}).
		Where("DATE(created_at) >= DATE_SUB(CURDATE(), INTERVAL 30 DAY)").
		Distinct("center_id").
		Count(&count).Error
	return count, err
}

func (r *RechargeRepository) GetYesterdayRechargeTotal(centerID string) (float64, error) {
	var total float64
	query := r.db.Model(&model.CRecharge{}).
		Where("DATE(created_at) = DATE_SUB(CURDATE(), INTERVAL 1 DAY)")
	if centerID != "" {
		query = query.Where("center_id = ?", centerID)
	}
	err := query.Select("COALESCE(SUM(amount), 0)").Scan(&total).Error
	return total, err
}

func (r *RechargeRepository) GetYesterdayConsumptionTotal(centerID string) (float64, error) {
	var total float64
	query := r.db.Model(&model.CardTransaction{}).
		Where("type = 'consume' AND DATE(card_transactions.created_at) = DATE_SUB(CURDATE(), INTERVAL 1 DAY)")
	if centerID != "" {
		query = query.Where("card_no IN (SELECT card_no FROM store_cards WHERE recharge_center_id = ?)", centerID)
	}
	err := query.Select("COALESCE(SUM(amount), 0)").Scan(&total).Error
	return total, err
}

func (r *RechargeRepository) CountPendingApprovals(centerID string) (int64, error) {
	var count int64
	query := r.db.Model(&model.RechargeApplication{}).Where("status = 'pending'")
	if centerID != "" {
		query = query.Where("center_id = ?", centerID)
	}
	err := query.Count(&count).Error
	return count, err
}

func (r *RechargeRepository) CountExpiringCards(centerID string) (int64, error) {
	var count int64
	query := r.db.Model(&model.StoreCard{}).
		Where("status NOT IN (5, 6)"). // 排除已过期和已作废
		Where("expired_at BETWEEN NOW() AND DATE_ADD(NOW(), INTERVAL 7 DAY)")
	if centerID != "" {
		query = query.Where("recharge_center_id = ?", centerID)
	}
	err := query.Count(&count).Error
	return count, err
}

func (r *RechargeRepository) GetRechargeTrends(days int, centerID string) ([]string, []float64, error) {
	type result struct {
		Date  string
		Total float64
	}
	var results []result

	query := r.db.Model(&model.CRecharge{}).
		Select("DATE(created_at) AS date, COALESCE(SUM(amount), 0) AS total").
		Where("DATE(created_at) >= DATE_SUB(CURDATE(), INTERVAL ? DAY)", days-1).
		Group("DATE(created_at)").
		Order("DATE(created_at) ASC")
	if centerID != "" {
		query = query.Where("center_id = ?", centerID)
	}
	err := query.Scan(&results).Error
	if err != nil {
		return nil, nil, err
	}

	// 确保连续：填充缺失日期
	dates := make([]string, days)
	values := make([]float64, days)
	now := time.Now()
	for i := 0; i < days; i++ {
		d := now.AddDate(0, 0, -(days - 1 - i))
		dates[i] = d.Format("01-02")
		values[i] = 0
	}

	resultMap := make(map[string]float64, len(results))
	for _, r := range results {
		// DB returns "YYYY-MM-DD", convert to "MM-DD"
		parts := strings.SplitN(r.Date, "-", 3)
		if len(parts) == 3 {
			resultMap[parts[1]+"-"+parts[2]] = r.Total
		}
	}
	for i, d := range dates {
		if v, ok := resultMap[d]; ok {
			values[i] = v
		}
	}

	return dates, values, nil
}

// ========== 月度消费 ==========

// UpsertMonthlyConsumption 录入/更新月度消费（ON DUPLICATE KEY UPDATE）
func (r *RechargeRepository) UpsertMonthlyConsumption(record *model.CenterMonthlyConsumption) error {
	return r.db.Save(record).Error
}

// GetMonthlyConsumption 查询指定中心指定月份的消费记录
func (r *RechargeRepository) GetMonthlyConsumption(centerID, month string) (*model.CenterMonthlyConsumption, error) {
	var record model.CenterMonthlyConsumption
	if err := r.db.Where("center_id = ? AND month = ?", centerID, month).First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &record, nil
}

// ListMonthlyConsumption 查询指定月份所有中心的消费记录
func (r *RechargeRepository) ListMonthlyConsumption(month string) ([]model.CenterMonthlyConsumption, error) {
	var list []model.CenterMonthlyConsumption
	query := r.db.Model(&model.CenterMonthlyConsumption{})
	if month != "" {
		query = query.Where("month = ?", month)
	}
	if err := query.Order("center_id").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
