package service

import (
	"context"
	"math"

	"github.com/docker/go-connections/nat"
	"github.com/kitanoyoru/golang-orchestrator/pkg/types"
	"github.com/pkg/errors"
)

var (
	defaultResourceLimitation = ResourceLimitation{
		CPU:    types.Float64(0.5),
		Memory: types.Int64(25),
		Disk:   types.Int64(0),
	}
	defaultRestartPolicy = RestartPolicyDisabled
)

type RestartPolicyMode string

const (
	RestartPolicyDisabled      RestartPolicyMode = "no"
	RestartPolicyAlways        RestartPolicyMode = "always"
	RestartPolicyOnFailure     RestartPolicyMode = "on-failure"
	RestartPolicyUnlessStopped RestartPolicyMode = "unless-stopped"
)

type ResourceLimitation struct {
	CPU    *float64
	Memory *int64
	Disk   *int64
}

func (r *ResourceLimitation) GetNanoCPUs() (int64, error) {
	if r.CPU == nil {
		return 0, errors.New("cpu is nil")
	}
	return int64(*r.CPU * math.Pow(10, 9)), nil
}

type RunOptions struct {
	Title         *string
	AttachStdin   bool
	AttachStdout  bool
	AttachStderr  bool
	ExposedPorts  nat.PortSet
	Cmd           []string
	Env           []string
	Limit         *ResourceLimitation
	RestartPolicy *RestartPolicyMode
}

type RunOption func(*RunOptions)

func WithTitle(title string) RunOption {
	return func(opts *RunOptions) {
		opts.Title = &title
	}
}

func WithResourceLimitation(limit *ResourceLimitation) RunOption {
	return func(opts *RunOptions) {
		opts.Limit = limit
	}
}

func WithCmd(cmd []string) RunOption {
	return func(opts *RunOptions) {
		opts.Cmd = cmd
	}
}

func WithExposedPorts(ports nat.PortSet) RunOption {
	return func(opts *RunOptions) {
		opts.ExposedPorts = ports 
	}
}

func WithRestartPolicy(policy RestartPolicyMode) RunOption {
	return func(opts *RunOptions) {
		opts.RestartPolicy = &policy
	}
}

func WithEnv(env []string) RunOption {
	return func(opts *RunOptions) {
		opts.Env = env 
	}
}

type CRI interface {
	Run(ctx context.Context, image string, options ...RunOption) (string, error)
	Stop(ctx context.Context, containerID string) error
}

