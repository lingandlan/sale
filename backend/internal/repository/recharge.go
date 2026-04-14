package repository

import (
	"marketplace/backend/internal/model"
	"strings"
	"time"

	"gorm.io/gorm"
)

// RechargeRepoInterface 充值仓库接口
type RechargeRepoInterface interface {
	CreateRechargeApplication(app *model.RechargeApplication) error
	GetRechargeApplications(status string, page, pageSize int) ([]model.RechargeApplication, int64, error)
	GetRechargeApplicationByID(id string) (*model.RechargeApplication, error)
	UpdateRechargeApplicationStatus(id, status, approvedBy, remark string) error
	CreateCRecharge(recharge *model.CRecharge) error
	GetCRechargeList(memberPhone, centerID string, page, pageSize int) ([]model.CRecharge, int64, error)
	GetCRechargeByID(id string) (*model.CRecharge, error)
	CreateCard(card *model.StoreCard) error
	GetCardList(status, cardNo, holderPhone string, page, pageSize int) ([]model.StoreCard, int64, error)
	GetCardByCardNo(cardNo string) (*model.StoreCard, error)
	UpdateCardBalance(cardNo string, balance float64) error
	UpdateCardStatus(cardNo, status string) error
	CreateCardTransaction(transaction *model.CardTransaction) error
	GetCardTransactions(cardNo string) ([]model.CardTransaction, error)
	GetCardStats() (map[string]int64, error)
	GetCenterByID(id string) (*model.RechargeCenter, error)
	DeductCenterBalance(id string, amount float64) (float64, error)
	GetCenters() ([]model.RechargeCenter, error)
	CreateCenter(center *model.RechargeCenter) error
	UpdateCenter(center *model.RechargeCenter) error
	DeleteCenter(id string) error
	GetOperators() ([]model.RechargeOperator, error)
	CreateOperator(operator *model.RechargeOperator) error
	UpdateOperator(operator *model.RechargeOperator) error
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
func (r *RechargeRepository) GetCRechargeList(memberPhone, centerID string, page, pageSize int) ([]model.CRecharge, int64, error) {
	var list []model.CRecharge
	var total int64

	query := r.db.Model(&model.CRecharge{})
	if memberPhone != "" {
		query = query.Where("member_phone = ?", memberPhone)
	}
	if centerID != "" {
		query = query.Where("center_id = ?", centerID)
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

// GetCardList 获取门店卡列表
func (r *RechargeRepository) GetCardList(status, cardNo, holderPhone string, page, pageSize int) ([]model.StoreCard, int64, error) {
	var list []model.StoreCard
	var total int64

	query := r.db.Model(&model.StoreCard{})
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if cardNo != "" {
		query = query.Where("card_no LIKE ?", "%"+cardNo+"%")
	}
	if holderPhone != "" {
		query = query.Where("holder_phone = ?", holderPhone)
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

// UpdateCardBalance 更新卡余额
func (r *RechargeRepository) UpdateCardBalance(cardNo string, balance float64) error {
	return r.db.Model(&model.StoreCard{}).Where("card_no = ?", cardNo).Update("balance", balance).Error
}

// UpdateCardStatus 更新卡状态
func (r *RechargeRepository) UpdateCardStatus(cardNo, status string) error {
	return r.db.Model(&model.StoreCard{}).Where("card_no = ?", cardNo).Update("status", status).Error
}

// CreateCardTransaction 创建卡交易记录
func (r *RechargeRepository) CreateCardTransaction(transaction *model.CardTransaction) error {
	return r.db.Create(transaction).Error
}

// GetCardTransactions 获取卡交易记录
func (r *RechargeRepository) GetCardTransactions(cardNo string) ([]model.CardTransaction, error) {
	var list []model.CardTransaction
	err := r.db.Where("card_no = ?", cardNo).Order("created_at DESC").Limit(50).Find(&list).Error
	return list, err
}

// GetCardStats 获取门店卡统计
func (r *RechargeRepository) GetCardStats() (map[string]int64, error) {
	stats := make(map[string]int64)

	// 总卡数
	var totalCards int64
	r.db.Model(&model.StoreCard{}).Count(&totalCards)
	stats["totalCards"] = totalCards

	// 活跃卡数
	var activeCards int64
	r.db.Model(&model.StoreCard{}).Where("status = ?", "active").Count(&activeCards)
	stats["activeCards"] = activeCards

	// 冻结卡数
	var frozenCards int64
	r.db.Model(&model.StoreCard{}).Where("status = ?", "inactive").Count(&frozenCards)
	stats["frozenCards"] = frozenCards

	// 过期卡数
	var expiredCards int64
	r.db.Model(&model.StoreCard{}).Where("expiry_date < ?", time.Now()).Count(&expiredCards)
	stats["expiredCards"] = expiredCards

	// 总余额
	var totalBalance struct {
		Total float64
	}
	r.db.Model(&model.StoreCard{}).Select("COALESCE(SUM(balance), 0) as total").Scan(&totalBalance)
	stats["totalBalance"] = int64(totalBalance.Total)

	return stats, nil
}

// ========== 充值中心 ==========

// GetCenterByID 根据ID获取充值中心
func (r *RechargeRepository) GetCenterByID(id string) (*model.RechargeCenter, error) {
	var center model.RechargeCenter
	err := r.db.Where("id = ?", id).First(&center).Error
	return &center, err
}

// DeductCenterBalance 扣减充值中心余额，返回扣减后余额
func (r *RechargeRepository) DeductCenterBalance(id string, amount float64) (float64, error) {
	var center model.RechargeCenter
	if err := r.db.Where("id = ?", id).First(&center).Error; err != nil {
		return 0, err
	}
	if center.Balance < amount {
		return 0, gorm.ErrRecordNotFound // 余额不足
	}
	newBalance := center.Balance - amount
	if err := r.db.Model(&model.RechargeCenter{}).Where("id = ?", id).Update("balance", newBalance).Error; err != nil {
		return 0, err
	}
	return newBalance, nil
}

// GetCenters 获取充值中心列表
func (r *RechargeRepository) GetCenters() ([]model.RechargeCenter, error) {
	var list []model.RechargeCenter
	err := r.db.Where("status = ?", "active").Find(&list).Error
	return list, err
}

// CreateCenter 创建充值中心
func (r *RechargeRepository) CreateCenter(center *model.RechargeCenter) error {
	return r.db.Create(center).Error
}

// UpdateCenter 更新充值中心
func (r *RechargeRepository) UpdateCenter(center *model.RechargeCenter) error {
	return r.db.Save(center).Error
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
func (r *RechargeRepository) UpdateOperator(operator *model.RechargeOperator) error {
	return r.db.Save(operator).Error
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
