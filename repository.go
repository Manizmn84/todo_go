package qtodo

import (
	"errors"
	"sync"
)

type Database interface {
	GetTaskList() []Task
	GetTask(string) (Task, error)
	SaveTask(Task) error
	DelTask(string) error
}

type MemoryDatabase struct {
	mu    sync.RWMutex
	tasks map[string]MyTask
}

func NewDatabase() *MemoryDatabase {
	return &MemoryDatabase{
		tasks: make(map[string]MyTask),
	}
}

func (db *MemoryDatabase) GetTaskList() []MyTask {
	db.mu.RLock()
	defer db.mu.RUnlock()

	list := make([]MyTask, 0, len(db.tasks))
	for _, t := range db.tasks {
		list = append(list, t)
	}
	return list
}

func (db *MemoryDatabase) GetTask(name string) (MyTask, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	task, exists := db.tasks[name]
	if !exists {
		return MyTask{}, errors.New("task not found")
	}
	return task, nil
}

func (db *MemoryDatabase) SaveTask(task MyTask) error {
	if task.name == "" {
		return errors.New("task name cannot be empty")
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	db.tasks[task.name] = task
	return nil
}

func (db *MemoryDatabase) DelTask(name string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	_, exists := db.tasks[name]
	if !exists {
		return errors.New("task not found")
	}
	delete(db.tasks, name)
	return nil
}
