package api

import (
	"encoding/json"
	"strconv"
	"time"
	"wxcloudrun-golang/config"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
	"wxcloudrun-golang/middleware"
	"wxcloudrun-golang/service"
	"wxcloudrun-golang/utils"

	"github.com/gin-gonic/gin"
)

// CreatePaymentRequest 创建支付请求
type CreatePaymentRequest struct {
	OrderID int32 `json:"orderId" binding:"required"`
}

// CreatePayment 创建支付（统一下单）
func CreatePayment(c *gin.Context) {
	userID := middleware.GetUserId(c)
	openID := middleware.GetOpenId(c)

	if openID == "" {
		utils.Error(400, "无法获取用户OpenID").WriteJSON(c.Writer)
		return
	}

	var req CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(400, "参数错误").WriteJSON(c.Writer)
		return
	}

	// 查询订单
	var order model.OrderModel
	if err := db.Get().Where("id = ? AND user_id = ?", req.OrderID, userID).First(&order).Error; err != nil {
		utils.Error(404, "订单不存在").WriteJSON(c.Writer)
		return
	}

	if order.PayStatus == 1 {
		utils.Error(400, "订单已支付").WriteJSON(c.Writer)
		return
	}

	// 检查配置
	if config.AppConfig.WxPaySubMchId == "" {
		utils.Error(500, "微信支付未配置").WriteJSON(c.Writer)
		return
	}

	// 创建或查询支付记录
	var payment model.PaymentModel
	db.Get().Where("order_id = ?", order.Id).First(&payment)
	if payment.Id == 0 {
		payment = model.PaymentModel{
			OrderNo:     order.OrderNo,
			OrderId:     order.Id,
			UserId:      userID,
			PaymentType: 1, // 微信支付
			Amount:      order.PayAmount,
			Status:      0, // 待支付
		}
		// 使用 Select 明确指定要插入的字段，排除 pay_time，避免 MySQL 报错
		fieldsToSelect := []string{
			"order_no", "order_id", "user_id", "payment_type", "amount", "trade_no", "status",
			"created_at", "updated_at",
		}
		if err := db.Get().Select(fieldsToSelect).Create(&payment).Error; err != nil {
			utils.Error(500, "创建支付记录失败").WriteJSON(c.Writer)
			return
		}
	}

	// 调用统一下单接口
	paymentService := service.NewPaymentService()
	totalFee := int(order.PayAmount * 100) // 转换为分
	clientIP := c.ClientIP()
	if clientIP == "" {
		clientIP = "127.0.0.1"
	}

	unifiedReq := &service.UnifiedOrderRequest{
		OpenID:         openID,
		Body:           "花店订单-" + order.OrderNo,
		OutTradeNo:    order.OrderNo,
		SpbillCreateIP: clientIP,
		SubMchID:       config.AppConfig.WxPaySubMchId,
		TotalFee:       totalFee,
		EnvID:          config.AppConfig.WxCloudEnvId,
		CallbackType:   service.CallbackTypeCloudRun,
	}
	unifiedReq.Container.Service = config.AppConfig.WxCloudServiceName
	unifiedReq.Container.Path = "/api/payment/notify"

	unifiedResp, err := paymentService.UnifiedOrder(unifiedReq)
	if err != nil {
		utils.Error(500, "统一下单失败: "+err.Error()).WriteJSON(c.Writer)
		return
	}

	if unifiedResp.RespData.ReturnCode != "SUCCESS" ||
		unifiedResp.RespData.ResultCode != "SUCCESS" {
		utils.Error(500, "统一下单失败: "+unifiedResp.RespData.ReturnMsg).WriteJSON(c.Writer)
		return
	}

	// 返回支付参数
	utils.Success(map[string]interface{}{
		"paymentId": payment.Id,
		"params":    unifiedResp.RespData.Payment,
	}).WriteJSON(c.Writer)
}

