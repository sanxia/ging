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
 * Action Result
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * action result interface
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
type IActionResult interface {
	Render()
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * view result data
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
 * rendering Html templates
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
 * rendering json strings
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (result *ActionResult) Json(args ...interface{}) {
	if result.Context.IsAborted() {
		return
	}
	render.Json(result.Context, args...)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * rendering text strings
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (result *ActionResult) String(args ...interface{}) {
	if result.Context.IsAborted() {
		return
	}

	render.String(result.Context, args...)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * rendering xml strings
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (result *ActionResult) Xml(args ...interface{}) {
	if result.Context.IsAborted() {
		return
	}

	render.Xml(result.Context, args...)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * rendering disk file
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (result *ActionResult) File(filepath string) {
	if result.Context.IsAborted() {
		return
	}

	render.File(result.Context, filepath)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * rendering bytes
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (result *ActionResult) Data(data []byte, args ...interface{}) {
	if result.Context.IsAborted() {
		return
	}

	render.Data(result.Context, data, args...)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * rendering io stream
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (result *ActionResult) Stream(step func(w io.Writer) bool) {
	if result.Context.IsAborted() {
		return
	}

	render.Stream(result.Context, step)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * rendering error strings
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (result *ActionResult) Error(msg string) {
	result.String(msg, 400)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * rendering redirect url
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (result *ActionResult) Redirect(url string) {
	render.Redirect(result.Context, url)
}
