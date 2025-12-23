package model

import "time"

// OrderItemModel 订单项模型
type OrderItemModel struct {
	Id        int32     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	OrderId   int32     `gorm:"column:order_id;index" json:"orderId"` // 订单ID
	ProductId int32     `gorm:"column:product_id;index" json:"productId"` // 商品ID
	ProductName string  `gorm:"column:product_name;size:128" json:"productName"` // 商品名称（快照）
	ProductImage string `gorm:"column:product_image;size:255" json:"productImage"` // 商品图片（快照）
	Price     float64   `gorm:"column:price;type:decimal(10,2)" json:"price"` // 单价
	Quantity  int       `gorm:"column:quantity" json:"quantity"` // 数量
	Amount    float64   `gorm:"column:amount;type:decimal(10,2)" json:"amount"` // 小计
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
}

func (OrderItemModel) TableName() string {
	return "order_item"
}

