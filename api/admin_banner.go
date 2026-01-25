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

// AdminCreateBanner 创建Banner
func AdminCreateBanner(c *gin.Context) {
	var banner model.BannerModel
	if err := c.ShouldBindJSON(&banner); err != nil {
		utils.Error(400, "参数错误: "+err.Error()).WriteJSON(c.Writer)
		return
	}

	if err := db.Get().Create(&banner).Error; err != nil {
		utils.Error(500, "创建失败: "+err.Error()).WriteJSON(c.Writer)
		return
	}

	bannerDTO := ToBannerDTO(banner)
	utils.Success(bannerDTO).WriteJSON(c.Writer)
}

// AdminUpdateBanner 更新Banner
func AdminUpdateBanner(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 32)

	var banner model.BannerModel
	if err := db.Get().Where("id = ?", id).First(&banner).Error; err != nil {
		utils.Error(404, "Banner不存在").WriteJSON(c.Writer)
		return
	}

	if err := c.ShouldBindJSON(&banner); err != nil {
		utils.Error(400, "参数错误: "+err.Error()).WriteJSON(c.Writer)
		return
	}

	banner.Id = int32(id)
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
