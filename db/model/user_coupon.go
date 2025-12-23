package model

import "time"

// UserCouponModel 用户优惠券关联模型
type UserCouponModel struct {
	Id        int32     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserId    int32     `gorm:"column:user_id;index" json:"userId"` // 用户ID
	CouponId  int32     `gorm:"column:coupon_id;index" json:"couponId"` // 优惠券ID
	Status    int       `gorm:"column:status;default:0" json:"status"` // 状态 0-未使用 1-已使用 2-已过期
	UsedTime  time.Time `gorm:"column:used_time" json:"usedTime"` // 使用时间
	OrderId   int32     `gorm:"column:order_id;default:0" json:"orderId"` // 使用的订单ID
	ExpireTime time.Time `gorm:"column:expire_time" json:"expireTime"` // 过期时间
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (UserCouponModel) TableName() string {
	return "user_coupon"
}

