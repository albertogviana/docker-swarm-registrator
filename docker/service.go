package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
)

// SwarmService implements SwarmServices
type SwarmService struct {
	DockerClient *client.Client
	FilterLabel  string
}

// SwarmServices defines interfaces with the required methods
type SwarmServices interface {
	GetServices(ctx context.Context) (*[]swarm.Service, error)
}

// NewSwarmService returns an instance of SwarmService
func NewSwarmService(client *client.Client, filterLabel string) *SwarmService {
	return &SwarmService{
		client,
		filterLabel,
	}
}

// GetServices all services running in the cluster
// You will find the available filters on https://docs.docker.com/engine/api/v1.32/#operation/ServiceList
func (s *SwarmService) GetServices(ctx context.Context) (*[]swarm.Service, error) {
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

	return &swarmServices, nil
}
