package docker

import (
	"context"
	"testing"

	"github.com/docker/go-connections/nat"
	"github.com/kitanoyoru/golang-orchestrator/pkg/types"
	"github.com/kitanoyoru/golang-orchestrator/task/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type DockerTestSuite struct {
	suite.Suite

	docker service.CRI
}

func (suite *DockerTestSuite) SetupSuite() {
	var err error
	suite.docker, err = New()
	if err != nil {
		suite.T().Fatal(err)
	}
}

func (suite *DockerTestSuite) TestRunContainer() {
	ctx := context.Background()

	opts := []service.RunOption{
		service.WithTitle("test-container"),
		service.WithCmd([]string{"sh", "-c", "echo 'Hello, World!'"}),
		service.WithExposedPorts(nat.PortSet{"8080/tcp": struct{}{}}),
		service.WithResourceLimitation(&service.ResourceLimitation{
			Memory: types.Int64(50 * 1024 * 1024), // 50MB
			CPU:    types.Float64(0.1),            // 10% of a CPU
		}),
		service.WithEnv([]string{"ENV_VAR=test"}),
	}

	id, err := suite.docker.Run(ctx, "alpine", opts...)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), id)

	err = suite.docker.Stop(ctx, id)
	assert.NoError(suite.T(), err)
}

func TestDockerTestSuite(t *testing.T) {
	suite.Run(t, new(DockerTestSuite))
}
