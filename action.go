package ging

import (
	"github.com/gin-gonic/gin"
)

/* ================================================================================
 * Http Action
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Http请求动作接口
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
type (
	ActionHandler func(*gin.Context) IActionResult

	IHttpAction interface {
		Get(path string, actionHandler ActionHandler, args ...interface{}) gin.IRoutes
		Post(path string, actionHandler ActionHandler, args ...interface{}) gin.IRoutes
		Delete(path string, actionHandler ActionHandler, args ...interface{}) gin.IRoutes
		Patch(path string, actionHandler ActionHandler, args ...interface{}) gin.IRoutes
		Options(path string, actionHandler ActionHandler, args ...interface{}) gin.IRoutes
		Head(path string, actionHandler ActionHandler, args ...interface{}) gin.IRoutes
	}
)
