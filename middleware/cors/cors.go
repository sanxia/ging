package cors

import (
	"strings"
)

import (
	"github.com/gin-gonic/gin"
)

/* ================================================================================
 * 跨域处理
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */
func Cors(corsOption *CorsOption) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !corsOption.IsAllow {
			return
		}

		if corsOption.IsAllDomain {
			c.Header("Access-Control-Allow-Origin", "*")
		} else {
			c.Header("Access-Control-Allow-Origin", corsOption.Domain)
		}

		if corsOption.IsCredentials {
			c.Header("Access-Control-Allow-Credentials", "true")
		} else {
			c.Header("Access-Control-Allow-Credentials", "false")
		}

		if len(corsOption.Methods) > 0 {
			method := strings.Join(corsOption.Methods, ",")
			c.Header("Access-Control-Allow-Methods", method)
		}

		if len(corsOption.Headers) > 0 {
			header := strings.Join(corsOption.Headers, ",")
			c.Header("Access-Control-Allow-Headers", header)
		}

		c.Next()
	}
}
