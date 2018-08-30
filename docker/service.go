package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
)

// SwarmServiceClient implements SwarmServices
type SwarmServiceClient struct {
	DockerClient *client.Client
	FilterLabel  string
}

// SwarmServices defines interfaces with the required methods
type SwarmServices interface {
	GetServices(ctx context.Context) ([]swarm.Service, error)
	// GetTask(filter filters.Args) ([]swarm.Task, error)
	// GetDeploymentStatus(serviceName string, image string) (ServiceStatus, error)
	// GetServiceStatus(serviceName string) (ServiceStatus, error)
}

// NewSwarmServiceClient return an instance of service
func NewSwarmServiceClient(client *client.Client, filterLabel string) *SwarmServiceClient {
	return &SwarmServiceClient{
		client,
		filterLabel,
	}
}

// GetServices all services running in the cluster
// You will find the available filters on https://docs.docker.com/engine/api/v1.32/#operation/ServiceList
func (s *SwarmServiceClient) GetServices(ctx context.Context) ([]swarm.Service, error) {
	filter := filters.NewArgs()
	filter.Add("label", s.FilterLabel)

	serviceList, err := s.DockerClient.ServiceList(ctx, types.ServiceListOptions{Filters: filter})

	if err != nil {
		return nil, err
	}

	swarmServices := []swarm.Service{}
	for _, service := range serviceList {
		swarmServices = append(swarmServices, service)
	}

	return swarmServices, nil
}
