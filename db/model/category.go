package model

import "time"

// CategoryModel 商品分类模型
type CategoryModel struct {
	Id        int32     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"column:name;size:64" json:"name"` // 分类名称
	Icon      string    `gorm:"column:icon;size:255" json:"icon"` // 分类图标
	Sort      int       `gorm:"column:sort;default:0" json:"sort"` // 排序
	Status    int       `gorm:"column:status;default:1" json:"status"` // 状态 1-启用 0-禁用
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (CategoryModel) TableName() string {
	return "category"
}

