package model

import "time"

// CouponModel 优惠券模型
type CouponModel struct {
	Id          int32     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"column:name;size:128" json:"name"` // 优惠券名称
	Type        int       `gorm:"column:type;default:1" json:"type"` // 类型 1-满减 2-折扣
	Amount      float64   `gorm:"column:amount;type:decimal(10,2)" json:"amount"` // 面额（满减金额或折扣比例）
	MinAmount   float64   `gorm:"column:min_amount;default:0;type:decimal(10,2)" json:"minAmount"` // 最低使用金额
	TotalCount  int       `gorm:"column:total_count;default:0" json:"totalCount"` // 总数量 0-不限
	UsedCount   int       `gorm:"column:used_count;default:0" json:"usedCount"` // 已使用数量
	StartTime   time.Time `gorm:"column:start_time" json:"startTime"` // 开始时间
	EndTime     time.Time `gorm:"column:end_time" json:"endTime"` // 结束时间
	Status      int       `gorm:"column:status;default:1" json:"status"` // 状态 1-启用 0-禁用
	CreatedAt   time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (CouponModel) TableName() string {
	return "coupon"
}

