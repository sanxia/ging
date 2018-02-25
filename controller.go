package ging

import (
	"log"
	"strings"
)

import (
	"github.com/gin-gonic/gin"
)

import (
	"github.com/sanxia/ging/middleware/session"
)

/* ================================================================================
 * 控制器数据结构
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 控制器接口
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
type (
	IController interface {
		Action(action func(ctx *gin.Context) IActionResult, args ...interface{}) func(*gin.Context)
		Filter(filters ...IActionFilter) IController

		SaveSession(ctx *gin.Context, name, value string)
		ValidateSession(ctx *gin.Context, name, value string, args ...bool) bool
		RemoveSession(ctx *gin.Context, name string)
		ClearSession(ctx *gin.Context)

		GetUserIdentity(ctx *gin.Context) *UserIdentity
	}
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 控制器数据结构
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
type (
	Controller struct {
		filters []IActionFilter
	}
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 控制器动作
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (ctrl *Controller) Action(action func(ctx *gin.Context) IActionResult, args ...interface{}) func(*gin.Context) {
	return func(context *gin.Context) {
		var actionFilters []IActionFilter
		var filterResult IActionResult

		isEnabled := true
		argsCount := len(args)
		if argsCount > 0 {
			if value, isOk := args[0].(bool); isOk {
				isEnabled = value
			} else {
				for _, actionFilter := range args {
					if actionFilter, isOk := actionFilter.(IActionFilter); isOk {
						if len(actionFilters) == 0 {
							actionFilters = make([]IActionFilter, argsCount)
						}
						actionFilters = append(actionFilters, actionFilter)
					}
				}
			}
		}

		if isEnabled {
			//动作执行之前的控制器动作过滤器
			//Before返回非空IActionResult即终止
			for _, filter := range ctrl.filters {
				if filter != nil {
					if filterResult = filter.Before(context); filterResult != nil {
						break
					}
				}
			}

			if filterResult == nil {
				//动作执行之前的动作过滤器
				for _, filter := range actionFilters {
					if filter != nil {
						if filterResult = filter.Before(context); filterResult != nil {
							break
						}
					}
				}
			}
		}

		//执行过滤器
		if filterResult != nil {
			filterResult.Render()
		} else {
			//执行动作
			action(context).Render()
		}

		if isEnabled {
			//动作执行之后的动作过滤器
			for _, filter := range actionFilters {
				if filter != nil {
					filter.After(context)
				}
			}

			//动作执行之后的控制器动作过滤器
			for _, filter := range ctrl.filters {
				if filter != nil {
					filter.After(context)
				}
			}
		}
	}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 设置控制器过滤器（控制器的方法执行前后都会执行过滤器接口，过滤器接口集合不支持排序）
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (ctrl *Controller) Filter(filters ...IActionFilter) IController {
	if len(filters) == 0 {
		return ctrl
	}

	if len(ctrl.filters) == 0 {
		ctrl.filters = make([]IActionFilter, 0)
	}

	for _, filter := range filters {
		ctrl.filters = append(ctrl.filters, filter)
	}

	return ctrl
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 获取用户标识
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (ctrl *Controller) GetUserIdentity(ctx *gin.Context) *UserIdentity {
	var userIdentity *UserIdentity
	if identity, ok := ctx.Get(UserIdentityKey); ok {
		user := identity.(UserIdentity)
		userIdentity = &user
	}

	return userIdentity
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 获取会话对象
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (ctrl *Controller) getSession(ctx *gin.Context) session.ISession {
	//log.Printf("ctrl GetSession: sid: %s, values: %v", session.Get(ctx).SessionId(), session.Get(ctx).Values())
	return session.Get(ctx)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 保存会话值
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (ctrl *Controller) SaveSession(ctx *gin.Context, name, value string) {
	log.Printf("SaveSession name:%s value: %s", name, value)

	if len(name) == 0 || len(value) == 0 {
		return
	}

	//保存手机验证码
	session := ctrl.getSession(ctx)
	session.Set(name, strings.ToLower(value))
	session.Save()
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 校验会话值
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (ctrl *Controller) ValidateSession(ctx *gin.Context, name, value string, args ...bool) bool {
	log.Printf("ValidateSession name:%s value: %s", name, value)

	if len(name) == 0 || len(value) == 0 {
		isSuccess := false
		return isSuccess
	}

	isSuccess := true
	isRemove := true

	session := ctrl.getSession(ctx)
	if sessionValue, isOk := session.Get(name).(string); isOk {
		log.Printf("ValidateSession sessionValue: %s", sessionValue)
		if sessionValue != strings.ToLower(value) {
			isSuccess = false
		}
	} else {
		isSuccess = false
	}

	//判断是否立即销毁会话数据
	if len(args) > 0 {
		isRemove = args[0]
	}
	if isRemove {
		session.Delete(name)
		session.Save()
	} else {
		if isSuccess {
			session.Delete(name)
			session.Save()
		}
	}

	return isSuccess
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 移除会话值
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (ctrl *Controller) RemoveSession(ctx *gin.Context, name string) {
	session := ctrl.getSession(ctx)
	session.Delete(name)
	session.Save()
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 清除会话
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (ctrl *Controller) ClearSession(ctx *gin.Context) {
	ctrl.getSession(ctx).Clear()
}
