package result

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
 * author: 美丽的地球啊
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
		context     *gin.Context
		data        interface{}
		contentType string
		statusCode  int
		isAbort     bool
	}
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 渲染Html模版
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (result *ActionResult) Html(template string, args ...interface{}) {
	if result.context.IsAborted() {
		return
	}

	//result.context.Writer.WriteHeader(200)
	//result.context.Writer.WriteString("hahah liuming")
	//result.context.Writer.Write([]byte("2016"))
	//result.context.WriteHeader(w.status)
	//result.context.Abort()

	render.Html(result.context, template, args...)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 渲染Json字符串
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (result *ActionResult) Json(args ...interface{}) {
	if result.context.IsAborted() {
		return
	}
	render.Json(result.context, args...)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 渲染字符串
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (result *ActionResult) String(args ...interface{}) {
	if result.context.IsAborted() {
		return
	}

	render.String(result.context, args...)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 渲染Xml
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (result *ActionResult) Xml(args ...interface{}) {
	if result.context.IsAborted() {
		return
	}

	render.Xml(result.context, args...)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 渲染磁盘物理文件
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (result *ActionResult) File(filepath string) {
	if result.context.IsAborted() {
		return
	}

	render.File(result.context, filepath)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 渲染字节数据
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (result *ActionResult) Data(data []byte, args ...interface{}) {
	if result.context.IsAborted() {
		return
	}

	render.Data(result.context, data, args...)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 渲染io数据流
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (result *ActionResult) Stream(step func(w io.Writer) bool) {
	if result.context.IsAborted() {
		return
	}

	render.Stream(result.context, step)
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
	render.Redirect(result.context, url)
}
