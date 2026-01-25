package model

import "time"

// BannerModel Banner模型
type BannerModel struct {
	Id        int32     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Title     string    `gorm:"column:title;size:128" json:"title"`           // Banner标题
	Image     string    `gorm:"column:image;size:255" json:"image"`          // Banner图片
	LinkType  int       `gorm:"column:link_type;default:1" json:"linkType"` // 链接类型 1-商品详情 2-活动详情 3-分类列表 4-外部链接
	LinkValue string    `gorm:"column:link_value;size:255" json:"linkValue"` // 链接值
	Sort      int       `gorm:"column:sort;default:0" json:"sort"`          // 排序
	Status    int       `gorm:"column:status;default:1" json:"status"`      // 状态 1-启用 0-禁用
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (BannerModel) TableName() string {
	return "banner"
}
