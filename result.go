package ging

import (
	"io"
)

import (
	"github.com/gin-gonic/gin"
)

import (
	"github.com/sanxia/ging/render"
)

/* ================================================================================
 * 动作结果
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊
 * ================================================================================ */

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 动作结果接口
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
type IActionResult interface {
	Render()
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 视图结果数据结构
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
type (
	ActionResult struct {
		Context     *gin.Context
		ContentData interface{}
		ContentType string
		StatusCode  int
		IsAbort     bool
	}
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 渲染Html模版
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (result *ActionResult) Html(template string, args ...interface{}) {
	if result.Context.IsAborted() {
		return
	}

	//result.Context.Writer.WriteHeader(200)
	//result.Context.Writer.WriteString("hahah liuming")
	//result.Context.Writer.Write([]byte("2016"))
	//result.Context.WriteHeader(w.status)
	//result.Context.Abort()

	render.Html(result.Context, template, args...)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 渲染Json字符串
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (result *ActionResult) Json(args ...interface{}) {
	if result.Context.IsAborted() {
		return
	}
	render.Json(result.Context, args...)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 渲染字符串
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (result *ActionResult) String(args ...interface{}) {
	if result.Context.IsAborted() {
		return
	}

	render.String(result.Context, args...)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 渲染Xml
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (result *ActionResult) Xml(args ...interface{}) {
	if result.Context.IsAborted() {
		return
	}

	render.Xml(result.Context, args...)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 渲染磁盘物理文件
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (result *ActionResult) File(filepath string) {
	if result.Context.IsAborted() {
		return
	}

	render.File(result.Context, filepath)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 渲染字节数据
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (result *ActionResult) Data(data []byte, args ...interface{}) {
	if result.Context.IsAborted() {
		return
	}

	render.Data(result.Context, data, args...)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 渲染io数据流
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (result *ActionResult) Stream(step func(w io.Writer) bool) {
	if result.Context.IsAborted() {
		return
	}

	render.Stream(result.Context, step)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 渲染io数据流
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (result *ActionResult) Error(msg string) {
	result.String(msg, 400)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 渲染err数据流
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (result *ActionResult) Redirect(url string) {
	render.Redirect(result.Context, url)
}
