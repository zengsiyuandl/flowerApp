package api

import (
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
	"wxcloudrun-golang/utils"

	"github.com/gin-gonic/gin"
)

// GetImage 根据图片ID获取图片数据
func GetImage(c *gin.Context) {
	// 解析图片ID
	imageID, err := utils.ParseInt32(c.Param("id"))
	if err != nil {
		utils.Error(400, "无效的图片ID").WriteJSON(c.Writer)
		return
	}

	// 从数据库查询图片
	var imageStorage model.ImageStorageModel
	if err := db.Get().Where("id = ?", imageID).First(&imageStorage).Error; err != nil {
		utils.Error(404, "图片不存在").WriteJSON(c.Writer)
		return
	}

	// 设置响应头
	c.Header("Content-Type", imageStorage.ContentType)
	c.Header("Content-Length", utils.IntToString(imageStorage.Size))
	c.Header("Cache-Control", "public, max-age=31536000") // 缓存1年

	// 返回图片数据
	c.Data(200, imageStorage.ContentType, imageStorage.Data)
}
