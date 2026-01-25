package api

import (
	"strconv"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
	"wxcloudrun-golang/utils"

	"github.com/gin-gonic/gin"
)

// GetProductList 获取商品列表
func GetProductList(c *gin.Context) {
	categoryId := c.Query("categoryId")
	keyword := c.Query("keyword")
	isHot := c.Query("isHot")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	query := db.Get().Model(&model.ProductModel{}).Where("status = ?", 1)

	if categoryId != "" {
		query = query.Where("category_id = ?", categoryId)
	}

	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}

	if isHot == "1" {
		query = query.Where("is_hot = ?", 1)
	}

	var total int64
	query.Count(&total)

	var products []model.ProductModel
	offset := (page - 1) * pageSize
	query.Order("sort DESC, id DESC").Offset(offset).Limit(pageSize).Find(&products)

	productDTOs := ToProductDTOList(products)
	utils.Success(map[string]interface{}{
		"list":     productDTOs,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	}).WriteJSON(c.Writer)
}

// GetProductDetail 获取商品详情
func GetProductDetail(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 32)

	var product model.ProductModel
	if err := db.Get().Where("id = ?", id).First(&product).Error; err != nil {
		utils.Error(404, "商品不存在").WriteJSON(c.Writer)
		return
	}

	// 获取商品图片
	var images []model.ProductImageModel
	db.Get().Where("product_id = ?", id).Order("sort ASC").Find(&images)

	utils.Success(map[string]interface{}{
		"product": product,
		"images":  images,
	}).WriteJSON(c.Writer)
}

// GetCategoryList 获取分类列表
func GetCategoryList(c *gin.Context) {
	var categories []model.CategoryModel
	db.Get().Where("status = ?", 1).Order("sort DESC, id ASC").Find(&categories)

	categoryDTOs := ToCategoryDTOList(categories)
	utils.Success(categoryDTOs).WriteJSON(c.Writer)
}

