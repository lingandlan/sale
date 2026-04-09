package errors

import "errors"

// 预定义错误
var (
	// ErrNotFound 资源不存在
	ErrNotFound = errors.New("resource not found")

	// ErrUnauthorized 未认证
	ErrUnauthorized = errors.New("unauthorized")

	// ErrForbidden 无权限
	ErrForbidden = errors.New("forbidden")

	// ErrDuplicate 资源重复
	ErrDuplicate = errors.New("resource already exists")

	// ErrInvalidParams 参数错误
	ErrInvalidParams = errors.New("invalid parameters")

	// ErrTokenExpired Token 过期
	ErrTokenExpired = errors.New("token expired")

	// ErrTokenInvalid Token 无效
	ErrTokenInvalid = errors.New("invalid token")

	// ErrPasswordIncorrect 密码错误
	ErrPasswordIncorrect = errors.New("password incorrect")

	// ErrUserDisabled 用户已禁用
	ErrUserDisabled = errors.New("user is disabled")
)

// Is 判断错误类型
func Is(err, target error) bool {
	return errors.Is(err, target)
}
