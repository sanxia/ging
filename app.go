package ging

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

/* ================================================================================
 * ging web integration framework
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */
type (
	IApp interface {
		GetName() string
		GetSetting() *Setting
		GetRouter() IHttpRouter

		RegisterTask(task ITask) error
		GetTasks() []ITask
	}
)

type (
	app struct {
		name       string
		setting    *Setting
		router     IHttpRouter
		tasks      []ITask
		IsActiving bool
	}
)

var (
	apps []IApp
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * init
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func init() {
	fmt.Printf("%v ging engine app init\n", time.Now())
	apps = make([]IApp, 0)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * instantiating app
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func NewApp(name string, setting *Setting, router IHttpRouter) IApp {
	return &app{
		name:    name,
		setting: setting,
		router:  router,
		tasks:   make([]ITask, 0),
	}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * register app
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func RegisterApp(args ...IApp) error {
	fmt.Printf("%v ging engine app register\n", time.Now())

	var err *CustomError

	for _, currentApp := range args {
		if currentApp == nil || len(currentApp.GetName()) == 0 {
			err = NewCustomError("app name is not found")
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
 * get app
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func GetApp(args ...string) IApp {
	var currentApp IApp

	if len(apps) != 0 {
		if len(args) == 0 {
			for _, _currentApp := range apps {
				if _currentApp, isOk := _currentApp.(*app); isOk && _currentApp.IsActiving {
					currentApp = _currentApp
					break
				}
			}
		} else {
			for _, _currentApp := range apps {
				if _currentApp.GetName() == strings.ToLower(args[0]) {
					currentApp = _currentApp
					break
				}
			}
		}
	}

	return currentApp
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * get app name
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *app) GetName() string {
	return s.name
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * get app setting
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *app) GetSetting() *Setting {
	return s.setting
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * get http router
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *app) GetRouter() IHttpRouter {
	return s.router
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * register task
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *app) RegisterTask(task ITask) error {
	if task == nil || len(task.GetName()) == 0 {
		return errors.New("argments error")
	}

	isExists := false
	for _, oldTask := range s.tasks {
		if oldTask.GetName() == task.GetName() {
			isExists = true
			break
		}
	}

	if !isExists {
		s.tasks = append(s.tasks, task)
	}

	return nil
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * get a task collection
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *app) GetTasks() []ITask {
	return s.tasks
}
