package worker

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/kitanoyoru/golang-orchestrator/pkg/queue"
	"github.com/kitanoyoru/golang-orchestrator/pkg/types"
	"github.com/kitanoyoru/golang-orchestrator/task"
	"github.com/kitanoyoru/golang-orchestrator/task/service/docker"
	"github.com/kitanoyoru/golang-orchestrator/task/state"
	"github.com/rs/zerolog/log"
)

var (
	ErrTaskContainerNotStarted     = errors.New("task container not started yet")
	ErrTaskContainerAlreadyStarted = errors.New("task container already started")
	ErrTaskInvalidStateTransition  = errors.New("task state cannot be transitioned to the specified state")
)

type Worker interface {
	AddTask(t *task.Task)
	RunTask(ctx context.Context) docker.DockerResult
	CollectStats()
}

type worker struct {
	Title     string
	Queue     queue.Queue
	DB        map[uuid.UUID]*task.Task
	TaskCount int
}

func NewWorker() Worker {
	return &worker{
		DB: make(map[uuid.UUID]*task.Task),
	}
}

func (w *worker) AddTask(t *task.Task) {
	w.Queue.Enqueue(t)
}

func (w *worker) RunTask(ctx context.Context) docker.DockerResult {
	t := w.Queue.Dequeue()
	if t == nil {
		log.Warn().Str("Worker", "RunTask").Msg("worker empty queue")
		return docker.DockerResult{}
	}

	taskQueued := t.(task.Task)

	taskPersisted, exists := w.DB[taskQueued.ID]
	if !exists {
		taskPersisted = &taskQueued
		w.DB[taskQueued.ID] = taskPersisted
	}

	if !state.ValidStateTransition(taskPersisted.State, taskQueued.State) {
		return docker.DockerResult{Error: ErrTaskInvalidStateTransition}
	}

	switch taskQueued.State {
	case state.Scheduled:
		return w.startTask(ctx, taskQueued)
	case state.Completed:
		return w.stopTask(ctx, taskQueued)
	default:
		return docker.DockerResult{Error: errors.New("unsupported task state")}
	}
}

func (w *worker) startTask(ctx context.Context, t task.Task) docker.DockerResult {
	if t.Runtime != nil && !t.StartTime.IsZero() && t.State == state.Running {
		return docker.DockerResult{Error: ErrTaskContainerAlreadyStarted}
	}

	d, err := setupDockerClientFromTask(ctx, &t)
	if err != nil {
		return docker.DockerResult{Error: err}
	}

	now := time.Now().UTC()

	res := d.Run(ctx)
	if res.Error != nil {
		w.updateTaskState(&t, state.Failed, res.ContainerID, types.Time(now), types.Time(now))
		return res
	}

	w.updateTaskState(&t, state.Running, res.ContainerID, types.Time(now), nil)

	return res
}

func (w *worker) stopTask(ctx context.Context, t task.Task) docker.DockerResult {
	if t.Runtime == nil && t.StartTime.IsZero() && t.State != state.Running {
		return docker.DockerResult{Error: ErrTaskContainerNotStarted}
	}

	d, err := setupDockerClientFromTask(ctx, &t)
	if err != nil {
		return docker.DockerResult{Error: err}
	}

	now := time.Now().UTC()

	res := d.Stop(ctx, t.Runtime.ContainerID)
	if res.Error != nil {
		w.updateTaskState(&t, state.Failed, res.ContainerID, types.Time(now), types.Time(now))
		return res
	}

	w.updateTaskState(&t, state.Completed, res.ContainerID, nil, types.Time(now))

	return res
}

func (w *worker) updateTaskState(t *task.Task, state state.State, containerID string, startedAt, finishedAt *time.Time) {
	t.State = state
	t.StartTime = time.Now().UTC()
	t.Runtime = &task.Runtime{ContainerID: containerID}
	if startedAt != nil {
		t.StartTime = *startedAt
	}
	if finishedAt != nil {
		t.StartTime = *finishedAt
	}
	w.DB[t.ID] = t
}

func setupDockerClientFromTask(ctx context.Context, t *task.Task) (docker.Docker, error) {
	config := docker.NewConfig(docker.FromTask(t))
	d, err := docker.NewDocker(config)
	if err != nil {
		log.Error().Ctx(ctx).Err(err).Send()
	}
	return d, err
}

func (w *worker) CollectStats() {}
