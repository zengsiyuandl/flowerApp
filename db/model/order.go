package model

import "time"

// OrderModel 订单模型
type OrderModel struct {
	Id            int32     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	OrderNo       string    `gorm:"column:order_no;uniqueIndex;size:64" json:"orderNo"` // 订单号
	UserId        int32     `gorm:"column:user_id;index" json:"userId"` // 用户ID
	TotalAmount   float64   `gorm:"column:total_amount;type:decimal(10,2)" json:"totalAmount"` // 订单总金额
	DiscountAmount float64 `gorm:"column:discount_amount;default:0;type:decimal(10,2)" json:"discountAmount"` // 优惠金额
	PayAmount     float64   `gorm:"column:pay_amount;type:decimal(10,2)" json:"payAmount"` // 实付金额
	DeliveryType  int       `gorm:"column:delivery_type;default:1" json:"deliveryType"` // 配送方式 1-外卖 2-自取
	DeliveryTime  time.Time `gorm:"column:delivery_time" json:"deliveryTime"` // 配送时间
	AddressId     int32     `gorm:"column:address_id" json:"addressId"` // 收货地址ID
	ReceiverName  string    `gorm:"column:receiver_name;size:32" json:"receiverName"` // 收货人姓名
	ReceiverPhone string    `gorm:"column:receiver_phone;size:20" json:"receiverPhone"` // 收货人电话
	ReceiverAddress string  `gorm:"column:receiver_address;size:255" json:"receiverAddress"` // 收货地址
	Remark        string    `gorm:"column:remark;size:255" json:"remark"` // 备注
	CouponId      int32     `gorm:"column:coupon_id;default:0" json:"couponId"` // 使用的优惠券ID
	Status        int       `gorm:"column:status;default:0" json:"status"` // 订单状态 0-待支付 1-待发货 2-配送中 3-已完成 4-已取消 5-已退款
	PayStatus     int       `gorm:"column:pay_status;default:0" json:"payStatus"` // 支付状态 0-未支付 1-已支付
	PayTime       time.Time `gorm:"column:pay_time" json:"payTime"` // 支付时间
	ShipTime      time.Time `gorm:"column:ship_time" json:"shipTime"` // 发货时间
	CompleteTime  time.Time `gorm:"column:complete_time" json:"completeTime"` // 完成时间
	CreatedAt     time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt     time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (OrderModel) TableName() string {
	return "order"
}

