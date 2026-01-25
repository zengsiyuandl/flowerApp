package middleware

import (
	"strings"
	"wxcloudrun-golang/utils"

	"github.com/gin-gonic/gin"
)

// AuthRequired 认证中间件（Gin版本）
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从header获取token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.Error(401, "未授权").WriteJSON(c.Writer)
			c.Abort()
			return
		}

		// 移除 "Bearer " 前缀
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// 解析token（简化实现，实际应使用JWT）
		// 这里暂时从token中提取用户信息，实际应验证token有效性
		if token == "" {
			utils.Error(401, "token无效").WriteJSON(c.Writer)
			c.Abort()
			return
		}

		// 优先从Header获取OpenID（微信云托管自动提供）
		openId := c.GetHeader("x-wx-openid")
		
		// 解析token获取用户ID
		userId, tokenOpenId, err := utils.ParseToken(token)
		if err != nil {
			utils.Error(401, "token无效").WriteJSON(c.Writer)
			c.Abort()
			return
		}

		// 如果Header中没有OpenID，使用token中的OpenID
		if openId == "" {
			openId = tokenOpenId
		}

		// 将用户信息存储到上下文
		c.Set("userId", int32(userId))
		c.Set("openId", openId)
		c.Set("token", token)

		c.Next()
	}
}

// GetUserId 从上下文获取用户ID
func GetUserId(c *gin.Context) int32 {
	if userId, exists := c.Get("userId"); exists {
		return userId.(int32)
	}
	return 0
}

// GetOpenId 从上下文获取OpenId
func GetOpenId(c *gin.Context) string {
	if openId, exists := c.Get("openId"); exists {
		return openId.(string)
	}
	return ""
}

