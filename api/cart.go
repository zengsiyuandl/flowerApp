package api

import (
	"strconv"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
	"wxcloudrun-golang/middleware"
	"wxcloudrun-golang/utils"

	"github.com/gin-gonic/gin"
)

// GetCart 获取购物车
func GetCart(c *gin.Context) {
	userId := middleware.GetUserId(c)

	var cartItems []model.CartModel
	db.Get().Where("user_id = ?", userId).Find(&cartItems)

	// 获取商品详情
	type CartItemWithProduct struct {
		model.CartModel
		Product model.ProductModel `json:"product"`
	}

	var result []CartItemWithProduct
	for _, item := range cartItems {
		var product model.ProductModel
		if err := db.Get().Where("id = ?", item.ProductId).First(&product).Error; err == nil {
			result = append(result, CartItemWithProduct{
				CartModel: item,
				Product:   product,
			})
		}
	}

	utils.Success(result).WriteJSON(c.Writer)
}

// AddToCartRequest 添加到购物车请求
type AddToCartRequest struct {
	ProductId int32 `json:"productId" binding:"required"`
	Quantity  int   `json:"quantity" binding:"required,min=1"`
}

// AddToCart 添加到购物车
func AddToCart(c *gin.Context) {
	userId := middleware.GetUserId(c)

	var req AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(400, "参数错误").WriteJSON(c.Writer)
		return
	}

	// 检查商品是否存在
	var product model.ProductModel
	if err := db.Get().Where("id = ? AND status = ?", req.ProductId, 1).First(&product).Error; err != nil {
		utils.Error(404, "商品不存在或已下架").WriteJSON(c.Writer)
		return
	}

	// 检查库存
	if product.Stock < req.Quantity {
		utils.Error(400, "库存不足").WriteJSON(c.Writer)
		return
	}

	// 检查购物车中是否已存在
	var cartItem model.CartModel
	err := db.Get().Where("user_id = ? AND product_id = ?", userId, req.ProductId).First(&cartItem).Error

	if err != nil {
		// 不存在，创建新记录
		cartItem = model.CartModel{
			UserId:    userId,
			ProductId: req.ProductId,
			Quantity:  req.Quantity,
		}
		if err := db.Get().Create(&cartItem).Error; err != nil {
			utils.Error(500, "添加失败").WriteJSON(c.Writer)
			return
		}
	} else {
		// 存在，更新数量
		newQuantity := cartItem.Quantity + req.Quantity
		if newQuantity > product.Stock {
			utils.Error(400, "库存不足").WriteJSON(c.Writer)
			return
		}
		cartItem.Quantity = newQuantity
		if err := db.Get().Save(&cartItem).Error; err != nil {
			utils.Error(500, "更新失败").WriteJSON(c.Writer)
			return
		}
	}

	utils.Success(cartItem).WriteJSON(c.Writer)
}

// UpdateCartRequest 更新购物车请求
type UpdateCartRequest struct {
	Quantity int `json:"quantity" binding:"required,min=1"`
}

// UpdateCart 更新购物车
func UpdateCart(c *gin.Context) {
	userId := middleware.GetUserId(c)
	id, _ := strconv.ParseInt(c.Param("id"), 10, 32)

	var req UpdateCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(400, "参数错误").WriteJSON(c.Writer)
		return
	}

	var cartItem model.CartModel
	if err := db.Get().Where("id = ? AND user_id = ?", id, userId).First(&cartItem).Error; err != nil {
		utils.Error(404, "购物车项不存在").WriteJSON(c.Writer)
		return
	}

	// 检查库存
	var product model.ProductModel
	if err := db.Get().Where("id = ?", cartItem.ProductId).First(&product).Error; err != nil {
		utils.Error(404, "商品不存在").WriteJSON(c.Writer)
		return
	}

	if req.Quantity > product.Stock {
		utils.Error(400, "库存不足").WriteJSON(c.Writer)
		return
	}

	cartItem.Quantity = req.Quantity
	if err := db.Get().Save(&cartItem).Error; err != nil {
		utils.Error(500, "更新失败").WriteJSON(c.Writer)
		return
	}

	utils.Success(cartItem).WriteJSON(c.Writer)
}

// DeleteCart 删除购物车项
func DeleteCart(c *gin.Context) {
	userId := middleware.GetUserId(c)
	id, _ := strconv.ParseInt(c.Param("id"), 10, 32)

	if err := db.Get().Where("id = ? AND user_id = ?", id, userId).Delete(&model.CartModel{}).Error; err != nil {
		utils.Error(500, "删除失败").WriteJSON(c.Writer)
		return
	}

	utils.Success(nil).WriteJSON(c.Writer)
}