// PaymentNotifyCallback 微信支付回调数据结构
type PaymentNotifyCallback struct {
	ReturnCode    string `json:"return_code"`
	ReturnMsg     string `json:"return_msg"`
	ResultCode    string `json:"result_code"`
	AppID         string `json:"appid"`
	MchID         string `json:"mch_id"`
	SubAppID      string `json:"sub_appid"`
	SubMchID      string `json:"sub_mch_id"`
	NonceStr      string `json:"nonce_str"`
	Sign          string `json:"sign"`
	TradeType     string `json:"trade_type"`
	OpenID        string `json:"openid"`
	TransactionID string `json:"transaction_id"`
	OutTradeNo    string `json:"out_trade_no"`
	TotalFee      int    `json:"total_fee"`
	TimeEnd       string `json:"time_end"`
}

// PaymentNotify 支付回调通知
func PaymentNotify(c *gin.Context) {
	var callback PaymentNotifyCallback
	if err := c.ShouldBindJSON(&callback); err != nil {
		// 尝试从body读取原始数据
		body, _ := c.GetRawData()
		if err := json.Unmarshal(body, &callback); err != nil {
			c.JSON(200, gin.H{"errcode": -1, "errmsg": "参数解析失败"})
			return
		}
	}

	// 验证回调数据
	if callback.ReturnCode != "SUCCESS" || callback.ResultCode != "SUCCESS" {
		c.JSON(200, gin.H{"errcode": 0}) // 仍然返回成功，避免重复回调
		return
	}

	// 查询支付记录
	var payment model.PaymentModel
	if err := db.Get().Where("order_no = ?", callback.OutTradeNo).First(&payment).Error; err != nil {
		c.JSON(200, gin.H{"errcode": 0}) // 返回成功，避免重复回调
		return
	}

	// 如果已处理，直接返回成功
	if payment.Status == 1 {
		c.JSON(200, gin.H{"errcode": 0})
		return
	}

	// 开始事务
	tx := db.Get().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 更新支付记录
	now := time.Now()
	payTime, _ := time.Parse("20060102150405", callback.TimeEnd)
	if payTime.IsZero() {
		payTime = now
	}

	if err := tx.Model(&payment).Updates(map[string]interface{}{
		"trade_no": callback.TransactionID,
		"status":   1,
		"pay_time": payTime,
	}).Error; err != nil {
		tx.Rollback()
		c.JSON(200, gin.H{"errcode": -1, "errmsg": "更新支付记录失败"})
		return
	}

	// 检查是否是充值订单
	var recharge model.RechargeModel
	isRecharge := tx.Where("order_no = ?", callback.OutTradeNo).First(&recharge).Error == nil

	if isRecharge {
		// 处理充值订单
		if err := tx.Model(&recharge).Update("status", 1).Error; err != nil {
			tx.Rollback()
			c.JSON(200, gin.H{"errcode": -1, "errmsg": "更新充值记录失败"})
			return
		}
		// 更新用户余额（如果有余额字段）
		// 这里可以根据实际业务逻辑处理
	} else {
		// 处理普通订单
		var order model.OrderModel
		if err := tx.Where("id = ?", payment.OrderId).First(&order).Error; err != nil {
			tx.Rollback()
			c.JSON(200, gin.H{"errcode": -1, "errmsg": "订单不存在"})
			return
		}

		if err := tx.Model(&order).Updates(map[string]interface{}{
			"pay_status": 1,
			"pay_time":   payTime,
			"status":     1, // 待发货
		}).Error; err != nil {
			tx.Rollback()
			c.JSON(200, gin.H{"errcode": -1, "errmsg": "更新订单状态失败"})
			return
		}

		// 如果使用了优惠券，标记为已使用
		if order.CouponId > 0 {
			tx.Model(&model.UserCouponModel{}).
				Where("user_id = ? AND coupon_id = ?", order.UserId, order.CouponId).
				Updates(map[string]interface{}{
					"status":    1,
					"used_time": payTime,
					"order_id":  order.Id,
				})
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(200, gin.H{"errcode": -1, "errmsg": "处理支付回调失败"})
		return
	}

	// 必须返回 {"errcode": 0}，否则会持续回调
	c.JSON(200, gin.H{"errcode": 0})
}

// QueryPaymentOrderRequest 查询支付订单请求
type QueryPaymentOrderRequest struct {
	OrderNo string `json:"orderNo" binding:"required"`
}

// QueryPaymentOrder 查询支付订单
func QueryPaymentOrder(c *gin.Context) {
	orderNo := c.Param("orderNo")
	if orderNo == "" {
		utils.Error(400, "订单号不能为空").WriteJSON(c.Writer)
		return
	}

	// 查询订单
	var order model.OrderModel
	if err := db.Get().Where("order_no = ?", orderNo).First(&order).Error; err != nil {
		utils.Error(404, "订单不存在").WriteJSON(c.Writer)
		return
	}

	// 调用微信支付查询接口
	paymentService := service.NewPaymentService()
	queryReq := &service.QueryOrderRequest{
		SubMchID:   config.AppConfig.WxPaySubMchId,
		OutTradeNo: orderNo,
		NonceStr:   utils.GenerateNonceStr(),
	}

	queryResp, err := paymentService.QueryOrder(queryReq)
	if err != nil {
		utils.Error(500, "查询订单失败: "+err.Error()).WriteJSON(c.Writer)
		return
	}

	if queryResp.ErrCode != 0 {
		utils.Error(500, "查询订单失败: "+queryResp.ErrMsg).WriteJSON(c.Writer)
		return
	}

	utils.Success(map[string]interface{}{
		"tradeState":    queryResp.RespData.TradeState,
		"transactionId": queryResp.RespData.TransactionID,
		"outTradeNo":   queryResp.RespData.OutTradeNo,
		"totalFee":     queryResp.RespData.TotalFee,
		"timeEnd":      queryResp.RespData.TimeEnd,
	}).WriteJSON(c.Writer)
}

// ClosePaymentOrderRequest 关闭支付订单请求
type ClosePaymentOrderRequest struct {
	OrderNo string `json:"orderNo" binding:"required"`
}

// ClosePaymentOrder 关闭支付订单
func ClosePaymentOrder(c *gin.Context) {
	var req ClosePaymentOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(400, "参数错误").WriteJSON(c.Writer)
		return
	}

	// 查询订单
	var order model.OrderModel
	if err := db.Get().Where("order_no = ?", req.OrderNo).First(&order).Error; err != nil {
		utils.Error(404, "订单不存在").WriteJSON(c.Writer)
		return
	}

	if order.PayStatus == 1 {
		utils.Error(400, "订单已支付，无法关闭").WriteJSON(c.Writer)
		return
	}

	// 调用微信支付关闭订单接口
	paymentService := service.NewPaymentService()
	closeReq := &service.CloseOrderRequest{
		SubMchID:   config.AppConfig.WxPaySubMchId,
		OutTradeNo: req.OrderNo,
		NonceStr:   utils.GenerateNonceStr(),
	}

	closeResp, err := paymentService.CloseOrder(closeReq)
	if err != nil {
		utils.Error(500, "关闭订单失败: "+err.Error()).WriteJSON(c.Writer)
		return
	}

	if closeResp.ErrCode != 0 {
		utils.Error(500, "关闭订单失败: "+closeResp.ErrMsg).WriteJSON(c.Writer)
		return
	}

	if closeResp.RespData.ReturnCode != "SUCCESS" ||
		closeResp.RespData.ResultCode != "SUCCESS" {
		utils.Error(500, "关闭订单失败").WriteJSON(c.Writer)
		return
	}

	utils.Success(nil).WriteJSON(c.Writer)
}

