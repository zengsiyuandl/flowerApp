package api

import (
	"io"
	"mime"
	"mime/multipart"
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

// UploadImageToDB 上传图片到数据库（内部函数，供其他接口复用）
// 返回图片ID和错误信息
func UploadImageToDB(file *multipart.FileHeader, category string) (int32, string) {
	if !isAllowedCategory(category) {
		return 0, "不支持的分类"
	}

	if file.Size > maxFileSizeBytes {
		return 0, "文件大小不能超过2MB"
	}

	ext := filepath.Ext(file.Filename)
	if !isAllowedFormat(ext) {
		return 0, "不支持的文件格式"
	}

	// 读取文件数据
	src, err := file.Open()
	if err != nil {
		return 0, "打开文件失败: " + err.Error()
	}
	defer src.Close()

	data, err := io.ReadAll(src)
	if err != nil {
		return 0, "读取文件失败: " + err.Error()
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
		return 0, "保存图片失败: " + err.Error()
	}

	return imageStorage.Id, ""
}

// AdminUploadImage 上传图片到数据库（独立上传接口，保留用于其他场景）
func AdminUploadImage(c *gin.Context) {
	category := c.DefaultQuery("category", "banner")

	file, err := c.FormFile("file")
	if err != nil {
		utils.Error(400, "获取文件失败: "+err.Error()).WriteJSON(c.Writer)
		return
	}

	imageID, errMsg := UploadImageToDB(file, category)
	if errMsg != "" {
		utils.Error(400, errMsg).WriteJSON(c.Writer)
		return
	}

	imageURL := utils.FormatImageURL(imageID)
	utils.Success(map[string]interface{}{
		"id":  imageID,
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
