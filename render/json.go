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
 * 输出Json字符串
 * args格式: data | data,statusCode | data,isAbort | data,statusCode,isAbort
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func Json(c *gin.Context, args ...interface{}) {
	var data interface{}
	statusCode := 200
	isAbort := false
	if len(args) == 1 {
		data = args[0]
	} else if len(args) == 2 {
		data = args[0]
		switch value := args[1].(type) {
		case int:
			statusCode = value
		case bool:
			isAbort = value
		}
	} else if len(args) == 3 {
		data = args[0]
		statusCode = args[1].(int)
		isAbort = args[2].(bool)
	}

	c.JSON(statusCode, data)

	if isAbort {
		c.Abort()
	}
}
