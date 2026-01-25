package api

import (
	"strconv"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
	"wxcloudrun-golang/utils"

	"github.com/gin-gonic/gin"
)

// AdminGetProductList 管理后台获取商品列表（包含所有状态）
func AdminGetProductList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	keyword := c.Query("keyword")
	categoryId := c.Query("categoryId")
	status := c.Query("status")

	query := db.Get().Model(&model.ProductModel{})

	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}
	if categoryId != "" {
		query = query.Where("category_id = ?", categoryId)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var products []model.ProductModel
	offset := (page - 1) * pageSize
	query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&products)

	// 获取请求的基础URL
	baseURL := getBaseURL(c)

	productDTOs := ToProductDTOList(products, baseURL)
	utils.Success(map[string]interface{}{
		"list":     productDTOs,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	}).WriteJSON(c.Writer)
}

// AdminCreateProduct 创建商品（支持multipart/form-data，可同时上传图片）
func AdminCreateProduct(c *gin.Context) {
	name, subtitle, description, detail, categoryId, stock, sort, status, isHot, isNew, price, originalPrice := parseProductFormData(c)

	if name == "" {
		utils.Error(400, "商品名称不能为空").WriteJSON(c.Writer)
		return
	}

	mainImageID, errMsg := processProductMainImage(c, 0)
	if errMsg != "" {
		utils.Error(400, errMsg).WriteJSON(c.Writer)
		return
	}

	product := buildProductModel(
		0, name, subtitle, description, detail,
		categoryId, int(mainImageID), stock, sort, status, isHot, isNew,
		price, originalPrice,
	)

	if err := db.Get().Create(&product).Error; err != nil {
		utils.Error(500, "创建失败: "+err.Error()).WriteJSON(c.Writer)
		return
	}

	baseURL := getBaseURL(c)
	productDTO := ToProductDTO(product, baseURL)
	utils.Success(productDTO).WriteJSON(c.Writer)
}

// AdminUpdateProduct 更新商品（支持multipart/form-data，可同时上传图片）
func AdminUpdateProduct(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 32)

	var existingProduct model.ProductModel
	if err := db.Get().Where("id = ?", id).First(&existingProduct).Error; err != nil {
		utils.Error(404, "商品不存在").WriteJSON(c.Writer)
		return
	}

	name, subtitle, description, detail, categoryId, stock, sort, status, isHot, isNew, price, originalPrice := parseProductFormData(c)

	if name == "" {
		utils.Error(400, "商品名称不能为空").WriteJSON(c.Writer)
		return
	}

	mainImageID, errMsg := processProductMainImage(c, existingProduct.MainImageId)
	if errMsg != "" {
		utils.Error(400, errMsg).WriteJSON(c.Writer)
		return
	}

	product := buildProductModel(
		int32(id), name, subtitle, description, detail,
		categoryId, int(mainImageID), stock, sort, status, isHot, isNew,
		price, originalPrice,
	)

	if err := db.Get().Save(&product).Error; err != nil {
		utils.Error(500, "更新失败: "+err.Error()).WriteJSON(c.Writer)
		return
	}

	baseURL := getBaseURL(c)
	productDTO := ToProductDTO(product, baseURL)
	utils.Success(productDTO).WriteJSON(c.Writer)
}

// AdminDeleteProduct 删除商品
func AdminDeleteProduct(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 32)

	if err := db.Get().Where("id = ?", id).Delete(&model.ProductModel{}).Error; err != nil {
		utils.Error(500, "删除失败: "+err.Error()).WriteJSON(c.Writer)
		return
	}

	utils.Success(nil).WriteJSON(c.Writer)
}

// AdminGetProductImages 获取商品图片列表
func AdminGetProductImages(c *gin.Context) {
	productId, _ := strconv.ParseInt(c.Param("productId"), 10, 32)

	var images []model.ProductImageModel
	db.Get().Where("product_id = ?", productId).Order("sort ASC").Find(&images)

	// 获取请求的基础URL
	baseURL := getBaseURL(c)

	imageDTOs := ToProductImageDTOList(images, baseURL)
	utils.Success(imageDTOs).WriteJSON(c.Writer)
}

// AdminAddProductImage 添加商品图片
func AdminAddProductImage(c *gin.Context) {
	productId, _ := strconv.ParseInt(c.Param("productId"), 10, 32)

	var image model.ProductImageModel
	if err := c.ShouldBindJSON(&image); err != nil {
		utils.Error(400, "参数错误").WriteJSON(c.Writer)
		return
	}

	image.ProductId = int32(productId)
	if err := db.Get().Create(&image).Error; err != nil {
		utils.Error(500, "添加失败").WriteJSON(c.Writer)
		return
	}

	// 获取请求的基础URL
	baseURL := getBaseURL(c)

	imageDTO := ToProductImageDTO(image, baseURL)
	utils.Success(imageDTO).WriteJSON(c.Writer)
}

// AdminDeleteProductImage 删除商品图片
func AdminDeleteProductImage(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 32)

	if err := db.Get().Where("id = ?", id).Delete(&model.ProductImageModel{}).Error; err != nil {
		utils.Error(500, "删除失败").WriteJSON(c.Writer)
		return
	}

	utils.Success(nil).WriteJSON(c.Writer)
}

