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
// 返还比例：上月净消费>=10万元返2%，否则返1%
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
	app, err := s.rechargeRepo.GetRechargeApplicationByID(id)
	if err != nil {
		return errors.New("充值申请不存在")
	}
	if app.Status != "pending" {
		return errors.New("该申请已被处理")
	}

	var status string
	if action == "approve" {
		status = "approved"
		// 增加充值中心积分（Points = 基础积分 + 返利积分）
		if err := s.rechargeRepo.AddCenterBalance(app.CenterID, float64(app.Points)); err != nil {
			return fmt.Errorf("增加充值中心余额失败: %w", err)
		}
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
func (s *RechargeService) GetCRechargeList(memberPhone, centerID, startDate, endDate string, page, pageSize int) (map[string]interface{}, error) {
	list, total, err := s.rechargeRepo.GetCRechargeList(memberPhone, centerID, startDate, endDate, page, pageSize)
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

// 状态转换白名单
var allowedTransitions = map[int][]int{
	model.CardStatusInStock: {model.CardStatusVoided},
	model.CardStatusIssued:  {model.CardStatusFrozen, model.CardStatusVoided},
	model.CardStatusActive:  {model.CardStatusFrozen, model.CardStatusExpired},
	model.CardStatusFrozen:  {model.CardStatusActive, model.CardStatusVoided},
}

// BatchImportCards 批量入库门店卡
func (s *RechargeService) BatchImportCards(startSeq, endSeq, cardType int, operatorID string) ([]string, error) {
	if startSeq <= 0 || endSeq <= 0 || startSeq > endSeq {
		return nil, errors.New("序号范围无效")
	}
	if endSeq-startSeq+1 > 1000 {
		return nil, errors.New("单次入库不能超过1000张")
	}

	// 获取当前最大序号
	maxSeq, err := s.rechargeRepo.GetMaxCardSequence()
	if err != nil {
		return nil, err
	}

	// 检查序号冲突
	if startSeq <= maxSeq {
		return nil, fmt.Errorf("起始序号 %d 与已有卡号冲突（当前最大序号 %d）", startSeq, maxSeq)
	}

	batchNo := fmt.Sprintf("B%s", time.Now().Format("20060102150405"))
	cards := make([]*model.StoreCard, 0, endSeq-startSeq+1)
	cardNos := make([]string, 0, endSeq-startSeq+1)

	for seq := startSeq; seq <= endSeq; seq++ {
		cardNo := fmt.Sprintf("TJ%08d", seq)
		cardNos = append(cardNos, cardNo)
		cards = append(cards, &model.StoreCard{
			ID:       uuid.New().String(),
			CardNo:   cardNo,
			CardType: cardType,
			Status:   model.CardStatusInStock,
			Balance:  1000,
			BatchNo:  batchNo,
		})
	}

	if err := s.rechargeRepo.BatchCreateCards(cards); err != nil {
		return nil, err
	}

	// 创建入库交易记录
	for _, cardNo := range cardNos {
		s.rechargeRepo.CreateCardTransaction(&model.CardTransaction{
			ID:            uuid.New().String(),
			CardNo:        cardNo,
			Type:          "stock_in",
			Amount:        0,
			BalanceBefore: 1000,
			BalanceAfter:  1000,
			Remark:        fmt.Sprintf("批量入库（批次%s）", batchNo),
			OperatorID:    operatorID,
		})
	}

	return cardNos, nil
}

// AllocateCards 将卡号段划拨到充值中心
func (s *RechargeService) AllocateCards(centerID, startCardNo, endCardNo string) (int, error) {
	// 校验充值中心存在
	center, err := s.rechargeRepo.GetCenterByID(centerID)
	if err != nil || center == nil {
		return 0, errors.New("充值中心不存在")
	}

	count, err := s.rechargeRepo.AllocateCardsToCenter(centerID, startCardNo, endCardNo)
	if err != nil {
		return 0, err
	}
	if count == 0 {
		return 0, errors.New("没有符合条件的卡可划拨")
	}
	return count, nil
}

// BindCardToUser 绑定卡号到用户
func (s *RechargeService) BindCardToUser(cardNo, userPhone, issueReason string, issueType int, rechargeCenterID, operatorID, relatedUserPhone, remark string) error {
	// 查卡
	card, err := s.rechargeRepo.GetCardByCardNo(cardNo)
	if err != nil {
		return errors.New("卡号不存在")
	}
	if card.Status != model.CardStatusInStock {
		return errors.New("该卡不在库存中，无法发放")
	}
	if card.RechargeCenterID != rechargeCenterID {
		return errors.New("该卡不在本中心库存中")
	}

	// TODO: 调商城接口验证用户存在
	userID := "" // 后续对接商城接口获取

	now := time.Now()
	updates := map[string]interface{}{
		"status":       model.CardStatusIssued,
		"user_id":      userID,
		"issue_reason": issueReason,
		"issued_at":    &now,
	}

	record := &model.CardIssueRecord{
		ID:               uuid.New().String(),
		CardNo:           cardNo,
		UserID:           userID,
		UserPhone:        userPhone,
		IssueReason:      issueReason,
		IssueType:        issueType,
		RechargeCenterID: rechargeCenterID,
		OperatorID:       operatorID,
		RelatedUserPhone: relatedUserPhone,
		Remark:           remark,
	}

	if err := s.rechargeRepo.BindCardToUser(cardNo, updates, record); err != nil {
		return err
	}

	// 创建发放交易记录
	s.rechargeRepo.CreateCardTransaction(&model.CardTransaction{
		ID:            uuid.New().String(),
		CardNo:        cardNo,
		Type:          "issue",
		Amount:        0,
		BalanceBefore: card.Balance,
		BalanceAfter:  card.Balance,
		Remark:        fmt.Sprintf("发放给用户 %s（%s）", userPhone, issueReason),
		OperatorID:    operatorID,
	})

	return nil
}

// VerifyCard 验证门店卡
func (s *RechargeService) VerifyCard(cardNo string) (*model.StoreCard, error) {
	card, err := s.rechargeRepo.GetCardByCardNo(cardNo)
	if err != nil {
		return nil, errors.New("卡号不存在")
	}
	if card.Status == model.CardStatusFrozen {
		return nil, errors.New("卡已冻结")
	}
	if card.Status == model.CardStatusExpired {
		return nil, errors.New("卡已过期")
	}
	if card.Status == model.CardStatusVoided {
		return nil, errors.New("卡已作废")
	}
	if card.Status == model.CardStatusInStock {
		return nil, errors.New("卡未发放")
	}
	// 已发放(2)和已激活(3)的卡可以查询
	// 已激活的卡检查是否过期
	if card.Status == model.CardStatusActive && card.ExpiredAt != nil && time.Now().After(*card.ExpiredAt) {
		// 自动标记为过期
		s.rechargeRepo.UpdateCardByMap(cardNo, map[string]interface{}{"status": model.CardStatusExpired})
		return nil, errors.New("卡已过期")
	}
	return card, nil
}

// ConsumeCard 核销门店卡（委托给事务方法）
func (s *RechargeService) ConsumeCard(cardNo string, amount int, operatorID, remark string) error {
	return s.rechargeRepo.ConsumeCardInTx(cardNo, amount, operatorID, remark)
}

// GetCardList 获取门店卡列表
func (s *RechargeService) GetCardList(status int, cardNo, centerID string, page, pageSize int) (map[string]interface{}, error) {
	list, total, err := s.rechargeRepo.GetCardList(status, cardNo, centerID, page, pageSize)
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
	card, err := s.rechargeRepo.GetCardByCardNo(cardNo)
	if err != nil {
		return nil, errors.New("卡号不存在")
	}
	transactions, _, err := s.rechargeRepo.GetCardTransactions(cardNo, 1, 50)
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
	result := make(map[string]interface{})
	for k, v := range stats {
		result[k] = v
	}
	return result, nil
}

// GetCardInventoryStats 获取门店卡库存统计
func (s *RechargeService) GetCardInventoryStats() (map[string]interface{}, error) {
	stats, err := s.rechargeRepo.GetCardInventoryStats()
	if err != nil {
		return nil, err
	}
	result := make(map[string]interface{})
	for k, v := range stats {
		result[k] = v
	}
	return result, nil
}

// transitionCardStatus 卡状态转换通用方法
func (s *RechargeService) transitionCardStatus(cardNo string, targetStatus int, operatorID, txnType, remark string) error {
	card, err := s.rechargeRepo.GetCardByCardNo(cardNo)
	if err != nil {
		return errors.New("卡号不存在")
	}

	// 校验状态转换合法性
	allowed, ok := allowedTransitions[card.Status]
	if !ok {
		return errors.New("当前状态不允许操作")
	}
	valid := false
	for _, s := range allowed {
		if s == targetStatus {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("不允许从状态%d转换到状态%d", card.Status, targetStatus)
	}

	if err := s.rechargeRepo.UpdateCardByMap(cardNo, map[string]interface{}{"status": targetStatus}); err != nil {
		return err
	}

	s.rechargeRepo.CreateCardTransaction(&model.CardTransaction{
		ID:            uuid.New().String(),
		CardNo:        cardNo,
		Type:          txnType,
		Amount:        0,
		BalanceBefore: card.Balance,
		BalanceAfter:  card.Balance,
		Remark:        remark,
		OperatorID:    operatorID,
	})
	return nil
}

// FreezeCard 冻结卡
func (s *RechargeService) FreezeCard(cardNo, operatorID string) error {
	return s.transitionCardStatus(cardNo, model.CardStatusFrozen, operatorID, "freeze", "冻结卡")
}

// UnfreezeCard 解冻卡
func (s *RechargeService) UnfreezeCard(cardNo, operatorID string) error {
	return s.transitionCardStatus(cardNo, model.CardStatusActive, operatorID, "unfreeze", "解冻卡")
}

// VoidCard 作废卡
func (s *RechargeService) VoidCard(cardNo, operatorID string) error {
	return s.transitionCardStatus(cardNo, model.CardStatusVoided, operatorID, "void", "作废卡")
}

// ========== 充值中心 ==========

// GetCenters 获取充值中心列表（含管理员、累计充值、已消耗）
func (s *RechargeService) GetCenters() ([]map[string]interface{}, error) {
	centers, err := s.rechargeRepo.GetCenters()
	if err != nil {
		return nil, err
	}

	operators, _ := s.rechargeRepo.GetOperators()
	opMap := make(map[string]*model.RechargeOperator)
	for i := range operators {
		opMap[operators[i].CenterID] = &operators[i]
	}

	result := make([]map[string]interface{}, 0, len(centers))
	for _, c := range centers {
		item := map[string]interface{}{
			"id":        c.ID,
			"name":      c.Name,
			"code":      c.Code,
			"province":  c.Province,
			"city":      c.City,
			"district":  c.District,
			"address":   c.Address,
			"managerId": c.ManagerID,
			"phone":     c.Phone,
			"balance":   c.Balance,
			"status":    c.Status,
			"createdAt": c.CreatedAt,
			"updatedAt": c.UpdatedAt,
		}

		// 管理员：从 operator 列表匹配
		if op, ok := opMap[c.ID]; ok {
			item["managerName"] = op.Name
			item["managerPhone"] = op.Phone
		}

			item["totalRecharge"] = s.rechargeRepo.GetCenterTotalRecharge(c.ID)
			item["totalConsumed"] = s.rechargeRepo.GetCenterTotalConsumed(c.ID)


		result = append(result, item)
	}
	return result, nil
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
	if err := s.rechargeRepo.UpdateCenter(id, data); err != nil {
		return nil, err
	}
	return s.rechargeRepo.GetCenterByID(id)
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
	updates := map[string]interface{}{}

	if v, ok := data["name"].(string); ok && v != "" {
		updates["name"] = v
	}
	if v, ok := data["phone"].(string); ok && v != "" {
		updates["phone"] = v
	}
	if v, ok := data["role"].(string); ok && v != "" {
		updates["role"] = v
	}
	if v, ok := data["status"].(string); ok && v != "" {
		updates["status"] = v
	}
	if v, ok := data["password"].(string); ok && v != "" {
		updates["password"] = v
	}
	// 前端发 center_id，数据库字段 center_id
	if v, ok := data["center_id"].(string); ok && v != "" {
		updates["center_id"] = v
	}
	if v, ok := data["centerId"].(string); ok && v != "" {
		updates["center_id"] = v
	}

	if len(updates) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	if err := s.rechargeRepo.UpdateOperator(id, updates); err != nil {
		return nil, err
	}

	return &model.RechargeOperator{ID: id}, nil
}

// DeleteOperator 删除操作员
func (s *RechargeService) DeleteOperator(id string) error {
	return s.rechargeRepo.DeleteOperator(id)
}
