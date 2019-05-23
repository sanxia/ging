package ging

import (
//"log"
)

import (
	"github.com/gin-gonic/gin"
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
		IHttpAction

		Action(actionHandler ActionHandler, args ...interface{}) func(*gin.Context)
		Filter(filters ...IActionFilter) IController

		SaveSession(ctx *gin.Context, name, value string)
		GetSession(ctx *gin.Context, name string) string
		ValidateSession(ctx *gin.Context, name, value string, args ...bool) bool
		RemoveSession(ctx *gin.Context, name string)
		ClearSession(ctx *gin.Context)

		GetToken(ctx *gin.Context) IToken
	}
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 控制器数据结构
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
type (
	Controller struct {
		GroupName string
		Engine    IHttpEngine
		filters   []IActionFilter
	}
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 初始化控制器
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func NewController(groupName string, engine IHttpEngine, args ...IActionFilter) IController {
	controller := &Controller{
		GroupName: groupName,
		Engine:    engine,
	}

	return controller.Filter(args...)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * IHttpAction接口实现 － Http Get请求
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (ctrl *Controller) Get(path string, actionHandler ActionHandler, args ...interface{}) gin.IRoutes {
	handlerFunc := ctrl.Action(actionHandler, args...)

	if group := ctrl.Engine.Group(ctrl.GroupName); group != nil {
		return group.GET(path, handlerFunc)
	}

	return ctrl.Engine.Engine().GET(path, handlerFunc)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * IHttpAction接口实现 － Http Post请求
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (ctrl *Controller) Post(path string, actionHandler ActionHandler, args ...interface{}) gin.IRoutes {
	handlerFunc := ctrl.Action(actionHandler, args...)

	if group := ctrl.Engine.Group(ctrl.GroupName); group != nil {
		return group.POST(path, handlerFunc)
	}

	return ctrl.Engine.Engine().POST(path, handlerFunc)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * IHttpAction接口实现 － Http Delete请求
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (ctrl *Controller) Delete(path string, actionHandler ActionHandler, args ...interface{}) gin.IRoutes {
	handlerFunc := ctrl.Action(actionHandler, args...)

	if group := ctrl.Engine.Group(ctrl.GroupName); group != nil {
		return group.DELETE(path, handlerFunc)
	}

	return ctrl.Engine.Engine().DELETE(path, handlerFunc)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * IHttpAction接口实现 － Http Patch请求
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (ctrl *Controller) Patch(path string, actionHandler ActionHandler, args ...interface{}) gin.IRoutes {
	handlerFunc := ctrl.Action(actionHandler, args...)

	if group := ctrl.Engine.Group(ctrl.GroupName); group != nil {
		return group.PATCH(path, handlerFunc)
	}

	return ctrl.Engine.Engine().PATCH(path, handlerFunc)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * IHttpAction接口实现 － Http Options请求
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (ctrl *Controller) Options(path string, actionHandler ActionHandler, args ...interface{}) gin.IRoutes {
	handlerFunc := ctrl.Action(actionHandler, args...)

	if group := ctrl.Engine.Group(ctrl.GroupName); group != nil {
		return group.OPTIONS(path, handlerFunc)
	}

	return ctrl.Engine.Engine().OPTIONS(path, handlerFunc)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * IHttpAction接口实现 － Http Head请求
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (ctrl *Controller) Head(path string, actionHandler ActionHandler, args ...interface{}) gin.IRoutes {
	handlerFunc := ctrl.Action(actionHandler, args...)

	if group := ctrl.Engine.Group(ctrl.GroupName); group != nil {
		return group.HEAD(path, handlerFunc)
	}

	return ctrl.Engine.Engine().HEAD(path, handlerFunc)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 控制器动作
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (ctrl *Controller) Action(actionHandler ActionHandler, args ...interface{}) func(*gin.Context) {
	return func(context *gin.Context) {
		var actionFilters IActionFilterList
		var filterResult IActionResult

		argsCount := len(args)
		isFilterEnabled := true

		if argsCount > 0 {
			for _, arg := range args {
				if arg == nil {
					continue
				}

				//判断是否禁用当前动作过滤器
				if actionFilter, isOk := arg.(IActionFilter); isOk {
					if len(actionFilters) == 0 {
						actionFilters = make(IActionFilterList, argsCount)
					}

					actionFilters = append(actionFilters, actionFilter)
				} else {
					if isFilterEnabledValue, isOk := arg.(bool); isOk {
						if !isFilterEnabledValue {
							isFilterEnabled = false
							break
						}
					}
				}
			}
		}

		if isFilterEnabled {
			//顺序执行之前的控制器过滤器
			//Before返回非空IActionResult即终止
			for _, filter := range ctrl.filters {
				if filter != nil {
					if filterResult = filter.Before(context); filterResult != nil {
						break
					}
				}
			}

			if filterResult == nil {
				//顺序执行之前的动作过滤器
				for _, filter := range actionFilters {
					if filter != nil {
						if filterResult = filter.Before(context); filterResult != nil {
							break
						}
					}
				}
			}
		}

		if filterResult != nil {
			filterResult.Render()
			context.Abort()
			return
		} else {
			actionHandler(context).Render()
		}

		if isFilterEnabled {
			//逆序执行之后的动作过滤器
			for i := len(actionFilters) - 1; i >= 0; i-- {
				actionFilters[i].After(context)
			}

			//逆序执行之后的控制器过滤器
			for i := len(ctrl.filters) - 1; i >= 0; i-- {
				ctrl.filters[i].After(context)
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
		if filter != nil {
			ctrl.filters = append(ctrl.filters, filter)
		}
	}

	return ctrl
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 获取会话对象
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (ctrl *Controller) getSession(ctx *gin.Context) session.ISession {
	newSession := session.NewSession(ctx)
	//log.Printf("ctrl getSession: id: %s, values: %v", newSession.SessionId(), newSession.Values())

	return newSession
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 保存会话值
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (ctrl *Controller) SaveSession(ctx *gin.Context, name, value string) {
	if len(name) == 0 || len(value) == 0 {
		return
	}

	session := ctrl.getSession(ctx)
	session.Set(name, value)
	session.Save()
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 获取会话值
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (ctrl *Controller) GetSession(ctx *gin.Context, name string) string {
	value := ""
	if len(name) == 0 {
		return value
	}

	if sessionValue, isOk := ctrl.getSession(ctx).Get(name).(string); isOk {
		value = sessionValue
	}

	return value
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 校验会话值
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (ctrl *Controller) ValidateSession(ctx *gin.Context, name, value string, args ...bool) bool {
	isSuccess := true
	isRemove := true

	if len(name) == 0 || len(value) == 0 {
		isSuccess = false
		return isSuccess
	}

	session := ctrl.getSession(ctx)
	if sessionValue, isOk := session.Get(name).(string); isOk {
		if sessionValue != value {
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
	session := ctrl.getSession(ctx)
	session.Clear()
	session.Save()
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 获取Token接口
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (ctrl *Controller) GetToken(ctx *gin.Context) IToken {
	var userToken IToken

	if ctx != nil {
		if userIdentity, isOk := ctx.Get(USER_IDENTITY); userIdentity != nil && isOk {
			if tokenIdentity, isOk := userIdentity.(*Token); isOk {
				userToken = tokenIdentity
			}
		}
	}

	return userToken
}
