package router

import (
	"wxcloudrun-golang/api"
	"wxcloudrun-golang/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter 设置路由
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// CORS中间件
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-User-Id, X-Open-Id")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		c.Next()
	})

	// 健康检查
	r.GET("/", api.IndexHandler)
	r.GET("/health", api.HealthHandler)

	// API路由组
	apiGroup := r.Group("/api")
	{
		// 用户相关（部分接口不需要认证）
		userGroup := apiGroup.Group("/user")
		{
			userGroup.POST("/login", api.UserLogin)
			userGroup.GET("/info", middleware.AuthRequired(), api.GetUserInfo)
			userGroup.PUT("/info", middleware.AuthRequired(), api.UpdateUserInfo)

			// 地址管理
			userGroup.GET("/address", middleware.AuthRequired(), api.GetAddressList)
			userGroup.POST("/address", middleware.AuthRequired(), api.AddAddress)
			userGroup.PUT("/address/:id", middleware.AuthRequired(), api.UpdateAddress)
			userGroup.DELETE("/address/:id", middleware.AuthRequired(), api.DeleteAddress)
		}

		// 商品相关
		productGroup := apiGroup.Group("/product")
		{
			productGroup.GET("/list", api.GetProductList)
			productGroup.GET("/:id", api.GetProductDetail)
		}

		// 分类相关
		apiGroup.GET("/category/list", api.GetCategoryList)

		// 购物车相关
		cartGroup := apiGroup.Group("/cart")
		cartGroup.Use(middleware.AuthRequired())
		{
			cartGroup.GET("", api.GetCart)
			cartGroup.POST("", api.AddToCart)
			cartGroup.PUT("/:id", api.UpdateCart)
			cartGroup.DELETE("/:id", api.DeleteCart)
		}

		// 订单相关
		orderGroup := apiGroup.Group("/order")
		orderGroup.Use(middleware.AuthRequired())
		{
			orderGroup.POST("/create", api.CreateOrder)
			orderGroup.GET("/list", api.GetOrderList)
			orderGroup.GET("/:id", api.GetOrderDetail)
			orderGroup.PUT("/:id/cancel", api.CancelOrder)
			orderGroup.PUT("/:id/confirm", api.ConfirmOrder)
		}

		// 支付相关
		paymentGroup := apiGroup.Group("/payment")
		paymentGroup.Use(middleware.AuthRequired())
		{
			paymentGroup.POST("/create", api.CreatePayment)
			paymentGroup.POST("/notify", api.PaymentNotify)
			paymentGroup.GET("/:orderId/status", api.GetPaymentStatus)
		}

		// 配送相关
		deliveryGroup := apiGroup.Group("/delivery")
		deliveryGroup.Use(middleware.AuthRequired())
		{
			deliveryGroup.GET("/:orderId/track", api.GetDeliveryTrack)
		}

		// 活动相关
		activityGroup := apiGroup.Group("/activity")
		{
			activityGroup.GET("/list", api.GetActivityList)
			activityGroup.GET("/:id", api.GetActivityDetail)
		}

		// 优惠券相关
		couponGroup := apiGroup.Group("/coupon")
		couponGroup.Use(middleware.AuthRequired())
		{
			couponGroup.GET("/list", api.GetCouponList)
			couponGroup.POST("/receive", api.ReceiveCoupon)
		}

		// 抽奖相关
		lotteryGroup := apiGroup.Group("/lottery")
		lotteryGroup.Use(middleware.AuthRequired())
		{
			lotteryGroup.GET("/info", api.GetLotteryInfo)
			lotteryGroup.POST("/draw", api.DrawLottery)
		}

		// 储值相关
		rechargeGroup := apiGroup.Group("/recharge")
		rechargeGroup.Use(middleware.AuthRequired())
		{
			rechargeGroup.POST("/create", api.CreateRecharge)
			rechargeGroup.GET("/list", api.GetRechargeList)
		}

		// ========== 管理后台API ==========
		// 注意：生产环境应该添加管理员权限验证
		adminGroup := apiGroup.Group("/admin")
		{
			// 商品管理
			adminGroup.GET("/product/list", api.AdminGetProductList)
			adminGroup.POST("/product", api.AdminCreateProduct)
			adminGroup.PUT("/product/:id", api.AdminUpdateProduct)
			adminGroup.DELETE("/product/:id", api.AdminDeleteProduct)
			adminGroup.GET("/product/:productId/images", api.AdminGetProductImages)
			adminGroup.POST("/product/:productId/images", api.AdminAddProductImage)
			adminGroup.DELETE("/product/image/:id", api.AdminDeleteProductImage)

			// 分类管理
			adminGroup.GET("/category/list", api.AdminGetCategoryList)
			adminGroup.POST("/category", api.AdminCreateCategory)
			adminGroup.PUT("/category/:id", api.AdminUpdateCategory)
			adminGroup.DELETE("/category/:id", api.AdminDeleteCategory)

			// 活动管理
			adminGroup.GET("/activity/list", api.AdminGetActivityList)
			adminGroup.POST("/activity", api.AdminCreateActivity)
			adminGroup.PUT("/activity/:id", api.AdminUpdateActivity)
			adminGroup.DELETE("/activity/:id", api.AdminDeleteActivity)

			// 优惠券管理
			adminGroup.GET("/coupon/list", api.AdminGetCouponList)
			adminGroup.POST("/coupon", api.AdminCreateCoupon)
			adminGroup.PUT("/coupon/:id", api.AdminUpdateCoupon)
			adminGroup.DELETE("/coupon/:id", api.AdminDeleteCoupon)

			// 订单管理
			adminGroup.GET("/order/list", api.AdminGetOrderList)
			adminGroup.PUT("/order/:id/ship", api.AdminShipOrder)

			// 用户管理
			adminGroup.GET("/user/list", api.AdminGetUserList)
			adminGroup.POST("/user/:userId/recharge", api.AdminRechargeUser)

			// 抽奖管理
			adminGroup.GET("/lottery/list", api.AdminGetLotteryList)
			adminGroup.POST("/lottery", api.AdminCreateLottery)
			adminGroup.PUT("/lottery/:id", api.AdminUpdateLottery)
			adminGroup.GET("/lottery/records", api.AdminGetLotteryRecords)
		}
	}

	return r
}

