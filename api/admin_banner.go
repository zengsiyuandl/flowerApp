package api

import (
	"strconv"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
	"wxcloudrun-golang/utils"

	"github.com/gin-gonic/gin"
)

// AdminGetBannerList 管理后台获取Banner列表
func AdminGetBannerList(c *gin.Context) {
	var banners []model.BannerModel
	db.Get().Order("sort ASC, id DESC").Find(&banners)

	bannerDTOs := ToBannerDTOList(banners)
	utils.Success(bannerDTOs).WriteJSON(c.Writer)
}

// AdminCreateBanner 创建Banner（支持multipart/form-data，可同时上传图片）
func AdminCreateBanner(c *gin.Context) {
	// 解析表单数据
	title, linkType, linkValue, sort, status := parseBannerFormData(c)

	// 验证必填字段
	if title == "" {
		utils.Error(400, "标题不能为空").WriteJSON(c.Writer)
		return
	}

	// 处理图片上传
	imageID, errMsg := processBannerImage(c, 0)
	if errMsg != "" {
		utils.Error(400, errMsg).WriteJSON(c.Writer)
		return
	}

	// 创建Banner
	banner := buildBannerModel(0, title, imageID, linkType, linkValue, sort, status)

	if err := db.Get().Create(&banner).Error; err != nil {
		utils.Error(500, "创建失败: "+err.Error()).WriteJSON(c.Writer)
		return
	}

	bannerDTO := ToBannerDTO(banner)
	utils.Success(bannerDTO).WriteJSON(c.Writer)
}

// AdminUpdateBanner 更新Banner（支持multipart/form-data，可同时上传图片）
func AdminUpdateBanner(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 32)

	// 检查Banner是否存在
	var existingBanner model.BannerModel
	if err := db.Get().Where("id = ?", id).First(&existingBanner).Error; err != nil {
		utils.Error(404, "Banner不存在").WriteJSON(c.Writer)
		return
	}

	// 解析表单数据
	title, linkType, linkValue, sort, status := parseBannerFormData(c)

	// 验证必填字段
	if title == "" {
		utils.Error(400, "标题不能为空").WriteJSON(c.Writer)
		return
	}

	// 处理图片上传
	imageID, errMsg := processBannerImage(c, existingBanner.ImageId)
	if errMsg != "" {
		utils.Error(400, errMsg).WriteJSON(c.Writer)
		return
	}

	// 更新Banner
	banner := buildBannerModel(int32(id), title, imageID, linkType, linkValue, sort, status)

	if err := db.Get().Save(&banner).Error; err != nil {
		utils.Error(500, "更新失败: "+err.Error()).WriteJSON(c.Writer)
		return
	}

	bannerDTO := ToBannerDTO(banner)
	utils.Success(bannerDTO).WriteJSON(c.Writer)
}

// AdminDeleteBanner 删除Banner
func AdminDeleteBanner(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 32)

	if err := db.Get().Where("id = ?", id).Delete(&model.BannerModel{}).Error; err != nil {
		utils.Error(500, "删除失败: "+err.Error()).WriteJSON(c.Writer)
		return
	}

	utils.Success(nil).WriteJSON(c.Writer)
}