// RefundRequest 申请退款请求
type RefundRequest struct {
	OrderNo    string  `json:"orderNo" binding:"required"`
	RefundNo   string  `json:"refundNo" binding:"required"`
	RefundFee  float64 `json:"refundFee" binding:"required"`
	RefundDesc string  `json:"refundDesc"`
}

// Refund 申请退款
func Refund(c *gin.Context) {
	var req RefundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(400, "参数错误").WriteJSON(c.Writer)
		return
	}

	// 查询订单和支付记录
	var order model.OrderModel
	if err := db.Get().Where("order_no = ?", req.OrderNo).First(&order).Error; err != nil {
		utils.Error(404, "订单不存在").WriteJSON(c.Writer)
		return
	}

	if order.PayStatus != 1 {
		utils.Error(400, "订单未支付，无法退款").WriteJSON(c.Writer)
		return
	}

	var payment model.PaymentModel
	if err := db.Get().Where("order_no = ?", req.OrderNo).First(&payment).Error; err != nil {
		utils.Error(404, "支付记录不存在").WriteJSON(c.Writer)
		return
	}

	if payment.TradeNo == "" {
		utils.Error(400, "支付交易号不存在").WriteJSON(c.Writer)
		return
	}

	// 调用微信支付退款接口
	paymentService := service.NewPaymentService()
	refundReq := &service.RefundRequest{
		SubMchID:     config.AppConfig.WxPaySubMchId,
		TransactionID: payment.TradeNo,
		OutRefundNo:  req.RefundNo,
		TotalFee:     int(order.PayAmount * 100),
		RefundFee:    int(req.RefundFee * 100),
		RefundDesc:   req.RefundDesc,
		EnvID:        config.AppConfig.WxCloudEnvId,
		CallbackType: service.CallbackTypeCloudRun,
	}
	refundReq.Container.Service = config.AppConfig.WxCloudServiceName
	refundReq.Container.Path = "/api/payment/refund/notify"

	refundResp, err := paymentService.Refund(refundReq)
	if err != nil {
		utils.Error(500, "申请退款失败: "+err.Error()).WriteJSON(c.Writer)
		return
	}

	if refundResp.ErrCode != 0 {
		utils.Error(500, "申请退款失败: "+refundResp.ErrMsg).WriteJSON(c.Writer)
		return
	}

	if refundResp.RespData.ReturnCode != "SUCCESS" ||
		refundResp.RespData.ResultCode != "SUCCESS" {
		utils.Error(500, "申请退款失败").WriteJSON(c.Writer)
		return
	}

	// 更新支付记录状态为已退款
	db.Get().Model(&payment).Update("status", 2) // 2-已退款

	utils.Success(map[string]interface{}{
		"refundId":   refundResp.RespData.RefundID,
		"outRefundNo": refundResp.RespData.OutRefundNo,
		"refundFee":  refundResp.RespData.RefundFee,
	}).WriteJSON(c.Writer)
}

