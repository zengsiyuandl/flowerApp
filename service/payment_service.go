package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	// WeChatPayAPIBase 微信支付开放接口服务基础地址
	WeChatPayAPIBase = "http://api.weixin.qq.com/_/pay"
	// CallbackTypeCloudRun 回调类型：云托管
	CallbackTypeCloudRun = 2
)

// PaymentService 支付服务
type PaymentService struct{}

// NewPaymentService 创建支付服务实例
func NewPaymentService() *PaymentService {
	return &PaymentService{}
}

// UnifiedOrderRequest 统一下单请求
type UnifiedOrderRequest struct {
	OpenID        string `json:"openid"`         // 用户openid
	Body          string `json:"body"`           // 商品描述
	OutTradeNo    string `json:"out_trade_no"`   // 商户订单号
	SpbillCreateIP string `json:"spbill_create_ip"` // 终端IP
	SubMchID      string `json:"sub_mch_id"`     // 子商户号
	TotalFee      int    `json:"total_fee"`      // 订单总金额（分）
	EnvID         string `json:"env_id"`         // 云托管环境ID
	CallbackType  int    `json:"callback_type"` // 回调类型
	Container     struct {
		Service string `json:"service"` // 服务名称
		Path    string `json:"path"`    // 回调路径
	} `json:"container"`
}

// UnifiedOrderResponse 统一下单响应
type UnifiedOrderResponse struct {
	ErrCode  int    `json:"errcode"`
	ErrMsg   string `json:"errmsg"`
	RespData struct {
		ReturnCode string `json:"return_code"`
		ReturnMsg  string `json:"return_msg"`
		ResultCode string `json:"result_code"`
		Payment    struct {
			AppID     string `json:"appId"`
			TimeStamp string `json:"timeStamp"`
			NonceStr  string `json:"nonceStr"`
			Package   string `json:"package"`
			SignType  string `json:"signType"`
			PaySign   string `json:"paySign"`
		} `json:"payment"`
		PrepayID string `json:"prepay_id"`
	} `json:"respdata"`
}

// QueryOrderRequest 查询订单请求
type QueryOrderRequest struct {
	SubMchID     string `json:"sub_mch_id"`
	OutTradeNo   string `json:"out_trade_no,omitempty"`   // 商户订单号
	TransactionID string `json:"transaction_id,omitempty"` // 微信订单号
	NonceStr     string `json:"nonce_str"`
}

// QueryOrderResponse 查询订单响应
type QueryOrderResponse struct {
	ErrCode  int    `json:"errcode"`
	ErrMsg   string `json:"errmsg"`
	RespData struct {
		ReturnCode    string `json:"return_code"`
		ReturnMsg     string `json:"return_msg"`
		ResultCode    string `json:"result_code"`
		TradeState    string `json:"trade_state"`     // 交易状态
		TransactionID string `json:"transaction_id"`   // 微信订单号
		OutTradeNo    string `json:"out_trade_no"`     // 商户订单号
		TotalFee      int    `json:"total_fee"`        // 订单金额（分）
		TimeEnd       string `json:"time_end"`        // 支付完成时间
	} `json:"respdata"`
}

// CloseOrderRequest 关闭订单请求
type CloseOrderRequest struct {
	SubMchID   string `json:"sub_mch_id"`
	OutTradeNo string `json:"out_trade_no"`
	NonceStr   string `json:"nonce_str"`
}

// CloseOrderResponse 关闭订单响应
type CloseOrderResponse struct {
	ErrCode  int    `json:"errcode"`
	ErrMsg   string `json:"errmsg"`
	RespData struct {
		ReturnCode string `json:"return_code"`
		ReturnMsg  string `json:"return_msg"`
		ResultCode string `json:"result_code"`
	} `json:"respdata"`
}

// RefundRequest 申请退款请求
type RefundRequest struct {
	SubMchID     string `json:"sub_mch_id"`
	OutTradeNo   string `json:"out_trade_no,omitempty"`
	TransactionID string `json:"transaction_id,omitempty"`
	OutRefundNo  string `json:"out_refund_no"`
	TotalFee     int    `json:"total_fee"`
	RefundFee    int    `json:"refund_fee"`
	RefundDesc   string `json:"refund_desc,omitempty"`
	EnvID        string `json:"env_id"`
	CallbackType int    `json:"callback_type"`
	Container    struct {
		Service string `json:"service"`
		Path    string `json:"path"`
	} `json:"container"`
}

