package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
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

// GetStoragePath 获取存储路径
// 优先使用环境变量 STORAGE_PATH，如果未设置则使用默认路径 "images"
// 在云托管中，可以设置 STORAGE_PATH 为持久化存储路径，如 "/data/images"
func GetStoragePath() string {
	storagePath := os.Getenv("STORAGE_PATH")
	if storagePath == "" {
		// 默认使用项目目录下的 images 文件夹（本地开发）
		storagePath = "images"
	}
	// 确保路径是绝对路径或相对于工作目录的路径
	if !filepath.IsAbs(storagePath) {
		// 如果是相对路径，确保相对于当前工作目录
		wd, err := os.Getwd()
		if err == nil {
			storagePath = filepath.Join(wd, storagePath)
		}
	}
	return storagePath
}

