package handler

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"marketplace/backend/internal/model"
	"marketplace/backend/internal/repository"
	"marketplace/backend/internal/service"
	"marketplace/backend/pkg/errmsg"
	"marketplace/backend/pkg/errno"
	"marketplace/backend/pkg/response"
)

type RechargeHandler struct {
	rechargeService service.RechargeServiceInterface
	userRepo        repository.UserRepositoryInterface
}

func NewRechargeHandler(rechargeService service.RechargeServiceInterface, userRepo repository.UserRepositoryInterface) *RechargeHandler {
	return &RechargeHandler{rechargeService: rechargeService, userRepo: userRepo}
}

// getOperatorInfo 从 JWT context 获取操作员信息（userID, role, centerID, name）
// 对于 super_admin/hq_admin/finance，centerID 返回 ""
func (h *RechargeHandler) getOperatorInfo(c *gin.Context) (int64, string, string, string, error) {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		return 0, "", "", "", fmt.Errorf("未获取到用户身份")
	}
	userID, ok := userIDRaw.(int64)
	if !ok {
		return 0, "", "", "", fmt.Errorf("用户ID类型错误")
	}
	roleRaw, _ := c.Get("role")
	role, _ := roleRaw.(string)

	// 先尝试查 DB 获取用户名（所有角色都需要）
	user, err := h.userRepo.GetByID(c.Request.Context(), userID)
	if err != nil || user == nil {
		// 查不到时降级为用户名未知
		userName := "未知用户"
		// super_admin/hq_admin/finance 不需要 center 限制
		if role == model.RoleSuperAdmin || role == model.RoleHQAdmin || role == model.RoleFinance {
			return userID, role, "", userName, nil
		}
		return userID, role, "", userName, nil
	}

	userName := user.Name
	if userName == "" {
		userName = user.Username
	}

	// super_admin/hq_admin/finance 不需要 center 限制
	if role == model.RoleSuperAdmin || role == model.RoleHQAdmin || role == model.RoleFinance {
		return userID, role, "", userName, nil
	}

	// center_admin/operator 需要 center_id
	if user.CenterID == nil || *user.CenterID == "" {
		return userID, role, "", userName, fmt.Errorf("当前用户未分配充值中心")
	}
	centerID := *user.CenterID
	return userID, role, centerID, userName, nil
}

// bizError 统一处理业务错误：业务错误返回400+具体信息，非业务错误返回500
func bizError(c *gin.Context, err error) {
	code, msg := errno.Resolve(err)
	if code != "" {
		response.Error(c, 400, msg)
	} else {
		response.InternalError(c, msg)
	}
}

// ========== B端充值申请 ==========

func (h *RechargeHandler) CreateBRechargeApplication(c *gin.Context) {
	var req model.CreateBRechargeApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamsError(c, err.Error())
		return
	}

	applicantID, _, _, applicantName, err := h.getOperatorInfo(c)
	if err != nil {
		response.Error(c, 401, err.Error())
		return
	}

	data := map[string]interface{}{
		"memberId":             req.MemberID,
		"centerId":             req.CenterID,
		"amount":               req.Amount,
		"paymentMethod":        req.PaymentMethod,
		"remark":               req.Remark,
		"lastMonthConsumption": req.LastMonthConsumption,
		"centerName":           req.CenterName,
		"transactionNo":        req.TransactionNo,
		"screenshot":           req.Screenshot,
		"memberName":           req.MemberName,
		"memberPhone":          req.MemberPhone,
		"applicantId":          fmt.Sprintf("%d", applicantID),
		"applicantName":        applicantName,
	}

	app, err := h.rechargeService.CreateBRechargeApplication(data)
	if err != nil {
		response.InternalError(c, errmsg.Get("recharge.apply_failed"))
		return
	}

	response.SuccessWithMessage(c, errmsg.Get("recharge.apply_success"), gin.H{
		"id":     app.ID,
		"status": app.Status,
	})
}

