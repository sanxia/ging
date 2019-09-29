package cookie

import (
	"github.com/gin-gonic/gin"
	"github.com/sanxia/ging"
	"github.com/sanxia/glib"
)

/* ================================================================================
 * cookie middleware
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */
var (
	cookieAuth *cookieAuthentication
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * cookie auth middleware
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func CookieAuthenticationMiddleware(extend CookieExtend) gin.HandlerFunc {
	cookieAuth = &cookieAuthentication{
		Extend: extend,
		Validate: func(ctx *gin.Context, extend CookieExtend, tokenIdentity ging.IToken) bool {
			return customValidate(ctx, extend, tokenIdentity)
		},
	}

	return cookieAuth.Validation()
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * custom validate
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func customValidate(ctx *gin.Context, extend CookieExtend, tokenIdentity ging.IToken) bool {
	isInRole := true

	if len(extend.Roles) > 0 {
		if roles := glib.StringInter(tokenIdentity.GetPayload().Roles, extend.Roles); len(roles) == 0 {
			isInRole = false
		}
	}

	return isInRole
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * logon
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func Logon(ctx *gin.Context, payload *ging.TokenPayload) bool {
	return cookieAuth.Logon(ctx, payload)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * logoff
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func Logoff(ctx *gin.Context) {
	cookieAuth.Logoff(ctx)
}
