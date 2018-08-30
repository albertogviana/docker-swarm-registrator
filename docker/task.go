package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
)

// SwarmTaskClient implements SwarmTasks
type SwarmTaskClient struct {
	DockerClient *client.Client
}

// SwarmTasks defines interfaces with the required methods
type SwarmTasks interface {
	GetTask(ctx context.Context, filter filters.Args) ([]swarm.Task, error)
}

// NewSwarmTaskClient return an instance of task
func NewSwarmTaskClient(client *client.Client) *SwarmTaskClient {
	return &SwarmTaskClient{
		client,
	}
}

// GetTask running in the cluster
func (t *SwarmTaskClient) GetTask(ctx context.Context, filter filters.Args) ([]swarm.Task, error) {
	tasks, err := t.DockerClient.TaskList(ctx, types.TaskListOptions{Filters: filter})

	if err != nil {
		return []swarm.Task{}, err
	}

	return tasks, nil
}
