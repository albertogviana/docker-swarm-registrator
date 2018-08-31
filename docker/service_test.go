package docker

import (
	"context"
	"testing"

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
	createTestService("nginx-registrator", []string{"registrator.enabled=true"}, []string{"mode=host,target=80"}, "", "dnsrr", "nginx:alpine")
	createTestService("nginx", []string{}, []string{}, "", "", "nginx:alpine")
}

func (s *SwarmServiceTestSuite) SetupTest() {
	client, err := NewDockerClient("", map[string]string{})

	s.Require().NoError(err)
	s.Require().NotNil(client)

	s.SwarmService = NewSwarmService(client, "registrator.enabled=true")
}

func (s *SwarmServiceTestSuite) TearDownSuite() {
	removeTestService("nginx-registrator")
	removeTestService("nginx")
}

func (s *SwarmServiceTestSuite) Test_GetServices() {
	services, err := s.SwarmService.GetServices(context.Background())

	s.Require().NoError(err)

	s.Len(services, 1)
}