func (h *RechargeHandler) GetRechargeApplicationList(c *gin.Context) {
	status := c.DefaultQuery("status", "")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	result, err := h.rechargeService.GetRechargeApplicationList(status, page, pageSize)
	if err != nil {
		response.InternalError(c, errmsg.Get("recharge.list_failed"))
		return
	}

	response.Success(c, result)
}

func (h *RechargeHandler) GetRechargeApplicationDetail(c *gin.Context) {
	id := c.Param("id")

	result, err := h.rechargeService.GetRechargeApplicationDetail(id)
	if err != nil {
		response.NotFound(c, errmsg.Get("recharge.not_found"))
		return
	}

	response.Success(c, result)
}

func (h *RechargeHandler) ApprovalRechargeApplication(c *gin.Context) {
	var req model.ApprovalRechargeApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamsError(c, err.Error())
		return
	}

	approvedByID, _, _, _, err := h.getOperatorInfo(c)
	if err != nil {
		response.Error(c, 401, err.Error())
		return
	}
	approvedBy := fmt.Sprintf("%d", approvedByID)

	if err := h.rechargeService.ApproveRechargeApplication(req.ID, req.Action, approvedBy, req.Reason); err != nil {
		response.InternalError(c, errmsg.Get("recharge.approval_failed"))
		return
	}

	response.Success(c, gin.H{"success": true})
}

// ========== C端充值 ==========

func (h *RechargeHandler) CreateCRecharge(c *gin.Context) {
	var req model.CreateCRechargeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamsError(c, err.Error())
		return
	}

	operatorID, _, _, operatorName, err := h.getOperatorInfo(c)
	if err != nil {
		response.Error(c, 401, err.Error())
		return
	}

	data := map[string]interface{}{
		"memberId":      req.MemberID,
		"centerId":      req.CenterID,
		"amount":        req.Amount,
		"paymentMethod": req.PaymentMethod,
		"remark":        req.Remark,
		"memberName":    req.MemberName,
		"memberPhone":   req.MemberPhone,
		"centerName":    req.CenterName,
		"operatorId":    fmt.Sprintf("%d", operatorID),
		"operatorName":  operatorName,
	}

	recharge, err := h.rechargeService.CreateCRecharge(data)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.SuccessWithMessage(c, errmsg.Get("recharge.c_create_success"), gin.H{
		"id":            recharge.ID,
		"transactionNo": recharge.ID,
		"balanceBefore": recharge.BalanceBefore,
		"balanceAfter":  recharge.BalanceAfter,
	})
}

func (h *RechargeHandler) GetCRechargeList(c *gin.Context) {
	memberPhone := c.Query("memberPhone")
	centerID := c.Query("centerId")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	result, err := h.rechargeService.GetCRechargeList(memberPhone, centerID, startDate, endDate, page, pageSize)
	if err != nil {
		response.InternalError(c, errmsg.Get("recharge.c_list_failed"))
		return
	}

	response.Success(c, result)
}

func (h *RechargeHandler) GetCRechargeDetail(c *gin.Context) {
	id := c.Param("id")

	result, err := h.rechargeService.GetCRechargeDetail(id)
	if err != nil {
		response.NotFound(c, errmsg.Get("recharge.c_not_found"))
		return
	}

	response.Success(c, result)
}

// ========== 门店卡 ==========

func (h *RechargeHandler) BatchImportCards(c *gin.Context) {
	// 限制文件大小 10MB
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 10<<20)

	file, err := c.FormFile("file")
	if err != nil {
		response.ParamsError(c, "请上传Excel文件（最大10MB）")
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".xlsx" && ext != ".xls" && ext != ".csv" {
		response.ParamsError(c, "仅支持.xlsx或.csv格式文件")
		return
	}

	// MIME type 校验
	contentType := file.Header.Get("Content-Type")
	allowedMIME := map[string]bool{
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": true, // .xlsx
		"application/vnd.ms-excel": true,                                          // .xls
		"text/csv":                  true,
		"application/octet-stream":  true, // 部分浏览器上传 .xlsx 时的 MIME
	}
	if !allowedMIME[contentType] {
		response.ParamsError(c, "文件类型不支持")
		return
	}

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

	userID, _, _, _, err := h.getOperatorInfo(c)
	if err != nil {
		response.Error(c, 401, err.Error())
		return
	}
	operatorID := fmt.Sprintf("%d", userID)

	count, cardNos, err := h.rechargeService.BatchImportCards(fileBytes, ext, operatorID)
	if err != nil {
		bizError(c, err)
		return
	}

	response.SuccessWithMessage(c, errmsg.Get("card.issue_success"), gin.H{
		"count":   count,
		"cardNos": cardNos,
	})
}

