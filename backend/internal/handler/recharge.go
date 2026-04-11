package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"marketplace/backend/internal/service"
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
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "申请提交成功", gin.H{
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
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, result)
}

func (h *RechargeHandler) GetRechargeApplicationDetail(c *gin.Context) {
	id := c.Param("id")

	result, err := h.rechargeService.GetRechargeApplicationDetail(id)
	if err != nil {
		response.NotFound(c, err.Error())
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
		response.InternalError(c, err.Error())
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
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "充值成功", gin.H{
		"id":            recharge.ID,
		"transactionNo": recharge.ID,
	})
}

func (h *RechargeHandler) GetCRechargeList(c *gin.Context) {
	memberPhone := c.Query("memberPhone")
	centerID := c.Query("centerId")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	result, err := h.rechargeService.GetCRechargeList(memberPhone, centerID, page, pageSize)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, result)
}

func (h *RechargeHandler) GetCRechargeDetail(c *gin.Context) {
	id := c.Param("id")

	result, err := h.rechargeService.GetCRechargeDetail(id)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, result)
}

// ========== 门店卡 ==========

func (h *RechargeHandler) IssueCard(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamsError(c, err.Error())
		return
	}

	// TODO: 从JWT获取操作员信息
	req["operatorId"] = "op123"

	card, err := h.rechargeService.IssueCard(req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "发放成功", gin.H{
		"id":     card.ID,
		"cardNo": card.CardNo,
	})
}

func (h *RechargeHandler) VerifyCard(c *gin.Context) {
	cardNo := c.Param("cardNo")

	result, err := h.rechargeService.VerifyCard(cardNo)
	if err != nil {
		response.NotFound(c, err.Error())
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
		response.ParamsError(c, "cardNo is required")
		return
	}
	amount, ok := req["amount"].(float64)
	if !ok {
		response.ParamsError(c, "amount is required")
		return
	}
	remark := ""
	if req["remark"] != nil {
		remark, _ = req["remark"].(string)
	}

	// TODO: 从JWT获取操作员ID
	operatorID := "op123"

	if err := h.rechargeService.ConsumeCard(cardNo, amount, remark, operatorID); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "核销成功", gin.H{"success": true})
}

func (h *RechargeHandler) UpdateCardStatus(c *gin.Context) {
	cardNo := c.Param("cardNo")
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamsError(c, err.Error())
		return
	}

	status, ok := req["status"].(string)
	if !ok {
		response.ParamsError(c, "status is required")
		return
	}

	if err := h.rechargeService.UpdateCardStatus(cardNo, status); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"success": true})
}

func (h *RechargeHandler) GetCardList(c *gin.Context) {
	status := c.Query("status")
	cardNo := c.Query("cardNo")
	holderPhone := c.Query("holderPhone")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	result, err := h.rechargeService.GetCardList(status, cardNo, holderPhone, page, pageSize)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, result)
}

func (h *RechargeHandler) GetCardDetail(c *gin.Context) {
	cardNo := c.Param("cardNo")

	result, err := h.rechargeService.GetCardDetail(cardNo)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, result)
}

func (h *RechargeHandler) GetCardStats(c *gin.Context) {
	result, err := h.rechargeService.GetCardStats()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, result)
}

// ========== 充值中心 ==========

func (h *RechargeHandler) GetCenters(c *gin.Context) {
	result, err := h.rechargeService.GetCenters()
	if err != nil {
		response.InternalError(c, err.Error())
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
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "创建成功", result)
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
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "更新成功", result)
}

func (h *RechargeHandler) DeleteCenter(c *gin.Context) {
	id := c.Param("id")

	if err := h.rechargeService.DeleteCenter(id); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "删除成功", gin.H{"success": true})
}

// ========== 操作员 ==========

func (h *RechargeHandler) GetOperators(c *gin.Context) {
	result, err := h.rechargeService.GetOperators()
	if err != nil {
		response.InternalError(c, err.Error())
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
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "创建成功", result)
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
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "更新成功", result)
}

func (h *RechargeHandler) DeleteOperator(c *gin.Context) {
	id := c.Param("id")

	if err := h.rechargeService.DeleteOperator(id); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "删除成功", gin.H{"success": true})
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