// QueryRefundRequest 查询退款请求
type QueryRefundRequest struct {
	RefundNo string `json:"refundNo" binding:"required"`
}

// QueryRefund 查询退款
func QueryRefund(c *gin.Context) {
	refundNo := c.Param("refundNo")
	if refundNo == "" {
		utils.Error(400, "退款单号不能为空").WriteJSON(c.Writer)
		return
	}

	// 调用微信支付查询退款接口
	paymentService := service.NewPaymentService()
	queryReq := &service.QueryRefundRequest{
		SubMchID:   config.AppConfig.WxPaySubMchId,
		OutRefundNo: refundNo,
	}

	queryResp, err := paymentService.QueryRefund(queryReq)
	if err != nil {
		utils.Error(500, "查询退款失败: "+err.Error()).WriteJSON(c.Writer)
		return
	}

	if queryResp.ErrCode != 0 {
		utils.Error(500, "查询退款失败: "+queryResp.ErrMsg).WriteJSON(c.Writer)
		return
	}

	utils.Success(map[string]interface{}{
		"transactionId": queryResp.RespData.TransactionID,
		"outTradeNo":   queryResp.RespData.OutTradeNo,
		"refundCount":  queryResp.RespData.RefundCount,
	}).WriteJSON(c.Writer)
}

// GetPaymentStatus 查询支付状态
func GetPaymentStatus(c *gin.Context) {
	orderIDStr := c.Param("orderId")
	orderID, _ := strconv.ParseInt(orderIDStr, 10, 32)

	var payment model.PaymentModel
	if err := db.Get().Where("order_id = ?", orderID).First(&payment).Error; err != nil {
		utils.Error(404, "支付记录不存在").WriteJSON(c.Writer)
		return
	}

	utils.Success(map[string]interface{}{
		"status":  payment.Status,
		"payTime": payment.PayTime,
		"tradeNo": payment.TradeNo,
	}).WriteJSON(c.Writer)
}
