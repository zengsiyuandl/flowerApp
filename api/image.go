package api

import (
	"fmt"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
	"wxcloudrun-golang/utils"

	"github.com/gin-gonic/gin"
)

const (
	cacheMaxAge = 31536000 // 缓存1年（秒）
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
		// 记录错误日志
		fmt.Printf("GetImage error: imageID=%d, err=%v\n", imageID, err)
		utils.Error(404, "图片不存在").WriteJSON(c.Writer)
		return
	}

	// 验证图片数据
	if len(imageStorage.Data) == 0 {
		fmt.Printf("GetImage error: imageID=%d, data is empty\n", imageID)
		utils.Error(500, "图片数据为空").WriteJSON(c.Writer)
		return
	}

	// 设置响应头
	c.Header("Content-Type", imageStorage.ContentType)
	c.Header("Content-Length", utils.IntToString(imageStorage.Size))
	c.Header("Cache-Control", fmt.Sprintf("public, max-age=%d", cacheMaxAge))

	// 返回图片数据
	c.Data(200, imageStorage.ContentType, imageStorage.Data)
}
