package ging

import (
	"fmt"
	"time"
)

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/sanxia/ging/plugin/pongo"
)

/* ================================================================================
 * http engine
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
 * instantiating http engine
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func NewHttpEngine(templatePath string, model string, isDebug bool) IHttpEngine {
	fmt.Printf("%v ging engine instantiating\n", time.Now())

	if isDebug {
		gin.SetMode(gin.DebugMode)
	} else {
		switch model {
		case DEBUG:
			gin.SetMode(gin.DebugMode)
		case RELEASE:
			gin.SetMode(gin.ReleaseMode)
		default:
			gin.SetMode(gin.TestMode)
		}
	}

	//initialization http engine
	httpEngine := &httpEngine{
		engine: gin.New(),
	}

	//register middleware
	httpEngine.Middleware(
		gin.Logger(),
		gin.Recovery(),
	)

	//register the rendering template engine
	httpEngine.Render(pongo.NewPongoTemplate(&pongo.PongoOption{
		TemplatePath: templatePath,
	}))

	return httpEngine
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * get gin.Engine
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *httpEngine) Engine() *gin.Engine {
	return s.engine
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * set gin middleware
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
 * set not found route processor
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *httpEngine) NoRoute(routeHandler gin.HandlerFunc) {
	s.engine.NoRoute(routeHandler)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * set route group
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *httpEngine) Group(path string) *gin.RouterGroup {
	if len(path) > 0 {
		return s.engine.Group(path)
	}

	return nil
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * http get action
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
 * http post action
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
 * set html renderer
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *httpEngine) Render(render render.HTMLRender) {
	s.engine.HTMLRender = render
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * set static file
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *httpEngine) Static(routerPath, filePath string) {
	s.engine.Static(routerPath, filePath)
}
