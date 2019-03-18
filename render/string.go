package render

import (
	"github.com/gin-gonic/gin"
)

/* ================================================================================
 * Render 工具模块
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 输出字符串
 * args格式: string | []byte | data,statusCode | data,isAbort | data,statusCode,isAbort
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func String(c *gin.Context, args ...interface{}) {
	var data string
	statusCode := 200
	isAbort := false
	argsCount := len(args)

	if argsCount == 1 {
		switch value := args[0].(type) {
		case string:
			data = value
		case []byte:
			data = string(value)
		}
	} else if argsCount > 1 {
		switch value := args[0].(type) {
		case string:
			data = value
		case []byte:
			data = string(value)
		}
		if argsCount == 2 {
			switch value := args[1].(type) {
			case int:
				statusCode = value
			case bool:
				isAbort = value
			}
		} else if argsCount == 3 {
			statusCode = args[1].(int)
			isAbort = args[2].(bool)
		}
	}

	c.String(statusCode, data)

	if isAbort {
		c.Abort()
	}
}
