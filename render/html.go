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
 * 输出Html
 * args格式: data | statusCode | data,statusCode | data,isAbort | data,statusCode,isAbort
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func Html(c *gin.Context, template string, args ...interface{}) {
	data := make(map[string]interface{}, 0)
	statusCode := 200
	isAbort := false
	argsCount := len(args)

	if argsCount == 1 {
		switch value := args[0].(type) {
		case int:
			statusCode = value
		case map[string]interface{}:
			data = value
		}
	} else if argsCount > 1 {
		data, _ = args[0].(map[string]interface{})
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

	c.HTML(statusCode, template, data)

	if isAbort {
		c.Abort()
	}
}
