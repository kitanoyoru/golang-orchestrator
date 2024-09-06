package docker

import (
	"context"
	"io"
	"os"

	"github.com/docker/docker/api/types/container"
	dockerImage "github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/google/uuid"
	"github.com/kitanoyoru/golang-orchestrator/task/service"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Docker struct {
	logger zerolog.Logger
	client *client.Client
}

func New() (service.CRI, error) {
	logger := log.With().Str("CRI", "Docker").Logger()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		logger.Error().Err(err)
		return nil, err
	}

	return &Docker{logger: logger, client: cli}, nil
}

func (d *Docker) Run(ctx context.Context, image string, options ...service.RunOption) (string, error) {
	d.logger.Info().Ctx(ctx).Str("method", "Run").Msg("call")

	opts := &service.RunOptions{}
	for _, o := range options {
		o(opts)
	}

	if err := d.pullImage(ctx, image); err != nil {
		return "", err
	}

	containerID, err := d.createContainer(ctx, image, opts)
	if err != nil {
		return "", err
	}

	if err = d.startContainer(ctx, containerID); err != nil {
		return "", err
	}

	if err := d.streamContainerLogs(ctx, containerID); err != nil {
		return "", err
	}

	return containerID, nil
}

func (d *Docker) pullImage(ctx context.Context, image string) error {
	d.logger.Info().Ctx(ctx).Str("method", "pullImage").Msg("call")

	reader, err := d.client.ImagePull(ctx, image, dockerImage.PullOptions{})
	if err != nil {
		return err
	}
	defer reader.Close()

	if _, err = io.Copy(os.Stdout, reader); err != nil {
		return err
	}
	return nil
}

func (d *Docker) createContainer(ctx context.Context, image string, opts *service.RunOptions) (string, error) {
	d.logger.Info().Ctx(ctx).Str("method", "createContainer").Msg("call")

	cc := container.Config{
		Image: image,
		Tty:   false,
	}

	if len(opts.Env) > 0 {
		cc.Env = opts.Env
	}

	if len(opts.ExposedPorts) > 0 {
		cc.ExposedPorts = opts.ExposedPorts
	}

	hc := container.HostConfig{
		RestartPolicy: container.RestartPolicy{
			Name: container.RestartPolicyDisabled,
		},
		PublishAllPorts: true,
	}

	if opts.RestartPolicy != nil {
		hc.RestartPolicy.Name = container.RestartPolicyMode(*opts.RestartPolicy)
	}

	if opts.Limit != nil {
		nanoCPUs, err := opts.Limit.GetNanoCPUs()
		if err != nil {
			return "", err
		}
		hc.Resources = container.Resources{
			Memory:   *opts.Limit.Memory,
			NanoCPUs: nanoCPUs,
		}
	}

  title := uuid.NewString()
	if opts.Title != nil {
    title = *opts.Title
	}

	resp, err := d.client.ContainerCreate(ctx, &cc, &hc, nil, nil, title)
	if err != nil {
		d.logger.Error().Ctx(ctx).Err(err)
		return "", err
	}

	return resp.ID, nil
}

func (d *Docker) startContainer(ctx context.Context, id string) error {
	d.logger.Info().Ctx(ctx).Str("method", "startContainer").Msg("call")

	if err := d.client.ContainerStart(ctx, id, container.StartOptions{}); err != nil {
		d.logger.Error().Ctx(ctx).Err(err)
		return err
	}

	return nil
}

func (d *Docker) streamContainerLogs(ctx context.Context, containerID string) error {
	out, err := d.client.ContainerLogs(ctx, containerID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	})
	if err != nil {
		d.logger.Error().Ctx(ctx).Err(err)
		return err
	}
	defer out.Close()

	if _, err := stdcopy.StdCopy(os.Stdout, os.Stderr, out); err != nil {
		d.logger.Error().Ctx(ctx).Err(err)
		return err
	}

	return nil
}

func (d *Docker) Stop(ctx context.Context, id string) error {
	err := d.client.ContainerStop(ctx, id, container.StopOptions{})
	if err != nil {
		log.Error().Ctx(ctx).Err(err)
		return err
	}

	err = d.client.ContainerRemove(ctx, id, container.RemoveOptions{})
	if err != nil {
		log.Error().Ctx(ctx).Err(err)
		return err
	}

	return nil
}
