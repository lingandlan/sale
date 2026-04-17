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
	GetRechargeApplications(status string, page, pageSize int) ([]model.RechargeApplication, int64, error)
	GetRechargeApplicationByID(id string) (*model.RechargeApplication, error)
	UpdateRechargeApplicationStatus(id, status, approvedBy, remark string) error
	CreateCRecharge(recharge *model.CRecharge) error
	GetCRechargeList(memberPhone, centerID, startDate, endDate string, page, pageSize int) ([]model.CRecharge, int64, error)
	GetCRechargeByID(id string) (*model.CRecharge, error)
	// 门店卡
	CreateCard(card *model.StoreCard) error
	BatchCreateCards(cards []*model.StoreCard) error
	GetCardList(status int, cardNo, centerID string, page, pageSize int) ([]model.StoreCard, int64, error)
	GetCardByCardNo(cardNo string) (*model.StoreCard, error)
	GetMaxCardSequence() (int, error)
	UpdateCardByMap(cardNo string, updates map[string]interface{}) error
	GetAllocatableCardCount() (int64, error)
	AllocateCardsByQuantity(centerID string, quantity int) (int, error)
	BindCardToUser(cardNo string, updates map[string]interface{}, record *model.CardIssueRecord) error
	ConsumeCardInTx(cardNo string, amount int, operatorID, remark string) error
	CreateCardTransaction(transaction *model.CardTransaction) error
	GetCardTransactions(cardNo string, page, pageSize int) ([]model.CardTransaction, int64, error)
	GetCardStats() (map[string]int64, error)
	GetCardInventoryStats() (map[string]int64, error)
	// 充值中心
	GetCenterByID(id string) (*model.RechargeCenter, error)
	AddCenterBalance(id string, amount float64) error
	DeductCenterBalance(id string, amount float64) (float64, error)
	GetCenterTotalRecharge(centerID string) int64
	GetCenterTotalConsumed(centerID string) float64
	GetCenters() ([]model.RechargeCenter, error)
	CreateCenter(center *model.RechargeCenter) error
	UpdateCenter(id string, updates map[string]interface{}) error
	DeleteCenter(id string) error
	// 操作员
	GetOperators() ([]model.RechargeOperator, error)
	CreateOperator(operator *model.RechargeOperator) error
	UpdateOperator(id string, updates map[string]interface{}) error
	DeleteOperator(id string) error
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
func (r *RechargeRepository) GetRechargeApplications(status string, page, pageSize int) ([]model.RechargeApplication, int64, error) {
	var list []model.RechargeApplication
	var total int64

	query := r.db.Model(&model.RechargeApplication{})
	if status != "" {
		statuses := strings.Split(status, ",")
		query = query.Where("status IN ?", statuses)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&list).Error

	return list, total, err
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

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&list).Error

	return list, total, err
}

// GetCRechargeByID 根据ID获取C端充值记录
func (r *RechargeRepository) GetCRechargeByID(id string) (*model.CRecharge, error) {
	var recharge model.CRecharge
	err := r.db.Where("id = ?", id).First(&recharge).Error
	return &recharge, err
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

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&list).Error

	return list, total, err
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
func (r *RechargeRepository) BindCardToUser(cardNo string, updates map[string]interface{}, record *model.CardIssueRecord) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 更新卡状态和绑定信息
		if err := tx.Model(&model.StoreCard{}).Where("card_no = ?", cardNo).Updates(updates).Error; err != nil {
			return err
		}
		// 创建发放记录
		if err := tx.Create(record).Error; err != nil {
			return err
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

// GetCardTransactions 获取卡交易记录（分页）
func (r *RechargeRepository) GetCardTransactions(cardNo string, page, pageSize int) ([]model.CardTransaction, int64, error) {
	var list []model.CardTransaction
	var total int64

	query := r.db.Model(&model.CardTransaction{}).Where("card_no = ?", cardNo)
	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&list).Error

	return list, total, err
}

// GetCardStats 获取门店卡统计
func (r *RechargeRepository) GetCardStats() (map[string]int64, error) {
	stats := make(map[string]int64)

	// 总卡数
	var total int64
	r.db.Model(&model.StoreCard{}).Count(&total)
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
		r.db.Model(&model.StoreCard{}).Where("status = ?", status).Count(&count)
		stats[field] = count
	}

	// 总余额（活跃+已冻结+已发放的卡）
	var totalBalance struct{ Total int }
	r.db.Model(&model.StoreCard{}).
		Where("status IN ?", []int{model.CardStatusActive, model.CardStatusFrozen, model.CardStatusIssued}).
		Select("COALESCE(SUM(balance), 0) as total").Scan(&totalBalance)
	stats["totalBalance"] = int64(totalBalance.Total)

	// 今日消费
	today := time.Now().Format("2006-01-02")
	var todayConsume struct{ Total int }
	r.db.Model(&model.CardTransaction{}).
		Where("type = ? AND DATE(created_at) = ?", "consume", today).
		Select("COALESCE(SUM(amount), 0) as total").Scan(&todayConsume)
	stats["todayConsume"] = int64(todayConsume.Total)

	// 7天内过期
	sevenDaysLater := time.Now().AddDate(0, 0, 7)
	var expireIn7Days int64
	r.db.Model(&model.StoreCard{}).
		Where("status = ? AND expired_at IS NOT NULL AND expired_at <= ?", model.CardStatusActive, sevenDaysLater).
		Count(&expireIn7Days)
	stats["expireIn7Days"] = expireIn7Days

	return stats, nil
}

// GetCardInventoryStats 总卡库统计
func (r *RechargeRepository) GetCardInventoryStats() (map[string]int64, error) {
	stats := make(map[string]int64)

	// 总卡数
	var total int64
	r.db.Model(&model.StoreCard{}).Count(&total)
	stats["totalCards"] = total

	// 已发放+已激活+已冻结+已过期的卡都是"已出库"的
	var issued int64
	r.db.Model(&model.StoreCard{}).Where("status IN ?", []int{model.CardStatusIssued, model.CardStatusActive, model.CardStatusFrozen, model.CardStatusExpired}).Count(&issued)
	stats["issuedCards"] = issued

	// 剩余库存 = 已入库且未划拨到充值中心的卡
	var inStock int64
	r.db.Model(&model.StoreCard{}).Where("status = ? AND (recharge_center_id IS NULL OR recharge_center_id = '')", model.CardStatusInStock).Count(&inStock)
	stats["inStockCards"] = inStock

	return stats, nil
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
	r.db.Where("id = ?", id).First(&center)
	return center.Balance, nil
}

// GetCenters 获取充值中心列表
func (r *RechargeRepository) GetCenters() ([]model.RechargeCenter, error) {
	var list []model.RechargeCenter
	err := r.db.Find(&list).Error
	return list, err
}

// GetCenterTotalRecharge 获取中心累计充值（approved 的申请 points 总和）
func (r *RechargeRepository) GetCenterTotalRecharge(centerID string) int64 {
	var total int64
	r.db.Model(&model.RechargeApplication{}).
		Where("center_id = ? AND status = ?", centerID, "approved").
		Select("COALESCE(SUM(points), 0)").
		Scan(&total)
	return total
}

// GetCenterTotalConsumed 获取中心已消耗（c_recharges 的 amount 总和）
func (r *RechargeRepository) GetCenterTotalConsumed(centerID string) float64 {
	var total float64
	r.db.Table("c_recharges").
		Where("center_id = ?", centerID).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&total)
	return total
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
