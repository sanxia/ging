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
 * qq: 2091938785
 * email: 2091938785@qq.com
 * author: 美丽的地球啊
 * ================================================================================ */
type logFilter struct {
	Filter
}

func LogFilter(args ...string) ging.IActionFilter {
	return &logFilter{
		Filter: Filter{
			Name: "log_filter",
		},
	}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 动作执行之前
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *logFilter) Before(ctx *gin.Context) ging.IActionResult {
	log.Printf("[%s] Before %v", s.Name, time.Now())

	return nil
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 动作执行之后
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *logFilter) After(ctx *gin.Context) {
	log.Printf("[%s] After %v", s.Name, time.Now())
}
