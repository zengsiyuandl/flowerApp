package api

import (
	"strconv"
	"time"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
	"wxcloudrun-golang/middleware"
	"wxcloudrun-golang/utils"

	"github.com/gin-gonic/gin"
)

// CreatePaymentRequest 创建支付请求
type CreatePaymentRequest struct {
	OrderId int32 `json:"orderId" binding:"required"`
}

// CreatePayment 创建支付（打桩实现）
func CreatePayment(c *gin.Context) {
	userId := middleware.GetUserId(c)

	var req CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(400, "参数错误").WriteJSON(c.Writer)
		return
	}

	// 查询订单
	var order model.OrderModel
	if err := db.Get().Where("id = ? AND user_id = ?", req.OrderId, userId).First(&order).Error; err != nil {
		utils.Error(404, "订单不存在").WriteJSON(c.Writer)
		return
	}

	if order.PayStatus == 1 {
		utils.Error(400, "订单已支付").WriteJSON(c.Writer)
		return
	}

	// 创建支付记录
	payment := model.PaymentModel{
		OrderNo:    order.OrderNo,
		OrderId:    order.Id,
		UserId:     userId,
		PaymentType: 1, // 微信支付
		Amount:     order.PayAmount,
		Status:     0, // 待支付
	}

	if err := db.Get().Create(&payment).Error; err != nil {
		utils.Error(500, "创建支付记录失败").WriteJSON(c.Writer)
		return
	}

	// 生成模拟支付参数（打桩实现）
	paymentParams := utils.GenerateMockPaymentParams(order.OrderNo, order.PayAmount)

	utils.Success(map[string]interface{}{
		"paymentId": payment.Id,
		"params":    paymentParams,
	}).WriteJSON(c.Writer)
}

// PaymentNotifyRequest 支付回调请求（打桩实现）
type PaymentNotifyRequest struct {
	OrderNo string `json:"orderNo"`
	TradeNo string `json:"tradeNo"`
}

// PaymentNotify 支付回调通知（打桩实现）
func PaymentNotify(c *gin.Context) {
	var req PaymentNotifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(400, "参数错误").WriteJSON(c.Writer)
		return
	}

	// 查询支付记录
	var payment model.PaymentModel
	if err := db.Get().Where("order_no = ?", req.OrderNo).First(&payment).Error; err != nil {
		utils.Error(404, "支付记录不存在").WriteJSON(c.Writer)
		return
	}

	if payment.Status == 1 {
		utils.Success(nil).WriteJSON(c.Writer)
		return
	}

	// 开始事务
	tx := db.Get().Begin()

	// 更新支付记录
	now := time.Now()
	if err := tx.Model(&payment).Updates(map[string]interface{}{
		"trade_no": req.TradeNo,
		"status":   1,
		"pay_time": now,
	}).Error; err != nil {
		tx.Rollback()
		utils.Error(500, "更新支付记录失败").WriteJSON(c.Writer)
		return
	}

	// 更新订单状态
	var order model.OrderModel
	if err := tx.Where("id = ?", payment.OrderId).First(&order).Error; err != nil {
		tx.Rollback()
		utils.Error(404, "订单不存在").WriteJSON(c.Writer)
		return
	}

	if err := tx.Model(&order).Updates(map[string]interface{}{
		"pay_status": 1,
		"pay_time":   now,
		"status":     1, // 待发货
	}).Error; err != nil {
		tx.Rollback()
		utils.Error(500, "更新订单状态失败").WriteJSON(c.Writer)
		return
	}

	// 如果使用了优惠券，标记为已使用
	if order.CouponId > 0 {
		tx.Model(&model.UserCouponModel{}).
			Where("user_id = ? AND coupon_id = ?", order.UserId, order.CouponId).
			Updates(map[string]interface{}{
				"status":   1,
				"used_time": now,
				"order_id": order.Id,
			})
	}

	if err := tx.Commit().Error; err != nil {
		utils.Error(500, "处理支付回调失败").WriteJSON(c.Writer)
		return
	}

	utils.Success(nil).WriteJSON(c.Writer)
}

// GetPaymentStatus 查询支付状态
func GetPaymentStatus(c *gin.Context) {
	orderIdStr := c.Param("orderId")
	orderId, _ := strconv.ParseInt(orderIdStr, 10, 32)

	var payment model.PaymentModel
	if err := db.Get().Where("order_id = ?", orderId).First(&payment).Error; err != nil {
		utils.Error(404, "支付记录不存在").WriteJSON(c.Writer)
		return
	}

	utils.Success(map[string]interface{}{
		"status":  payment.Status,
		"payTime": payment.PayTime,
	}).WriteJSON(c.Writer)
}

