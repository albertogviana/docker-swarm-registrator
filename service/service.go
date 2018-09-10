package service

import (
	"context"

	"github.com/albertogviana/docker-swarm-registrator/docker"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
)

// Service defines the based structure
type Service struct {
	SwarmService *docker.SwarmService
	SwarmTask    *docker.SwarmTask
	SwarmNode    *docker.SwarmNode
}

// Services defines interface with mandatory methods
type Services interface {
	GetServices(ctx context.Context) (*[]SwarmService, error)
}

// NewService returns a new instance of the Service structure
func NewService(swarmService *docker.SwarmService, swarmTask *docker.SwarmTask, swarmNode *docker.SwarmNode) *Service {
	return &Service{
		swarmService,
		swarmTask,
		swarmNode,
	}
}

// GetServices returns all services running in the cluster that has the label `registrator.enabled=true`
func (s *Service) GetServices(ctx context.Context) (*[]SwarmService, error) {
	svc, err := s.SwarmService.GetServices(ctx)
	if err != nil {
		return nil, err
	}

	services := []SwarmService{}
	for _, sm := range *svc {
		ss := SwarmService{}

		ss.ID = sm.ID
		ss.Name = sm.Spec.Name

		filter := filters.NewArgs()
		filter.Add("service", sm.Spec.Name)
		filter.Add("desired-state", "running")

		tasks, err := s.getTasksByService(ctx, sm, filter)
		if err != nil {
			return nil, err
		}
		ss.Task = tasks

		services = append(services, ss)
	}

	return &services, nil
}

func (s *Service) getTasksByService(ctx context.Context, ss swarm.Service, filters filters.Args) ([]SwarmTask, error) {
	tasks, err := s.SwarmTask.GetTask(ctx, filters)
	if err != nil {
		return nil, err
	}

	swarmTasks := []SwarmTask{}
	for _, task := range *tasks {
		t := SwarmTask{}
		t.ID = task.ID
		t.Name = ss.Spec.Name

		node, err := s.SwarmNode.GetNodeByID(ctx, task.NodeID)
		if err != nil {
			return nil, err
		}

		t.Address = node.Status.Addr

		for _, p := range task.Status.PortStatus.Ports {
			t.Port = int(p.PublishedPort)
		}

		swarmTasks = append(swarmTasks, t)
	}

	return swarmTasks, nil
}
