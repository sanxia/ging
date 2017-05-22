package result

import (
	"github.com/gin-gonic/gin"
)

import (
	"github.com/sanxia/ging"
)

/* ================================================================================
 * Json结果
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊
 * ================================================================================ */

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 视图结果数据结构
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
type (
	dataResult struct {
		ging.ActionResult
	}
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Data结果
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func DataResult(context *gin.Context, data interface{}, args ...interface{}) ging.IActionResult {
	result := &dataResult{}

	result.context = context
	result.data = data
	result.contentType = "image/png"
	result.statusCode = 200

	argsCount := len(args)
	if argsCount == 1 {
		switch value := args[0].(type) {
		case int:
			result.statusCode = value
		case bool:
			result.isAbort = value
		}
	} else if argsCount == 2 {
		if statusCode, ok := args[0].(int); ok {
			result.statusCode = statusCode
		} else if contentType, ok := args[0].(string); ok {
			result.contentType = contentType
		}
		if isAbort, ok := args[1].(bool); ok {
			result.isAbort = isAbort
		}
	} else if argsCount == 3 {
		if contentType, ok := args[0].(string); ok {
			result.contentType = contentType
		}
		if statusCode, ok := args[1].(int); ok {
			result.statusCode = statusCode
		}
		if isAbort, ok := args[2].(bool); ok {
			result.isAbort = isAbort
		}
	}

	return result
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 渲染
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (r *dataResult) Render() {
	data, _ := r.data.([]byte)
	r.Data(data, r.contentType, r.statusCode, r.isAbort)
}
