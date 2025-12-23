package model

import "time"

// DeliveryModel 配送信息模型
type DeliveryModel struct {
	Id            int32     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	OrderId       int32     `gorm:"column:order_id;uniqueIndex" json:"orderId"` // 订单ID
	OrderNo       string    `gorm:"column:order_no;index;size:64" json:"orderNo"` // 订单号
	DeliveryType  int       `gorm:"column:delivery_type;default:1" json:"deliveryType"` // 配送方式 1-自营 2-达达 3-美团
	DeliveryManId int32     `gorm:"column:delivery_man_id;default:0" json:"deliveryManId"` // 配送员ID（自营配送）
	DeliveryManName string  `gorm:"column:delivery_man_name;size:32" json:"deliveryManName"` // 配送员姓名
	DeliveryManPhone string `gorm:"column:delivery_man_phone;size:20" json:"deliveryManPhone"` // 配送员电话
	ThirdPartyOrderNo string `gorm:"column:third_party_order_no;size:64" json:"thirdPartyOrderNo"` // 第三方配送单号
	Status          int     `gorm:"column:status;default:0" json:"status"` // 配送状态 0-待分配 1-已分配 2-配送中 3-已完成 4-已取消
	StartTime      time.Time `gorm:"column:start_time" json:"startTime"` // 开始配送时间
	CompleteTime   time.Time `gorm:"column:complete_time" json:"completeTime"` // 完成时间
	CreatedAt      time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt      time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (DeliveryModel) TableName() string {
	return "delivery"
}

