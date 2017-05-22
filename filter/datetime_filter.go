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
	name := "datetime"
	if len(args) == 1 {
		name = args[0]
	}
	return &datetimeFilter{
		Filter: Filter{
			Name: name,
		},
	}
}

func (s *datetimeFilter) Before(ctx *gin.Context) ging.IActionResult {
	s.beginDate = time.Now()
	log.Printf("[filter: %s] Before %v", s.Name, s.beginDate)

	return nil
}

func (s *datetimeFilter) After(ctx *gin.Context) {
	if !ctx.IsAborted() {
		s.endDate = time.Now()
		nanoseconds := s.endDate.Sub(s.beginDate).Nanoseconds()
		msg := fmt.Sprintf("[filter: %s] After total time: %v", s.Name, nanoseconds)

		ctx.Writer.Write([]byte(msg))
	}

}
