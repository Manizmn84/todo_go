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
	tasks map[string]Task
}

func NewDatabase() *MemoryDatabase {
	return &MemoryDatabase{
		tasks: make(map[string]Task),
	}
}

func (db *MemoryDatabase) GetTaskList() []Task {
	db.mu.RLock()
	defer db.mu.RUnlock()

	list := make([]Task, 0, len(db.tasks))
	for _, t := range db.tasks {
		list = append(list, t)
	}
	return list
}

func (db *MemoryDatabase) GetTask(name string) (Task, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	task, exists := db.tasks[name]
	if !exists {
		return nil, errors.New("task not found")
	}
	return task, nil
}

func (db *MemoryDatabase) SaveTask(task Task) error {
	if task.GetName() == "" {
		return errors.New("task name cannot be empty")
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	db.tasks[task.GetName()] = task
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
