package utils

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

// MockPaymentParams 模拟支付参数（打桩实现）
type MockPaymentParams struct {
	TimeStamp string `json:"timeStamp"`
	NonceStr  string `json:"nonceStr"`
	Package   string `json:"package"`
	SignType  string `json:"signType"`
	PaySign   string `json:"paySign"`
}

// GenerateMockPaymentParams 生成模拟支付参数
func GenerateMockPaymentParams(orderNo string, amount float64) *MockPaymentParams {
	rand.Seed(time.Now().UnixNano())
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	nonceStr := fmt.Sprintf("%d%d", time.Now().UnixNano(), rand.Intn(10000))
	packageStr := fmt.Sprintf("prepay_id=MOCK_%s", orderNo)
	signType := "MD5"

	// 生成模拟签名
	signStr := fmt.Sprintf("appId=%s&nonceStr=%s&package=%s&signType=%s&timeStamp=%s&key=MOCK_KEY",
		"wx319f1d79bcfa3aed", nonceStr, packageStr, signType, timestamp)
	paySign := fmt.Sprintf("%x", md5.Sum([]byte(signStr)))

	return &MockPaymentParams{
		TimeStamp: timestamp,
		NonceStr:  nonceStr,
		Package:   packageStr,
		SignType:  signType,
		PaySign:   paySign,
	}
}

// TODO: 后续接入真实微信支付时，在此文件中实现真实的统一下单逻辑

