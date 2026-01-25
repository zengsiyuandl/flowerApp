package api

import (
	"strconv"
	"time"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
	"wxcloudrun-golang/middleware"
	"wxcloudrun-golang/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateOrderRequest 创建订单请求
type CreateOrderRequest struct {
	CartIds      []int32 `json:"cartIds"`      // 购物车ID列表（从购物车创建）
	ProductId    int32   `json:"productId"`   // 商品ID（立即购买）
	Quantity     int     `json:"quantity"`     // 数量（立即购买）
	DeliveryType int     `json:"deliveryType" binding:"required"` // 配送方式 1-外卖 2-自取
	AddressId    int32   `json:"addressId"`   // 地址ID（外卖必填）
	DeliveryTime string  `json:"deliveryTime"` // 配送时间
	CouponId     int32   `json:"couponId"`     // 优惠券ID
	Remark       string  `json:"remark"`       // 备注
}

// CreateOrder 创建订单
func CreateOrder(c *gin.Context) {
	userId := middleware.GetUserId(c)

	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(400, "参数错误").WriteJSON(c.Writer)
		return
	}

	// 开始事务
	tx := db.Get().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var orderItems []model.OrderItemModel
	var totalAmount float64

	// 处理订单项
	if len(req.CartIds) > 0 {
		// 从购物车创建
		var cartItems []model.CartModel
		if err := tx.Where("user_id = ? AND id IN ?", userId, req.CartIds).Find(&cartItems).Error; err != nil {
			tx.Rollback()
			utils.Error(500, "获取购物车失败").WriteJSON(c.Writer)
			return
		}

		for _, cartItem := range cartItems {
			var product model.ProductModel
			if err := tx.Where("id = ? AND status = ?", cartItem.ProductId, 1).First(&product).Error; err != nil {
				tx.Rollback()
				utils.Error(404, "商品不存在或已下架").WriteJSON(c.Writer)
				return
			}

			if product.Stock < cartItem.Quantity {
				tx.Rollback()
				utils.Error(400, "商品库存不足").WriteJSON(c.Writer)
				return
			}

			itemAmount := float64(cartItem.Quantity) * product.Price
			totalAmount += itemAmount

			// 获取请求的基础URL
			baseURL := getBaseURL(c)
			productImageURL := ""
			if product.MainImageId > 0 && baseURL != "" {
				productImageURL = utils.FormatImageURL(product.MainImageId, baseURL)
			}
			orderItems = append(orderItems, model.OrderItemModel{
				ProductId:    product.Id,
				ProductName:  product.Name,
				ProductImage: productImageURL,
				Price:        product.Price,
				Quantity:     cartItem.Quantity,
				Amount:       itemAmount,
			})
		}
	} else if req.ProductId > 0 {
		// 立即购买
		var product model.ProductModel
		if err := tx.Where("id = ? AND status = ?", req.ProductId, 1).First(&product).Error; err != nil {
			tx.Rollback()
			utils.Error(404, "商品不存在或已下架").WriteJSON(c.Writer)
			return
		}

		if product.Stock < req.Quantity {
			tx.Rollback()
			utils.Error(400, "商品库存不足").WriteJSON(c.Writer)
			return
		}

		itemAmount := float64(req.Quantity) * product.Price
		totalAmount += itemAmount

		// 获取请求的基础URL
		baseURL := getBaseURL(c)
		productImageURL := ""
		if product.MainImageId > 0 && baseURL != "" {
			productImageURL = utils.FormatImageURL(product.MainImageId, baseURL)
		}
		orderItems = append(orderItems, model.OrderItemModel{
			ProductId:    product.Id,
			ProductName:  product.Name,
			ProductImage: productImageURL,
			Price:        product.Price,
			Quantity:     req.Quantity,
			Amount:       itemAmount,
		})
	} else {
		tx.Rollback()
		utils.Error(400, "请选择商品").WriteJSON(c.Writer)
		return
	}

	// 处理优惠券
	var discountAmount float64
	if req.CouponId > 0 {
		var userCoupon model.UserCouponModel
		if err := tx.Where("user_id = ? AND coupon_id = ? AND status = ?", userId, req.CouponId, 0).First(&userCoupon).Error; err == nil {
			var coupon model.CouponModel
			if err := tx.Where("id = ?", req.CouponId).First(&coupon).Error; err == nil {
				now := time.Now()
				if now.After(coupon.StartTime) && now.Before(coupon.EndTime) && totalAmount >= coupon.MinAmount {
					if coupon.Type == 1 {
						// 满减
						discountAmount = coupon.Amount
					} else if coupon.Type == 2 {
						// 折扣
						discountAmount = totalAmount * (1 - coupon.Amount/100)
					}
				}
			}
		}
	}

	payAmount := totalAmount - discountAmount
	if payAmount < 0 {
		payAmount = 0
	}

	// 处理地址信息
	var receiverName, receiverPhone, receiverAddress string
	if req.DeliveryType == 1 {
		// 外卖需要地址
		var address model.AddressModel
		if err := tx.Where("id = ? AND user_id = ?", req.AddressId, userId).First(&address).Error; err != nil {
			tx.Rollback()
			utils.Error(404, "地址不存在").WriteJSON(c.Writer)
			return
		}
		receiverName = address.Name
		receiverPhone = address.Phone
		receiverAddress = address.Province + address.City + address.District + address.Address
	}

	// 解析配送时间
	var deliveryTime time.Time
	var hasDeliveryTime bool
	if req.DeliveryTime != "" {
		if parsedTime, err := time.Parse("2006-01-02 15:04:05", req.DeliveryTime); err == nil {
			deliveryTime = parsedTime
			hasDeliveryTime = true
		}
	}

	// 创建订单
	orderNo := utils.GenerateOrderNo()
	order := model.OrderModel{
		OrderNo:        orderNo,
		UserId:         userId,
		TotalAmount:    totalAmount,
		DiscountAmount: discountAmount,
		PayAmount:      payAmount,
		DeliveryType:   req.DeliveryType,
		AddressId:      req.AddressId,
		ReceiverName:   receiverName,
		ReceiverPhone:  receiverPhone,
		ReceiverAddress: receiverAddress,
		Remark:        req.Remark,
		CouponId:      req.CouponId,
		Status:        0, // 待支付
		PayStatus:     0,
	}

	// 只有在有配送时间时才设置，避免零值时间导致数据库错误
	if hasDeliveryTime {
		order.DeliveryTime = deliveryTime
	}

	// 使用 Omit 排除零值时间字段，避免 MySQL 报错
	createQuery := tx.Omit("pay_time", "ship_time", "complete_time")
	if !hasDeliveryTime {
		createQuery = createQuery.Omit("delivery_time")
	}

	if err := createQuery.Create(&order).Error; err != nil {
		tx.Rollback()
		utils.Error(500, "创建订单失败").WriteJSON(c.Writer)
		return
	}

	// 创建订单项
	for i := range orderItems {
		orderItems[i].OrderId = order.Id
	}
	if err := tx.Create(&orderItems).Error; err != nil {
		tx.Rollback()
		utils.Error(500, "创建订单项失败").WriteJSON(c.Writer)
		return
	}

	// 扣减库存
	for _, item := range orderItems {
		if err := tx.Model(&model.ProductModel{}).Where("id = ?", item.ProductId).
			Update("stock", gorm.Expr("stock - ?", item.Quantity)).Error; err != nil {
			tx.Rollback()
			utils.Error(500, "扣减库存失败").WriteJSON(c.Writer)
			return
		}
	}

	// 删除购物车（如果从购物车创建）
	if len(req.CartIds) > 0 {
		tx.Where("user_id = ? AND id IN ?", userId, req.CartIds).Delete(&model.CartModel{})
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		utils.Error(500, "提交订单失败").WriteJSON(c.Writer)
		return
	}

	utils.Success(order).WriteJSON(c.Writer)
}

