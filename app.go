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
	App struct {
		Name     string
		Router   IHttpRouter
		Settings *Settings
	}
)

var (
	apps []*App
)

func init() {
	fmt.Printf("%v ging app init\n", time.Now())
	apps = make([]*App, 0)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 注册App
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func RegisterApp(app *App) error {
	var err *CustomError

	fmt.Printf("%v ging app register\n", time.Now())
	if app == nil {
		err = NewCustomError("App不能为空")
	}

	if len(app.Name) == 0 {
		err = NewCustomError("App名称不能为空")
	}

	isExists := false
	for _, app := range apps {
		if app.Name == strings.ToLower(app.Name) {
			isExists = true
			break
		}
	}

	if !isExists {
		apps = append(apps, app)
	}

	return err
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 设置App配置
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func SetAppSettings(name string, settings *Settings) error {
	app, err := GetApp(name)
	if err == nil {
		app.Settings = settings
	}

	return err
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 获取App
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func GetApp(appName string) (*App, error) {
	if len(apps) == 0 {
		return nil, NewCustomError("未找到有效的App")
	}

	if len(appName) == 0 {
		return nil, NewCustomError("App名称参数不能为空")
	}

	for _, app := range apps {
		if app.Name == strings.ToLower(appName) {
			return app, nil
		}
	}

	return nil, NewCustomError("未找到有效的App")
}
