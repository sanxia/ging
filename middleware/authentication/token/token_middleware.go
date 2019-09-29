package token

import (
	"github.com/gin-gonic/gin"
	"github.com/sanxia/ging"
	"github.com/sanxia/glib"
)

/* ================================================================================
 * Token Auth Middleware
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */
var (
	tokenAuth *tokenAuthentication
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * token auth middleware
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func TokenAuthenticationMiddleware(extend TokenExtend) gin.HandlerFunc {
	tokenAuth = &tokenAuthentication{
		Extend: extend,
		Validate: func(ctx *gin.Context, tokenExtend TokenExtend, tokenIdentity ging.IToken) bool {
			return customValidate(ctx, tokenExtend, tokenIdentity)
		},
	}

	return tokenAuth.Validation()
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * custom validate
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func customValidate(ctx *gin.Context, extend TokenExtend, tokenIdentity ging.IToken) bool {
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
func Logon(ctx *gin.Context, payload *ging.TokenPayload) {
	tokenAuth.Logon(ctx, payload)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * logoff
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func Logoff(ctx *gin.Context) {
	tokenAuth.Logoff(ctx)
}
