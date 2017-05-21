package gzip

import (
	"github.com/gin-gonic/gin"
)

/* ================================================================================
 * Gzip中间件模块
 * author: 美丽的地球啊
 * ================================================================================ */
func GzipMiddleware() gin.HandlerFunc {
	return Gzip(BestSpeed)
}
