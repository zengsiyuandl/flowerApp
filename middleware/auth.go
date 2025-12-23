package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"wxcloudrun-golang/utils"
)

// AuthMiddleware 认证中间件
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从header获取token
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.Error(401, "未授权").WriteJSON(w)
			return
		}

		// 移除 "Bearer " 前缀
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// 解析token
		userId, openId, err := utils.ParseToken(token)
		if err != nil {
			utils.Error(401, "token无效").WriteJSON(w)
			return
		}

		// 将用户信息存储到请求header
		r.Header.Set("X-User-Id", strconv.FormatInt(int64(userId), 10))
		r.Header.Set("X-Open-Id", openId)

		next(w, r)
	}
}

