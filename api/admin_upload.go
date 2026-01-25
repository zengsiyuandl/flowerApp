package api

import (
	"io"
	"mime"
	"path/filepath"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
	"wxcloudrun-golang/utils"

	"github.com/gin-gonic/gin"
)

const (
	maxFileSizeBytes = 2 * 1024 * 1024 // 2MB
)

var (
	allowedCategories = []string{"banner", "goods", "category", "activity", "user"}
	allowedExtensions = []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
)

// AdminUploadImage 上传图片到数据库
func AdminUploadImage(c *gin.Context) {
	category := c.DefaultQuery("category", "banner")

	if !isAllowedCategory(category) {
		utils.Error(400, "不支持的分类").WriteJSON(c.Writer)
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		utils.Error(400, "获取文件失败: "+err.Error()).WriteJSON(c.Writer)
		return
	}

	if file.Size > maxFileSizeBytes {
		utils.Error(400, "文件大小不能超过2MB").WriteJSON(c.Writer)
		return
	}

	ext := filepath.Ext(file.Filename)
	if !isAllowedFormat(ext) {
		utils.Error(400, "不支持的文件格式").WriteJSON(c.Writer)
		return
	}

	// 读取文件数据
	src, err := file.Open()
	if err != nil {
		utils.Error(500, "打开文件失败: "+err.Error()).WriteJSON(c.Writer)
		return
	}
	defer src.Close()

	data, err := io.ReadAll(src)
	if err != nil {
		utils.Error(500, "读取文件失败: "+err.Error()).WriteJSON(c.Writer)
		return
	}

	// 确定Content-Type
	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// 保存到数据库
	imageStorage := model.ImageStorageModel{
		Category:    category,
		Filename:    file.Filename,
		ContentType: contentType,
		Data:        data,
		Size:        int(file.Size),
	}

	if err := db.Get().Create(&imageStorage).Error; err != nil {
		utils.Error(500, "保存图片失败: "+err.Error()).WriteJSON(c.Writer)
		return
	}

	// 返回图片ID和URL
	imageURL := utils.FormatImageURL(imageStorage.Id)
	utils.Success(map[string]interface{}{
		"id":  imageStorage.Id,
		"url": imageURL,
	}).WriteJSON(c.Writer)
}

// isAllowedFormat 检查文件格式是否允许
func isAllowedFormat(ext string) bool {
	for _, format := range allowedExtensions {
		if ext == format {
			return true
		}
	}
	return false
}

// isAllowedCategory 检查分类是否允许
func isAllowedCategory(category string) bool {
	for _, cat := range allowedCategories {
		if cat == category {
			return true
		}
	}
	return false
}
