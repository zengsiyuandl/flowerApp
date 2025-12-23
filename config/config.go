package config

import (
	"fmt"
	"os"
)

// Config 配置结构
type Config struct {
	WxAppId    string
	WxSecret   string
	WxPayMchId string
	WxPayKey   string
}

var AppConfig *Config

// Init 初始化配置
func Init() {
	// 从环境变量读取，如果未设置则使用默认值
	wxAppId := os.Getenv("WX_APPID")
	if wxAppId == "" {
		wxAppId = "wx319f1d79bcfa3aed" // 默认值
		fmt.Println("WARNING: WX_APPID not set, using default value")
	}

	wxSecret := os.Getenv("WX_SECRET")
	if wxSecret == "" {
		wxSecret = "329528f1f8ee6881e3d8a71f99c73ec9" // 默认值
		fmt.Println("WARNING: WX_SECRET not set, using default value")
	}

	AppConfig = &Config{
		WxAppId:    wxAppId,
		WxSecret:   wxSecret,
		WxPayMchId: os.Getenv("WX_PAY_MCHID"),
		WxPayKey:   os.Getenv("WX_PAY_KEY"),
	}

	// 验证必要配置
	if AppConfig.WxAppId == "" {
		panic("WX_APPID is required")
	}
	if AppConfig.WxSecret == "" {
		panic("WX_SECRET is required")
	}

	fmt.Printf("Config initialized: WX_APPID=%s\n", AppConfig.WxAppId)
}

