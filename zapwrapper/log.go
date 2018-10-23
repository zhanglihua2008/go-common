package zapwrapper

import (
	"fmt"

	"go.uber.org/zap"
)

func getLogParams(params ...interface{}) []zap.Field {

	fields := make([]zap.Field, 0, 10)

	key := ""
	for i, param := range params {
		if i%2 == 0 {
			switch param.(type) {
			case string:
				key = param.(string)
			case error:
				key = param.(error).Error()
			default:
				key = fmt.Sprint(param)
			}
		} else {
			value := ""
			switch param.(type) {
			case string:
				value = param.(string)
			case error:
				value = param.(error).Error()
			default:
				value = fmt.Sprint(param)
			}
			fields = append(fields, zap.String(key, value))
		}
	}

	return fields
}

func Info(msg string, params ...interface{}) {
	if len(params) >= 2 {
		fields := getLogParams(params...)
		zap.L().Info(msg, fields...)
	} else {
		zap.L().Info(msg)
	}
}

func Error(msg string, params ...interface{}) {
	if len(params) >= 2 {
		fields := getLogParams(params...)
		zap.L().Error(msg, fields...)
	} else {
		zap.L().Error(msg)
	}
}
