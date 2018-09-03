package jrs

import (
	"errors"
	"github.com/shopspring/decimal"
)

func (ctx *Context) GetString(key string, def ...string) string {
	defLen := len(def)
	if ctx.Body.Params[key] == nil {
		if defLen>0 {
			return def[0]
		}
		return ""
	}
	result, ok := ctx.Body.Params[key].(string)
	if ok== true {
		return result
	}
	if defLen>0 {
		return def[0]
	}
	return ""
}

func (ctx *Context) GetInt(key string, def ...int) (int, error) {
	errMsg := "Invalid int value"
	defLen := len(def)
	if ctx.Body.Params[key] != nil {
		f64, ok := ctx.Body.Params[key].(float64)
		if ok== true {
			if int(decimal.NewFromFloat(f64).IntPart()) >0 {
				return int(decimal.NewFromFloat(f64).IntPart()), nil
			}
		}
	}
	if defLen>0 {
		return def[0],  nil
	}
	return 0, errors.New(errMsg)
}

func (ctx *Context) GetInt64(key string, def ...int64) (int64, error) {
	errMsg := "Invalid int64 value"
	defLen := len(def)
	if ctx.Body.Params[key] == nil {
		if defLen>0 {
			return def[0], nil
		}
		return 0, errors.New(errMsg)
	}
	f64, ok := ctx.Body.Params[key].(float64)
	if ok== true {
		return decimal.NewFromFloat(f64).IntPart(), nil
	}
	if defLen>0 {
		return def[0],  nil
	}
	return 0, errors.New(errMsg)
}

func (ctx *Context) GetFloat32(key string, def ...float32) (float32, error) {
	errMsg := "Invalid float32 value"
	defLen := len(def)
	if ctx.Body.Params[key] == nil {
		if defLen>0 {
			return def[0], nil
		}
		return 0, errors.New(errMsg)
	}
	result, ok := ctx.Body.Params[key].(float64)
	if ok== true {
		return float32(result), nil
	}
	if defLen>0 {
		return def[0],  nil
	}
	return 0, errors.New(errMsg)
}

func (ctx *Context) GetFloat64(key string, def ...float64) (float64, error) {
	errMsg := "Invalid float64 value"
	defLen := len(def)
	if ctx.Body.Params[key] == nil {
		if defLen>0 {
			return def[0], nil
		}
		return 0, errors.New(errMsg)
	}
	result, ok := ctx.Body.Params[key].(float64)
	if ok== true {
		return result, nil
	}
	if defLen>0 {
		return def[0],  nil
	}
	return 0, errors.New(errMsg)
}

func (ctx *Context) GetBool(key string, def ...bool) (bool, error) {
	errMsg := "Invalid bool value"
	defLen := len(def)
	if ctx.Body.Params[key] == nil {
		if defLen>0 {
			return def[0], nil
		}
		return false, errors.New(errMsg)
	}
	result, ok := ctx.Body.Params[key].(bool)
	if ok== true {
		return result, nil
	}
	if defLen>0 {
		return def[0],  nil
	}
	return false, errors.New(errMsg)
}