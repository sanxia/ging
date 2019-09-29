package security

import (
	"github.com/gin-gonic/gin"
	"github.com/sanxia/ging"
)

/* ================================================================================
 * 安全中间件模块
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */
type Security struct {
	ErrorHandler func(ctx *gin.Context, code string, isAjax bool) ging.IActionResult //error result
	IsDisabled   bool
}

func SecurityMiddleware(securityOption Security) gin.HandlerFunc {
	return securityHandler(securityOption)
}
