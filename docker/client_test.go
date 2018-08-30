package docker

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type DockerClientTestSuite struct {
	suite.Suite
}

func TestDockerClientTestSuite(t *testing.T) {
	suite.Run(t, new(DockerClientTestSuite))
}

func (c *DockerClientTestSuite) Test_NewDockerClient_ReturnClient() {
	client, err := NewDockerClient("", map[string]string{})

	c.Require().NoError(err)
	c.Require().NotNil(client)
}
