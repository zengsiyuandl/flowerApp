package utils

import (
	"fmt"
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

