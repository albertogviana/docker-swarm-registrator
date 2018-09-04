package docker

import (
	"context"
	"testing"

	"github.com/albertogviana/docker-swarm-registrator/tests"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/stretchr/testify/suite"
)

type SwarmTaskTestSuite struct {
	suite.Suite
	DockerClient *client.Client
	SwarmService *SwarmService
	SwarmTask    *SwarmTask
}

func TestSwarmTaskTestSuite(t *testing.T) {
	suite.Run(t, new(SwarmTaskTestSuite))
}

func (s *SwarmTaskTestSuite) SetupSuite() {
	tests.CreateTestService("nginx-registrator", []string{"registrator.enabled=true"}, []string{"mode=host,target=80"}, "", "dnsrr", "nginx:alpine")
	tests.CreateTestService("nginx", []string{}, []string{}, "", "", "nginx:alpine")
	tests.ScaleTestService("nginx-registrator", 5)
}

func (s *SwarmTaskTestSuite) SetupTest() {
	client, err := NewDockerClient("", map[string]string{})

	s.Require().NoError(err)
	s.Require().NotNil(client)

	s.SwarmService = NewSwarmService(client, "registrator.enabled=true")
	s.SwarmTask = NewSwarmTask(client)
}

func (s *SwarmTaskTestSuite) TearDownSuite() {
	tests.RemoveTestService("nginx-registrator")
	tests.RemoveTestService("nginx")
}

func (s *SwarmTaskTestSuite) Test_GetTasks() {
	ctx := context.Background()
	services, err := s.SwarmService.GetServices(ctx)

	s.Require().NoError(err)
	s.Len(services, 1)

	filter := filters.NewArgs()
	filter.Add("service", services[0].ID)
	filter.Add("desired-state", "running")

	task, err := s.SwarmTask.GetTask(ctx, filter)

	s.Require().NoError(err)
	s.Len(task, 5)
}
