package model

import "time"

// ProductImageModel 商品图片模型
type ProductImageModel struct {
	Id        int32     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ProductId int32     `gorm:"column:product_id;index" json:"productId"` // 商品ID
	ImageId   int32     `gorm:"column:image_id" json:"imageId"`            // 图片ID
	Sort      int       `gorm:"column:sort;default:0" json:"sort"` // 排序
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
}

func (ProductImageModel) TableName() string {
	return "product_image"
}

