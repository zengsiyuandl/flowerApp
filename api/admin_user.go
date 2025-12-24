package api

import (
	"strconv"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
	"wxcloudrun-golang/utils"

	"github.com/gin-gonic/gin"
)

// AdminGetUserList 管理后台获取用户列表
func AdminGetUserList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	keyword := c.Query("keyword")

	query := db.Get().Model(&model.UserModel{})

	if keyword != "" {
		query = query.Where("nickname LIKE ? OR phone LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	var total int64
	query.Count(&total)

	var users []model.UserModel
	offset := (page - 1) * pageSize
	query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&users)

	utils.Success(map[string]interface{}{
		"list":     users,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	}).WriteJSON(c.Writer)
}

// AdminRechargeUser 给用户充值
func AdminRechargeUser(c *gin.Context) {
	userId, _ := strconv.ParseInt(c.Param("userId"), 10, 32)

	var req struct {
		Amount float64 `json:"amount" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(400, "参数错误").WriteJSON(c.Writer)
		return
	}

	// 更新用户余额
	var user model.UserModel
	if err := db.Get().Where("id = ?", userId).First(&user).Error; err != nil {
		utils.Error(404, "用户不存在").WriteJSON(c.Writer)
		return
	}

	user.Balance += req.Amount
	if err := db.Get().Save(&user).Error; err != nil {
		utils.Error(500, "充值失败").WriteJSON(c.Writer)
		return
	}

	// 创建充值记录
	recharge := model.RechargeModel{
		UserId:      int32(userId),
		OrderNo:     utils.GenerateRechargeNo(),
		Amount:      req.Amount,
		BonusAmount: 0,
		Status:      1, // 已支付
	}
	db.Get().Create(&recharge)

	utils.Success(nil).WriteJSON(c.Writer)
}

