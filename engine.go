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
		Render(render render.HTMLRender)
		Group(path string) *gin.RouterGroup
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

	//初始化httpEngine
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
 * gin渲染器设置
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *httpEngine) Render(render render.HTMLRender) {
	s.engine.HTMLRender = render
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * gin路由组设置
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *httpEngine) Group(path string) *gin.RouterGroup {
	return s.engine.Group(path)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * gin静态文件设置
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *httpEngine) Static(routerPath, filePath string) {
	s.engine.Static(routerPath, filePath)
}