// GetOrderList 获取订单列表
func GetOrderList(c *gin.Context) {
	userId := middleware.GetUserId(c)
	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	query := db.Get().Model(&model.OrderModel{}).Where("user_id = ?", userId)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var orders []model.OrderModel
	offset := (page - 1) * pageSize
	query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&orders)

	utils.Success(map[string]interface{}{
		"list":     orders,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	}).WriteJSON(c.Writer)
}

// GetOrderDetail 获取订单详情
func GetOrderDetail(c *gin.Context) {
	userId := middleware.GetUserId(c)
	id, _ := strconv.ParseInt(c.Param("id"), 10, 32)

	var order model.OrderModel
	if err := db.Get().Where("id = ? AND user_id = ?", id, userId).First(&order).Error; err != nil {
		utils.Error(404, "订单不存在").WriteJSON(c.Writer)
		return
	}

	var orderItems []model.OrderItemModel
	db.Get().Where("order_id = ?", id).Find(&orderItems)

	utils.Success(map[string]interface{}{
		"order":      order,
		"orderItems": orderItems,
	}).WriteJSON(c.Writer)
}

// CancelOrder 取消订单
func CancelOrder(c *gin.Context) {
	userId := middleware.GetUserId(c)
	id, _ := strconv.ParseInt(c.Param("id"), 10, 32)

	var order model.OrderModel
	if err := db.Get().Where("id = ? AND user_id = ?", id, userId).First(&order).Error; err != nil {
		utils.Error(404, "订单不存在").WriteJSON(c.Writer)
		return
	}

	if order.Status != 0 {
		utils.Error(400, "订单状态不允许取消").WriteJSON(c.Writer)
		return
	}

	// 开始事务
	tx := db.Get().Begin()

	// 更新订单状态
	if err := tx.Model(&order).Update("status", 4).Error; err != nil {
		tx.Rollback()
		utils.Error(500, "取消订单失败").WriteJSON(c.Writer)
		return
	}

	// 恢复库存
	var orderItems []model.OrderItemModel
	tx.Where("order_id = ?", id).Find(&orderItems)
	for _, item := range orderItems {
		tx.Model(&model.ProductModel{}).Where("id = ?", item.ProductId).
			Update("stock", gorm.Expr("stock + ?", item.Quantity))
	}

	if err := tx.Commit().Error; err != nil {
		utils.Error(500, "取消订单失败").WriteJSON(c.Writer)
		return
	}

	utils.Success(nil).WriteJSON(c.Writer)
}

// ConfirmOrder 确认收货
func ConfirmOrder(c *gin.Context) {
	userId := middleware.GetUserId(c)
	id, _ := strconv.ParseInt(c.Param("id"), 10, 32)

	var order model.OrderModel
	if err := db.Get().Where("id = ? AND user_id = ?", id, userId).First(&order).Error; err != nil {
		utils.Error(404, "订单不存在").WriteJSON(c.Writer)
		return
	}

	if order.Status != 2 {
		utils.Error(400, "订单状态不允许确认收货").WriteJSON(c.Writer)
		return
	}

	now := time.Now()
	if err := db.Get().Model(&order).Updates(map[string]interface{}{
		"status":       3,
		"complete_time": now,
	}).Error; err != nil {
		utils.Error(500, "确认收货失败").WriteJSON(c.Writer)
		return
	}

	utils.Success(nil).WriteJSON(c.Writer)
}

