package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"wxcloudrun-golang/config"
)

// WxLoginResponse 微信登录响应
type WxLoginResponse struct {
	OpenId     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionId    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

// Code2Session 微信登录 code2Session
func Code2Session(code string) (*WxLoginResponse, error) {
	// 验证配置
	if config.AppConfig == nil {
		return nil, fmt.Errorf("配置未初始化")
	}
	if config.AppConfig.WxAppId == "" {
		return nil, fmt.Errorf("WX_APPID 未配置")
	}
	if config.AppConfig.WxSecret == "" {
		return nil, fmt.Errorf("WX_SECRET 未配置")
	}

	url := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		config.AppConfig.WxAppId, config.AppConfig.WxSecret, code)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("请求微信API失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	var result WxLoginResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	if result.ErrCode != 0 {
		return nil, fmt.Errorf("微信登录失败: %s (errcode: %d)", result.ErrMsg, result.ErrCode)
	}

	return &result, nil
}

