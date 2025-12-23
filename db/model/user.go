package model

import "time"

// UserModel 用户模型
type UserModel struct {
	Id        int32     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	OpenId    string    `gorm:"column:openid;uniqueIndex;size:64" json:"openid"`
	UnionId   string    `gorm:"column:unionid;size:64;index" json:"unionid"`
	NickName  string    `gorm:"column:nickname;size:64" json:"nickname"`
	AvatarUrl string    `gorm:"column:avatar_url;size:255" json:"avatarUrl"`
	Phone     string    `gorm:"column:phone;size:20" json:"phone"`
	Gender    int       `gorm:"column:gender;default:0" json:"gender"` // 0-未知 1-男 2-女
	MemberLevel int    `gorm:"column:member_level;default:0" json:"memberLevel"` // 会员等级 0-普通 1-银卡 2-金卡 3-钻石
	Balance   float64   `gorm:"column:balance;default:0;type:decimal(10,2)" json:"balance"` // 余额
	Points    int       `gorm:"column:points;default:0" json:"points"` // 积分
	Status    int       `gorm:"column:status;default:1" json:"status"` // 状态 1-正常 0-禁用
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (UserModel) TableName() string {
	return "user"
}

