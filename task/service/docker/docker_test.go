package docker

import (
	"context"
	"testing"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type DockerTestSuite struct {
	suite.Suite

	docker Docker
	config Config
}

func (suite *DockerTestSuite) SetupSuite() {
	suite.config = Config{
		Title:         "test-container",
		Image:         "alpine", // Use a lightweight image for testing
		Cmd:           []string{"sh", "-c", "echo 'Hello, World!'"},
		ExposedPorts:  nat.PortSet{"8080/tcp": struct{}{}},
		Memory:        50 * 1024 * 1024, // 50MB
		Cpu:           0.1,              // 10% of a CPU
		Env:           []string{"ENV_VAR=test"},
		RestartPolicy: container.RestartPolicyMode("no"),
	}

	var err error
	suite.docker, err = NewDocker(suite.config)
	if err != nil {
		suite.T().Fatal(err)
	}
}

func (suite *DockerTestSuite) TestRunContainer() {
	ctx := context.Background()

	// Run the container
	result := suite.docker.Run(ctx)
	assert.NoError(suite.T(), result.Error)
	assert.NotEmpty(suite.T(), result.ContainerID)

	// Stop the container
	stopResult := suite.docker.Stop(ctx, result.ContainerID)
	assert.NoError(suite.T(), stopResult.Error)
	assert.Equal(suite.T(), "stop", stopResult.Action)
	assert.Equal(suite.T(), "success", stopResult.Result)
}

func TestDockerTestSuite(t *testing.T) {
	suite.Run(t, new(DockerTestSuite))
}
