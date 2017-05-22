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

	result.Context = context
	result.ContentData = data
	result.ContentType = "image/png"
	result.StatusCode = 200

	argsCount := len(args)
	if argsCount == 1 {
		switch value := args[0].(type) {
		case int:
			result.StatusCode = value
		case bool:
			result.IsAbort = value
		}
	} else if argsCount == 2 {
		if statusCode, ok := args[0].(int); ok {
			result.StatusCode = statusCode
		} else if contentType, ok := args[0].(string); ok {
			result.ContentType = contentType
		}
		if isAbort, ok := args[1].(bool); ok {
			result.IsAbort = isAbort
		}
	} else if argsCount == 3 {
		if contentType, ok := args[0].(string); ok {
			result.ContentType = contentType
		}
		if statusCode, ok := args[1].(int); ok {
			result.StatusCode = statusCode
		}
		if isAbort, ok := args[2].(bool); ok {
			result.IsAbort = isAbort
		}
	}

	return result
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 渲染
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (r *dataResult) Render() {
	data, _ := r.ContentData.([]byte)
	r.Data(data, r.ContentType, r.StatusCode, r.IsAbort)
}