// RefundResponse 申请退款响应
type RefundResponse struct {
	ErrCode  int    `json:"errcode"`
	ErrMsg   string `json:"errmsg"`
	RespData struct {
		ReturnCode    string `json:"return_code"`
		ReturnMsg     string `json:"return_msg"`
		ResultCode    string `json:"result_code"`
		TransactionID string `json:"transaction_id"`
		OutTradeNo    string `json:"out_trade_no"`
		OutRefundNo   string `json:"out_refund_no"`
		RefundID      string `json:"refund_id"`
		RefundFee     int    `json:"refund_fee"`
	} `json:"respdata"`
}

// QueryRefundRequest 查询退款请求
type QueryRefundRequest struct {
	SubMchID     string `json:"sub_mch_id"`
	OutTradeNo   string `json:"out_trade_no,omitempty"`
	TransactionID string `json:"transaction_id,omitempty"`
	OutRefundNo  string `json:"out_refund_no,omitempty"`
	RefundID     string `json:"refund_id,omitempty"`
}

// QueryRefundResponse 查询退款响应
type QueryRefundResponse struct {
	ErrCode  int    `json:"errcode"`
	ErrMsg   string `json:"errmsg"`
	RespData struct {
		ReturnCode    string `json:"return_code"`
		ReturnMsg     string `json:"return_msg"`
		ResultCode    string `json:"result_code"`
		TransactionID string `json:"transaction_id"`
		OutTradeNo    string `json:"out_trade_no"`
		RefundCount   int    `json:"refund_count"`
	} `json:"respdata"`
}

// UnifiedOrder 统一下单
func (s *PaymentService) UnifiedOrder(req *UnifiedOrderRequest) (*UnifiedOrderResponse, error) {
	url := fmt.Sprintf("%s/unifiedOrder", WeChatPayAPIBase)
	return s.callWeChatPayAPI(url, req)
}

// QueryOrder 查询订单
func (s *PaymentService) QueryOrder(req *QueryOrderRequest) (*QueryOrderResponse, error) {
	url := fmt.Sprintf("%s/queryorder", WeChatPayAPIBase)
	resp, err := s.callWeChatPayAPIGeneric(url, req)
	if err != nil {
		return nil, err
	}

	var queryResp QueryOrderResponse
	if err := json.Unmarshal(resp, &queryResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}
	return &queryResp, nil
}

// CloseOrder 关闭订单
func (s *PaymentService) CloseOrder(req *CloseOrderRequest) (*CloseOrderResponse, error) {
	url := fmt.Sprintf("%s/closeorder", WeChatPayAPIBase)
	resp, err := s.callWeChatPayAPIGeneric(url, req)
	if err != nil {
		return nil, err
	}

	var closeResp CloseOrderResponse
	if err := json.Unmarshal(resp, &closeResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}
	return &closeResp, nil
}

// Refund 申请退款
func (s *PaymentService) Refund(req *RefundRequest) (*RefundResponse, error) {
	url := fmt.Sprintf("%s/refund", WeChatPayAPIBase)
	resp, err := s.callWeChatPayAPIGeneric(url, req)
	if err != nil {
		return nil, err
	}

	var refundResp RefundResponse
	if err := json.Unmarshal(resp, &refundResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}
	return &refundResp, nil
}

// QueryRefund 查询退款
func (s *PaymentService) QueryRefund(req *QueryRefundRequest) (*QueryRefundResponse, error) {
	url := fmt.Sprintf("%s/queryrefund", WeChatPayAPIBase)
	resp, err := s.callWeChatPayAPIGeneric(url, req)
	if err != nil {
		return nil, err
	}

	var queryRefundResp QueryRefundResponse
	if err := json.Unmarshal(resp, &queryRefundResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}
	return &queryRefundResp, nil
}

// callWeChatPayAPI 调用微信支付API（统一下单专用）
func (s *PaymentService) callWeChatPayAPI(url string, payload interface{}) (*UnifiedOrderResponse, error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal request failed: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response failed: %w", err)
	}

	var result UnifiedOrderResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}

	if result.ErrCode != 0 {
		return nil, fmt.Errorf("wechat pay api error: %s (errcode: %d)", result.ErrMsg, result.ErrCode)
	}

	return &result, nil
}

// callWeChatPayAPIGeneric 调用微信支付API（通用方法，返回原始JSON）
func (s *PaymentService) callWeChatPayAPIGeneric(url string, payload interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal request failed: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response failed: %w", err)
	}

	// 检查是否有错误
	var errorResp struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}
	if err := json.Unmarshal(body, &errorResp); err == nil && errorResp.ErrCode != 0 {
		return nil, fmt.Errorf("wechat pay api error: %s (errcode: %d)", errorResp.ErrMsg, errorResp.ErrCode)
	}

	return body, nil
}

