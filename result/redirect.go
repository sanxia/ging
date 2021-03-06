package result

import (
	"github.com/gin-gonic/gin"
)

import (
	"github.com/sanxia/ging"
)

/* ================================================================================
 * Redirect结果
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Redirect视图结果数据结构
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
type (
	redirectResult struct {
		ging.ActionResult
		Url string
	}
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Redirect视图结果
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func RedirectResult(context *gin.Context, args ...string) ging.IActionResult {
	redirectUrl := "/"
	if len(args) > 0 {
		redirectUrl = args[0]
	}

	result := &redirectResult{
		Url: redirectUrl,
	}

	result.Context = context
	return result
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 渲染
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (r *redirectResult) Render() {
	r.Redirect(r.Url)
}
