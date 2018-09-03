package jrs

import (
	"github.com/json-iterator/go"
	"io/ioutil"
	"log"
	"net/http"
)

var JSON = jsoniter.ConfigCompatibleWithStandardLibrary

type Func func(*Context) (interface{}, int, string)

type ReqBody struct {
	ID     int64                  `json:"id"`
	Method string                 `json:"method"`
	Params map[string]interface{} `json:"params"`
}

type ErrorBody struct {
	ID  int64 `sql:"id"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type SuccessBody struct {
	ID     int64       `json:"id"`
	Result interface{} `json:"result"`
}

type Context struct {
	Resp http.ResponseWriter
	Req  *http.Request
	Body *ReqBody
}

type Service struct {
	EndPoint string
	mapper   map[string]Func
}

func New(endPoint string) *Service {
	var service Service
	service.EndPoint = endPoint
	service.mapper = make(map[string]Func)
	return &service
}

func (this *Service) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	if req.RequestURI != this.EndPoint {
		resp.WriteHeader(http.StatusNotFound)
		resp.Write([]byte(http.StatusText(http.StatusNotFound)))
		return
	}
	if req.Method != "POST" {
		resp.WriteHeader(http.StatusMethodNotAllowed)
		resp.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
		return
	}

	result, err := ioutil.ReadAll(req.Body)
	if err != nil {
		returnError(resp, ParseErrorCode, ParseErrorMessage)
		return
	}
	err = req.Body.Close()
	if err != nil {
		returnError(resp, ParseErrorCode, ParseErrorMessage)
		return
	}
	if JSON.Valid(result) == false {
		returnError(resp, ParseErrorCode, ParseErrorMessage)
		return
	}

	switch result[0] {
	case 123:
		var reqBody ReqBody
		reqBody.Params = make(map[string]interface{})
		err = JSON.Unmarshal(result, &reqBody)
		if err != nil {
			returnError(resp, ParseErrorCode, ParseErrorMessage)
			return
		}

		// 创建ctx
		var ctx Context
		ctx.Req = req
		ctx.Resp = resp
		ctx.Body = &reqBody

		// 取得对应的函数，并判断函数是否存在
		handler := this.mapper[ctx.Body.Method]
		if handler == nil {
			makeError(InvalidMethodCode, InvalidMethodMessage)
			return
		}

		// 执行对应的函数，并取得返回结果
		result, code, msg := handler(&ctx)
		if code != 0 {
			returnError(resp, code, msg)
			return
		}

		var successBody SuccessBody
		successBody.ID = reqBody.ID
		successBody.Result = result
		resultBytes, err := JSON.Marshal(successBody)
		if err != nil {
			returnError(resp, InternalErrorCode, InternalErrorMessage)
			return
		}
		returnSuccess(resp, resultBytes)
	case 91:
		var reqBodys []ReqBody
		err = JSON.Unmarshal(result, &reqBodys)
		if err != nil {
			log.Println(err.Error())
			returnError(resp, ParseErrorCode, ParseErrorMessage)
			return
		}

		var results []interface{}

		for k := range reqBodys {
			reqBody := reqBodys[k]
			// 创建ctx
			var ctx Context
			ctx.Req = req
			ctx.Resp = resp
			ctx.Body = &reqBody

			// 取得对应的函数，并判断函数是否存在
			handler := this.mapper[ctx.Body.Method]
			if handler == nil {
				results = append(results, ErrorBody{
					ID : reqBody.ID,
					Code : InvalidMethodCode,
					Message : InvalidMethodMessage,
				})
			} else {
				// 执行对应的函数，并取得返回结果
				result, code, msg := handler(&ctx)
				if code == 0 {
					results = append(results, SuccessBody{
						ID : reqBody.ID,
						Result : result,
					})
				} else {
					results = append(results, ErrorBody{
						ID:reqBody.ID,
						Code : code,
						Message : msg,
					})
				}
			}
		}
		result, err := JSON.Marshal(results)
		if err != nil {
			returnError(resp, InternalErrorCode, InternalErrorMessage)
			return
		}
		returnSuccess(resp, result)
	default:
		returnError(resp, ParseErrorCode, ParseErrorMessage)
		return
	}

}

// 设置函数到映射器
func (this *Service) SetFunc(name string, f Func) {
	this.mapper[name] = f
}
