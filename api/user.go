package api

import (
	"strconv"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
	"wxcloudrun-golang/middleware"
	"wxcloudrun-golang/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UserLoginRequest 登录请求
type UserLoginRequest struct {
	Code string `json:"code" binding:"required"`
}

// UserLoginResponse 登录响应
type UserLoginResponse struct {
	Token  string           `json:"token"`
	User   *model.UserModel `json:"user"`
}

// UserLogin 微信登录
func UserLogin(c *gin.Context) {
	var req UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(400, "参数错误").WriteJSON(c.Writer)
		return
	}

	// 调用微信API获取openid
	wxResp, err := utils.Code2Session(req.Code)
	if err != nil {
		utils.Error(500, err.Error()).WriteJSON(c.Writer)
		return
	}

	// 查询或创建用户
	var user model.UserModel
	err = db.Get().Where("openid = ?", wxResp.OpenId).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		// 创建新用户
		user = model.UserModel{
			OpenId:   wxResp.OpenId,
			UnionId:  wxResp.UnionId,
			Status:   1,
		}
		if err := db.Get().Create(&user).Error; err != nil {
			utils.Error(500, "创建用户失败").WriteJSON(c.Writer)
			return
		}
	}

	// 生成token
	token := utils.GenerateToken(user.Id, user.OpenId)

	utils.Success(UserLoginResponse{
		Token: token,
		User:  &user,
	}).WriteJSON(c.Writer)
}

// GetUserInfo 获取用户信息
func GetUserInfo(c *gin.Context) {
	userId := middleware.GetUserId(c)

	var user model.UserModel
	if err := db.Get().Where("id = ?", userId).First(&user).Error; err != nil {
		utils.Error(404, "用户不存在").WriteJSON(c.Writer)
		return
	}

	utils.Success(user).WriteJSON(c.Writer)
}

// UpdateUserInfoRequest 更新用户信息请求
type UpdateUserInfoRequest struct {
	NickName  string `json:"nickname"`
	AvatarUrl string `json:"avatarUrl"`
	Phone     string `json:"phone"`
	Gender    int    `json:"gender"`
}

// UpdateUserInfo 更新用户信息
func UpdateUserInfo(c *gin.Context) {
	userId := middleware.GetUserId(c)

	var req UpdateUserInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(400, "参数错误").WriteJSON(c.Writer)
		return
	}

	updates := make(map[string]interface{})
	if req.NickName != "" {
		updates["nickname"] = req.NickName
	}
	if req.AvatarUrl != "" {
		updates["avatar_url"] = req.AvatarUrl
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.Gender > 0 {
		updates["gender"] = req.Gender
	}

	if err := db.Get().Model(&model.UserModel{}).Where("id = ?", userId).Updates(updates).Error; err != nil {
		utils.Error(500, "更新失败").WriteJSON(c.Writer)
		return
	}

	utils.Success(nil).WriteJSON(c.Writer)
}

// GetAddressList 获取地址列表
func GetAddressList(c *gin.Context) {
	userId := middleware.GetUserId(c)

	var addresses []model.AddressModel
	db.Get().Where("user_id = ?", userId).Order("is_default DESC, id DESC").Find(&addresses)

	utils.Success(addresses).WriteJSON(c.Writer)
}

// AddAddressRequest 添加地址请求
type AddAddressRequest struct {
	Name      string  `json:"name" binding:"required"`
	Phone     string  `json:"phone" binding:"required"`
	Province  string  `json:"province" binding:"required"`
	City      string  `json:"city" binding:"required"`
	District  string  `json:"district" binding:"required"`
	Address   string  `json:"address" binding:"required"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	IsDefault int     `json:"isDefault"`
}

// AddAddress 添加地址
func AddAddress(c *gin.Context) {
	userId := middleware.GetUserId(c)

	var req AddAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(400, "参数错误").WriteJSON(c.Writer)
		return
	}

	// 如果设置为默认地址，先取消其他默认地址
	if req.IsDefault == 1 {
		db.Get().Model(&model.AddressModel{}).Where("user_id = ?", userId).Update("is_default", 0)
	}

	address := model.AddressModel{
		UserId:    userId,
		Name:      req.Name,
		Phone:     req.Phone,
		Province:  req.Province,
		City:      req.City,
		District:  req.District,
		Address:   req.Address,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		IsDefault: req.IsDefault,
	}

	if err := db.Get().Create(&address).Error; err != nil {
		utils.Error(500, "添加失败").WriteJSON(c.Writer)
		return
	}

	utils.Success(address).WriteJSON(c.Writer)
}

// UpdateAddress 更新地址
func UpdateAddress(c *gin.Context) {
	userId := middleware.GetUserId(c)
	id, _ := strconv.ParseInt(c.Param("id"), 10, 32)

	var req AddAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(400, "参数错误").WriteJSON(c.Writer)
		return
	}

	// 检查地址是否属于当前用户
	var address model.AddressModel
	if err := db.Get().Where("id = ? AND user_id = ?", id, userId).First(&address).Error; err != nil {
		utils.Error(404, "地址不存在").WriteJSON(c.Writer)
		return
	}

	// 如果设置为默认地址，先取消其他默认地址
	if req.IsDefault == 1 {
		db.Get().Model(&model.AddressModel{}).Where("user_id = ? AND id != ?", userId, id).Update("is_default", 0)
	}

	updates := map[string]interface{}{
		"name":       req.Name,
		"phone":      req.Phone,
		"province":    req.Province,
		"city":        req.City,
		"district":    req.District,
		"address":     req.Address,
		"latitude":    req.Latitude,
		"longitude":   req.Longitude,
		"is_default":  req.IsDefault,
	}

	if err := db.Get().Model(&address).Updates(updates).Error; err != nil {
		utils.Error(500, "更新失败").WriteJSON(c.Writer)
		return
	}

	utils.Success(nil).WriteJSON(c.Writer)
}

// DeleteAddress 删除地址
func DeleteAddress(c *gin.Context) {
	userId := middleware.GetUserId(c)
	id, _ := strconv.ParseInt(c.Param("id"), 10, 32)

	if err := db.Get().Where("id = ? AND user_id = ?", id, userId).Delete(&model.AddressModel{}).Error; err != nil {
		utils.Error(500, "删除失败").WriteJSON(c.Writer)
		return
	}

	utils.Success(nil).WriteJSON(c.Writer)
}

