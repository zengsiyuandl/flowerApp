package model

import "time"

// CartModel 购物车模型
type CartModel struct {
	Id        int32     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserId    int32     `gorm:"column:user_id;index" json:"userId"` // 用户ID
	ProductId int32     `gorm:"column:product_id;index" json:"productId"` // 商品ID
	Quantity  int       `gorm:"column:quantity;default:1" json:"quantity"` // 数量
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (CartModel) TableName() string {
	return "cart"
}