func (h *RechargeHandler) AllocateCards(c *gin.Context) {
	var req struct {
		CenterID string `json:"centerId"`
		Quantity  int    `json:"quantity"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamsError(c, err.Error())
		return
	}

	if req.Quantity <= 0 {
		response.ParamsError(c, "划拨数量必须大于0")
		return
	}

	count, err := h.rechargeService.AllocateCards(req.CenterID, req.Quantity)
	if err != nil {
		bizError(c, err)
		return
	}

	response.SuccessWithMessage(c, errmsg.Get("card.issue_success"), gin.H{
		"count": count,
	})
}

func (h *RechargeHandler) BindCardToUser(c *gin.Context) {
	var req struct {
		CardNo           string `json:"cardNo"`
		UserPhone        string `json:"userPhone"`
		IssueReason      string `json:"issueReason"`
		IssueType        int    `json:"issueType"`
		RechargeCenterID string `json:"rechargeCenterId"`
		RelatedUserPhone string `json:"relatedUserPhone"`
		Remark           string `json:"remark"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamsError(c, err.Error())
		return
	}

	userID, _, centerID, _, err := h.getOperatorInfo(c)
	if err != nil {
		response.Error(c, 401, err.Error())
		return
	}
	operatorID := fmt.Sprintf("%d", userID)

	// center permission check: center_admin/operator can only bind cards from their own center
	if centerID != "" && req.RechargeCenterID != centerID {
		response.Error(c, 403, "无权操作其他充值中心的卡片")
		return
	}

	if err := h.rechargeService.BindCardToUser(req.CardNo, req.UserPhone, req.IssueReason, req.IssueType, req.RechargeCenterID, operatorID, req.RelatedUserPhone, req.Remark); err != nil {
		bizError(c, err)
		return
	}

	response.SuccessWithMessage(c, errmsg.Get("card.issue_success"), gin.H{"success": true})
}

func (h *RechargeHandler) VerifyCard(c *gin.Context) {
	cardNo := c.Param("cardNo")

	result, err := h.rechargeService.VerifyCard(cardNo)
	if err != nil {
		bizError(c, err)
		return
	}

	response.Success(c, result)
}

func (h *RechargeHandler) ConsumeCard(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamsError(c, err.Error())
		return
	}

	cardNo, ok := req["cardNo"].(string)
	if !ok {
		response.ParamsError(c, errmsg.Get("card.consume_no_card"))
		return
	}
	amountFloat, ok := req["amount"].(float64)
	if !ok {
		response.ParamsError(c, errmsg.Get("card.consume_no_amount"))
		return
	}
	amount := int(amountFloat)
	remark := ""
	if req["remark"] != nil {
		remark, _ = req["remark"].(string)
	}

	userID, _, _, _, err := h.getOperatorInfo(c)
	if err != nil {
		response.Error(c, 401, err.Error())
		return
	}
	operatorID := fmt.Sprintf("%d", userID)

	if err := h.rechargeService.ConsumeCard(cardNo, amount, operatorID, remark); err != nil {
		bizError(c, err)
		return
	}

	response.SuccessWithMessage(c, errmsg.Get("card.consume_success"), gin.H{"success": true})
}

func (h *RechargeHandler) FreezeCard(c *gin.Context) {
	cardNo := c.Param("cardNo")

	userID, _, _, _, err := h.getOperatorInfo(c)
	if err != nil {
		response.Error(c, 401, err.Error())
		return
	}
	operatorID := fmt.Sprintf("%d", userID)

	if err := h.rechargeService.FreezeCard(cardNo, operatorID); err != nil {
		bizError(c, err)
		return
	}

	response.Success(c, gin.H{"success": true})
}

