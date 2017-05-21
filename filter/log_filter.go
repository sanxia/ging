package filter

import (
	"log"
	"time"
)

import (
	"github.com/gin-gonic/gin"
)

import (
	"github.com/sanxia/ging/result"
)

/* ================================================================================
 * 日志过滤器
 * author: 美丽的地球啊
 * ================================================================================ */
type logFilter struct {
	Filter
}

func LogFilter(args ...string) IActionFilter {
	name := "log"
	if len(args) == 1 {
		name = args[0]
	}
	return &logFilter{
		Filter: Filter{
			Name: name,
		},
	}
}

func (s *logFilter) Before(ctx *gin.Context) result.IActionResult {
	log.Printf("[filter: %s] Before %v", s.Name, time.Now())

	return nil
}

func (s *logFilter) After(ctx *gin.Context) {
	log.Printf("[filter: %s] After %v", s.Name, time.Now())
}
