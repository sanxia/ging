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
 * 计时过滤器
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */
type timerFilter struct {
	ging.Filter
	beginDate time.Time
	endDate   time.Time
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 实例化计时过滤器
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func NewTimerFilter(args ...string) ging.IActionFilter {
	return &timerFilter{
		Filter: ging.Filter{
			Name: "timer_filter",
		},
	}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 动作执行之前
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *timerFilter) Before(ctx *gin.Context) ging.IActionResult {
	s.beginDate = time.Now()
	log.Printf("[%s] Before %v", s.Name, s.beginDate)

	return nil
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 动作执行之后
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *timerFilter) After(ctx *gin.Context) {
	if !ctx.IsAborted() {
		s.endDate = time.Now()
		nanoseconds := s.endDate.Sub(s.beginDate).Nanoseconds()

		msg := fmt.Sprintf("total time: %v", nanoseconds)
		log.Printf("[%s] After %v %s ", s.Name, s.endDate, msg)
	}
}
