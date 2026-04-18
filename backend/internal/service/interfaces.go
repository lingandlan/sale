package service

import (
	"context"

	"marketplace/backend/internal/model"
)

type UserServiceInterface interface {
	GetByID(ctx context.Context, id int64) (*model.User, error)
	GetByPhone(ctx context.Context, phone string) (*model.User, error)
	Create(ctx context.Context, req *model.CreateUserRequest) (*model.User, error)
	Update(ctx context.Context, id int64, req *model.UpdateUserRequest) (*model.User, error)
	List(ctx context.Context, page, pageSize int) ([]*model.User, int64, error)
	ListWithFilters(ctx context.Context, req *model.ListUsersRequest) (*model.ListUsersResponse, error)
	ResetPassword(ctx context.Context, id int64, req *model.ResetPasswordRequest) error
	UpdateStatus(ctx context.Context, id int64, req *model.UpdateUserStatusRequest) error
	Delete(ctx context.Context, id int64) error
	UpdatePassword(ctx context.Context, id int64, hashedPassword string) error
}

var _ UserServiceInterface = (*UserService)(nil)

type AuthServiceInterface interface {
	Login(ctx context.Context, phone, password string) (*model.LoginResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (string, string, error)
	Logout(ctx context.Context, userID int64) error
	GenerateToken(user *model.User) (string, error)
	GenerateRefreshToken(user *model.User) (string, error)
	ParseToken(tokenString string) (*Claims, error)
}

var _ AuthServiceInterface = (*AuthService)(nil)

// RechargeServiceInterface 充值服务接口
type RechargeServiceInterface interface {
	CalculatePoints(amount float64, lastMonthConsumption float64) (int, int, int)
	CreateBRechargeApplication(data map[string]interface{}) (*model.RechargeApplication, error)
	GetRechargeApplicationList(status string, page, pageSize int) (map[string]interface{}, error)
	GetRechargeApplicationDetail(id string) (*model.RechargeApplication, error)
	ApproveRechargeApplication(id, action, approvedBy, remark string) error
	CreateCRecharge(data map[string]interface{}) (*model.CRecharge, error)
	GetCRechargeList(memberPhone, centerID, startDate, endDate string, page, pageSize int) (map[string]interface{}, error)
	GetCRechargeDetail(id string) (*model.CRecharge, error)
	// 门店卡
	BatchImportCards(file []byte, ext string, operatorID string) (int, []string, error)
	AllocateCards(centerID string, quantity int) (int, error)
	BindCardToUser(cardNo, userPhone, issueReason string, issueType int, rechargeCenterID, operatorID, relatedUserPhone, remark string) error
	VerifyCard(cardNo string) (*model.StoreCard, error)
	ConsumeCard(cardNo string, amount int, operatorID, remark string) error
	GetCardList(status int, cardNo, centerID string, page, pageSize int) (map[string]interface{}, error)
	GetCardDetail(cardNo string) (map[string]interface{}, error)
	GetCardStats() (map[string]interface{}, error)
	GetCardInventoryStats() (map[string]interface{}, error)
	GetAvailableCards(centerID string, keyword string) ([]string, error)
	FreezeCard(cardNo, operatorID string) error
	UnfreezeCard(cardNo, operatorID string) error
	VoidCard(cardNo, operatorID string) error
	// 充值中心
	GetCenters() ([]map[string]interface{}, error)
	GetCenterDetail(id string) (*model.RechargeCenter, error)
	CreateCenter(data map[string]interface{}) (*model.RechargeCenter, error)
	UpdateCenter(id string, data map[string]interface{}) (*model.RechargeCenter, error)
	DeleteCenter(id string) error
	// 操作员
	GetOperators() ([]model.RechargeOperator, error)
	CreateOperator(data map[string]interface{}) (*model.RechargeOperator, error)
	UpdateOperator(id string, data map[string]interface{}) (*model.RechargeOperator, error)
	DeleteOperator(id string) error
		// Dashboard
		GetDashboardStatistics(role, centerID string) (map[string]interface{}, error)
		GetDashboardTodos(role, centerID string) (map[string]interface{}, error)
		GetDashboardRechargeTrends(days int, role, centerID string) (map[string]interface{}, error)
}

var _ RechargeServiceInterface = (*RechargeService)(nil)
