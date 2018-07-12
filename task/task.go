package task

import (
	"errors"
)

import (
	"github.com/sanxia/ging"
)

/* ================================================================================
 * ITask接口
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */
type (
	ITaskCenter interface {
		RegisterTask(channel, name string, task ITask) error
		Start(channel string)
	}

	ITask interface {
		Run(settings *ging.Settings)
	}
)

type (
	TaskCenter struct {
		tasks    map[string]map[string]ITask
		settings *ging.Settings
	}
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 实例化TaskCenter
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func NewTaskCenter(settings *ging.Settings) ITaskCenter {
	taskCenter := &TaskCenter{
		tasks:    make(map[string]map[string]ITask, 0),
		settings: settings,
	}

	return taskCenter
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 注册任务
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (t *TaskCenter) RegisterTask(channel, name string, task ITask) error {
	if len(channel) == 0 || len(name) == 0 {
		return errors.New("argments error")
	}

	if _, isOk := t.tasks[channel]; !isOk {
		t.tasks[channel] = make(map[string]ITask, 0)
	}

	if tasks, isOk := t.tasks[channel]; isOk {
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
func (t *TaskCenter) Start(channel string) {
	tasks, isOk := t.tasks[channel]
	if !isOk {
		return
	}

	for _, task := range tasks {
		go task.Run(t.settings)
	}
}
