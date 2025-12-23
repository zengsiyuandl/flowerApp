package model

import "time"

// LotteryRecordModel 抽奖记录模型
type LotteryRecordModel struct {
	Id        int32     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	LotteryId int32     `gorm:"column:lottery_id;index" json:"lotteryId"` // 抽奖活动ID
	UserId    int32     `gorm:"column:user_id;index" json:"userId"` // 用户ID
	PrizeType int       `gorm:"column:prize_type" json:"prizeType"` // 奖品类型 1-优惠券 2-积分 3-余额 4-实物
	PrizeId   int32     `gorm:"column:prize_id;default:0" json:"prizeId"` // 奖品ID（如优惠券ID）
	PrizeName string    `gorm:"column:prize_name;size:128" json:"prizeName"` // 奖品名称
	PrizeValue string   `gorm:"column:prize_value;size:128" json:"prizeValue"` // 奖品值（如金额、积分）
	Status    int       `gorm:"column:status;default:0" json:"status"` // 状态 0-未领取 1-已领取
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (LotteryRecordModel) TableName() string {
	return "lottery_record"
}

