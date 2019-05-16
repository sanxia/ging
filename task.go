package ging

import (
	"errors"
)

/* ================================================================================
 * ITask接口
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */
type (
	ITaskCenter interface {
		RegisterTask(appName, name string, task ITask) error
		Start(app IApp)
	}

	ITask interface {
		Run(app IApp)
	}
)

type (
	taskCenter struct {
		tasks map[string]map[string]ITask
	}
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 实例化TaskCenter
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func NewTaskCenter() ITaskCenter {
	return &taskCenter{
		tasks: make(map[string]map[string]ITask, 0),
	}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 注册任务
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (t *taskCenter) RegisterTask(appName, name string, task ITask) error {
	if len(appName) == 0 || len(name) == 0 {
		return errors.New("argments error")
	}

	if _, isOk := t.tasks[appName]; !isOk {
		t.tasks[appName] = make(map[string]ITask, 0)
	}

	if tasks, isOk := t.tasks[appName]; isOk {
		if _, isOk := tasks[name]; isOk {
			return errors.New("task name is exists")
		} else {
			tasks[name] = task
		}
	}

	return nil
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 启动任务
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (t *taskCenter) Start(app IApp) {
	tasks, isOk := t.tasks[app.GetName()]
	if !isOk {
		return
	}

	for _, task := range tasks {
		go task.Run(app)
	}
}
