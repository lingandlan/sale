package service

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"marketplace/backend/internal/model"
	"marketplace/backend/internal/repository"
	"marketplace/backend/pkg/errno"
	"marketplace/backend/pkg/logger"

	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
)

type RechargeService struct {
	rechargeRepo  repository.RechargeRepoInterface
	memberService MemberServiceInterface
}

func NewRechargeService(rechargeRepo repository.RechargeRepoInterface, memberService MemberServiceInterface) *RechargeService {
	return &RechargeService{rechargeRepo: rechargeRepo, memberService: memberService}
}

// ========== Map 参数安全解析 helper ==========

func getFloat64(data map[string]interface{}, key string) (float64, error) {
	v, ok := data[key]
	if !ok {
		return 0, errno.New(errno.CodeInvalidParam, fmt.Sprintf("参数 %s 不能为空", key))
	}
	val, ok := v.(float64)
	if !ok {
		return 0, errno.New(errno.CodeInvalidParam, fmt.Sprintf("参数 %s 类型错误", key))
	}
	return val, nil
}

func getString(data map[string]interface{}, key string) (string, error) {
	v, ok := data[key]
	if !ok {
		return "", errno.New(errno.CodeInvalidParam, fmt.Sprintf("参数 %s 不能为空", key))
	}
	val, ok := v.(string)
	if !ok {
		return "", errno.New(errno.CodeInvalidParam, fmt.Sprintf("参数 %s 类型错误", key))
	}
	return val, nil
}

// ========== B端充值申请 ==========

// CalculatePoints 计算积分
// 规则：基础积分 = 金额，返还积分 = 金额 * 返还比例
// 返还比例：上月净消费>=10万元返2%，否则返1%
func (s *RechargeService) CalculatePoints(amount float64, lastMonthConsumption float64) (basePoints, rebatePoints, totalPoints, rebateRate int) {
	basePoints = int(amount)
	if lastMonthConsumption >= 100000 {
		rebateRate = 2
	} else {
		rebateRate = 1
	}
	rebatePoints = int(float64(basePoints) * float64(rebateRate) / 100)
	totalPoints = basePoints + rebatePoints
	return
}

// CreateBRechargeApplication 创建B端充值申请
func (s *RechargeService) CreateBRechargeApplication(data map[string]interface{}) (*model.RechargeApplication, error) {
	// 计算积分
	amount, err := getFloat64(data, "amount")
	if err != nil {
		return nil, err
	}
	lastMonthConsumption, err := getFloat64(data, "lastMonthConsumption")
	if err != nil {
		return nil, err
	}
	basePoints, rebatePoints, totalPoints, rebateRate := s.CalculatePoints(amount, lastMonthConsumption)

	centerID, err := getString(data, "centerId")
	if err != nil {
		return nil, err
	}
	centerName, err := getString(data, "centerName")
	if err != nil {
		return nil, err
	}
	applicantID, err := getString(data, "applicantId")
	if err != nil {
		return nil, err
	}
	applicantName, err := getString(data, "applicantName")
	if err != nil {
		return nil, err
	}
	transactionNo, err := getString(data, "transactionNo")
	if err != nil {
		return nil, err
	}
	screenshot, err := getString(data, "screenshot")
	if err != nil {
		return nil, err
	}
	remark, err := getString(data, "remark")
	if err != nil {
		return nil, err
	}

	// 创建申请记录
	app := &model.RechargeApplication{
		ID:            uuid.New().String(),
		CenterID:      centerID,
		CenterName:    centerName,
		Amount:        amount,
		BasePoints:    basePoints,
		RebatePoints:  rebatePoints,
		Points:        totalPoints,
		RebateRate:    rebateRate,
		ApplicantID:   applicantID,
		ApplicantName: applicantName,
		TransactionNo: transactionNo,
		Screenshot:    screenshot,
		Remark:        remark,
		LastMonthConsumption: lastMonthConsumption,
		Status:        "pending",
	}

	if err := s.rechargeRepo.CreateRechargeApplication(app); err != nil {
		return nil, err
	}

	return app, nil
}

