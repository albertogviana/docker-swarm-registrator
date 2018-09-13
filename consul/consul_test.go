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
	tests.CreateTestService("consul", []string{}, []string{"8300:8300", "8500:8500"}, "", "", "consul", []string{"CONSUL_BIND_INTERFACE=eth0"})
	tests.CreateTestService("nginx-registrator", []string{"registrator.enabled=true"}, []string{"mode=host,target=80"}, "", "dnsrr", "nginx:alpine", []string{})
	tests.CreateTestService("service-1", []string{"registrator.enabled=true", "registrator.checks.1.name=service-health", "registrator.checks.1.id=service-health", "registrator.checks.1.interval=10s", "registrator.checks.1.timeout=10s", "registrator.checks.1.path=/", "registrator.checks.1.http=true", "registrator.checks.1.removefailedserviceafter=30s"}, []string{"mode=host,target=80"}, "", "dnsrr", "nginx:alpine", []string{})
	tests.ScaleTestService("nginx-registrator", 3)
	tests.ScaleTestService("service-1", 5)
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
	tests.RemoveTestService("consul")
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

	service1Catalog, _, err := c.ConsulClient.Catalog().Service("service-1", "", &consulapi.QueryOptions{})
	c.Require().NoError(err)
	c.Len(service1Catalog, 5)

	nginxCatalog, _, err := c.ConsulClient.Catalog().Service("nginx-registrator", "", &consulapi.QueryOptions{})
	c.Require().NoError(err)
	c.Len(nginxCatalog, 3)

	for _, svc := range *services {
		for _, task := range svc.Task {
			err = consul.Deregister(task.ID)
			c.Require().NoError(err)
		}
	}

	service1Catalog, _, err = c.ConsulClient.Catalog().Service("service-1", "", &consulapi.QueryOptions{})
	c.Require().NoError(err)
	c.Len(service1Catalog, 0)

	nginxCatalog, _, err = c.ConsulClient.Catalog().Service("nginx-registrator", "", &consulapi.QueryOptions{})
	c.Require().NoError(err)
	c.Len(nginxCatalog, 0)
}
