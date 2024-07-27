package manager

import (
	"github.com/google/uuid"
	"github.com/kitanoyoru/golang-orchestrator/pkg/queue"
	"github.com/kitanoyoru/golang-orchestrator/task"
)

type Manager struct {
	Pending       queue.Queue
	TaskDB        map[string][]*task.Task
	EventDB       map[string][]*task.TaskEvent
	Workers       []string
	WorkerTaskMap map[string][]uuid.UUID
	TaskWorkerMap map[uuid.UUID]string
}

func (m *Manager) SelectWorker() {}

func (m *Manager) UpdateTasks() {}

func (m *Manager) SendWork() {}
