package utils

import (
	"encoding/json"
	"net/http"
)

// JsonResult 统一响应结构
type JsonResult struct {
	Code     int         `json:"code"`
	ErrorMsg string      `json:"errorMsg,omitempty"`
	Data     interface{} `json:"data,omitempty"`
}

// Success 成功响应
func Success(data interface{}) *JsonResult {
	return &JsonResult{
		Code: 0,
		Data: data,
	}
}

// Error 错误响应
func Error(code int, errorMsg string) *JsonResult {
	return &JsonResult{
		Code:     code,
		ErrorMsg: errorMsg,
	}
}

// WriteJSON 写入JSON响应
func (jr *JsonResult) WriteJSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jr)
}

