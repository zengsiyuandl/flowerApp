package api

import (
	"encoding/json"
	"fmt"
	"os"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
	"wxcloudrun-golang/utils"

	"github.com/gin-gonic/gin"
)

// GetImage 根据图片ID获取图片数据
func GetImage(c *gin.Context) {
	// #region agent log
	logEntry := map[string]interface{}{
		"sessionId":    "debug-session",
		"runId":        "run1",
		"hypothesisId": "A",
		"location":     "image.go:14",
		"message":      "GetImage called",
		"data":         map[string]interface{}{"id": c.Param("id")},
		"timestamp":    fmt.Sprintf("%d", 0),
	}
	if logFile, err := os.OpenFile("d:\\code\\workspace\\.cursor\\debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
		json.NewEncoder(logFile).Encode(logEntry)
		logFile.Close()
	}
	// #endregion

	// 解析图片ID
	imageID, err := utils.ParseInt32(c.Param("id"))
	if err != nil {
		// #region agent log
		logEntry = map[string]interface{}{
			"sessionId":    "debug-session",
			"runId":        "run1",
			"hypothesisId": "A",
			"location":     "image.go:20",
			"message":      "ParseInt32 failed",
			"data":         map[string]interface{}{"error": err.Error(), "param": c.Param("id")},
			"timestamp":    fmt.Sprintf("%d", 0),
		}
		if logFile, err2 := os.OpenFile("d:\\code\\workspace\\.cursor\\debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err2 == nil {
			json.NewEncoder(logFile).Encode(logEntry)
			logFile.Close()
		}
		// #endregion
		utils.Error(400, "无效的图片ID").WriteJSON(c.Writer)
		return
	}

	// #region agent log
	logEntry = map[string]interface{}{
		"sessionId":    "debug-session",
		"runId":        "run1",
		"hypothesisId": "B",
		"location":     "image.go:25",
		"message":      "Before database query",
		"data":         map[string]interface{}{"imageID": imageID},
		"timestamp":    fmt.Sprintf("%d", 0),
	}
	if logFile, err2 := os.OpenFile("d:\\code\\workspace\\.cursor\\debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err2 == nil {
		json.NewEncoder(logFile).Encode(logEntry)
		logFile.Close()
	}
	// #endregion

	// 从数据库查询图片
	var imageStorage model.ImageStorageModel
	dbErr := db.Get().Where("id = ?", imageID).First(&imageStorage).Error
	if dbErr != nil {
		// #region agent log
		logEntry = map[string]interface{}{
			"sessionId":    "debug-session",
			"runId":        "run1",
			"hypothesisId": "B",
			"location":     "image.go:32",
			"message":      "Database query failed",
			"data":         map[string]interface{}{"error": dbErr.Error(), "imageID": imageID},
			"timestamp":    fmt.Sprintf("%d", 0),
		}
		if logFile, err2 := os.OpenFile("d:\\code\\workspace\\.cursor\\debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err2 == nil {
			json.NewEncoder(logFile).Encode(logEntry)
			logFile.Close()
		}
		// #endregion
		utils.Error(404, "图片不存在").WriteJSON(c.Writer)
		return
	}

	// #region agent log
	logEntry = map[string]interface{}{
		"sessionId":    "debug-session",
		"runId":        "run1",
		"hypothesisId": "C",
		"location":     "image.go:40",
		"message":      "After database query",
		"data": map[string]interface{}{
			"id":          imageStorage.Id,
			"contentType": imageStorage.ContentType,
			"size":        imageStorage.Size,
			"dataLength":  len(imageStorage.Data),
		},
		"timestamp": fmt.Sprintf("%d", 0),
	}
	if logFile, err2 := os.OpenFile("d:\\code\\workspace\\.cursor\\debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err2 == nil {
		json.NewEncoder(logFile).Encode(logEntry)
		logFile.Close()
	}
	// #endregion

	// 设置响应头
	if imageStorage.ContentType == "" {
		imageStorage.ContentType = "image/jpeg" // 默认类型
	}
	c.Header("Content-Type", imageStorage.ContentType)
	c.Header("Content-Length", utils.IntToString(imageStorage.Size))
	c.Header("Cache-Control", "public, max-age=31536000") // 缓存1年

	// #region agent log
	logEntry = map[string]interface{}{
		"sessionId":    "debug-session",
		"runId":        "run1",
		"hypothesisId": "D",
		"location":     "image.go:52",
		"message":      "Before c.Data",
		"data": map[string]interface{}{
			"contentType": imageStorage.ContentType,
			"size":        imageStorage.Size,
			"dataLength":  len(imageStorage.Data),
		},
		"timestamp": fmt.Sprintf("%d", 0),
	}
	if logFile, err2 := os.OpenFile("d:\\code\\workspace\\.cursor\\debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err2 == nil {
		json.NewEncoder(logFile).Encode(logEntry)
		logFile.Close()
	}
	// #endregion

	// 检查数据是否为空
	if len(imageStorage.Data) == 0 {
		// #region agent log
		logEntry = map[string]interface{}{
			"sessionId":    "debug-session",
			"runId":        "run1",
			"hypothesisId": "C",
			"location":     "image.go:55",
			"message":      "Image data is empty",
			"data":         map[string]interface{}{"imageID": imageID, "id": imageStorage.Id},
			"timestamp":    fmt.Sprintf("%d", 0),
		}
		if logFile, err2 := os.OpenFile("d:\\code\\workspace\\.cursor\\debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err2 == nil {
			json.NewEncoder(logFile).Encode(logEntry)
			logFile.Close()
		}
		// #endregion
		utils.Error(500, "图片数据为空").WriteJSON(c.Writer)
		return
	}

	// 返回图片数据
	c.Data(200, imageStorage.ContentType, imageStorage.Data)

	// #region agent log
	logEntry = map[string]interface{}{
		"sessionId":    "debug-session",
		"runId":        "run1",
		"hypothesisId": "D",
		"location":     "image.go:70",
		"message":      "After c.Data",
		"data":         map[string]interface{}{"status": "completed"},
		"timestamp":    fmt.Sprintf("%d", 0),
	}
	if logFile, err2 := os.OpenFile("d:\\code\\workspace\\.cursor\\debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err2 == nil {
		json.NewEncoder(logFile).Encode(logEntry)
		logFile.Close()
	}
	// #endregion
}
