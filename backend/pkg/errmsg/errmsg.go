package errmsg

// Msg 错误提示语言包
// key 格式: 模块.操作.错误类型
var Msg = map[string]string{
	// ========== 通用 ==========
	"common.param_error":      "请求参数有误",
	"common.internal_error":   "服务器内部错误",
	"common.not_found":        "请求的资源不存在",
	"common.forbidden":        "无权限访问",
	"common.unauthorized":     "请先登录",
	"common.conflict":         "数据冲突，请刷新后重试",

	// ========== 认证 auth ==========
	"auth.login_failed":       "手机号或密码错误",
	"auth.user_disabled":      "用户已被禁用",
	"auth.token_invalid":      "Token 无效或已过期",
	"auth.token_missing":      "缺少认证信息",
	"auth.token_format_error": "Token 格式错误",
	"auth.logout_failed":      "登出失败",
	"auth.register_failed":    "注册失败",

	// ========== 用户 user ==========
	"user.not_found":          "用户不存在",
	"user.update_failed":      "更新失败",
	"user.list_failed":        "获取列表失败",
	"user.password_encrypt":   "密码加密失败",
	"user.password_change":    "密码修改失败",
	"user.id_format_error":    "用户ID格式错误",

	// ========== 管理员 admin ==========
	"admin.user_list_failed":    "获取用户列表失败",
	"admin.user_not_found":      "用户不存在",
	"admin.reset_password":      "重置密码失败",
	"admin.update_status":       "更新用户状态失败",
	"admin.get_user_failed":     "获取用户信息失败",
	"admin.create_user_failed":  "创建用户失败",
	"admin.update_user_failed":  "更新用户失败",
	"admin.delete_user_failed":  "删除用户失败",

	// ========== 充值申请 B端 ==========
	"recharge.apply_failed":     "提交充值申请失败",
	"recharge.apply_success":    "申请提交成功",
	"recharge.list_failed":      "获取充值申请列表失败",
	"recharge.not_found":        "充值申请不存在",
	"recharge.approval_failed":  "审批操作失败",

	// ========== 充值记录 C端 ==========
	"recharge.c_create_failed":  "提交充值记录失败",
	"recharge.c_create_success": "充值成功",
	"recharge.c_list_failed":    "获取充值记录列表失败",
	"recharge.c_not_found":      "充值记录不存在",

	// ========== 门店卡 ==========
	"card.issue_failed":         "发放门店卡失败",
	"card.issue_success":        "发放成功",
	"card.verify_not_found":     "门店卡不存在或状态异常",
	"card.consume_failed":       "核销失败",
	"card.consume_success":      "核销成功",
	"card.consume_no_card":      "缺少卡号参数",
	"card.consume_no_amount":    "缺少金额参数",
	"card.status_failed":        "更新卡状态失败",
	"card.status_no_param":      "缺少状态参数",
	"card.list_failed":          "获取门店卡列表失败",
	"card.detail_not_found":     "门店卡不存在",
	"card.stats_failed":         "获取门店卡统计失败",

	// ========== 充值中心 ==========
	"center.list_failed":        "获取充值中心列表失败",
	"center.create_failed":      "创建充值中心失败",
	"center.create_success":     "创建成功",
	"center.update_failed":      "更新充值中心失败",
	"center.update_success":     "更新成功",
	"center.delete_failed":      "删除充值中心失败",
	"center.delete_success":     "删除成功",

	// ========== 操作员 ==========
	"operator.list_failed":      "获取操作员列表失败",
	"operator.create_failed":    "创建操作员失败",
	"operator.create_success":   "创建成功",
	"operator.update_failed":    "更新操作员失败",
	"operator.update_success":   "更新成功",
	"operator.delete_failed":    "删除操作员失败",
	"operator.delete_success":   "删除成功",

	// ========== RBAC ==========
	"rbac.check_failed":         "权限检查失败",
	"rbac.no_role":              "无有效角色",
	"rbac.no_permission":        "无权限访问",
	"rbac.role_mismatch":        "角色不匹配",
}

// Get 获取错误提示，不存在时返回兜底文案
func Get(key string) string {
	if msg, ok := Msg[key]; ok {
		return msg
	}
	return "操作失败"
}
