package consul

import (
	"context"
	"testing"

	"github.com/albertogviana/docker-swarm-registrator/docker"
	"github.com/albertogviana/docker-swarm-registrator/service"
	"github.com/albertogviana/docker-swarm-registrator/tests"
	"github.com/docker/docker/client"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/suite"
)

type ConsulTestSuite struct {
	suite.Suite
	DockerClient *client.Client
	SwarmService *docker.SwarmService
	SwarmTask    *docker.SwarmTask
	SwarmNode    *docker.SwarmNode
	Service      *service.Service
	ConsulClient *consulapi.Client
}

func TestConsulTestSuite(t *testing.T) {
	suite.Run(t, new(ConsulTestSuite))
}

func (c *ConsulTestSuite) SetupSuite() {
	tests.CreateTestService("nginx-registrator", []string{"registrator.enabled=true"}, []string{"mode=host,target=80"}, "", "dnsrr", "nginx:alpine")
	tests.CreateTestService("service-1", []string{"registrator.enabled=true"}, []string{"mode=host,target=80"}, "", "dnsrr", "nginx:alpine")
	tests.ScaleTestService("nginx-registrator", 3)
}

func (c *ConsulTestSuite) SetupTest() {
	client, err := docker.NewDockerClient("", map[string]string{})

	c.Require().NoError(err)
	c.Require().NotNil(client)

	c.SwarmService = docker.NewSwarmService(client, "registrator.enabled=true")
	c.SwarmTask = docker.NewSwarmTask(client)
	c.SwarmNode = docker.NewSwarmNode(client)
	c.Service = service.NewService(c.SwarmService, c.SwarmTask, c.SwarmNode)
	c.DockerClient = client

	consulClient, err := consulapi.NewClient(consulapi.DefaultConfig())
	c.Require().NoError(err)
	c.ConsulClient = consulClient
}

func (c *ConsulTestSuite) TearDownSuite() {
	tests.RemoveTestService("nginx-registrator")
	tests.RemoveTestService("service-1")
}

func (c *ConsulTestSuite) Test_Consul() {
	ctx := context.Background()

	services, err := c.Service.GetServices(ctx)
	c.Require().NoError(err)
	c.Len(*services, 2)

	consul := NewConsulClient(c.ConsulClient)
	for _, svc := range *services {
		for _, task := range svc.Task {
			err = consul.Register(task)
			c.Require().NoError(err)
		}
	}

	for _, svc := range *services {
		for _, task := range svc.Task {
			err = consul.Deregister(task.ID)
			c.Require().NoError(err)
		}
	}
}
