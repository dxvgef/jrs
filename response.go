package jrs

import (
	"net/http"
)
const (
	ParseErrorCode        = -32700 // 服务端解析JSON失败
	InvalidRequestCode    = -32600 // 无效的请求
	InvalidMethodCode     = -32601 // 无效的方法
	InvalidParamsCode     = -32602 // 无效的参数
	InternalErrorCode     = -32603 // 内部错误
	ParseErrorMessage     = "Parse error"
	InvalidRequestMessage = "Invalid Request"
	InvalidMethodMessage  = "Invalid Method"
	InvalidParamsMessage  = "Invalid Params"
	InternalErrorMessage  = "Internal Error"
)

// 构建错误
func makeError(code int, msg string) ([]byte, bool) {
	var errorBody ErrorBody
	errorBody.Code = code
	errorBody.Message = msg
	result, err := JSON.Marshal(&errorBody)
	if err != nil {
		return nil, false
	}
	return result, true
}

// 立即在http.response中输出错误
func returnError(resp http.ResponseWriter, code int, msg string) {
	result, ok := makeError(code, msg)
	if ok == true {
		resp.Header().Set("Content-Type", "application/json; charset=UTF-8")
		resp.WriteHeader(http.StatusOK)
		resp.Write(result)
	}
}

// 立即在http.response中输出成功
func returnSuccess(resp http.ResponseWriter, result []byte) {
	resp.Header().Set("Content-Type", "application/json; charset=UTF-8")
	resp.WriteHeader(http.StatusOK)
	resp.Write(result)
}