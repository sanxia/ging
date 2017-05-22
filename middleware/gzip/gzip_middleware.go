package gzip

import (
	"github.com/gin-gonic/gin"
)

/* ================================================================================
 * Gzip中间件模块
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊
 * ================================================================================ */
func GzipMiddleware() gin.HandlerFunc {
	return Gzip(BestSpeed)
}
