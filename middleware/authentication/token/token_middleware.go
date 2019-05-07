package token

import (
	"github.com/gin-gonic/gin"
	"github.com/sanxia/ging"
	"github.com/sanxia/glib"
)

/* ================================================================================
 * Token认证中间件模块
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */
var (
	tokenAuth *tokenAuthentication
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Token认证中间件
 * extend: Token扩展数据
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func TokenAuthenticationMiddleware(extend TokenExtend) gin.HandlerFunc {
	tokenAuth = &tokenAuthentication{
		Extend: extend,
		Validate: func(ctx *gin.Context, tokenExtend TokenExtend, tokenIdentity ging.IToken) bool {
			return customValidate(ctx, tokenExtend, tokenIdentity)
		},
	}

	//身份验证
	return tokenAuth.Validation()
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 自定义验证扩展点
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
 * 登入
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func Logon(ctx *gin.Context, payload *ging.TokenPayload) {
	tokenAuth.Logon(ctx, payload)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 登出
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func Logoff(ctx *gin.Context) {
	tokenAuth.Logoff(ctx)
}
