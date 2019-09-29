package cors

import (
	"github.com/gin-gonic/gin"
)

/* ================================================================================
 * 跨域中间件模块
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */
type CorsOption struct {
	IsAllowDomain bool //是否允许跨域
	Domains       []string
}

func CorsMiddleware(corsOption *CorsOption) gin.HandlerFunc {
	return Cors(corsOption)
}
