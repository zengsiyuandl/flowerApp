package model

import "time"

// PaymentModel 支付记录模型
type PaymentModel struct {
	Id          int32     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	OrderNo     string    `gorm:"column:order_no;index;size:64" json:"orderNo"` // 订单号
	OrderId     int32     `gorm:"column:order_id;index" json:"orderId"` // 订单ID
	UserId      int32     `gorm:"column:user_id;index" json:"userId"` // 用户ID
	PaymentType int       `gorm:"column:payment_type;default:1" json:"paymentType"` // 支付方式 1-微信支付 2-余额支付
	Amount      float64   `gorm:"column:amount;type:decimal(10,2)" json:"amount"` // 支付金额
	TradeNo     string    `gorm:"column:trade_no;size:64;index" json:"tradeNo"` // 微信交易号
	Status      int       `gorm:"column:status;default:0" json:"status"` // 支付状态 0-待支付 1-已支付 2-已退款
	PayTime     time.Time `gorm:"column:pay_time" json:"payTime"` // 支付时间
	CreatedAt   time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (PaymentModel) TableName() string {
	return "payment"
}

