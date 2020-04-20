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
	Domain        string
	Headers       []string
	Methods       []string
	IsCredentials bool //是否cookie
	IsAllDomain   bool //是否所有域
	IsAllow       bool //是否允许跨域

}

func CorsMiddleware(corsOption *CorsOption) gin.HandlerFunc {
	return Cors(corsOption)
}
