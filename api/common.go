package api

import (
	"net/http"
	"wxcloudrun-golang/utils"

	"github.com/gin-gonic/gin"
)

// IndexHandler 首页
func IndexHandler(c *gin.Context) {
	c.String(http.StatusOK, "Flower Shop API")
}

// HealthHandler 健康检查
func HealthHandler(c *gin.Context) {
	utils.Success(map[string]string{"status": "ok"}).WriteJSON(c.Writer)
}

