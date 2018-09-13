package service

import (
	"context"
	"testing"

	"github.com/albertogviana/docker-swarm-registrator/docker"
	"github.com/albertogviana/docker-swarm-registrator/tests"
	"github.com/docker/docker/client"
	"github.com/stretchr/testify/suite"
)

type ServiceTestSuite struct {
	suite.Suite
	Service      *Service
	DockerClient *client.Client
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

func (s *ServiceTestSuite) SetupSuite() {
	tests.CreateTestService("service-1", []string{"registrator.enabled=true", "registrator.checks.1.name=service-1-health", "registrator.checks.1.id=service-1-health", "registrator.checks.1.interval=10s", "registrator.checks.1.timeout=10s", "registrator.checks.1.path=/", "registrator.checks.1.http=true", "registrator.checks.1.removefailedserviceafter=30s"}, []string{"mode=host,target=80"}, "", "dnsrr", "nginx:alpine", []string{})
	tests.CreateTestService("service-2", []string{"registrator.enabled=true"}, []string{"mode=host,target=80"}, "", "dnsrr", "nginx:alpine", []string{})
	tests.ScaleTestService("service-1", 3)
}

func (s *ServiceTestSuite) SetupTest() {
	client, err := docker.NewDockerClient("", map[string]string{"Cache-Control": "no-cache, no-store, must-revalidate"})

	s.Require().NoError(err)
	s.Require().NotNil(client)

	swarmService := docker.NewSwarmService(client, "registrator.enabled=true")
	swarmTask := docker.NewSwarmTask(client)
	swarmNode := docker.NewSwarmNode(client)
	s.Service = NewService(swarmService, swarmTask, swarmNode)
	s.DockerClient = client
}

func (s *ServiceTestSuite) TearDownSuite() {
	tests.RemoveTestService("service-1")
	tests.RemoveTestService("service-2")
	s.DockerClient.Close()
}

func (s *ServiceTestSuite) Test_GetServices() {
	ctx := context.Background()

	svc, err := s.Service.GetServices(ctx)

	s.Require().NoError(err)
	s.Len(*svc, 2)
}
