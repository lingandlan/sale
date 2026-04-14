package service

import (
	"errors"
	"fmt"
	"time"

	"marketplace/backend/internal/model"
	"marketplace/backend/internal/repository"

	"github.com/google/uuid"
)

type RechargeService struct {
	rechargeRepo repository.RechargeRepoInterface
}

func NewRechargeService(rechargeRepo repository.RechargeRepoInterface) *RechargeService {
	return &RechargeService{rechargeRepo: rechargeRepo}
}

// ========== B端充值申请 ==========

// CalculatePoints 计算积分
// 规则：基础积分 = 金额，返还积分 = 金额 * 返还比例
// 返还比例：上月净消费≥10万元返2%，否则返1%
func (s *RechargeService) CalculatePoints(amount float64, lastMonthConsumption float64) (int, int, int) {
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

// CreateBRechargeApplication 创建B端充值申请
func (s *RechargeService) CreateBRechargeApplication(data map[string]interface{}) (*model.RechargeApplication, error) {
	// 计算积分
	amount := data["amount"].(float64)
	lastMonthConsumption := data["lastMonthConsumption"].(float64)
	basePoints, rebatePoints, totalPoints := s.CalculatePoints(amount, lastMonthConsumption)

	// 创建申请记录
	app := &model.RechargeApplication{
		ID:            uuid.New().String(),
		CenterID:      data["centerId"].(string),
		CenterName:    data["centerName"].(string),
		Amount:        amount,
		BasePoints:    basePoints,
		RebatePoints:  rebatePoints,
		Points:        totalPoints,
		RebateRate:    rebatePoints * 100 / basePoints,
		ApplicantID:   data["applicantId"].(string),
		ApplicantName: data["applicantName"].(string),
		TransactionNo: data["transactionNo"].(string),
		Screenshot:     data["screenshot"].(string),
		Remark:        data["remark"].(string),
		Status:        "pending",
	}

	if err := s.rechargeRepo.CreateRechargeApplication(app); err != nil {
		return nil, err
	}

	return app, nil
}

// GetRechargeApplicationList 获取充值申请列表
func (s *RechargeService) GetRechargeApplicationList(status string, page, pageSize int) (map[string]interface{}, error) {
	list, total, err := s.rechargeRepo.GetRechargeApplications(status, page, pageSize)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"list":  list,
		"total": total,
	}, nil
}

// GetRechargeApplicationDetail 获取充值申请详情
func (s *RechargeService) GetRechargeApplicationDetail(id string) (*model.RechargeApplication, error) {
	app, err := s.rechargeRepo.GetRechargeApplicationByID(id)
	if err != nil {
		return nil, err
	}
	return app, nil
}

// ApproveRechargeApplication 审批充值申请
func (s *RechargeService) ApproveRechargeApplication(id, action, approvedBy, remark string) error {
	_, err := s.rechargeRepo.GetRechargeApplicationByID(id)
	if err != nil {
		return errors.New("充值申请不存在")
	}

	var status string
	if action == "approve" {
		status = "approved"
		// TODO: 实际增加积分到充值中心
	} else if action == "reject" {
		status = "rejected"
	} else {
		return errors.New("invalid action")
	}

	return s.rechargeRepo.UpdateRechargeApplicationStatus(id, status, approvedBy, remark)
}

// ========== C端充值 ==========

// CreateCRecharge C端充值
func (s *RechargeService) CreateCRecharge(data map[string]interface{}) (*model.CRecharge, error) {
	memberPhone := data["memberPhone"].(string)
	amount := data["amount"].(float64)
	centerID := data["centerId"].(string)
	points := int(amount)

	// 获取充值中心信息及余额
	center, err := s.rechargeRepo.GetCenterByID(centerID)
	if err != nil {
		return nil, errors.New("充值中心不存在")
	}
	if center.Balance < amount {
		return nil, errors.New("充值中心余额不足")
	}

	// TODO: 获取会员当前余额
	memberBalanceBefore := 0
	memberBalanceAfter := memberBalanceBefore + points

	recharge := &model.CRecharge{
		ID:            uuid.New().String(),
		MemberID:      data["memberId"].(string),
		MemberName:    data["memberName"].(string),
		MemberPhone:   memberPhone,
		CenterID:      centerID,
		CenterName:    data["centerName"].(string),
		Amount:        amount,
		Points:        points,
		PaymentMethod: data["paymentMethod"].(string),
		OperatorID:    data["operatorId"].(string),
		OperatorName:  data["operatorName"].(string),
		Remark:        data["remark"].(string),
		BalanceBefore: memberBalanceBefore,
		BalanceAfter:  memberBalanceAfter,
	}

	if err := s.rechargeRepo.CreateCRecharge(recharge); err != nil {
		return nil, err
	}

	// 扣减充值中心余额
	if _, err := s.rechargeRepo.DeductCenterBalance(centerID, amount); err != nil {
		return nil, errors.New("扣减中心余额失败")
	}

	return recharge, nil
}

