package config

import (
	"fmt"
	"os"
)

// Config 配置结构
type Config struct {
	WxAppId            string
	WxSecret           string
	WxPaySubMchId      string // 微信支付子商户号
	WxCloudEnvId       string // 云托管环境ID
	WxCloudServiceName string // 云托管服务名称
}

var AppConfig *Config

// Init 初始化配置
func Init() {
	wxAppId := os.Getenv("WX_APPID")
	if wxAppId == "" {
		panic("WX_APPID environment variable is required")
	}

	wxSecret := os.Getenv("WX_SECRET")
	if wxSecret == "" {
		panic("WX_SECRET environment variable is required")
	}

	wxPaySubMchId := os.Getenv("WX_PAY_SUB_MCH_ID")
	if wxPaySubMchId == "" {
		panic("WX_PAY_SUB_MCH_ID environment variable is required")
	}

	wxCloudEnvId := os.Getenv("WX_CLOUD_ENV_ID")
	if wxCloudEnvId == "" {
		panic("WX_CLOUD_ENV_ID environment variable is required")
	}

	wxCloudServiceName := os.Getenv("WX_CLOUD_SERVICE_NAME")
	if wxCloudServiceName == "" {
		panic("WX_CLOUD_SERVICE_NAME environment variable is required")
	}

	AppConfig = &Config{
		WxAppId:            wxAppId,
		WxSecret:           wxSecret,
		WxPaySubMchId:      wxPaySubMchId,
		WxCloudEnvId:       wxCloudEnvId,
		WxCloudServiceName: wxCloudServiceName,
	}

	fmt.Println("Config initialized successfully")
}