func (h *RechargeHandler) UnfreezeCard(c *gin.Context) {
	cardNo := c.Param("cardNo")

	userID, _, _, _, err := h.getOperatorInfo(c)
	if err != nil {
		response.Error(c, 401, err.Error())
		return
	}
	operatorID := fmt.Sprintf("%d", userID)

	if err := h.rechargeService.UnfreezeCard(cardNo, operatorID); err != nil {
		bizError(c, err)
		return
	}

	response.Success(c, gin.H{"success": true})
}

func (h *RechargeHandler) VoidCard(c *gin.Context) {
	cardNo := c.Param("cardNo")

	userID, _, _, _, err := h.getOperatorInfo(c)
	if err != nil {
		response.Error(c, 401, err.Error())
		return
	}
	operatorID := fmt.Sprintf("%d", userID)

	if err := h.rechargeService.VoidCard(cardNo, operatorID); err != nil {
		bizError(c, err)
		return
	}

	response.Success(c, gin.H{"success": true})
}

func (h *RechargeHandler) GetAvailableCards(c *gin.Context) {
	centerID := c.Query("centerId")
	keyword := c.Query("keyword")

	_, _, opCenterID, _, err := h.getOperatorInfo(c)
	if err != nil {
		response.Error(c, 401, err.Error())
		return
	}

	// center_admin/operator: force use own center_id
	if opCenterID != "" {
		centerID = opCenterID
	}

	cardNos, err := h.rechargeService.GetAvailableCards(centerID, keyword)
	if err != nil {
		bizError(c, err)
		return
	}
	response.Success(c, gin.H{"cardNos": cardNos})
}

func (h *RechargeHandler) GetAvailableCardCount(c *gin.Context) {
	centerID := c.Query("centerId")
	if centerID == "" {
		response.Error(c, 400, "centerId is required")
		return
	}

	_, _, opCenterID, _, err := h.getOperatorInfo(c)
	if err != nil {
		response.Error(c, 401, err.Error())
		return
	}

	if opCenterID != "" {
		centerID = opCenterID
	}

	count, err := h.rechargeService.GetAvailableCardCount(centerID)
	if err != nil {
		bizError(c, err)
		return
	}
	response.Success(c, gin.H{"count": count})
}

func (h *RechargeHandler) GetCardList(c *gin.Context) {
	status, _ := strconv.Atoi(c.Query("status"))
	cardNo := c.Query("cardNo")
	centerID := c.Query("centerId")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	_, _, opCenterID, _, err := h.getOperatorInfo(c)
	if err != nil {
		response.Error(c, 401, err.Error())
		return
	}
	if opCenterID != "" {
		centerID = opCenterID
	}

	result, err := h.rechargeService.GetCardList(status, cardNo, centerID, page, pageSize)
	if err != nil {
		response.InternalError(c, errmsg.Get("card.list_failed"))
		return
	}

	response.Success(c, result)
}

func (h *RechargeHandler) GetCardDetail(c *gin.Context) {
	cardNo := c.Param("cardNo")

	result, err := h.rechargeService.GetCardDetail(cardNo)
	if err != nil {
		response.NotFound(c, errmsg.Get("card.detail_not_found"))
		return
	}

	response.Success(c, result)
}

func (h *RechargeHandler) GetCardStats(c *gin.Context) {
	_, role, opCenterID, _, err := h.getOperatorInfo(c)
	if err != nil {
		response.Error(c, 401, err.Error())
		return
	}

	centerID := c.Query("centerId")
	// 中心角色强制使用自己的中心
	if opCenterID != "" {
		centerID = opCenterID
	} else if centerID == "" && (role == model.RoleSuperAdmin || role == model.RoleHQAdmin || role == model.RoleFinance) {
		// 总部角色不传 centerId 则看全部
		_ = role
	}

	result, err := h.rechargeService.GetCardStats(centerID)
	if err != nil {
		response.InternalError(c, errmsg.Get("card.stats_failed"))
		return
	}
	response.Success(c, result)
}

