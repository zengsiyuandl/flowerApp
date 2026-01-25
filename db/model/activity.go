package model

import "time"

// ActivityModel 活动模型
type ActivityModel struct {
	Id        int32     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Type      int       `gorm:"column:type" json:"type"` // 活动类型 1-满减 2-折扣 3-赠品 4-新用户
	Title     string    `gorm:"column:title;size:128" json:"title"` // 活动标题
	Content   string    `gorm:"column:content;type:text" json:"content"` // 活动内容
	ImageId   int32     `gorm:"column:image_id" json:"imageId"`     // 活动图片ID
	StartTime time.Time `gorm:"column:start_time" json:"startTime"` // 开始时间
	EndTime   time.Time `gorm:"column:end_time" json:"endTime"` // 结束时间
	Status    int       `gorm:"column:status;default:1" json:"status"` // 状态 1-进行中 0-已结束
	Sort      int       `gorm:"column:sort;default:0" json:"sort"` // 排序
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (ActivityModel) TableName() string {
	return "activity"
}