// GetRechargeApplicationList 获取充值申请列表
func (s *RechargeService) GetRechargeApplicationList(status string, centerID string, page, pageSize int) (map[string]interface{}, error) {
	list, total, err := s.rechargeRepo.GetRechargeApplications(status, centerID, page, pageSize)
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
		return errno.New(errno.CodeRechargeNotFound, "充值申请不存在")
	}
	if app.Status != "pending" {
		return errno.New(errno.CodeAlreadyProcessed, "该申请已被处理")
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
	memberPhone, err := getString(data, "memberPhone")
	if err != nil {
		return nil, err
	}
	amount, err := getFloat64(data, "amount")
	if err != nil {
		return nil, err
	}
	centerID, err := getString(data, "centerId")
	if err != nil {
		return nil, err
	}
	points := int(amount)

	// 获取充值中心信息及余额
	center, err := s.rechargeRepo.GetCenterByID(centerID)
	if err != nil {
		return nil, errno.New(errno.CodeCenterNotFound, "充值中心不存在")
	}
	if center.Balance < amount {
		return nil, errno.New(errno.CodeCenterNoBalance, "充值中心余额不足")
	}

	// 获取会员当前余额（WSY积分）
	memberBalanceBefore := 0
	if info, err := s.memberService.SearchByPhone(memberPhone); err == nil {
		memberBalanceBefore = int(info.Balance)
	}

	memberID, err := getString(data, "memberId")
	if err != nil {
		return nil, err
	}
	memberName, err := getString(data, "memberName")
	if err != nil {
		return nil, err
	}
	centerName, err := getString(data, "centerName")
	if err != nil {
		return nil, err
	}
	paymentMethod, err := getString(data, "paymentMethod")
	if err != nil {
		return nil, err
	}
	operatorID, err := getString(data, "operatorId")
	if err != nil {
		return nil, err
	}
	operatorName, err := getString(data, "operatorName")
	if err != nil {
		return nil, err
	}
	remark, err := getString(data, "remark")
	if err != nil {
		return nil, err
	}

	recharge := &model.CRecharge{
		ID:            uuid.New().String(),
		MemberID:      memberID,
		MemberName:    memberName,
		MemberPhone:   memberPhone,
		CenterID:      centerID,
		CenterName:    centerName,
		Amount:        amount,
		Points:        points,
		PaymentMethod: paymentMethod,
		OperatorID:    operatorID,
		OperatorName:  operatorName,
		Remark:        remark,

		BalanceBefore: memberBalanceBefore,
		BalanceAfter:  memberBalanceBefore + points,
		Status:        "pending",
	}

	if err := s.rechargeRepo.CreateCRecharge(recharge); err != nil {
		return nil, err
	}

	// 扣减充值中心余额
	if _, err := s.rechargeRepo.DeductCenterBalance(centerID, amount); err != nil {
		return nil, errors.New("扣减中心余额失败")
	}

	// WSY加积分
	batchcode := fmt.Sprintf("%d%s", time.Now().Unix(), recharge.ID)
	if len(batchcode) > 30 {
		batchcode = batchcode[:30]
	}
	if afterIntegral, apiErr := s.memberService.AddIntegral(memberPhone, float64(points), batchcode, fmt.Sprintf("充值中心充值 %s", centerID)); apiErr == nil {
		recharge.BalanceAfter = int(afterIntegral)
		s.rechargeRepo.UpdateCRecharge(recharge.ID, map[string]interface{}{
			"balance_before": memberBalanceBefore,
			"balance_after":  int(afterIntegral),
			"status":         "success",
		})
		recharge.Status = "success"
	} else {
		// WSY加积分失败，记录错误，保持 status = pending
		logger.Error("WSY AddIntegral failed, recharge pending",
			zap.String("rechargeId", recharge.ID),
			zap.String("memberPhone", memberPhone),
			zap.Int("points", points),
			zap.String("centerID", centerID),
			zap.Error(apiErr),
		)
		return nil, errno.New(errno.CodeInvalidParam, fmt.Sprintf("充值记录已创建，但积分发放失败，请联系管理员处理（单号：%s）", recharge.ID))
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

// 状态转换规则：
// - 非冻结状态 → 可冻结
// - 冻结状态 → 可解冻（转为已激活）

// cardTypeMap 卡类型中文映射
var cardTypeMap = map[string]int{
	"实体":   model.CardTypePhysical,
	"实体卡":  model.CardTypePhysical,
	"虚拟":   model.CardTypeVirtual,
	"虚拟卡":  model.CardTypeVirtual,
}

// parseExcel 解析 Excel 文件为行数据
func (s *RechargeService) parseExcel(file []byte) ([][]string, error) {
	f, err := excelize.OpenReader(bytes.NewReader(file))
	if err != nil {
		return nil, errors.New("无法解析Excel文件")
	}
	defer f.Close()

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, errors.New("Excel文件没有工作表")
	}
	rows, err := f.GetRows(sheets[0])
	if err != nil {
		return nil, errors.New("无法读取Excel内容")
	}
	return rows, nil
}

// parseCSV 解析 CSV 文件为行数据（逗号分隔，UTF-8 BOM 兼容）
func (s *RechargeService) parseCSV(file []byte) ([][]string, error) {
	text := string(file)
	// 去掉 UTF-8 BOM
	text = strings.TrimPrefix(text, "\uFEFF")

	var rows [][]string
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		rows = append(rows, strings.Split(line, ","))
	}
	if len(rows) == 0 {
		return nil, errors.New("CSV文件没有数据")
	}
	return rows, nil
}

