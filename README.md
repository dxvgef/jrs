# JRS(JSON-RPC Server)

使用Go语言开发的JSON-RPC Server，仿遵JSON-RPC 2约定，仅做了以下改动：
* 去除了JSON-RPC协议版本字段
* 为了兼容批量请求，所有请求和返回都必须带上ID

### 示例：
```Go
package main

import (
	"log"
	"net/http"

	"github.com/dxvgef/jrs"
)

func main() {
	// 获得一个jrs实例，指定唯一入口路径
	svr := jrs.New("/rpc")
	// 设置映射
	svr.SetFunc("user.add", hello)
    // 启动HTTP服务
	http.ListenAndServe(":10080", svr)
}

/*
 用于映射的处理函数
 入参*jrs.Context里包含了请求的http.Context和JSON-RPC请求数据
 第1个出参是返回成功时的result字段数据
 第2个出参是返回错误时的代码
 第3个出参是返回错误时的消息字符串
 如果第2个参数值不为0，则自动视为出错，不对第1个出参进行处理
*/
func hello(ctx *jrs.Context) (interface, int, string) {
	// 准备一个map做为返回数据的result部份
	data := make(map[string]interface{})
	/* 
	 用jrs.Context的GetString、GetInt、GetInt64、GetFloat32、GetFloat64、GetBool函数
	 把请求参数值转为对应的类型
    */
    data["name"] = ctx.GetString("name")
    data["size"], _ = ctx.GetInt("size")
    if data["size"] == 0 {
    	// 如果出参的第二个
        return nil, jrs.InvalidParamsCode, "size不能等于0"
    }
    data["enabled"], _ = ctx.GetBool("enabled")
    data["price"], _ = ctx.GetFloat64("price")
    if data["price"] == 0 {
        return nil, jrs.InvalidParamsCode, "price不能等于0"
    }
    return data, 0, ""
}
```

### 请求示例
```JSON
[
	{
		"id": 1,
		"method": "add",
		"params": {
			"name": "苹果",
			"size": 1,
			"enabled": true,
			"price": 3.131415927
		}
	},
	{
		"id": 2,
		"method": "add",
		"params": {
			"name": "梨子",
			"size": 0,
			"enabled": false,
			"price": 100
		}
	}
]
```

### 返回示例
```JSON
[
    {
        "id": 1,
        "result": {
            "name": "苹果",
            "size": 1,
            "enabled": true,
            "price": 3.131415927
        }
    },
    {
        "ID": 2,
        "code": -32602,
        "message": "size不能等于0"
    }
]
```