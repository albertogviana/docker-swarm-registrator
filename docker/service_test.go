package docker

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type SwarmServiceClientTestSuite struct {
	suite.Suite
	SwarmServiceClient *SwarmServiceClient
}

func TestSwarmServiceClientTestSuite(t *testing.T) {
	suite.Run(t, new(SwarmServiceClientTestSuite))
}

func (s *SwarmServiceClientTestSuite) SetupSuite() {
	createTestService("nginx-registrator", []string{"registrator.enabled=true"}, []string{"mode=host,target=80"}, "", "dnsrr", "nginx:alpine")
	createTestService("nginx", []string{}, []string{}, "", "", "nginx:alpine")
}

func (s *SwarmServiceClientTestSuite) SetupTest() {
	client, err := NewDockerClient("", map[string]string{})

	s.Require().NoError(err)
	s.Require().NotNil(client)

	s.SwarmServiceClient = NewSwarmServiceClient(client, "registrator.enabled=true")
}

func (s *SwarmServiceClientTestSuite) TearDownSuite() {
	removeTestService("nginx-registrator")
	removeTestService("nginx")
}

func (s *SwarmServiceClientTestSuite) Test_GetServices() {
	services, err := s.SwarmServiceClient.GetServices(context.Background())

	s.Require().NoError(err)

	s.Len(services, 1)
}
