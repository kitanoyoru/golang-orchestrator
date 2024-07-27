package task

import (
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
	"github.com/kitanoyoru/golang-orchestrator/task/state"
)

type Task struct {
	ID            uuid.UUID
	Title         string
	State         state.State
	Image         string
	Cpu           float64
	Memory        int
	Disk          int
	ExposedPorts  nat.PortSet
	PortBindings  map[string]string
	RestartPolicy string
	Runtime       *Runtime
	StartTime     time.Time
	FinishTime    time.Time
}

type Runtime struct {
	ContainerID string
}

type TaskEvent struct {
	ID        uuid.UUID
	State     state.State
	Timestamp time.Time
	Task      Task
}
