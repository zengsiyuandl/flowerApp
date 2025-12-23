package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// GenerateToken 生成简单的token（实际生产环境应使用JWT库）
func GenerateToken(userId int32, openId string) string {
	timestamp := time.Now().Unix()
	randomBytes := make([]byte, 16)
	rand.Read(randomBytes)
	randomStr := base64.URLEncoding.EncodeToString(randomBytes)
	return fmt.Sprintf("%d_%s_%d_%s", userId, openId, timestamp, randomStr)
}

// ParseToken 解析token（简单实现，生产环境应使用JWT库）
func ParseToken(token string) (userId int32, openId string, err error) {
	if token == "" {
		return 0, "", fmt.Errorf("token为空")
	}
	
	// 简单实现：token格式为 userId_openId_timestamp_randomStr
	parts := strings.Split(token, "_")
	if len(parts) < 2 {
		return 0, "", fmt.Errorf("token格式错误")
	}
	
	userIdInt, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return 0, "", fmt.Errorf("token解析失败")
	}
	
	return int32(userIdInt), parts[1], nil
}

