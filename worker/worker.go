package worker

import (
	"github.com/google/uuid"
	"github.com/kitanoyoru/golang-orchestrator/pkg/queue"
	"github.com/kitanoyoru/golang-orchestrator/task"
)

type Worker struct {
	Title     string
	Queue     queue.Queue
	DB        map[uuid.UUID]*task.Task
	TaskCount int
}

func (w *Worker) RunTask() {}

func (w *Worker) StartTask() {}

func (w *Worker) StopTask() {}

func (w *Worker) CollectStats() {}
