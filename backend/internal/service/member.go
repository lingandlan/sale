package service

import (
	"fmt"
)

// MemberInfo 会员信息（来自WSY商城服务）
type MemberInfo struct {
	UserID   string  `json:"userId"`
	Name     string  `json:"name"`
	Phone    string  `json:"phone"`
	Balance  float64 `json:"balance"`
	Level    string  `json:"level"`
	NickName string  `json:"nickName"`
}

// MemberServiceInterface 会员服务接口
type MemberServiceInterface interface {
	SearchByPhone(phone string) (*MemberInfo, error)
}

var _ MemberServiceInterface = (*MemberService)(nil)

// MemberService 会员服务
type MemberService struct {
	// TODO: 后续注入 WSY HTTP client
}

func NewMemberService() *MemberService {
	return &MemberService{}
}

// SearchByPhone 通过手机号查询会员信息
// 调用链路：先调 WSY 10000_phone_get_user_info 获取 user_id，
// 再调 WSY 10000_integral_user_integral 获取积分余额
func (s *MemberService) SearchByPhone(phone string) (*MemberInfo, error) {
	if phone == "" {
		return nil, fmt.Errorf("手机号不能为空")
	}

	// TODO: 对接 WSY 服务
	// 步骤1: POST wsy_pub/third_api/index.php?m=third_application_authorization&a=third_function
	//        act=10000_phone_get_user_info, phone=xxx → 获取 user_id
	// 步骤2: act=10000_integral_user_integral, user_id=xxx → 获取积分余额

	// Mock 数据 — 后续替换为真实 WSY 调用
	return &MemberInfo{
		UserID:   "370955981",
		Name:     "张三",
		Phone:    phone,
		Balance:  5000,
		Level:    "普通会员",
		NickName: "张三",
	}, nil
}
