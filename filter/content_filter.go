package filter

import (
	"github.com/gin-gonic/gin"
)

import (
	"github.com/sanxia/ging"
)

/* ================================================================================
 * 内容过滤器
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */
type contentFilter struct {
	ging.Filter
	header string
	footer string
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 实例化内容过滤器
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func NewContentFilter(header, footer string) ging.IActionFilter {
	return &contentFilter{
		Filter: ging.Filter{
			Name: "content_filter",
		},
		header: header,
		footer: footer,
	}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 动作执行之前
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *contentFilter) Before(ctx *gin.Context) ging.IActionResult {
	if len(s.header) > 0 {
		ctx.Writer.Write([]byte(s.header))
	}

	return nil
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 动作执行之后
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *contentFilter) After(ctx *gin.Context) {
	if !ctx.IsAborted() {
		if len(s.footer) > 0 {
			ctx.Writer.Write([]byte(s.footer))
		}
	}
}
