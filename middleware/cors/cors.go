package cors

import (
	"github.com/gin-gonic/gin"
)

import (
	"huilibao.com/core/common"
)

/* ================================================================================
 * 跨域处理
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊
 * ================================================================================ */
func Cors(corsOption *CorsOption) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !corsOption.IsAllowDomain {
			return
		}

		for _, domain := range corsOption.Domains {
			c.Header("Access-Control-Allow-Origin", domain)
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		}

		c.Next()
	}
}
