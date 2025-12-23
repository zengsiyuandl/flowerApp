package config

import (
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
	AppConfig = &Config{
		WxAppId:    os.Getenv("WX_APPID"),
		WxSecret:   os.Getenv("WX_SECRET"),
		WxPayMchId: os.Getenv("WX_PAY_MCHID"),
		WxPayKey:   os.Getenv("WX_PAY_KEY"),
	}
}

