package model

import "time"

// ImageStorageModel 图片存储模型
type ImageStorageModel struct {
	Id          int32     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Category    string    `gorm:"column:category;size:50;index" json:"category"` // 分类
	Filename    string    `gorm:"column:filename;size:255" json:"filename"`     // 原始文件名
	ContentType string    `gorm:"column:content_type;size:100" json:"contentType"` // MIME类型
	Data        []byte    `gorm:"column:data;type:longblob" json:"-"`            // 图片二进制数据
	Size        int       `gorm:"column:size" json:"size"`                       // 文件大小（字节）
	CreatedAt   time.Time `gorm:"column:created_at" json:"createdAt"`
}

func (ImageStorageModel) TableName() string {
	return "image_storage"
}
