package docker

import (
	"context"
	"testing"

	"github.com/albertogviana/docker-swarm-registrator/tests"
	"github.com/stretchr/testify/suite"
)

type SwarmServiceTestSuite struct {
	suite.Suite
	SwarmService *SwarmService
}

func TestSwarmServiceTestSuite(t *testing.T) {
	suite.Run(t, new(SwarmServiceTestSuite))
}

func (s *SwarmServiceTestSuite) SetupSuite() {
	tests.CreateTestService("nginx-registrator", []string{"registrator.enabled=true"}, []string{"mode=host,target=80"}, "", "dnsrr", "nginx:alpine", []string{})
	tests.CreateTestService("nginx", []string{}, []string{}, "", "", "nginx:alpine", []string{})
}

func (s *SwarmServiceTestSuite) SetupTest() {
	client, err := NewDockerClient("", map[string]string{})

	s.Require().NoError(err)
	s.Require().NotNil(client)

	s.SwarmService = NewSwarmService(client, "registrator.enabled=true")
}

func (s *SwarmServiceTestSuite) TearDownSuite() {
	tests.RemoveTestService("nginx-registrator")
	tests.RemoveTestService("nginx")
}

func (s *SwarmServiceTestSuite) Test_GetServices() {
	services, err := s.SwarmService.GetServices(context.Background())

	s.Require().NoError(err)
	s.Len(*services, 1)
}
