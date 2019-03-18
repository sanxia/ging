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
 * 输出Data
 * args格式: contentType | statusCode | contentType,statusCode |
 *          contentType,isAbort | contentType,statusCode,isAbort
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func Data(c *gin.Context, data []byte, args ...interface{}) {
	var contentType string
	statusCode := 200
	isAbort := false
	argsCount := len(args)

	if argsCount == 1 {
		switch value := args[0].(type) {
		case int:
			statusCode = value
		case string:
			contentType = value
		}
	} else if argsCount > 1 {
		contentType = args[0].(string)
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

	c.Data(statusCode, contentType, data)

	if isAbort {
		c.Abort()
	}
}
