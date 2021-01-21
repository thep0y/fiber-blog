/*
* @Author: thepoy
* @Email:  thepoy@aliyun.com
* @File Name: responses.go
* @Created:   2021-01-17 18:58:31
* @Modified:  2021-01-19 08:33:34
 */

package models

// ErrorResponse 错误响应
type ErrorResponse struct {
	Error string `json:"error"`
}

// BaseResponse 基本响应
type BaseResponse struct {
	Code    uint   `json:"status"`
	Message string `json:"msg"`
}
