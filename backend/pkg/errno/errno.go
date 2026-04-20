package errno

import "fmt"

// Code 业务错误码，前端可据此做逻辑分支
type Code string

const (
	// 通用
	CodeInvalidParam Code = "INVALID_PARAM"

	// 门店卡
	CodeCardNotFound      Code = "CARD_NOT_FOUND"
	CodeCardFrozen        Code = "CARD_FROZEN"
	CodeCardExpired       Code = "CARD_EXPIRED"
	CodeCardVoided        Code = "CARD_VOIDED"
	CodeCardNotIssued     Code = "CARD_NOT_ISSUED"
	CodeCardNotInStock    Code = "CARD_NOT_IN_STOCK"
	CodeCardNotInCenter   Code = "CARD_NOT_IN_CENTER"
	CodeInsufficientStock Code = "INSUFFICIENT_STOCK"
	CodeInvalidCardType   Code = "INVALID_CARD_TYPE"
	CodeInvalidAction     Code = "INVALID_ACTION"
	CodeInvalidTransition Code = "INVALID_TRANSITION"
	CodeAlreadyProcessed  Code = "ALREADY_PROCESSED"
	CodeDuplicateCardNo   Code = "DUPLICATE_CARD_NO"

	// 充值中心
	CodeCenterNotFound  Code = "CENTER_NOT_FOUND"
	CodeCenterNoBalance Code = "CENTER_NO_BALANCE"

	// 充值申请
	CodeRechargeNotFound Code = "RECHARGE_NOT_FOUND"
)

// BizError 业务错误，携带 code + message
type BizError struct {
	Code    Code   `json:"code"`
	Message string `json:"message"`
}

func (e *BizError) Error() string { return e.Message }

// New 创建业务错误
func New(code Code, msg string) *BizError {
	return &BizError{Code: code, Message: msg}
}

// Newf 创建带格式化的业务错误
func Newf(code Code, format string, args ...interface{}) *BizError {
	return &BizError{Code: code, Message: fmt.Sprintf(format, args...)}
}

// Resolve 从 error 中提取 BizError，非 BizError 返回空 code
func Resolve(err error) (Code, string) {
	if biz, ok := err.(*BizError); ok {
		return biz.Code, biz.Message
	}
	return "", err.Error()
}
