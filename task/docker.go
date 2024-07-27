package task

import (
	"context"
	"io"
	"math"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Title         string
	AttachStdin   bool
	AttachStdout  bool
	AttachStderr  bool
	ExposedPorts  nat.PortSet
	Cmd           []string
	Image         string
	Cpu           float64
	Memory        int64
	Disk          int64
	Env           []string
	RestartPolicy container.RestartPolicyMode
}

type Docker struct {
	Client *client.Client
	Config Config
}

type DockerResult struct {
	Error       error
	Action      string
	ContainerID string
	Result      string
}

func (d *Docker) Run(ctx context.Context) DockerResult {
	reader, err := d.Client.ImagePull(ctx, d.Config.Image, image.PullOptions{})
	if err != nil {
		log.Error().Ctx(ctx).Err(err)
		return DockerResult{Error: err}
	}

	io.Copy(os.Stdout, reader)

	cc := container.Config{
		Image:        d.Config.Image,
		Tty:          false,
		Env:          d.Config.Env,
		ExposedPorts: d.Config.ExposedPorts,
	}

	hc := container.HostConfig{
		RestartPolicy: container.RestartPolicy{
			Name: d.Config.RestartPolicy,
		},
		Resources: container.Resources{
			Memory:   d.Config.Memory,
			NanoCPUs: int64(d.Config.Cpu * math.Pow(10, 9)),
		},
		PublishAllPorts: true,
	}

	resp, err := d.Client.ContainerCreate(ctx, &cc, &hc, nil, nil, d.Config.Title)
	if err != nil {
		log.Error().Ctx(ctx).Err(err)
		return DockerResult{Error: err}
	}

	err = d.Client.ContainerStart(ctx, resp.ID, container.StartOptions{})
	if err != nil {
		log.Error().Ctx(ctx).Err(err)
		return DockerResult{Error: err}
	}

	out, err := d.Client.ContainerLogs(ctx, resp.ID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	})
	if err != nil {
		log.Error().Ctx(ctx).Err(err)
		return DockerResult{Error: err}
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)

	return DockerResult{
		ContainerID: resp.ID,
		Action:      "start",
		Result:      "success",
	}
}

func (d *Docker) Stop(ctx context.Context, id string) DockerResult {
	err := d.Client.ContainerStop(ctx, id, container.StopOptions{})
	if err != nil {
		log.Error().Ctx(ctx).Err(err)
		return DockerResult{Error: err}
	}

	err = d.Client.ContainerRemove(ctx, id, container.RemoveOptions{})
	if err != nil {
		log.Error().Ctx(ctx).Err(err)
		return DockerResult{Error: err}
	}

	return DockerResult{
		ContainerID: id,
		Action:      "stop",
		Result:      "success",
	}
}