func (h *RechargeHandler) GetCardInventoryStats(c *gin.Context) {
	result, err := h.rechargeService.GetCardInventoryStats()
	if err != nil {
		response.InternalError(c, errmsg.Get("card.stats_failed"))
		return
	}
	response.Success(c, result)
}

func (h *RechargeHandler) GetMonthlyTrend(c *gin.Context) {
	_, _, opCenterID, _, err := h.getOperatorInfo(c)
	if err != nil {
		response.Error(c, 401, err.Error())
		return
	}
	centerID := c.Query("centerId")
	if opCenterID != "" {
		centerID = opCenterID
	}
	result, err := h.rechargeService.GetMonthlyTrend(centerID)
	if err != nil {
		response.InternalError(c, errmsg.Get("card.stats_failed"))
		return
	}
	response.Success(c, result)
}

func (h *RechargeHandler) GetCenterCardStats(c *gin.Context) {
	_, _, opCenterID, _, err := h.getOperatorInfo(c)
	if err != nil {
		response.Error(c, 401, err.Error())
		return
	}
	centerID := ""
	if opCenterID != "" {
		centerID = opCenterID
	}
	result, err := h.rechargeService.GetCenterCardStats(centerID)
	if err != nil {
		response.InternalError(c, errmsg.Get("card.stats_failed"))
		return
	}
	response.Success(c, result)
}

// ========== 充值中心 ==========

func (h *RechargeHandler) GetCenters(c *gin.Context) {
	result, err := h.rechargeService.GetCenters()
	if err != nil {
		response.InternalError(c, errmsg.Get("center.list_failed"))
		return
	}
	response.Success(c, result)
}

func (h *RechargeHandler) GetCenterDetail(c *gin.Context) {
	id := c.Param("id")
	result, err := h.rechargeService.GetCenterDetail(id)
	if err != nil {
		response.NotFound(c, errmsg.Get("center.list_failed"))
		return
	}
	response.Success(c, result)
}

func (h *RechargeHandler) CreateCenter(c *gin.Context) {
	var req model.CreateCenterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamsError(c, err.Error())
		return
	}

	data := map[string]interface{}{
		"name":      req.Name,
		"code":      req.Code,
		"address":   req.Address,
		"phone":     req.Phone,
		"province":  req.Province,
		"city":      req.City,
		"district":  req.District,
		"managerId": req.ManagerID,
	}

	result, err := h.rechargeService.CreateCenter(data)
	if err != nil {
		response.InternalError(c, errmsg.Get("center.create_failed"))
		return
	}

	response.SuccessWithMessage(c, errmsg.Get("center.create_success"), result)
}

func (h *RechargeHandler) UpdateCenter(c *gin.Context) {
	id := c.Param("id")
	var req model.UpdateCenterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamsError(c, err.Error())
		return
	}

	data := map[string]interface{}{}
	if req.Name != "" {
		data["name"] = req.Name
	}
	if req.Address != "" {
		data["address"] = req.Address
	}
	if req.Phone != "" {
		data["phone"] = req.Phone
	}

	result, err := h.rechargeService.UpdateCenter(id, data)
	if err != nil {
		response.InternalError(c, errmsg.Get("center.update_failed"))
		return
	}

	response.SuccessWithMessage(c, errmsg.Get("center.update_success"), result)
}

func (h *RechargeHandler) DeleteCenter(c *gin.Context) {
	id := c.Param("id")

	if err := h.rechargeService.DeleteCenter(id); err != nil {
		response.InternalError(c, errmsg.Get("center.delete_failed"))
		return
	}

	response.SuccessWithMessage(c, errmsg.Get("center.delete_success"), gin.H{"success": true})
}

// ========== 操作员 ==========

func (h *RechargeHandler) GetOperators(c *gin.Context) {
	result, err := h.rechargeService.GetOperators()
	if err != nil {
		response.InternalError(c, errmsg.Get("operator.list_failed"))
		return
	}
	response.Success(c, result)
}

