package filter

import (
	"fmt"
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
 * 日期过滤器
 * author: 美丽的地球啊
 * ================================================================================ */
type datetimeFilter struct {
	Filter
	beginDate time.Time
	endDate   time.Time
}

func DatetimeFilter(args ...string) ging.IActionFilter {
	return &datetimeFilter{
		Filter: Filter{
			Name: "datetime_filter",
		},
	}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 动作执行之前
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *datetimeFilter) Before(ctx *gin.Context) ging.IActionResult {
	s.beginDate = time.Now()
	log.Printf("[%s] Before %v", s.Name, s.beginDate)

	return nil
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 动作执行之后
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *datetimeFilter) After(ctx *gin.Context) {
	if !ctx.IsAborted() {
		s.endDate = time.Now()
		nanoseconds := s.endDate.Sub(s.beginDate).Nanoseconds()
		msg := fmt.Sprintf("[%s] After total time: %v", s.Name, nanoseconds)

		ctx.Writer.Write([]byte(msg))
	}

}
