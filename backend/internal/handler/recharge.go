package handler

import (
	"io"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"marketplace/backend/internal/service"
	"marketplace/backend/pkg/errmsg"
	"marketplace/backend/pkg/response"
)

type RechargeHandler struct {
	rechargeService service.RechargeServiceInterface
}

func NewRechargeHandler(rechargeService service.RechargeServiceInterface) *RechargeHandler {
	return &RechargeHandler{rechargeService: rechargeService}
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

	// TODO: 从JWT获取操作员信息
	req["operatorId"] = "op123"
	req["operatorName"] = "张出纳"

	recharge, err := h.rechargeService.CreateCRecharge(req)
	if err != nil {
		response.InternalError(c, errmsg.Get("recharge.c_create_failed"))
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

	// TODO: 从JWT获取操作员信息
	operatorID := "op123"

	count, cardNos, err := h.rechargeService.BatchImportCards(fileBytes, ext, operatorID)
	if err != nil {
		response.InternalError(c, errmsg.Get("card.issue_failed")+":"+err.Error())
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
		response.InternalError(c, errmsg.Get("card.issue_failed")+":"+err.Error())
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

	// TODO: 从JWT获取操作员信息
	operatorID := "op123"

	if err := h.rechargeService.BindCardToUser(req.CardNo, req.UserPhone, req.IssueReason, req.IssueType, req.RechargeCenterID, operatorID, req.RelatedUserPhone, req.Remark); err != nil {
		response.InternalError(c, errmsg.Get("card.issue_failed")+":"+err.Error())
		return
	}

	response.SuccessWithMessage(c, errmsg.Get("card.issue_success"), gin.H{"success": true})
}

func (h *RechargeHandler) VerifyCard(c *gin.Context) {
	cardNo := c.Param("cardNo")

	result, err := h.rechargeService.VerifyCard(cardNo)
	if err != nil {
		response.NotFound(c, errmsg.Get("card.verify_not_found")+":"+err.Error())
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

	// TODO: 从JWT获取操作员ID
	operatorID := "op123"

	if err := h.rechargeService.ConsumeCard(cardNo, amount, operatorID, remark); err != nil {
		response.InternalError(c, errmsg.Get("card.consume_failed")+":"+err.Error())
		return
	}

	response.SuccessWithMessage(c, errmsg.Get("card.consume_success"), gin.H{"success": true})
}

func (h *RechargeHandler) FreezeCard(c *gin.Context) {
	cardNo := c.Param("cardNo")

	// TODO: 从JWT获取操作员ID
	operatorID := "op123"

	if err := h.rechargeService.FreezeCard(cardNo, operatorID); err != nil {
		response.InternalError(c, errmsg.Get("card.status_failed")+":"+err.Error())
		return
	}

	response.Success(c, gin.H{"success": true})
}

func (h *RechargeHandler) UnfreezeCard(c *gin.Context) {
	cardNo := c.Param("cardNo")

	// TODO: 从JWT获取操作员ID
	operatorID := "op123"

	if err := h.rechargeService.UnfreezeCard(cardNo, operatorID); err != nil {
		response.InternalError(c, errmsg.Get("card.status_failed")+":"+err.Error())
		return
	}

	response.Success(c, gin.H{"success": true})
}

func (h *RechargeHandler) VoidCard(c *gin.Context) {
	cardNo := c.Param("cardNo")

	// TODO: 从JWT获取操作员ID
	operatorID := "op123"

	if err := h.rechargeService.VoidCard(cardNo, operatorID); err != nil {
		response.InternalError(c, errmsg.Get("card.status_failed")+":"+err.Error())
		return
	}

	response.Success(c, gin.H{"success": true})
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
