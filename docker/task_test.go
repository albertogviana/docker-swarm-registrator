package docker

import (
	"context"
	"testing"

	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/stretchr/testify/suite"
)

type SwarmTaskClientTestSuite struct {
	suite.Suite
	DockerClient       *client.Client
	SwarmServiceClient *SwarmServiceClient
	SwarmTaskClient    *SwarmTaskClient
}

func TestSwarmTaskClientTestSuite(t *testing.T) {
	suite.Run(t, new(SwarmTaskClientTestSuite))
}

func (s *SwarmTaskClientTestSuite) SetupSuite() {
	createTestService("nginx-registrator", []string{"registrator.enabled=true"}, []string{"mode=host,target=80"}, "", "dnsrr", "nginx:alpine")
	createTestService("nginx", []string{}, []string{}, "", "", "nginx:alpine")
	scaleTestService("nginx-registrator", 5)
}

func (s *SwarmTaskClientTestSuite) SetupTest() {
	client, err := NewDockerClient("", map[string]string{})

	s.Require().NoError(err)
	s.Require().NotNil(client)

	s.SwarmServiceClient = NewSwarmServiceClient(client, "registrator.enabled=true")
	s.SwarmTaskClient = NewSwarmTaskClient(client)
}

func (s *SwarmTaskClientTestSuite) TearDownSuite() {
	removeTestService("nginx-registrator")
	removeTestService("nginx")
}

func (s *SwarmTaskClientTestSuite) Test_GetTasks() {
	ctx := context.Background()
	services, err := s.SwarmServiceClient.GetServices(ctx)

	s.Require().NoError(err)
	s.Len(services, 1)

	filter := filters.NewArgs()
	filter.Add("service", services[0].ID)
	filter.Add("desired-state", "running")

	task, err := s.SwarmTaskClient.GetTask(ctx, filter)

	s.Require().NoError(err)
	s.Len(task, 5)
}
