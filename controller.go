package ging

import (
	"fmt"
	"log"
	"time"
)

import (
	"github.com/gin-gonic/gin"
	"github.com/sanxia/ging/middleware/session"
)

/* ================================================================================
 * controller
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * controller interface
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
type (
	IController interface {
		IHttpAction

		Action(actionHandler ActionHandler, args ...interface{}) func(*gin.Context)
		Filter(filters ...IActionFilter) IController

		SaveSession(ctx *gin.Context, name, value string)
		GetSession(ctx *gin.Context, name string) string
		ValidateSession(ctx *gin.Context, name, value string, args ...interface{}) bool
		RemoveSession(ctx *gin.Context, name string)
		ClearSession(ctx *gin.Context)

		GetToken(ctx *gin.Context) IToken
		GetApp() IApp
	}
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * controller data structure
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
type (
	Controller struct {
		GroupName string
		Engine    IHttpEngine
		filters   []IActionFilter
	}
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * instantiating the controller
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func NewController(groupName string, engine IHttpEngine, args ...IActionFilter) IController {
	fmt.Printf("%v ging engine controller %s instantiating\n", time.Now(), groupName)

	controller := &Controller{
		GroupName: groupName,
		Engine:    engine,
	}

	return controller.Filter(args...)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * IHttpAction interface implementation － Http Get
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (ctrl *Controller) Get(path string, actionHandler ActionHandler, args ...interface{}) gin.IRoutes {
	handlerFunc := ctrl.Action(actionHandler, args...)

	if group := ctrl.Engine.Group(ctrl.GroupName); group != nil {
		return group.GET(path, handlerFunc)
	}

	return ctrl.Engine.Engine().GET(path, handlerFunc)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * IHttpAction interface implementation － Http Post
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (ctrl *Controller) Post(path string, actionHandler ActionHandler, args ...interface{}) gin.IRoutes {
	handlerFunc := ctrl.Action(actionHandler, args...)

	if group := ctrl.Engine.Group(ctrl.GroupName); group != nil {
		return group.POST(path, handlerFunc)
	}

	return ctrl.Engine.Engine().POST(path, handlerFunc)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * IHttpAction interface implementation － Http Delete
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (ctrl *Controller) Delete(path string, actionHandler ActionHandler, args ...interface{}) gin.IRoutes {
	handlerFunc := ctrl.Action(actionHandler, args...)

	if group := ctrl.Engine.Group(ctrl.GroupName); group != nil {
		return group.DELETE(path, handlerFunc)
	}

	return ctrl.Engine.Engine().DELETE(path, handlerFunc)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * IHttpAction interface implementation － Http Patch
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (ctrl *Controller) Patch(path string, actionHandler ActionHandler, args ...interface{}) gin.IRoutes {
	handlerFunc := ctrl.Action(actionHandler, args...)

	if group := ctrl.Engine.Group(ctrl.GroupName); group != nil {
		return group.PATCH(path, handlerFunc)
	}

	return ctrl.Engine.Engine().PATCH(path, handlerFunc)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * IHttpAction interface implementation － Http Options
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (ctrl *Controller) Options(path string, actionHandler ActionHandler, args ...interface{}) gin.IRoutes {
	handlerFunc := ctrl.Action(actionHandler, args...)

	if group := ctrl.Engine.Group(ctrl.GroupName); group != nil {
		return group.OPTIONS(path, handlerFunc)
	}

	return ctrl.Engine.Engine().OPTIONS(path, handlerFunc)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * IHttpAction interface implementation － Http Head
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (ctrl *Controller) Head(path string, actionHandler ActionHandler, args ...interface{}) gin.IRoutes {
	handlerFunc := ctrl.Action(actionHandler, args...)

	if group := ctrl.Engine.Group(ctrl.GroupName); group != nil {
		return group.HEAD(path, handlerFunc)
	}

	return ctrl.Engine.Engine().HEAD(path, handlerFunc)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * controller action
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (ctrl *Controller) Action(actionHandler ActionHandler, args ...interface{}) func(*gin.Context) {
	return func(ctx *gin.Context) {
		var actionFilters IActionFilterList
		var filterResult IActionResult

		argsCount := len(args)
		isFilterEnabled := true

		if argsCount > 0 {
			for _, arg := range args {
				if arg == nil {
					continue
				}

				//determine if the current motion filter is disabled
				if actionFilter, isOk := arg.(IActionFilter); isOk {
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
			//the controller filter before sequential execution
			//Before returns non-empty IActionResult is terminated
			for _, filter := range ctrl.filters {
				if filter != nil {
					if filterResult = filter.Before(ctx); filterResult != nil {
						break
					}
				}
			}

			if filterResult == nil {
				//the action filter before the sequential execution
				for _, filter := range actionFilters {
					if filter != nil {
						if filterResult = filter.Before(ctx); filterResult != nil {
							break
						}
					}
				}
			}
		}

		if filterResult != nil {
			filterResult.Render()
			ctx.Abort()
			return
		} else {
			actionHandler(ctx).Render()
		}

		if isFilterEnabled {
			//action filter after reverse sequence execution
			for i := len(actionFilters) - 1; i >= 0; i-- {
				actionFilters[i].After(ctx)
			}

			//controller filter after reverse sequence execution
			for i := len(ctrl.filters) - 1; i >= 0; i-- {
				ctrl.filters[i].After(ctx)
			}
		}
	}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * set up the controller filter
 * The filter interface is executed before and after the controller's method is executed
 * and the filter interface collection does not support sorting
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
 * get session objects
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (ctrl *Controller) getSession(ctx *gin.Context) session.ISession {
	newSession := session.NewSession(ctx)

	log.Printf("ctrl getSession values: %v", newSession.Values())
	//log.Printf("ctrl getSession: id: %s, values: %v", newSession.SessionId(), newSession.Values())

	return newSession
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * save session values
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (ctrl *Controller) SaveSession(ctx *gin.Context, name, value string) {
	if len(name) == 0 || len(value) == 0 {
		return
	}

	currentSession := ctrl.getSession(ctx)
	currentSession.Set(name, value)
	currentSession.Save()
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * get session values
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
 * check session value
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (ctrl *Controller) ValidateSession(ctx *gin.Context, name, value string, args ...interface{}) bool {
	isSuccess := true
	isRemove := true

	if len(name) == 0 || len(value) == 0 {
		isSuccess = false
		return isSuccess
	}

	currentSession := ctrl.getSession(ctx)
	sessionValue, isSessionValueOk := currentSession.Get(name).(string)
	if !isSessionValueOk {
		isSuccess = false
	}

	//determine whether session data is destroyed immediately
	if len(args) > 0 {
		if _isRemove, isOk := args[0].(bool); isOk {
			if sessionValue != value {
				isSuccess = false
			}

			isRemove = _isRemove
		} else if validateHandler, isOk := args[0].(func(string) (bool, error)); isOk {
			if _isRemove, isError := validateHandler(sessionValue); isError != nil {
				isSuccess = false
			} else {
				isRemove = _isRemove
			}
		} else {
			isSuccess = false
		}
	} else {
		if sessionValue != value {
			isSuccess = false
		}
	}

	if isSuccess {
		currentSession.Delete(name)
		currentSession.Save()
	} else {
		if isRemove {
			currentSession.Delete(name)
			currentSession.Save()
		}
	}

	return isSuccess
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * remove session values
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (ctrl *Controller) RemoveSession(ctx *gin.Context, name string) {
	currentSession := ctrl.getSession(ctx)
	currentSession.Delete(name)
	currentSession.Save()
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * clear session
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (ctrl *Controller) ClearSession(ctx *gin.Context) {
	currentSession := ctrl.getSession(ctx)
	currentSession.Clear()
	currentSession.Save()
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * get IToken interface
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (ctrl *Controller) GetToken(ctx *gin.Context) IToken {
	return GetToken(ctx)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Get the IApp interface
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (ctrl *Controller) GetApp() IApp {
	return GetApp()
}
