package model

import "time"

// AddressModel 收货地址模型
type AddressModel struct {
	Id        int32     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserId    int32     `gorm:"column:user_id;index" json:"userId"`
	Name      string    `gorm:"column:name;size:32" json:"name"` // 收货人姓名
	Phone     string    `gorm:"column:phone;size:20" json:"phone"` // 收货人电话
	Province  string    `gorm:"column:province;size:32" json:"province"` // 省份
	City      string    `gorm:"column:city;size:32" json:"city"` // 城市
	District  string    `gorm:"column:district;size:32" json:"district"` // 区县
	Address   string    `gorm:"column:address;size:255" json:"address"` // 详细地址
	Latitude  float64   `gorm:"column:latitude;type:decimal(10,7)" json:"latitude"` // 纬度
	Longitude float64   `gorm:"column:longitude;type:decimal(10,7)" json:"longitude"` // 经度
	IsDefault int       `gorm:"column:is_default;default:0" json:"isDefault"` // 是否默认地址 0-否 1-是
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (AddressModel) TableName() string {
	return "address"
}

