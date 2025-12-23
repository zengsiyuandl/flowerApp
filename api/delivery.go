package api

import (
	"strconv"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
	"wxcloudrun-golang/middleware"
	"wxcloudrun-golang/utils"

	"github.com/gin-gonic/gin"
)

// GetDeliveryTrack 获取配送跟踪信息
func GetDeliveryTrack(c *gin.Context) {
	userId := middleware.GetUserId(c)
	orderIdStr := c.Param("orderId")
	orderId, _ := strconv.ParseInt(orderIdStr, 10, 32)

	// 验证订单是否属于当前用户
	var order model.OrderModel
	if err := db.Get().Where("id = ? AND user_id = ?", orderId, userId).First(&order).Error; err != nil {
		utils.Error(404, "订单不存在").WriteJSON(c.Writer)
		return
	}

	// 获取配送信息
	var delivery model.DeliveryModel
	if err := db.Get().Where("order_id = ?", orderId).First(&delivery).Error; err != nil {
		utils.Error(404, "配送信息不存在").WriteJSON(c.Writer)
		return
	}

	// 获取配送跟踪记录
	var tracks []model.DeliveryTrackModel
	db.Get().Where("order_id = ?", orderId).Order("created_at DESC").Find(&tracks)

	utils.Success(map[string]interface{}{
		"delivery": delivery,
		"tracks":   tracks,
	}).WriteJSON(c.Writer)
}

