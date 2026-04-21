package model

import (
	"time"
)

// RechargeApplication B端充值申请
type RechargeApplication struct {
	ID            string    `json:"id" gorm:"primaryKey;size:64"`
	CenterID      string    `json:"centerId" gorm:"index;size:64"`    // 充值中心ID
	CenterName    string    `json:"centerName"`                        // 充值中心名称
	Amount        float64   `json:"amount"`                            // 充值金额
	Points        int       `json:"points"`                            // 预计积分
	BasePoints    int       `json:"basePoints"`                        // 基础积分
	RebatePoints  int       `json:"rebatePoints"`                      // 返还积分
	RebateRate    int       `json:"rebateRate"`                        // 返还比例
	ApplicantID   string    `json:"applicantId" gorm:"index;size:64"`  // 申请人ID
	ApplicantName string    `json:"applicantName"`                     // 申请人姓名
	TransactionNo string    `json:"transactionNo"`                     // 银行流水单号
	Screenshot    string    `json:"screenshot"`                        // 付款截图
	Remark        string    `json:"remark"`                            // 备注
	Status        string    `json:"status" gorm:"default:'pending';size:32"` // pending/approved/rejected
	ApprovedBy    string    `json:"approvedBy"`                        // 审批人
	ApprovedAt    *time.Time `json:"approvedAt"`                       // 审批时间
	ApprovalRemark string   `json:"approvalRemark"`                   // 审批备注
	LastMonthConsumption float64 `json:"lastMonthConsumption"`

	CreatedAt     time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}

// CRecharge C端充值记录
type CRecharge struct {
	ID            string    `json:"id" gorm:"primaryKey;size:64"`
	MemberID      string    `json:"memberId" gorm:"index;size:64"`    // 会员ID
	MemberName    string    `json:"memberName"`                        // 会员姓名
	MemberPhone   string    `json:"memberPhone" gorm:"index;size:32"`  // 会员手机
	CenterID      string    `json:"centerId"`                          // 充值中心ID
	CenterName    string    `json:"centerName"`                        // 充值中心名称
	Amount        float64   `json:"amount"`                            // 充值金额
	Points        int       `json:"points"`                            // 获得积分
	PaymentMethod string    `json:"paymentMethod"`                     // 支付方式 cash/wechat/alipay/card
	OperatorID    string    `json:"operatorId"`                        // 操作员ID
	OperatorName  string    `json:"operatorName"`                      // 操作员姓名
	Remark        string    `json:"remark"`                            // 备注
	BalanceBefore int       `json:"balanceBefore"`                     // 充值前余额
	BalanceAfter  int       `json:"balanceAfter"`                      // 充值后余额
	Status        string    `json:"status" gorm:"default:'success';size:32"` // pending/success/failed
	CreatedAt     time.Time `json:"createdAt" gorm:"autoCreateTime"`
}

// 卡状态常量
const (
	CardStatusInStock = 1 // 已入库
	CardStatusIssued  = 2 // 已发放
	CardStatusActive  = 3 // 已激活
	CardStatusFrozen  = 4 // 已冻结
	CardStatusExpired = 5 // 已过期
	CardStatusVoided  = 6 // 已作废
)

// 卡类型
const (
	CardTypePhysical = 1 // 实体卡
	CardTypeVirtual  = 2 // 虚拟卡
)

