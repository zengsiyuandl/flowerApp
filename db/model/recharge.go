package model

import "time"

// RechargeModel 储值充值记录模型
type RechargeModel struct {
	Id          int32     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserId      int32     `gorm:"column:user_id;index" json:"userId"`
	OrderNo     string    `gorm:"column:order_no;uniqueIndex;size:64" json:"orderNo"` // 充值订单号
	Amount      float64   `gorm:"column:amount;type:decimal(10,2)" json:"amount"` // 充值金额
	BonusAmount float64   `gorm:"column:bonus_amount;default:0;type:decimal(10,2)" json:"bonusAmount"` // 赠送金额
	PaymentId   int32     `gorm:"column:payment_id;index" json:"paymentId"` // 支付记录ID
	Status      int       `gorm:"column:status;default:0" json:"status"` // 状态 0-待支付 1-已支付 2-已取消
	CreatedAt   time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (RechargeModel) TableName() string {
	return "recharge"
}