// GetCRechargeList 获取C端充值列表
func (s *RechargeService) GetCRechargeList(memberPhone, centerID string, page, pageSize int) (map[string]interface{}, error) {
	list, total, err := s.rechargeRepo.GetCRechargeList(memberPhone, centerID, page, pageSize)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"list":  list,
		"total": total,
	}, nil
}

// GetCRechargeDetail 获取C端充值详情
func (s *RechargeService) GetCRechargeDetail(id string) (*model.CRecharge, error) {
	recharge, err := s.rechargeRepo.GetCRechargeByID(id)
	if err != nil {
		return nil, err
	}
	return recharge, nil
}

// ========== 门店卡 ==========

// IssueCard 发放门店卡
func (s *RechargeService) IssueCard(data map[string]interface{}) (*model.StoreCard, error) {
	// 生成卡号
	cardNo := s.generateCardNo()

	// 计算过期日期（1年后）
	issueDate := time.Now()
	expiryDate := issueDate.AddDate(1, 0, 0)

	card := &model.StoreCard{
		ID:              uuid.New().String(),
		CardNo:          cardNo,
		HolderID:        data["holderId"].(string),
		HolderName:      data["holderName"].(string),
		HolderPhone:     data["holderPhone"].(string),
		Balance:         data["amount"].(float64),
		Status:          "active",
		IssueCenterID:   data["centerId"].(string),
		IssueCenterName: data["centerName"].(string),
		IssueDate:       issueDate,
		ExpiryDate:      expiryDate,
	}

	if err := s.rechargeRepo.CreateCard(card); err != nil {
		return nil, err
	}

	// 创建交易记录
	transaction := &model.CardTransaction{
		ID:          uuid.New().String(),
		CardNo:      cardNo,
		Type:        "issue",
		Amount:      card.Balance,
		BalanceAfter: card.Balance,
		Remark:      "发卡",
		OperatorID:  data["operatorId"].(string),
	}
	s.rechargeRepo.CreateCardTransaction(transaction)

	return card, nil
}

// generateCardNo 生成卡号
func (s *RechargeService) generateCardNo() string {
	// TJ + 年份后2位 + 5位序号
	// TODO: 从数据库获取最大序号，保证唯一性
	year := time.Now().Format("06")
	sequence := 12345 // 临时序号
	return fmt.Sprintf("TJ%s%05d", year, sequence)
}

// VerifyCard 验证门店卡
func (s *RechargeService) VerifyCard(cardNo string) (*model.StoreCard, error) {
	card, err := s.rechargeRepo.GetCardByCardNo(cardNo)
	if err != nil {
		return nil, errors.New("卡号不存在")
	}

	if card.Status != "active" {
		return nil, errors.New("卡状态异常")
	}

	// 检查是否过期
	if time.Now().After(card.ExpiryDate) {
		return nil, errors.New("卡已过期")
	}

	return card, nil
}

// ConsumeCard 核销门店卡
func (s *RechargeService) ConsumeCard(cardNo string, amount float64, remark, operatorID string) error {
	card, err := s.rechargeRepo.GetCardByCardNo(cardNo)
	if err != nil {
		return errors.New("卡号不存在")
	}

	// 检查卡状态
	if card.Status != "active" {
		return errors.New("卡状态异常，无法核销")
	}

	// 检查余额
	if amount > card.Balance {
		return errors.New("余额不足")
	}

	// 更新余额
	newBalance := card.Balance - amount
	if err := s.rechargeRepo.UpdateCardBalance(cardNo, newBalance); err != nil {
		return err
	}

	// 创建交易记录
	transaction := &model.CardTransaction{
		ID:           uuid.New().String(),
		CardNo:       cardNo,
		Type:         "consume",
		Amount:       amount,
		BalanceAfter: newBalance,
		Remark:       remark,
		OperatorID:   operatorID,
	}

	return s.rechargeRepo.CreateCardTransaction(transaction)
}

// GetCardList 获取门店卡列表
func (s *RechargeService) GetCardList(status, cardNo, holderPhone string, page, pageSize int) (map[string]interface{}, error) {
	list, total, err := s.rechargeRepo.GetCardList(status, cardNo, holderPhone, page, pageSize)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"list":  list,
		"total": total,
	}, nil
}

// GetCardDetail 获取门店卡详情
func (s *RechargeService) GetCardDetail(cardNo string) (map[string]interface{}, error) {
	card, err := s.VerifyCard(cardNo)
	if err != nil {
		return nil, err
	}

	transactions, err := s.rechargeRepo.GetCardTransactions(cardNo)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"card":         card,
		"transactions": transactions,
	}, nil
}

