package api

import (
	"strconv"
	"wxcloudrun-golang/config"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
	"wxcloudrun-golang/middleware"
	"wxcloudrun-golang/service"
	"wxcloudrun-golang/utils"

	"github.com/gin-gonic/gin"
)

// CreateRechargeRequest 创建储值订单请求
type CreateRechargeRequest struct {
	Amount float64 `json:"amount" binding:"required,min=0.01"`
}

// CreateRecharge 创建储值订单
func CreateRecharge(c *gin.Context) {
	userId := middleware.GetUserId(c)

	var req CreateRechargeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(400, "参数错误").WriteJSON(c.Writer)
		return
	}

	// 计算赠送金额（示例：充100送10）
	bonusAmount := float64(0)
	if req.Amount >= 100 {
		bonusAmount = req.Amount * 0.1
	}

	// 生成充值订单号
	orderNo := utils.GenerateRechargeNo()

	// 创建充值记录
	recharge := model.RechargeModel{
		UserId:      userId,
		OrderNo:     orderNo,
		Amount:      req.Amount,
		BonusAmount: bonusAmount,
		Status:      0, // 待支付
	}

	if err := db.Get().Create(&recharge).Error; err != nil {
		utils.Error(500, "创建充值订单失败").WriteJSON(c.Writer)
		return
	}

	// 创建支付记录
	payment := model.PaymentModel{
		OrderNo:    orderNo,
		UserId:     userId,
		PaymentType: 1, // 微信支付
		Amount:     req.Amount,
		Status:     0, // 待支付
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

	recharge.PaymentId = payment.Id
	db.Get().Save(&recharge)

	// 调用统一下单接口
	openID := middleware.GetOpenId(c)
	if openID == "" {
		utils.Error(400, "无法获取用户OpenID").WriteJSON(c.Writer)
		return
	}

	if config.AppConfig.WxPaySubMchId == "" {
		utils.Error(500, "微信支付未配置").WriteJSON(c.Writer)
		return
	}

	paymentService := service.NewPaymentService()
	totalFee := int(req.Amount * 100) // 转换为分
	clientIP := c.ClientIP()
	if clientIP == "" {
		clientIP = "127.0.0.1"
	}

	unifiedReq := &service.UnifiedOrderRequest{
		OpenID:         openID,
		Body:           "储值充值-" + orderNo,
		OutTradeNo:    orderNo,
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

	utils.Success(map[string]interface{}{
		"recharge":  recharge,
		"paymentId": payment.Id,
		"params":    unifiedResp.RespData.Payment,
	}).WriteJSON(c.Writer)
}

// GetRechargeList 获取储值记录列表
func GetRechargeList(c *gin.Context) {
	userId := middleware.GetUserId(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	query := db.Get().Model(&model.RechargeModel{}).Where("user_id = ?", userId)

	var total int64
	query.Count(&total)

	var recharges []model.RechargeModel
	offset := (page - 1) * pageSize
	query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&recharges)

	utils.Success(map[string]interface{}{
		"list":     recharges,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	}).WriteJSON(c.Writer)
}

