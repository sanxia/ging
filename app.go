package ging

import (
	"fmt"
	"strings"
	"time"
)

/* ================================================================================
 * ging web framework
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */

type (
	IApp interface {
		GetName() string
		GetSettings() *Settings
		GetRouter() IHttpRouter
	}
)

type (
	app struct {
		name     string
		settings *Settings
		router   IHttpRouter
	}
)

var (
	apps []IApp
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * app 初始化
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func init() {
	fmt.Printf("%v ging app init\n", time.Now())
	apps = make([]IApp, 0)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 初始化App
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func NewApp(name string, settings *Settings, router IHttpRouter) IApp {
	return &app{
		name:     name,
		settings: settings,
		router:   router,
	}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 获取App Name
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *app) GetName() string {
	return s.name
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 获取App Settings
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *app) GetSettings() *Settings {
	return s.settings
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 获取Http Router
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *app) GetRouter() IHttpRouter {
	return s.router
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 注册App
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func RegisterApp(args ...IApp) error {
	fmt.Printf("%v ging app register\n", time.Now())

	var err *CustomError

	for _, currentApp := range args {
		if currentApp == nil || len(currentApp.GetName()) == 0 {
			err = NewCustomError("App名称不能为空")
			break
		}

		isExists := false
		for _, app := range apps {
			if strings.ToLower(app.GetName()) == strings.ToLower(currentApp.GetName()) {
				isExists = true
				break
			}
		}

		if !isExists {
			apps = append(apps, currentApp)
		}
	}

	return err
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 获取App
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func GetApp(name string) IApp {
	fmt.Printf("%v ging get app\n", time.Now())

	if len(apps) != 0 && len(name) != 0 {
		for _, app := range apps {
			if app.GetName() == strings.ToLower(name) {
				return app
			}
		}
	}

	return nil
}
