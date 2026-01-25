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

// getBaseURL 从请求中获取基础URL
func getBaseURL(c *gin.Context) string {
	scheme := c.GetHeader("X-Forwarded-Proto")
	if scheme == "" {
		scheme = c.Request.URL.Scheme
	}
	host := c.GetHeader("X-Forwarded-Host")
	if host == "" {
		host = c.GetHeader("Host")
	}
	if host == "" {
		host = c.Request.Host
	}
	return utils.GetBaseURLFromRequest(scheme, host)
}

