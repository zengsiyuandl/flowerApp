package api

import (
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
	"wxcloudrun-golang/utils"

	"github.com/gin-gonic/gin"
)

// GetBannerList 获取Banner列表（小程序端使用，仅返回启用的）
func GetBannerList(c *gin.Context) {
	var banners []model.BannerModel
	db.Get().Where("status = ?", 1).
		Order("sort ASC, id DESC").
		Find(&banners)

	utils.Success(banners).WriteJSON(c.Writer)
}
