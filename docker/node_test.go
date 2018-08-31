package docker

import (
	"context"
	"testing"

	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/stretchr/testify/suite"
)

type SwarmNodeTestSuite struct {
	suite.Suite
	DockerClient *client.Client
	SwarmNode    *SwarmNode
}

func TestSwarmNodeTestSuite(t *testing.T) {
	suite.Run(t, new(SwarmNodeTestSuite))
}

func (s *SwarmNodeTestSuite) SetupTest() {
	client, err := NewDockerClient("", map[string]string{})

	s.Require().NoError(err)
	s.Require().NotNil(client)

	s.SwarmNode = NewSwarmNode(client)
}

func (s *SwarmNodeTestSuite) Test_GetNodes() {
	ctx := context.Background()
	nodes, err := s.SwarmNode.GetNodes(ctx, filters.NewArgs())

	s.Require().NoError(err)
	s.Require().NotNil(nodes)

	s.Len(nodes, 1)
}
