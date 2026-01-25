package api

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
	"wxcloudrun-golang/utils"

	"github.com/gin-gonic/gin"
)

const (
	maxFileSize    = 5 * 1024 * 1024 // 5MB
	allowedFormats = ".jpg,.jpeg,.png,.gif,.webp"
)

// 允许的分类
var allowedCategories = []string{"banner", "goods", "category", "activity", "user"}

// AdminUploadImage 上传图片
func AdminUploadImage(c *gin.Context) {
	// 获取分类参数，默认为banner
	category := c.DefaultQuery("category", "banner")
	
	// 验证分类是否合法
	if !isAllowedCategory(category) {
		utils.Error(400, "不支持的分类，仅支持: banner, goods, category, activity, user").WriteJSON(c.Writer)
		return
	}

	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		utils.Error(400, "获取文件失败: "+err.Error()).WriteJSON(c.Writer)
		return
	}

	// 检查文件大小
	if file.Size > maxFileSize {
		utils.Error(400, "文件大小不能超过5MB").WriteJSON(c.Writer)
		return
	}

	// 检查文件格式
	ext := filepath.Ext(file.Filename)
	if !isAllowedFormat(ext) {
		utils.Error(400, "不支持的文件格式，仅支持: "+allowedFormats).WriteJSON(c.Writer)
		return
	}

	// 获取存储路径（支持环境变量配置，用于持久化存储）
	baseUploadDir := utils.GetStoragePath()
	
	// 创建分类目录
	uploadDir := filepath.Join(baseUploadDir, category)
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		utils.Error(500, "创建上传目录失败: "+err.Error()).WriteJSON(c.Writer)
		return
	}

	// 生成唯一文件名
	filename := generateFilename(ext)
	filePath := filepath.Join(uploadDir, filename)

	// 保存文件
	src, err := file.Open()
	if err != nil {
		utils.Error(500, "打开文件失败: "+err.Error()).WriteJSON(c.Writer)
		return
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		utils.Error(500, "创建文件失败: "+err.Error()).WriteJSON(c.Writer)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		utils.Error(500, "保存文件失败: "+err.Error()).WriteJSON(c.Writer)
		return
	}

	// 返回文件URL
	fileURL := fmt.Sprintf("/images/%s/%s", category, filename)
	utils.Success(map[string]interface{}{
		"url": fileURL,
	}).WriteJSON(c.Writer)
}

// isAllowedFormat 检查文件格式是否允许
func isAllowedFormat(ext string) bool {
	allowed := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	ext = filepath.Ext(ext)
	for _, format := range allowed {
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

// generateFilename 生成唯一文件名
func generateFilename(ext string) string {
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("%d%s", timestamp, ext)
}
