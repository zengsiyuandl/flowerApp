package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strconv"
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

// FormatImageURL 格式化图片URL
func FormatImageURL(imageID int32) string {
	return "/images/" + strconv.FormatInt(int64(imageID), 10)
}

