package model

import "time"

// DeliveryTrackModel 配送跟踪模型
type DeliveryTrackModel struct {
	Id          int32     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	DeliveryId  int32     `gorm:"column:delivery_id;index" json:"deliveryId"` // 配送信息ID
	OrderId     int32     `gorm:"column:order_id;index" json:"orderId"` // 订单ID
	Status      int       `gorm:"column:status" json:"status"` // 状态
	StatusDesc  string    `gorm:"column:status_desc;size:128" json:"statusDesc"` // 状态描述
	Latitude    float64   `gorm:"column:latitude;type:decimal(10,7)" json:"latitude"` // 纬度
	Longitude   float64   `gorm:"column:longitude;type:decimal(10,7)" json:"longitude"` // 经度
	Address     string    `gorm:"column:address;size:255" json:"address"` // 地址描述
	CreatedAt   time.Time `gorm:"column:created_at" json:"createdAt"`
}

func (DeliveryTrackModel) TableName() string {
	return "delivery_track"
}

