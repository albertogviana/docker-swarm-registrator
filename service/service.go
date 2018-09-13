package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/albertogviana/docker-swarm-registrator/docker"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/mitchellh/mapstructure"
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

	checks, err := s.getServiceLabels(ss)
	if err != nil {
		return nil, err
	}

	swarmTasks := []SwarmTask{}
	for _, task := range *tasks {
		t := SwarmTask{}
		t.ID = task.ID
		t.Name = ss.Spec.Name

		if len(checks) > 0 {
			t.Checks = checks
		}

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

func (s *Service) getServiceLabels(ss swarm.Service) ([]Check, error) {
	checks := []Check{}

	params := map[string]map[string]string{}
	for k, v := range ss.Spec.Labels {
		if strings.HasPrefix(k, "registrator.checks.") && k != "registrator.enabled" {
			index := strings.TrimPrefix(k, "registrator.checks.")[:1]
			if params[index] == nil {
				params[index] = map[string]string{}
			}

			params[index][strings.TrimPrefix(k, fmt.Sprintf("registrator.checks.%s.", index))] = v
		}
	}

	if len(params) > 0 {
		for k := range params {
			if len(params[k]) > 0 {
				var check Check
				mapstructure.Decode(params[k], &check)
				checks = append(checks, check)
			}
		}
	}

	return checks, nil
}
