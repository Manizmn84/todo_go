package qtodo

import (
	"errors"
	"strings"
	"time"
)

type Task interface {
	DoAction()
	GetAlarmTime() time.Time
	GetAction() func()
	GetName() string
	GetDescription() string
}

type MyTask struct {
	action      func()
	alarmTime   time.Time
	name        string
	description string
}

func NewTask(action func(), alarm time.Time, name, desc string) (Task, error) {
	if action == nil {
		return nil, errors.New("action can not be nil")
	}

	if alarm.IsZero() || alarm.Before(time.Now()) {
		return nil, errors.New("alarm time must be valid and in the future")
	}

	if strings.TrimSpace(name) == "" {
		return nil, errors.New("name can not be empty")
	}

	if strings.TrimSpace(desc) == "" {
		return nil, errors.New("description can not be empty")
	}

	return &MyTask{
		action:      action,
		alarmTime:   alarm,
		name:        name,
		description: desc,
	}, nil
}

func (t *MyTask) DoAction() {
	t.action()
}

func (t *MyTask) GetAlarmTime() time.Time {
	return t.alarmTime
}

func (t *MyTask) GetAction() func() {
	return t.action
}

func (t *MyTask) GetName() string {
	return t.name
}

func (t *MyTask) GetDescription() string {
	return t.description
}
