package model

import "time"

// ProductModel 商品模型
type ProductModel struct {
	Id          int32     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	CategoryId  int32     `gorm:"column:category_id;index" json:"categoryId"` // 分类ID
	Name        string    `gorm:"column:name;size:128" json:"name"` // 商品名称
	Subtitle    string    `gorm:"column:subtitle;size:255" json:"subtitle"` // 副标题
	MainImageId int32     `gorm:"column:main_image_id" json:"mainImageId"`    // 主图ID
	Price       float64   `gorm:"column:price;type:decimal(10,2)" json:"price"` // 价格
	OriginalPrice float64 `gorm:"column:original_price;type:decimal(10,2)" json:"originalPrice"` // 原价
	Stock       int       `gorm:"column:stock;default:0" json:"stock"` // 库存
	Sales       int       `gorm:"column:sales;default:0" json:"sales"` // 销量
	Description string    `gorm:"column:description;type:text" json:"description"` // 商品描述
	Detail      string    `gorm:"column:detail;type:text" json:"detail"` // 商品详情
	Sort        int       `gorm:"column:sort;default:0" json:"sort"` // 排序
	Status      int       `gorm:"column:status;default:1" json:"status"` // 状态 1-上架 0-下架
	IsHot       int       `gorm:"column:is_hot;default:0" json:"isHot"` // 是否热门 0-否 1-是
	IsNew       int       `gorm:"column:is_new;default:0" json:"isNew"` // 是否新品 0-否 1-是
	CreatedAt   time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (ProductModel) TableName() string {
	return "product"
}

