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

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * app 初始化
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func init() {
	fmt.Printf("%v ging app init\n", time.Now())
	apps = make([]*App, 0)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 注册App
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func RegisterApp(args ...*App) error {
	var err *CustomError

	for _, currentApp := range args {
		fmt.Printf("%v ging app register\n", time.Now())
		if currentApp == nil || len(currentApp.Name) == 0 {
			err = NewCustomError("App名称不能为空")
			break
		}

		isExists := false
		for _, app := range apps {
			if strings.ToLower(app.Name) == strings.ToLower(currentApp.Name) {
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
