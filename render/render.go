package render

import (
	"io"
	"net/http"
)

import (
	"github.com/gin-gonic/gin"
)

/* ================================================================================
 * Render 工具模块
 * qq: 2091938785
 * email: 2091938785@qq.com
 * author: 美丽的地球啊
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

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 输出Xml字符串
 * args格式: data | data,statusCode | data,isAbort | data,statusCode,isAbort
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func Xml(c *gin.Context, args ...interface{}) {
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
	c.XML(statusCode, data)
	if isAbort {
		c.Abort()
	}
}

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

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 输出磁盘物理文件
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func File(c *gin.Context, filepath string) {
	c.File(filepath)
	c.Abort()
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 输出二进制流
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func Stream(c *gin.Context, step func(w io.Writer) bool) {
	c.Stream(step)
	c.Abort()
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 输出错误消息
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func Error(c *gin.Context, msg string) {
	String(c, msg, 400)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 跳转
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func Redirect(c *gin.Context, url string) {
	c.Redirect(http.StatusSeeOther, url)
	c.Abort()
}
