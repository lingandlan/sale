package service

import (
	"fmt"

	"marketplace/backend/pkg/mall"
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
	wsyClient *mall.WSYClient
}

func NewMemberService(wsyClient *mall.WSYClient) *MemberService {
	return &MemberService{wsyClient: wsyClient}
}

// SearchByPhone 通过手机号查询会员信息
// 调用链路：先调 WSY 10000_phone_get_user_info 获取 user_id，
// 再调 WSY 10000_integral_user_integral 获取积分余额
func (s *MemberService) SearchByPhone(phone string) (*MemberInfo, error) {
	if phone == "" {
		return nil, fmt.Errorf("手机号不能为空")
	}

	// 步骤1: 通过手机号获取 user_id
	userID, err := s.wsyClient.PhoneToUserID(phone)
	if err != nil {
		return nil, fmt.Errorf("查询会员失败: %w", err)
	}

	// 步骤2: 获取积分余额
	integral, err := s.wsyClient.GetUserIntegral(userID)
	if err != nil {
		return nil, fmt.Errorf("查询积分失败: %w", err)
	}

	return &MemberInfo{
		UserID:   userID,
		Name:     "",
		Phone:    phone,
		Balance:  integral,
		Level:    "普通会员",
		NickName: "",
	}, nil
}
