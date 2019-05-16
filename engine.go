package ging

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/sanxia/ging/plugin/pongo"
)

const (
	DEBUG   string = "debug"
	RELEASE string = "release"
)

/* ================================================================================
* Http引擎数据结构
* qq group: 582452342
* email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
* ================================================================================ */
type (
	IHttpEngine interface {
		Engine() *gin.Engine
		Middleware(args ...gin.HandlerFunc)
		Get(groupName, path string, actionHandler ActionHandler) gin.IRoutes
		Post(groupName, path string, actionHandler ActionHandler) gin.IRoutes
		NoRoute(routeHandler gin.HandlerFunc)
		Group(path string) *gin.RouterGroup
		Render(render render.HTMLRender)
		Static(routerPath, filePath string)
	}

	httpEngine struct {
		engine *gin.Engine
	}
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 实例化HttpEngine
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func NewHttpEngine(templatePath string, model string, isDebug bool) IHttpEngine {
	if isDebug {
		gin.SetMode(gin.DebugMode)
	} else {
		switch model {
		case "debug":
			gin.SetMode(gin.DebugMode)
		case "release":
			gin.SetMode(gin.ReleaseMode)
		default:
			gin.SetMode(gin.TestMode)
		}
	}

	//初始化HttpEngine
	httpEngine := &httpEngine{
		engine: gin.New(),
	}

	//注册中间件
	httpEngine.Middleware(
		gin.Logger(),
		gin.Recovery(),
	)

	//注册渲染引擎
	httpEngine.Render(pongo.NewRender(templatePath))

	return httpEngine
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 获取gin.Engine
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *httpEngine) Engine() *gin.Engine {
	return s.engine
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * gin中间件设置
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *httpEngine) Middleware(args ...gin.HandlerFunc) {
	count := len(args)
	if count > 0 {
		for _, middleware := range args {
			s.engine.Use(middleware)
		}
	}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * gin未命中路由处理器设置
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *httpEngine) NoRoute(routeHandler gin.HandlerFunc) {
	s.engine.NoRoute(routeHandler)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * gin路由组设置
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *httpEngine) Group(path string) *gin.RouterGroup {
	if len(path) > 0 {
		return s.engine.Group(path)
	}

	return nil
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Http Get 动作
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *httpEngine) Get(groupName, path string, actionHandler ActionHandler) gin.IRoutes {
	if group := s.Group(groupName); group != nil {
		return group.GET(path, func(ctx *gin.Context) {
			actionHandler(ctx).Render()
		})
	}

	return s.engine.GET(path, func(ctx *gin.Context) {
		actionHandler(ctx).Render()
	})
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Http POST 动作
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *httpEngine) Post(groupName, path string, actionHandler ActionHandler) gin.IRoutes {
	if group := s.Group(groupName); group != nil {
		return group.POST(path, func(ctx *gin.Context) {
			actionHandler(ctx).Render()
		})
	}

	return s.engine.POST(path, func(ctx *gin.Context) {
		actionHandler(ctx).Render()
	})
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * gin渲染器设置
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *httpEngine) Render(render render.HTMLRender) {
	s.engine.HTMLRender = render
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * gin静态文件设置
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *httpEngine) Static(routerPath, filePath string) {
	s.engine.Static(routerPath, filePath)
}
