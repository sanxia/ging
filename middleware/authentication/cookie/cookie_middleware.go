package cookie

import (
	"github.com/gin-gonic/gin"
)

import (
	"github.com/sanxia/ging"
	"github.com/sanxia/ging/result"
	"github.com/sanxia/ging/util"
)

/* ================================================================================
 * 表单认证中间件模块
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */
var (
	cookieAuth *CookieAuthentication
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Cookie认证中间件
 * extend: 表单扩展数据
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func CookieAuthenticationMiddleware(extend CookieExtend) gin.HandlerFunc {
	//初始化表单验证
	var err error
	cookieAuth, err = NewCookieAuthentication(CookieAuthentication{
		Extend: extend,
		Validate: func(ctx *gin.Context, extend CookieExtend, userIdentity *ging.UserIdentity) bool {
			return customValidate(ctx, extend, userIdentity)
		},
	})

	if err != nil {
		return func(ctx *gin.Context) {
			errorResult := map[string]interface{}{
				"Code": 111,
				"Msg":  "参数错误",
				"Data": nil,
			}
			viewResult := result.JsonResult(ctx, errorResult)
			viewResult.Render()

			ctx.Abort()
		}
	}

	//身份验证
	return cookieAuth.Validation()
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 自定义验证扩展点
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func customValidate(ctx *gin.Context, extend CookieExtend, userIdentity *ging.UserIdentity) bool {
	//用户角色是否匹配
	if len(extend.Role) > 0 {
		isInRole := util.IsInRole(userIdentity.Role, extend.Role)
		return isInRole
	}

	return true
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 登陆
 * userModel: 用户数据模型
 * isPersistence: 是否持久化登陆信息
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func Logon(ctx *gin.Context, userIdentity *ging.UserIdentity, isRemember bool) bool {
	return cookieAuth.Logon(ctx, userIdentity, isRemember)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 登出
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func Logoff(ctx *gin.Context) {
	cookieAuth.Logoff(ctx)
}