func (h *RechargeHandler) CreateOperator(c *gin.Context) {
	var req model.CreateOperatorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamsError(c, err.Error())
		return
	}

	data := map[string]interface{}{
		"name":     req.Name,
		"phone":    req.Phone,
		"password": req.Password,
		"centerId": req.CenterID,
		"role":     req.Role,
	}

	result, err := h.rechargeService.CreateOperator(data)
	if err != nil {
		response.InternalError(c, errmsg.Get("operator.create_failed"))
		return
	}

	response.SuccessWithMessage(c, errmsg.Get("operator.create_success"), result)
}

func (h *RechargeHandler) UpdateOperator(c *gin.Context) {
	id := c.Param("id")
	var req model.UpdateOperatorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamsError(c, err.Error())
		return
	}

	data := map[string]interface{}{}
	if req.Name != "" {
		data["name"] = req.Name
	}
	if req.Phone != "" {
		data["phone"] = req.Phone
	}
	if req.Password != "" {
		data["password"] = req.Password
	}
	if req.Role != "" {
		data["role"] = req.Role
	}
	if req.Status != "" {
		data["status"] = req.Status
	}
	if req.CenterID != "" {
		data["centerId"] = req.CenterID
	}

	result, err := h.rechargeService.UpdateOperator(id, data)
	if err != nil {
		response.InternalError(c, errmsg.Get("operator.update_failed"))
		return
	}

	response.SuccessWithMessage(c, errmsg.Get("operator.update_success"), result)
}

func (h *RechargeHandler) DeleteOperator(c *gin.Context) {
	id := c.Param("id")

	if err := h.rechargeService.DeleteOperator(id); err != nil {
		response.InternalError(c, errmsg.Get("operator.delete_failed"))
		return
	}

	response.SuccessWithMessage(c, errmsg.Get("operator.delete_success"), gin.H{"success": true})
}

// ========== Dashboard ==========

func (h *RechargeHandler) GetDashboardStatistics(c *gin.Context) {
	_, role, centerID, _, err := h.getOperatorInfo(c)
	if err != nil {
		response.Error(c, 401, err.Error())
		return
	}

	// center_admin/operator 限定本中心
	filterCenterID := ""
	if role != model.RoleSuperAdmin && role != model.RoleHQAdmin && role != model.RoleFinance {
		filterCenterID = centerID
	}

	data, err := h.rechargeService.GetDashboardStatistics(role, filterCenterID)
	if err != nil {
		response.InternalError(c, "统计数据加载失败")
		return
	}
	response.Success(c, data)
}

func (h *RechargeHandler) GetDashboardTodos(c *gin.Context) {
	_, role, centerID, _, err := h.getOperatorInfo(c)
	if err != nil {
		response.Error(c, 401, err.Error())
		return
	}

	filterCenterID := ""
	if role != model.RoleSuperAdmin && role != model.RoleHQAdmin && role != model.RoleFinance {
		filterCenterID = centerID
	}

	data, err := h.rechargeService.GetDashboardTodos(role, filterCenterID)
	if err != nil {
		response.InternalError(c, "待办事项加载失败")
		return
	}
	response.Success(c, data)
}

func (h *RechargeHandler) GetDashboardRechargeTrends(c *gin.Context) {
	_, role, centerID, _, err := h.getOperatorInfo(c)
	if err != nil {
		response.Error(c, 401, err.Error())
		return
	}

	filterCenterID := ""
	if role != model.RoleSuperAdmin && role != model.RoleHQAdmin && role != model.RoleFinance {
		filterCenterID = centerID
	}

	days := 7
	if d := c.DefaultQuery("days", "7"); d != "" {
		if parsed, e := strconv.Atoi(d); e == nil && parsed > 0 {
			days = parsed
		}
	}

	data, err := h.rechargeService.GetDashboardRechargeTrends(days, role, filterCenterID)
	if err != nil {
		response.InternalError(c, "趋势数据加载失败")
		return
	}
	response.Success(c, data)
}