// BatchImportCards 批量入库门店卡（支持 xlsx/csv）
func (s *RechargeService) BatchImportCards(file []byte, ext string, operatorID string) (int, []string, error) {
	// 1. 根据文件类型解析为行数据
	var rows [][]string
	var err error

	switch ext {
	case ".csv":
		rows, err = s.parseCSV(file)
	default: // .xlsx, .xls
		rows, err = s.parseExcel(file)
	}
	if err != nil {
		return 0, nil, err
	}

	// 跳过表头，至少需要2行（1行表头+1行数据）
	if len(rows) <= 1 {
		return 0, nil, errors.New("Excel文件没有数据行（至少需要表头+1行数据）")
	}

	dataRows := rows[1:] // 跳过表头
	if len(dataRows) > 1000 {
		return 0, nil, errors.New("单次入库不能超过1000张")
	}

	// 3. 解析并校验每行
	cardNoSet := make(map[string]bool)
	cards := make([]*model.StoreCard, 0, len(dataRows))
	cardNos := make([]string, 0, len(dataRows))
	batchNo := fmt.Sprintf("B%s", time.Now().Format("20060102150405"))

	for i, row := range dataRows {
		rowNum := i + 2 // Excel行号（1-based，跳过表头）
		if len(row) < 3 {
			return 0, nil, fmt.Errorf("第%d行列数不足，需要3列（卡号、卡类型、面值）", rowNum)
		}

		cardNo := strings.TrimSpace(row[0])
		cardTypeStr := strings.TrimSpace(row[1])
		balanceStr := strings.TrimSpace(row[2])

		// 校验卡号非空
		if cardNo == "" {
			return 0, nil, fmt.Errorf("第%d行卡号为空", rowNum)
		}

		// 校验卡类型
		cardType, ok := cardTypeMap[cardTypeStr]
		if !ok {
			return 0, nil, fmt.Errorf("第%d行卡类型错误（需为实体/实体卡/虚拟/虚拟卡）: %s", rowNum, cardTypeStr)
		}

		// 校验面值（正整数）
		balance, err := strconv.Atoi(balanceStr)
		if err != nil || balance <= 0 {
			return 0, nil, fmt.Errorf("第%d行面值必须为正整数: %s", rowNum, balanceStr)
		}

		// Excel内卡号去重
		if cardNoSet[cardNo] {
			return 0, nil, fmt.Errorf("第%d行卡号重复: %s", rowNum, cardNo)
		}
		cardNoSet[cardNo] = true

		cardNos = append(cardNos, cardNo)
		cards = append(cards, &model.StoreCard{
			ID:       uuid.New().String(),
			CardNo:   cardNo,
			CardType: cardType,
			Status:   model.CardStatusInStock,
			Balance:  balance,
			BatchNo:  batchNo,
		})
	}

	// 6. 入库
	if err := s.rechargeRepo.BatchCreateCards(cards); err != nil {
		return 0, nil, err
	}

	// 7. 创建入库交易记录（事务）
	txns := make([]*model.CardTransaction, 0, len(cards))
	for _, card := range cards {
		txns = append(txns, &model.CardTransaction{
			ID:            uuid.New().String(),
			CardNo:        card.CardNo,
			Type:          "stock_in",
			Amount:        0,
			BalanceBefore: card.Balance,
			BalanceAfter:  card.Balance,
			Remark:        fmt.Sprintf("批量入库（批次%s）", batchNo),
			OperatorID:    operatorID,
		})
	}
	if err := s.rechargeRepo.BatchCreateCardTransactions(txns); err != nil {
		return 0, nil, fmt.Errorf("创建入库交易记录失败: %w", err)
	}

	return len(cards), cardNos, nil
}

