package qtodo

import (
	"errors"
	"sync"
	"time"
)

type runningTask struct {
	stopChan chan bool
	isActive bool
	temp     bool
}

type App interface {
	StartTask(string) error
	StopTask(string)
	AddTask(name string, desc string, alarm time.Time, action func(), isTemp bool) error
	DelTask(string) error
	GetTaskList() []Task
	GetActiveTaskList() []Task
	GetTask(string) (Task, error)
}

type MyApp struct {
	db     Database
	mu     sync.RWMutex
	runner map[string]*runningTask
}

func NewApp(db Database) *MyApp {
	return &MyApp{
		db:     db,
		runner: make(map[string]*runningTask),
	}
}

func (a *MyApp) AddTask(name, desc string, alarm time.Time, action func(), isTemp bool) error {
	task, err := NewTask(action, alarm, name, desc)
	if err != nil {
		return err
	}

	err = a.db.SaveTask(task)
	if err != nil {
		return err
	}

	a.mu.Lock()
	defer a.mu.Unlock()
	if _, ok := a.runner[name]; !ok {
		a.runner[name] = &runningTask{
			stopChan: make(chan bool),
			isActive: false,
			temp:     isTemp,
		}
	}

	return nil
}

func (a *MyApp) DelTask(name string) error {
	a.StopTask(name)

	err := a.db.DelTask(name)
	if err != nil {
		return err
	}

	a.mu.Lock()
	delete(a.runner, name)
	a.mu.Unlock()

	return nil
}

func (a *MyApp) GetTaskList() []Task {
	return a.db.GetTaskList()
}

func (a *MyApp) GetActiveTaskList() []Task {
	a.mu.RLock()
	defer a.mu.RUnlock()

	out := []Task{}
	for name, rt := range a.runner {
		if rt.isActive {
			if t, err := a.db.GetTask(name); err == nil {
				out = append(out, t)
			}
		}
	}
	return out
}

func (a *MyApp) GetTask(name string) (Task, error) {
	return a.db.GetTask(name)
}

func (a *MyApp) StartTask(name string) error {
	task, err := a.db.GetTask(name)
	if err != nil {
		return err
	}

	a.mu.Lock()
	r, ok := a.runner[name]
	if !ok {
		a.mu.Unlock()
		return errors.New("task runner not found")
	}

	if r.isActive {
		a.mu.Unlock()
		return errors.New("task already running")
	}

	r.isActive = true
	r.stopChan = make(chan bool)
	stopChan := r.stopChan
	isTemp := r.temp
	a.mu.Unlock()

	go func() {
		now := time.Now()
		wait := task.GetAlarmTime().Sub(now)

		if wait < 0 {
			wait = 0
		}

		timer := time.NewTimer(wait)
		defer timer.Stop()

		select {
		case <-timer.C:
			a.mu.RLock()
			active := r.isActive
			a.mu.RUnlock()

			if active {
				task.DoAction()
				if isTemp {
					a.DelTask(name)
				}
			}

		case <-stopChan:
			return
		}
	}()

	return nil
}

func (a *MyApp) StopTask(name string) {
	a.mu.Lock()
	defer a.mu.Unlock()

	r, ok := a.runner[name]
	if !ok {
		return
	}

	if r.isActive {
		r.isActive = false
		close(r.stopChan)
	}
}
