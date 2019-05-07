package cookie

import (
	"github.com/gin-gonic/gin"
	"github.com/sanxia/ging"
	"github.com/sanxia/glib"
)

/* ================================================================================
 * 表单认证中间件模块
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */
var (
	cookieAuth *cookieAuthentication
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Cookie认证中间件
 * extend: Cookie扩展数据
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func CookieAuthenticationMiddleware(extend CookieExtend) gin.HandlerFunc {
	cookieAuth = &cookieAuthentication{
		Extend: extend,
		Validate: func(ctx *gin.Context, extend CookieExtend, tokenIdentity ging.IToken) bool {
			return customValidate(ctx, extend, tokenIdentity)
		},
	}

	//身份验证
	return cookieAuth.Validation()
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 自定义验证扩展点
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
 * 登入
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func Logon(ctx *gin.Context, payload *ging.TokenPayload) bool {
	return cookieAuth.Logon(ctx, payload)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 登出
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func Logoff(ctx *gin.Context) {
	cookieAuth.Logoff(ctx)
}