// AllocateCards 按数量划拨卡到充值中心
func (s *RechargeService) AllocateCards(centerID string, quantity int) (int, error) {
	// 校验充值中心存在
	center, err := s.rechargeRepo.GetCenterByID(centerID)
	if err != nil || center == nil {
		return 0, errno.New(errno.CodeCenterNotFound, "充值中心不存在")
	}

	// 校验库存充足
	available, err := s.rechargeRepo.GetAllocatableCardCount()
	if err != nil {
		return 0, err
	}
	if available < int64(quantity) {
		return 0, errno.Newf(errno.CodeInsufficientStock, "库存不足，当前可划拨%d张，需要%d张", available, quantity)
	}

	count, err := s.rechargeRepo.AllocateCardsByQuantity(centerID, quantity)
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
		return errno.New(errno.CodeCardNotFound, "卡号不存在")
	}
	if card.Status != model.CardStatusInStock {
		return errno.New(errno.CodeCardNotInStock, "该卡不在库存中，无法发放")
	}
	if card.RechargeCenterID != rechargeCenterID {
		return errno.New(errno.CodeCardNotInCenter, "该卡不在本中心库存中")
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

	// 发放交易记录（在同一事务中创建）
	txn := &model.CardTransaction{
		ID:            uuid.New().String(),
		CardNo:        cardNo,
		Type:          "issue",
		Amount:        0,
		BalanceBefore: card.Balance,
		BalanceAfter:  card.Balance,
		Remark:        fmt.Sprintf("发放给用户 %s（%s）", userPhone, issueReason),
		OperatorID:    operatorID,
	}

	if err := s.rechargeRepo.BindCardToUser(cardNo, updates, record, txn); err != nil {
		return err
	}

	return nil
}

