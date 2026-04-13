package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

// ListResponse 分页列表响应
type ListResponse struct {
	Items    interface{} `json:"items"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

// 业务错误码
const (
	CodeSuccess       = 0
	CodeParamsError   = 400
	CodeUnauthorized  = 401
	CodeForbidden     = 403
	CodeNotFound      = 404
	CodeConflict      = 409
	CodeInternalError = 500
)

// errorMessages 错误码 → 默认中文提示
var errorMessages = map[int]string{
	CodeParamsError:   "请求参数有误",
	CodeUnauthorized:  "请先登录",
	CodeForbidden:     "无权限访问",
	CodeNotFound:      "请求的资源不存在",
	CodeConflict:      "数据冲突，请刷新后重试",
	CodeInternalError: "服务器内部错误",
}

// ErrorMessage 根据错误码返回默认中文提示
func ErrorMessage(code int) string {
	if msg, ok := errorMessages[code]; ok {
		return msg
	}
	return "未知错误"
}

// Success 返回成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: "success",
		Data:    data,
	})
}

// SuccessWithMessage 返回成功响应（自定义消息）
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: message,
		Data:    data,
	})
}

// Created 返回创建成功响应
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Code:    CodeSuccess,
		Message: "created",
		Data:    data,
	})
}

// NoContent 返回无内容响应
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// Error 返回错误响应
func Error(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, ErrorResponse{
		Code:    code,
		Message: message,
	})
}

// ErrorWithDetail 返回错误响应（带详情）
func ErrorWithDetail(c *gin.Context, code int, message string, detail string) {
	c.JSON(http.StatusOK, ErrorResponse{
		Code:    code,
		Message: message,
		Detail:  detail,
	})
}

// ParamsError 返回参数错误
func ParamsError(c *gin.Context, detail string) {
	ErrorWithDetail(c, CodeParamsError, "参数错误", detail)
}

// Unauthorized 返回未认证
func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, ErrorResponse{
		Code:    CodeUnauthorized,
		Message: message,
	})
}

// Forbidden 返回无权限
func Forbidden(c *gin.Context, message string) {
	Error(c, CodeForbidden, message)
}

// NotFound 返回资源不存在
func NotFound(c *gin.Context, message string) {
	Error(c, CodeNotFound, message)
}

// Conflict 返回资源冲突
func Conflict(c *gin.Context, message string) {
	Error(c, CodeConflict, message)
}

// InternalError 返回内部错误
func InternalError(c *gin.Context, message string) {
	Error(c, CodeInternalError, message)
}
