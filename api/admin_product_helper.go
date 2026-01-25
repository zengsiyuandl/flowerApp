package api

import (
	"strconv"
	"wxcloudrun-golang/db/model"
	"wxcloudrun-golang/utils"

	"github.com/gin-gonic/gin"
)

// parseProductFormData 解析商品表单数据
func parseProductFormData(c *gin.Context) (
	name, subtitle, description, detail string,
	categoryId, stock, sort, status, isHot, isNew int,
	price, originalPrice float64,
) {
	name = c.PostForm("name")
	subtitle = c.PostForm("subtitle")
	description = c.PostForm("description")
	detail = c.PostForm("detail")
	categoryIdStr := c.PostForm("categoryId")
	stockStr := c.PostForm("stock")
	sortStr := c.PostForm("sort")
	statusStr := c.PostForm("status")
	isHotStr := c.PostForm("isHot")
	isNewStr := c.PostForm("isNew")
	priceStr := c.PostForm("price")
	originalPriceStr := c.PostForm("originalPrice")

	categoryId, _ = strconv.Atoi(categoryIdStr)
	stock, _ = strconv.Atoi(stockStr)
	sort, _ = strconv.Atoi(sortStr)
	status, _ = strconv.Atoi(statusStr)
	isHot, _ = strconv.Atoi(isHotStr)
	isNew, _ = strconv.Atoi(isNewStr)
	price, _ = strconv.ParseFloat(priceStr, 64)
	originalPrice, _ = strconv.ParseFloat(originalPriceStr, 64)

	return
}

// processProductMainImage 处理商品主图上传
func processProductMainImage(c *gin.Context, existingImageID int32) (int32, string) {
	if file, err := c.FormFile("mainImage"); err == nil {
		imageID, errMsg := UploadImageToDB(file, "goods")
		if errMsg != "" {
			return 0, errMsg
		}
		return imageID, ""
	}

	mainImageIdStr := c.PostForm("mainImageId")
	if mainImageIdStr != "" {
		if id, err := utils.ParseInt32(mainImageIdStr); err == nil {
			return id, ""
		}
	}

	return existingImageID, ""
}

// buildProductModel 构建商品模型
func buildProductModel(
	id int32,
	name, subtitle, description, detail string,
	categoryId, mainImageId, stock, sort, status, isHot, isNew int,
	price, originalPrice float64,
) model.ProductModel {
	return model.ProductModel{
		Id:            id,
		Name:          name,
		Subtitle:      subtitle,
		CategoryId:    int32(categoryId),
		MainImageId:   int32(mainImageId),
		Price:         price,
		OriginalPrice: originalPrice,
		Stock:         stock,
		Description:   description,
		Detail:        detail,
		Sort:          sort,
		Status:        status,
		IsHot:         isHot,
		IsNew:         isNew,
	}
}
