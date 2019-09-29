package security

import (
	"strings"
)

import (
	"github.com/gin-gonic/gin"
	"github.com/sanxia/ging"
)

/* ================================================================================
 * 安全处理
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */
func securityHandler(securityOption Security) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if securityOption.IsDisabled {
			ctx.Next()
			return
		}

		appSetting := ging.GetApp().GetSetting()
		currentUserId := ""
		if userToken := ging.GetToken(ctx); userToken != nil {
			currentUserId = userToken.GetPayload().UserId
		}

		//白名单
		if !appSetting.Security.White.IsDisabled {
			if appSetting.Security.White.IsInIps(ctx.ClientIP()) ||
				appSetting.Security.White.IsInUsers(currentUserId) {
				ctx.Next()
				return
			}
		}

		//黑名单
		if !appSetting.Security.Black.IsDisabled {
			if appSetting.Security.Black.IsInIps(ctx.ClientIP()) ||
				appSetting.Security.Black.IsInUsers(currentUserId) {
				securityOption.errorHandler(ctx, "blacklist")
				return
			}
		}

		//应用运行状态
		if appSetting.Domain.StatusCode != "runing" {
			securityOption.errorHandler(ctx, appSetting.Domain.StatusCode)
			return
		}

		//访问时段
		requestMethod := strings.ToLower(ctx.Request.Method)
		inMethods := map[string]bool{
			"post":   true,
			"put":    true,
			"delete": true,
		}

		if _, isOk := inMethods[requestMethod]; isOk {
			// 写时段
			if !appSetting.Security.InTime.IsDisabled && !appSetting.Security.InTime.IsTime() {
				securityOption.errorHandler(ctx, "time")
				return
			}
		} else {
			// 读时段
			if !appSetting.Security.OutTime.IsDisabled && !appSetting.Security.OutTime.IsTime() {
				securityOption.errorHandler(ctx, "time")
				return
			}
		}
	}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * error process
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s Security) errorHandler(ctx *gin.Context, code string) {
	isAjax := ging.IsAjax(ctx)

	if s.ErrorHandler != nil {
		s.ErrorHandler(ctx, code, isAjax).Render()
	}

	ctx.Abort()
}
