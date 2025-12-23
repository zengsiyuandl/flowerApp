package model

import "time"

// LotteryModel 抽奖活动模型
type LotteryModel struct {
	Id        int32     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"column:name;size:128" json:"name"` // 活动名称
	Image     string    `gorm:"column:image;size:255" json:"image"` // 活动图片
	Rule      string    `gorm:"column:rule;type:text" json:"rule"` // 活动规则
	StartTime time.Time `gorm:"column:start_time" json:"startTime"` // 开始时间
	EndTime   time.Time `gorm:"column:end_time" json:"endTime"` // 结束时间
	DailyLimit int      `gorm:"column:daily_limit;default:1" json:"dailyLimit"` // 每日抽奖次数限制
	TotalLimit int      `gorm:"column:total_limit;default:0" json:"totalLimit"` // 总抽奖次数限制 0-不限
	Status    int       `gorm:"column:status;default:1" json:"status"` // 状态 1-进行中 0-已结束
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (LotteryModel) TableName() string {
	return "lottery"
}