// StoreCard 门店卡
type StoreCard struct {
	ID              string     `json:"id" gorm:"primaryKey;size:64"`
	CardNo          string     `json:"cardNo" gorm:"uniqueIndex;size:32;not null"`  // 卡号 TJ00000001
	CardType        int        `json:"cardType" gorm:"default:1"`                    // 1=实体卡,2=虚拟卡
	Status          int        `json:"status" gorm:"default:1"`                      // 1=已入库,2=已发放,3=已激活,4=已冻结,5=已过期,6=已作废
	Balance         int        `json:"balance" gorm:"default:1000"`                  // 余额（元），固定面值1000
	RechargeCenterID string    `json:"rechargeCenterId" gorm:"index;size:64"`        // 划拨到的充值中心ID
	UserID          string     `json:"userId" gorm:"index;size:64"`                  // 绑定的用户ID
	BatchNo         string     `json:"batchNo" gorm:"size:64"`                       // 批次号
	IssueReason     string     `json:"issueReason" gorm:"size:64"`                   // 发放原因:购买套餐包/推荐奖励/其他
	IssuedAt        *time.Time `json:"issuedAt"`                                      // 发放时间
	ActivatedAt     *time.Time `json:"activatedAt"`                                   // 激活时间（首次核销）
	ExpiredAt       *time.Time `json:"expiredAt"`                                     // 过期时间（激活日+1年）
	CreatedAt       time.Time  `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt       time.Time  `json:"updatedAt" gorm:"autoUpdateTime"`
}

// CardIssueRecord 门店卡发放记录
type CardIssueRecord struct {
	ID               string    `json:"id" gorm:"primaryKey;size:64"`
	CardNo           string    `json:"cardNo" gorm:"index;size:32;not null"`         // 卡号
	UserID           string    `json:"userId" gorm:"size:64;not null"`               // 用户ID
	UserPhone        string    `json:"userPhone" gorm:"size:32;not null"`            // 用户手机号
	IssueReason      string    `json:"issueReason" gorm:"size:64;not null"`          // 发放原因
	IssueType        int       `json:"issueType" gorm:"not null"`                    // 1=实体卡,2=虚拟卡
	RechargeCenterID string    `json:"rechargeCenterId" gorm:"size:64;not null"`     // 充值中心ID
	OperatorID       string    `json:"operatorId" gorm:"size:64;not null"`           // 操作员ID
	RelatedUserPhone string    `json:"relatedUserPhone" gorm:"size:32"`              // 推荐奖励时关联购买人手机号
	Remark           string    `json:"remark" gorm:"size:500"`                       // 备注
	CreatedAt        time.Time `json:"createdAt" gorm:"autoCreateTime"`
}

// CenterMonthlyConsumption 充值中心月度消费记录（手动录入商城消费）
type CenterMonthlyConsumption struct {
	ID          string    `json:"id" gorm:"primaryKey;size:64"`
	CenterID    string    `json:"centerId" gorm:"uniqueIndex:uk_center_month;size:64"`
	Month       string    `json:"month" gorm:"uniqueIndex:uk_center_month;size:7"`
	Consumption float64   `json:"consumption"`
	CreatedAt   time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}

// CardTransaction 门店卡交易记录
type CardTransaction struct {
	ID            string    `json:"id" gorm:"primaryKey;size:64"`
	CardNo        string    `json:"cardNo" gorm:"index;size:32;not null"`  // 卡号
	Type          string    `json:"type" gorm:"size:32;not null"`          // issue/consume/freeze/unfreeze/activate/void
	Amount        int       `json:"amount" gorm:"default:0"`               // 金额（元）
	BalanceBefore int       `json:"balanceBefore" gorm:"default:0"`        // 交易前余额
	BalanceAfter  int       `json:"balanceAfter" gorm:"default:0"`         // 交易后余额
	Remark        string    `json:"remark" gorm:"size:500"`                // 备注
	OperatorID    string    `json:"operatorId" gorm:"size:64"`             // 操作员ID
	CreatedAt     time.Time `json:"createdAt" gorm:"autoCreateTime"`
}

// RechargeCenter 充值中心
type RechargeCenter struct {
	ID        string    `json:"id" gorm:"primaryKey;size:64"`
	Name      string    `json:"name" gorm:"uniqueIndex;size:128"`   // 中心名称
	Code      string    `json:"code" gorm:"uniqueIndex;size:64"`    // 中心编码
	Province  string    `json:"province" gorm:"size:32"`            // 省
	City      string    `json:"city" gorm:"size:32"`                // 市
	District  string    `json:"district" gorm:"size:32"`            // 区/县
	Address   string    `json:"address"`                            // 具体位置
	ManagerID string    `json:"managerId" gorm:"index;size:64"`     // 管理员(操作员)ID
	Phone     string    `json:"phone"`                              // 联系电话
	Balance   float64   `json:"balance" gorm:"default:0"`           // 积分余额
	Status    string    `json:"status" gorm:"default:'active';size:32"` // active/inactive
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}

// RechargeOperator 充值操作员
type RechargeOperator struct {
	ID        string    `json:"id" gorm:"primaryKey;size:64"`
	Name      string    `json:"name"`                               // 姓名
	Phone     string    `json:"phone" gorm:"uniqueIndex;size:32"`   // 手机号
	Password  string    `json:"-"`                                  // 密码
	CenterID  string    `json:"centerId" gorm:"index;size:64"`      // 所属中心ID
	Role      string    `json:"role"`                               // 角色
	Status    string    `json:"status" gorm:"default:'active';size:32"` // active/inactive
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}

// CreateBRechargeApplicationRequest B端充值申请
type CreateBRechargeApplicationRequest struct {
	CenterID             string  `json:"centerId" binding:"required"`
	CenterName           string  `json:"centerName"`
	Amount               float64 `json:"amount" binding:"required,gt=0"`
	LastMonthConsumption float64 `json:"lastMonthConsumption"`
	TransactionNo        string  `json:"transactionNo"`
	Screenshot           string  `json:"screenshot"`
	Remark               string  `json:"remark"`
}

// ApprovalRechargeApplicationRequest B端充值审批
type ApprovalRechargeApplicationRequest struct {
	ID     string `json:"id" binding:"required"`
	Action string `json:"action" binding:"required,oneof=approve reject"`
	Reason string `json:"reason"`
}

// CreateCRechargeRequest C端充值
type CreateCRechargeRequest struct {
	MemberID      string  `json:"memberId" binding:"required"`
	CenterID      string  `json:"centerId" binding:"required"`
	Amount        float64 `json:"amount" binding:"required,gt=0"`
	PaymentMethod string  `json:"paymentMethod"`
	Remark        string  `json:"remark"`
	MemberName    string  `json:"memberName"`
	MemberPhone   string  `json:"memberPhone"`
	CenterName    string  `json:"centerName"`
}

// CreateCenterRequest 创建充值中心
type CreateCenterRequest struct {
	Name      string `json:"name" binding:"required,min=1,max=100"`
	Code      string `json:"code"`
	Address   string `json:"address"`
	Phone     string `json:"phone"`
	Province  string `json:"province"`
	City      string `json:"city"`
	District  string `json:"district"`
	ManagerID string `json:"managerId"`
}

// UpdateCenterRequest 更新充值中心
type UpdateCenterRequest struct {
	Name     string `json:"name" binding:"omitempty,min=1,max=100"`
	Code     string `json:"code" binding:"omitempty"`
	Address  string `json:"address" binding:"omitempty"`
	Phone    string `json:"phone" binding:"omitempty"`
	Status   string `json:"status" binding:"omitempty,oneof=active frozen"`
	Province string `json:"province" binding:"omitempty"`
	City     string `json:"city" binding:"omitempty"`
	District string `json:"district" binding:"omitempty"`
}

// CreateOperatorRequest 创建操作员
type CreateOperatorRequest struct {
	Name     string `json:"name" binding:"required,min=1,max=50"`
	Phone    string `json:"phone" binding:"required,len=11"`
	Password string `json:"password" binding:"required,min=6,max=32"`
	CenterID string `json:"centerId" binding:"required"`
	Role     string `json:"role" binding:"required,oneof=center_admin operator"`
}

// UpdateOperatorRequest 更新操作员
type UpdateOperatorRequest struct {
	Name     string `json:"name" binding:"omitempty,min=1,max=50"`
	Phone    string `json:"phone" binding:"omitempty,len=11"`
	Password string `json:"password" binding:"omitempty,min=6,max=32"`
	Role     string `json:"role" binding:"omitempty,oneof=center_admin operator"`
	Status   string `json:"status" binding:"omitempty,oneof=active inactive"`
	CenterID string `json:"centerId" binding:"omitempty"`
}

// TableName 显式指定表名
func (RechargeApplication) TableName() string { return "recharge_applications" }
func (CRecharge) TableName() string           { return "c_recharges" }
func (StoreCard) TableName() string           { return "store_cards" }
func (CardIssueRecord) TableName() string     { return "card_issue_records" }
func (CardTransaction) TableName() string     { return "card_transactions" }
func (RechargeCenter) TableName() string      { return "recharge_centers" }
func (CenterMonthlyConsumption) TableName() string { return "center_monthly_consumption" }
