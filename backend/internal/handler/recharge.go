package handler

import (
	"fmt"
	"io"
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
	if user.CenterID == nil {
		return userID, role, "", userName, fmt.Errorf("当前用户未分配充值中心")
	}
	centerID := strconv.FormatUint(uint64(*user.CenterID), 10)
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
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamsError(c, err.Error())
		return
	}

	// TODO: 从JWT获取申请人信息
	req["applicantId"] = "user123"
	req["applicantName"] = "张财务"

	app, err := h.rechargeService.CreateBRechargeApplication(req)
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
	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamsError(c, err.Error())
		return
	}

	id := req["id"]
	action := req["action"]
	remark := req["remark"]

	// TODO: 从JWT获取审批人信息
	approvedBy := "admin"

	if err := h.rechargeService.ApproveRechargeApplication(id, action, approvedBy, remark); err != nil {
		response.InternalError(c, errmsg.Get("recharge.approval_failed"))
		return
	}

	response.Success(c, gin.H{"success": true})
}

// ========== C端充值 ==========

func (h *RechargeHandler) CreateCRecharge(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamsError(c, err.Error())
		return
	}

	operatorID, _, _, operatorName, err := h.getOperatorInfo(c)
	if err != nil {
		response.Error(c, 401, err.Error())
		return
	}
	req["operatorId"] = fmt.Sprintf("%d", operatorID)
	req["operatorName"] = operatorName

	recharge, err := h.rechargeService.CreateCRecharge(req)
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
	file, err := c.FormFile("file")
	if err != nil {
		response.ParamsError(c, "请上传Excel文件")
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".xlsx" && ext != ".xls" && ext != ".csv" {
		response.ParamsError(c, "仅支持.xlsx或.csv格式文件")
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

func (h *RechargeHandler) GetCardList(c *gin.Context) {
	status, _ := strconv.Atoi(c.Query("status"))
	cardNo := c.Query("cardNo")
	centerID := c.Query("centerId")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

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
	result, err := h.rechargeService.GetCardStats()
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
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamsError(c, err.Error())
		return
	}

	result, err := h.rechargeService.CreateCenter(req)
	if err != nil {
		response.InternalError(c, errmsg.Get("center.create_failed"))
		return
	}

	response.SuccessWithMessage(c, errmsg.Get("center.create_success"), result)
}

func (h *RechargeHandler) UpdateCenter(c *gin.Context) {
	id := c.Param("id")
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamsError(c, err.Error())
		return
	}

	result, err := h.rechargeService.UpdateCenter(id, req)
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
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamsError(c, err.Error())
		return
	}

	result, err := h.rechargeService.CreateOperator(req)
	if err != nil {
		response.InternalError(c, errmsg.Get("operator.create_failed"))
		return
	}

	response.SuccessWithMessage(c, errmsg.Get("operator.create_success"), result)
}

func (h *RechargeHandler) UpdateOperator(c *gin.Context) {
	id := c.Param("id")
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamsError(c, err.Error())
		return
	}

	result, err := h.rechargeService.UpdateOperator(id, req)
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
	// TODO: 从数据库获取实际统计数据
	response.Success(c, gin.H{
		"todayRecharge":    15000,
		"monthRecharge":    450000,
		"totalCards":       500,
		"pendingApprovals": 5,
	})
}

func (h *RechargeHandler) GetDashboardTodos(c *gin.Context) {
	// TODO: 从数据库获取实际待办事项
	response.Success(c, gin.H{
		"pendingApprovals": gin.H{
			"count":       3,
			"description": "3笔充值申请待审批",
		},
		"expiringCards": gin.H{
			"count":       5,
			"description": "5张门店卡将在7天内到期",
		},
	})
}

func (h *RechargeHandler) GetDashboardRechargeTrends(c *gin.Context) {
	// TODO: 从数据库获取实际趋势数据
	_ = c.DefaultQuery("days", "7")
	response.Success(c, gin.H{
		"dates":  []string{"04-05", "04-06", "04-07", "04-08", "04-09", "04-10", "04-11"},
		"values": []int{32000, 45000, 28000, 52000, 41000, 38000, 52800},
	})
}
