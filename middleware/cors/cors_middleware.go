package cors

import (
	"github.com/gin-gonic/gin"
)

import (
	"huilibao.com/core/common"
)

/* ================================================================================
 * 跨域中间件模块
 * author: 美丽的地球啊
 * ================================================================================ */
type CorsOption struct {
	IsAllowDomain bool //是否允许跨域
	Domains       []string
}

func CorsMiddleware(corsOption *CorsOption) gin.HandlerFunc {
	return Cors(corsOption)
}