// GetCardStats 获取门店卡统计
func (s *RechargeService) GetCardStats() (map[string]interface{}, error) {
	stats, err := s.rechargeRepo.GetCardStats()
	if err != nil {
		return nil, err
	}

	// 转换为 map[string]interface{}
	result := make(map[string]interface{})
	for k, v := range stats {
		result[k] = v
	}

	// 今日消费和7天内过期需要从数据库查询
	// TODO: 在Repository层实现这些统计
	result["todayConsume"] = int64(0)
	result["expireIn7Days"] = int64(0)

	return result, nil
}

// ========== 充值中心 ==========

// GetCenters 获取充值中心列表
func (s *RechargeService) GetCenters() ([]model.RechargeCenter, error) {
	return s.rechargeRepo.GetCenters()
}

// GetCenterDetail 获取充值中心详情
func (s *RechargeService) GetCenterDetail(id string) (*model.RechargeCenter, error) {
	return s.rechargeRepo.GetCenterByID(id)
}

// CreateCenter 创建充值中心
func (s *RechargeService) CreateCenter(data map[string]interface{}) (*model.RechargeCenter, error) {
	center := &model.RechargeCenter{
		ID:      uuid.New().String(),
		Name:    data["name"].(string),
		Code:    data["code"].(string),
		Status:  "active",
	}

	if v, ok := data["province"].(string); ok {
		center.Province = v
	}
	if v, ok := data["city"].(string); ok {
		center.City = v
	}
	if v, ok := data["district"].(string); ok {
		center.District = v
	}
	if v, ok := data["address"].(string); ok {
		center.Address = v
	}
	if v, ok := data["phone"].(string); ok {
		center.Phone = v
	}
	if v, ok := data["managerId"].(string); ok {
		center.ManagerID = v
	}

	if err := s.rechargeRepo.CreateCenter(center); err != nil {
		return nil, err
	}

	return center, nil
}

// UpdateCenter 更新充值中心
func (s *RechargeService) UpdateCenter(id string, data map[string]interface{}) (*model.RechargeCenter, error) {
	center := &model.RechargeCenter{
		ID:     id,
		Status: data["status"].(string),
	}

	if v, ok := data["name"].(string); ok {
		center.Name = v
	}
	if v, ok := data["code"].(string); ok {
		center.Code = v
	}
	if v, ok := data["province"].(string); ok {
		center.Province = v
	}
	if v, ok := data["city"].(string); ok {
		center.City = v
	}
	if v, ok := data["district"].(string); ok {
		center.District = v
	}
	if v, ok := data["address"].(string); ok {
		center.Address = v
	}
	if v, ok := data["phone"].(string); ok {
		center.Phone = v
	}
	if v, ok := data["managerId"].(string); ok {
		center.ManagerID = v
	}

	if err := s.rechargeRepo.UpdateCenter(center); err != nil {
		return nil, err
	}

	return center, nil
}

// DeleteCenter 删除充值中心
func (s *RechargeService) DeleteCenter(id string) error {
	return s.rechargeRepo.DeleteCenter(id)
}

// ========== 操作员 ==========

// GetOperators 获取操作员列表
func (s *RechargeService) GetOperators() ([]model.RechargeOperator, error) {
	return s.rechargeRepo.GetOperators()
}

// CreateOperator 创建操作员
func (s *RechargeService) CreateOperator(data map[string]interface{}) (*model.RechargeOperator, error) {
	operator := &model.RechargeOperator{
		ID:       uuid.New().String(),
		Name:     data["name"].(string),
		Phone:    data["phone"].(string),
		Password: data["password"].(string),
		CenterID: data["centerId"].(string),
		Role:     data["role"].(string),
		Status:   "active",
	}

	if err := s.rechargeRepo.CreateOperator(operator); err != nil {
		return nil, err
	}

	return operator, nil
}

// UpdateOperator 更新操作员
func (s *RechargeService) UpdateOperator(id string, data map[string]interface{}) (*model.RechargeOperator, error) {
	operator := &model.RechargeOperator{
		ID:       id,
		Name:     data["name"].(string),
		Phone:    data["phone"].(string),
		CenterID: data["centerId"].(string),
		Role:     data["role"].(string),
		Status:   data["status"].(string),
	}

	// 如果有新密码
	if password, ok := data["password"]; ok && password != "" {
		operator.Password = password.(string)
	}

	if err := s.rechargeRepo.UpdateOperator(operator); err != nil {
		return nil, err
	}

	return operator, nil
}

// UpdateCardStatus 更新卡状态
func (s *RechargeService) UpdateCardStatus(cardNo, status string) error {
	return s.rechargeRepo.UpdateCardStatus(cardNo, status)
}

// DeleteOperator 删除操作员
func (s *RechargeService) DeleteOperator(id string) error {
	return s.rechargeRepo.DeleteOperator(id)
}
