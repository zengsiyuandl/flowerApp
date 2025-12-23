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
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		config.AppConfig.WxAppId, config.AppConfig.WxSecret, code)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result WxLoginResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if result.ErrCode != 0 {
		return nil, fmt.Errorf("微信登录失败: %s", result.ErrMsg)
	}

	return &result, nil
}

