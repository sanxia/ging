package result

import (
	"github.com/gin-gonic/gin"
)

import (
	"github.com/sanxia/ging"
)

/* ================================================================================
 * Html结果
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊
 * ================================================================================ */

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Html视图结果数据结构
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
type (
	htmlResult struct {
		ging.ActionResult
		tmpl string
	}
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Html视图结果
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func HtmlResult(context *gin.Context, tmpl string, args ...interface{}) ging.IActionResult {
	result := &htmlResult{
		tmpl: tmpl,
	}

	result.context = context
	result.contentType = "html"
	result.statusCode = 200

	argsCount := len(args)
	if argsCount > 0 {
		result.data = args[0]
		if argsCount == 2 {
			switch value := args[1].(type) {
			case int:
				result.statusCode = value
			case bool:
				result.isAbort = value
			}
		} else if argsCount == 3 {
			if statusCode, ok := args[1].(int); ok {
				result.statusCode = statusCode
			}
			if isAbort, ok := args[2].(bool); ok {
				result.isAbort = isAbort
			}
		}
	}
	return result
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 渲染
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (r *htmlResult) Render() {
	r.Html(r.tmpl, r.data, r.statusCode, r.isAbort)
}
