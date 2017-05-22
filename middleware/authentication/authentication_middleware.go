package authentication

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
 * author  : 美丽的地球啊
 * ================================================================================ */
var (
	forms *FormsAuthentication
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 登陆
 * extend: 表单扩展数据
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func FormsAuthenticationMiddleware(extend FormsAuthenticationExtend) gin.HandlerFunc {
	//初始化表单验证
	var err error
	forms, err = NewFormAuthentication(FormsAuthentication{
		Extend: extend,
		Validate: func(ctx *gin.Context, formExtend FormsAuthenticationExtend, userIdentity *ging.UserIdentity) bool {
			return customValidate(ctx, formExtend, userIdentity)
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
	return forms.Validation()
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 自定义验证扩展点
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func customValidate(ctx *gin.Context, formExtend FormsAuthenticationExtend, userIdentity *ging.UserIdentity) bool {
	//用户角色是否匹配
	if len(formExtend.Role) > 0 {
		isInRole := util.IsInRole(userIdentity.Role, formExtend.Role)
		return isInRole
	}

	return true
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 登陆
 * userModel: 用户数据模型
 * isPersistence: 是否持久化登陆信息
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func Logon(ctx *gin.Context, userIdentity *ging.UserIdentity) bool {
	return forms.Logon(ctx, userIdentity)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 登出
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func Logoff(ctx *gin.Context) bool {
	return forms.Logoff(ctx)
}
