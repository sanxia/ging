package filter

import (
	"log"
	"time"
)

import (
	"github.com/gin-gonic/gin"
)

import (
	"github.com/sanxia/ging"
)

/* ================================================================================
 * 日志过滤器
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */
type logFilter struct {
	ging.Filter
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 实例化日志过滤器
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func NewLogFilter(args ...string) ging.IActionFilter {
	return &logFilter{
		Filter: ging.Filter{
			Name: "log_filter",
		},
	}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 动作执行之前
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *logFilter) Before(ctx *gin.Context) ging.IActionResult {
	url := ctx.Request.RequestURI
	method := ctx.Request.Method

	log.Printf("[%s] Before %s %s %v", s.Name, method, url, time.Now())

	return nil
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 动作执行之后
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *logFilter) After(ctx *gin.Context) {
	log.Printf("[%s] After %v", s.Name, time.Now())
}
