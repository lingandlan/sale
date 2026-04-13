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
	CreatedAt     time.Time `json:"createdAt" gorm:"autoCreateTime"`
}

// StoreCard 门店卡
type StoreCard struct {
	ID              string    `json:"id" gorm:"primaryKey;size:64"`
	CardNo          string    `json:"cardNo" gorm:"uniqueIndex;size:64"` // 卡号
	HolderID        string    `json:"holderId"`                          // 持卡人ID
	HolderName      string    `json:"holderName"`                        // 持卡人姓名
	HolderPhone     string    `json:"holderPhone" gorm:"index;size:32"`  // 持卡人手机
	Balance         float64   `json:"balance"`                           // 卡余额
	Status          string    `json:"status" gorm:"default:'active';size:32"` // active/inactive/expired
	IssueCenterID   string    `json:"issueCenterId"`                     // 发放中心ID
	IssueCenterName string    `json:"issueCenterName"`                   // 发放中心名称
	IssueDate       time.Time `json:"issueDate"`                         // 发放日期
	ExpiryDate      time.Time `json:"expiryDate"`                        // 过期日期
	CreatedAt       time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt       time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}

// CardTransaction 门店卡交易记录
type CardTransaction struct {
	ID           string    `json:"id" gorm:"primaryKey;size:64"`
	CardNo       string    `json:"cardNo" gorm:"index;size:64"`    // 卡号
	Type         string    `json:"type"`                           // issue/consume/recharge
	Amount       float64   `json:"amount"`                         // 金额
	BalanceAfter float64   `json:"balanceAfter"`                   // 交易后余额
	Remark       string    `json:"remark"`                         // 备注
	OperatorID   string    `json:"operatorId"`                     // 操作员ID
	CreatedAt    time.Time `json:"createdAt" gorm:"autoCreateTime"`
}

// RechargeCenter 充值中心
type RechargeCenter struct {
	ID        string    `json:"id" gorm:"primaryKey;size:64"`
	Name      string    `json:"name" gorm:"uniqueIndex;size:128"`   // 中心名称
	Code      string    `json:"code" gorm:"uniqueIndex;size:64"`    // 中心编码
	Address   string    `json:"address"`                            // 地址
	Phone     string    `json:"phone"`                              // 联系电话
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

// TableName 显式指定表名
func (RechargeApplication) TableName() string { return "recharge_applications" }
func (CRecharge) TableName() string           { return "c_recharges" }
func (StoreCard) TableName() string           { return "store_cards" }
func (CardTransaction) TableName() string     { return "card_transactions" }
func (RechargeCenter) TableName() string      { return "recharge_centers" }
func (RechargeOperator) TableName() string    { return "recharge_operators" }