// VerifyCard 验证门店卡
func (s *RechargeService) VerifyCard(cardNo string) (*model.StoreCard, error) {
	card, err := s.rechargeRepo.GetCardByCardNo(cardNo)
	if err != nil {
		return nil, errno.New(errno.CodeCardNotFound, "卡号不存在")
	}
	if card.Status == model.CardStatusFrozen {
		return nil, errno.New(errno.CodeCardFrozen, "卡已冻结")
	}
	if card.Status == model.CardStatusExpired {
		return nil, errno.New(errno.CodeCardExpired, "卡已过期")
	}
	if card.Status == model.CardStatusVoided {
		return nil, errno.New(errno.CodeCardVoided, "卡已作废")
	}
	if card.Status == model.CardStatusInStock {
		return nil, errno.New(errno.CodeCardNotIssued, "卡未发放")
	}
	// 已发放(2)和已激活(3)的卡可以查询
	// 已激活的卡检查是否过期
	if card.Status == model.CardStatusActive && card.ExpiredAt != nil && time.Now().After(*card.ExpiredAt) {
		// 自动标记为过期
		s.rechargeRepo.UpdateCardByMap(cardNo, map[string]interface{}{"status": model.CardStatusExpired})
		return nil, errno.New(errno.CodeCardExpired, "卡已过期")
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

	// 附加充值中心名称
	type CardWithCenter struct {
		model.StoreCard
		RechargeCenterName string `json:"rechargeCenterName"`
	}
	centerCache := make(map[string]string)
	resultList := make([]CardWithCenter, 0, len(list))
	for _, card := range list {
		name := ""
		if card.RechargeCenterID != "" {
			if n, ok := centerCache[card.RechargeCenterID]; ok {
				name = n
			} else {
				if c, err := s.rechargeRepo.GetCenterByID(card.RechargeCenterID); err == nil && c != nil {
					name = c.Name
					centerCache[card.RechargeCenterID] = name
				}
			}
		}
		resultList = append(resultList, CardWithCenter{StoreCard: card, RechargeCenterName: name})
	}

	return map[string]interface{}{
		"list":  resultList,
		"total": total,
	}, nil
}

// GetCardDetail 获取门店卡详情
func (s *RechargeService) GetCardDetail(cardNo string) (map[string]interface{}, error) {
	card, err := s.rechargeRepo.GetCardByCardNo(cardNo)
	if err != nil {
		return nil, errno.New(errno.CodeCardNotFound, "卡号不存在")
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
func (s *RechargeService) GetCardStats(centerID string) (map[string]interface{}, error) {
	stats, err := s.rechargeRepo.GetCardStats(centerID)
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

// GetMonthlyTrend 获取月度发放/核销趋势
func (s *RechargeService) GetMonthlyTrend(centerID string) (interface{}, error) {
	return s.rechargeRepo.GetMonthlyTrend(centerID)
}

// GetCenterCardStats 获取充值中心维度卡统计
func (s *RechargeService) GetCenterCardStats(centerID string) (interface{}, error) {
	return s.rechargeRepo.GetCenterCardStats(centerID)
}

// transitionCardStatus 卡状态转换通用方法
func (s *RechargeService) transitionCardStatus(cardNo string, targetStatus int, operatorID, txnType, remark string) error {
	card, err := s.rechargeRepo.GetCardByCardNo(cardNo)
	if err != nil {
		return errno.New(errno.CodeCardNotFound, "卡号不存在")
	}

	// 校验状态转换合法性
	if targetStatus == model.CardStatusFrozen {
		// 任何非冻结状态都可以冻结
		if card.Status == model.CardStatusFrozen {
			return errno.New(errno.CodeInvalidAction, "卡已冻结，不能重复冻结")
		}
	} else if targetStatus == model.CardStatusActive {
		// 只有冻结状态可以解冻
		if card.Status != model.CardStatusFrozen {
			return errno.New(errno.CodeInvalidAction, "只有冻结状态的卡才能解冻")
		}
	} else {
		return errno.New(errno.CodeInvalidAction, "不支持的状态转换")
	}

	if err := s.rechargeRepo.TransitionCardStatusTX(cardNo, map[string]interface{}{"status": targetStatus}, &model.CardTransaction{
		ID:            uuid.New().String(),
		CardNo:        cardNo,
		Type:          txnType,
		Amount:        0,
		BalanceBefore: card.Balance,
		BalanceAfter:  card.Balance,
		Remark:        remark,

		OperatorID:    operatorID,
	}); err != nil {
		return err
	}
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

// GetAvailableCards 获取指定充值中心的可用卡号列表
func (s *RechargeService) GetAvailableCards(centerID string, keyword string) ([]string, error) {
	if centerID == "" {
		return nil, errors.New("充值中心ID不能为空")
	}
	return s.rechargeRepo.GetAvailableCardNos(centerID, keyword)
}

// GetAvailableCardCount 获取指定充值中心的可用卡数量
func (s *RechargeService) GetAvailableCardCount(centerID string) (int64, error) {
	if centerID == "" {
		return 0, errors.New("充值中心ID不能为空")
	}
	return s.rechargeRepo.GetAvailableCardCount(centerID)
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

			if tr, err := s.rechargeRepo.GetCenterTotalRecharge(c.ID); err == nil {
			item["totalRecharge"] = tr
		} else {
			item["totalRecharge"] = int64(0)
		}
			if tc, err := s.rechargeRepo.GetCenterTotalConsumed(c.ID); err == nil {
			item["totalConsumed"] = tc
		} else {
			item["totalConsumed"] = float64(0)
		}


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
	name, err := getString(data, "name")
	if err != nil {
		return nil, err
	}
	code, err := getString(data, "code")
	if err != nil {
		return nil, err
	}

	center := &model.RechargeCenter{
		ID:     uuid.New().String(),
		Name:   name,
		Code:   code,
		Status: "active",
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
	name, err := getString(data, "name")
	if err != nil {
		return nil, err
	}
	phone, err := getString(data, "phone")
	if err != nil {
		return nil, err
	}
	password, err := getString(data, "password")
	if err != nil {
		return nil, err
	}
	centerID, err := getString(data, "centerId")
	if err != nil {
		return nil, err
	}
	role, err := getString(data, "role")
	if err != nil {
		return nil, err
	}

	// 密码哈希存储
	hashedPwd, err := HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("密码加密失败: %w", err)
	}

	operator := &model.RechargeOperator{
		ID:       uuid.New().String(),
		Name:     name,
		Phone:    phone,
		Password: hashedPwd,
		CenterID: centerID,
		Role:     role,
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
		hashed, err := HashPassword(v)
		if err != nil {
			return nil, fmt.Errorf("密码加密失败: %w", err)
		}
		updates["password"] = hashed
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

// ========== Dashboard ==========

func calcTrend(current, yesterday float64) string {
	if yesterday == 0 {
		if current > 0 {
			return "+100%"
		}
		return "—"
	}
	pct := (current - yesterday) / yesterday * 100
	if pct >= 0 {
		return fmt.Sprintf("+%.0f%%", pct)
	}
	return fmt.Sprintf("%.0f%%", pct)
}

func (s *RechargeService) GetDashboardStatistics(role, centerID string) (map[string]interface{}, error) {
	todayRecharge, err := s.rechargeRepo.GetTodayRechargeTotal(centerID)
	if err != nil {
		return nil, fmt.Errorf("统计数据加载失败: %w", err)
	}
	todayConsumption, err := s.rechargeRepo.GetTodayConsumptionTotal(centerID)
	if err != nil {
		return nil, fmt.Errorf("统计数据加载失败: %w", err)
	}
	activeCenters, err := s.rechargeRepo.GetActiveCenterCount(centerID)
	if err != nil {
		return nil, fmt.Errorf("统计数据加载失败: %w", err)
	}
	yesterdayRecharge, _ := s.rechargeRepo.GetYesterdayRechargeTotal(centerID)
	yesterdayConsumption, _ := s.rechargeRepo.GetYesterdayConsumptionTotal(centerID)

	return map[string]interface{}{
		"memberCount":       0,
		"memberTrend":       "—",
		"todayRecharge":     todayRecharge,
		"rechargeTrend":     calcTrend(todayRecharge, yesterdayRecharge),
		"todayConsumption":  todayConsumption,
		"consumptionTrend":  calcTrend(todayConsumption, yesterdayConsumption),
		"activeCenters":     activeCenters,
		"centerTrend":       "—",
	}, nil
}

func (s *RechargeService) GetDashboardTodos(role, centerID string) (map[string]interface{}, error) {
	pendingCount, err := s.rechargeRepo.CountPendingApprovals(centerID)
	if err != nil {
		return nil, fmt.Errorf("待办事项加载失败: %w", err)
	}
	expiringCount, err := s.rechargeRepo.CountExpiringCards(centerID)
	if err != nil {
		return nil, fmt.Errorf("待办事项加载失败: %w", err)
	}

	pendingDesc := ""
	if pendingCount > 0 {
		pendingDesc = fmt.Sprintf("%d笔充值申请待审批", pendingCount)
	}
	expiringDesc := ""
	if expiringCount > 0 {
		expiringDesc = fmt.Sprintf("%d张门店卡将在7天内到期", expiringCount)
	}

	return map[string]interface{}{
		"pendingApprovals": map[string]interface{}{
			"count":       pendingCount,
			"description": pendingDesc,
		},
		"expiringCards": map[string]interface{}{
			"count":       expiringCount,
			"description": expiringDesc,
		},
	}, nil
}

func (s *RechargeService) GetDashboardRechargeTrends(days int, role, centerID string) (map[string]interface{}, error) {
	dates, values, err := s.rechargeRepo.GetRechargeTrends(days, centerID)
	if err != nil {
		return nil, fmt.Errorf("趋势数据加载失败: %w", err)
	}
	return map[string]interface{}{
		"dates":  dates,
		"values": values,
	}, nil
}

// ========== 月度消费 ==========

// GetCenterLastMonthConsumption 查询充值中心上月消费金额及返还比例
func (s *RechargeService) GetCenterLastMonthConsumption(centerID string) (map[string]interface{}, error) {
	// 上月 YYYY-MM
	now := time.Now()
	lastMonth := time.Date(now.Year(), now.Month()-1, 1, 0, 0, 0, 0, time.UTC)
	monthStr := lastMonth.Format("2006-01")

	record, err := s.rechargeRepo.GetMonthlyConsumption(centerID, monthStr)
	if err != nil {
		return nil, err
	}

	consumption := float64(0)
	if record != nil {
		consumption = record.Consumption
	}

	rebateRate := 1
	if consumption >= 100000 {
		rebateRate = 2
	}

	return map[string]interface{}{
		"consumption": consumption,
		"rebateRate":  rebateRate,
		"month":       monthStr,
	}, nil
}

// UpsertMonthlyConsumption 录入/更新月度消费
func (s *RechargeService) UpsertMonthlyConsumption(centerID, month string, consumption float64) error {
	// 查已有记录
	existing, err := s.rechargeRepo.GetMonthlyConsumption(centerID, month)
	if err != nil {
		return err
	}

	if existing != nil {
		existing.Consumption = consumption
		return s.rechargeRepo.UpsertMonthlyConsumption(existing)
	}

	record := &model.CenterMonthlyConsumption{
		ID:          uuid.New().String(),
		CenterID:    centerID,
		Month:       month,
		Consumption: consumption,
	}
	return s.rechargeRepo.UpsertMonthlyConsumption(record)
}

// ListMonthlyConsumption 查询月度消费列表
func (s *RechargeService) ListMonthlyConsumption(month string) ([]model.CenterMonthlyConsumption, error) {
	return s.rechargeRepo.ListMonthlyConsumption(month)
}
