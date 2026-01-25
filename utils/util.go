package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// GenerateOrderNo 生成订单号
func GenerateOrderNo() string {
	return fmt.Sprintf("ORD%d%06d", time.Now().Unix(), time.Now().Nanosecond()%1000000)
}

// GenerateRechargeNo 生成充值订单号
func GenerateRechargeNo() string {
	return fmt.Sprintf("RCH%d%06d", time.Now().Unix(), time.Now().Nanosecond()%1000000)
}

// FormatTime 格式化时间
func FormatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// GenerateNonceStr 生成随机字符串（32位）
func GenerateNonceStr() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// ParseInt32 解析字符串为int32
func ParseInt32(s string) (int32, error) {
	val, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(val), nil
}

// IntToString 将int转换为string
func IntToString(i int) string {
	return strconv.Itoa(i)
}

// FormatImageURL 格式化图片URL（返回完整URL）
func FormatImageURL(imageID int32, baseURL string) string {
	// 确保baseURL不以/结尾
	if len(baseURL) > 0 && baseURL[len(baseURL)-1] == '/' {
		baseURL = baseURL[:len(baseURL)-1]
	}
	return baseURL + "/images/" + strconv.FormatInt(int64(imageID), 10)
}

// GetBaseURLFromRequest 从请求中获取基础URL
func GetBaseURLFromRequest(scheme, host string) string {
	if host == "" {
		return ""
	}
	
	// 清理host，移除可能包含的协议信息
	host = strings.TrimSpace(host)
	if strings.HasPrefix(host, "http://") {
		host = strings.TrimPrefix(host, "http://")
		scheme = "http"
	} else if strings.HasPrefix(host, "https://") {
		host = strings.TrimPrefix(host, "https://")
		scheme = "https"
	}
	
	// 验证和清理scheme，确保只能是http或https
	scheme = strings.ToLower(strings.TrimSpace(scheme))
	if scheme != "http" && scheme != "https" {
		scheme = "https" // 默认使用https（云托管通常使用https）
	}
	
	return scheme + "://" + host
}

