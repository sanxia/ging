package token

import (
	"github.com/gin-gonic/gin"
)

import (
	"github.com/sanxia/ging"
	"github.com/sanxia/ging/result"
	"github.com/sanxia/ging/util"
)

/* ================================================================================
 * Token认证中间件模块
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊
 * ================================================================================ */
var (
	tokenAuth *TokenAuthentication
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Token认证中间件
 * extend: Token扩展数据
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func TokenAuthenticationMiddleware(extend TokenExtend) gin.HandlerFunc {
	//初始化Token验证
	var err error
	tokenAuth, err = NewTokenAuthentication(TokenAuthentication{
		Extend: extend,
		Validate: func(ctx *gin.Context, tokenExtend TokenExtend, userIdentity *ging.UserIdentity) bool {
			return customValidate(ctx, tokenExtend, userIdentity)
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
	return tokenAuth.Validation()
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 自定义验证扩展点
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func customValidate(ctx *gin.Context, tokenExtend TokenExtend, userIdentity *ging.UserIdentity) bool {
	//用户角色是否匹配
	if len(tokenExtend.Role) > 0 {
		isInRole := util.IsInRole(userIdentity.Role, tokenExtend.Role)
		return isInRole
	}

	return true
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 登陆
 * userModel: 用户数据模型
 * isPersistence: 是否持久化登陆信息
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func Logon(ctx *gin.Context, userIdentity *ging.UserIdentity) string {
	return tokenAuth.Logon(ctx, userIdentity)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 登出
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func Logoff(ctx *gin.Context) {
	tokenAuth.Logoff(ctx)
}
