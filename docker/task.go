package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
)

// SwarmTask implements SwarmTasks
type SwarmTask struct {
	DockerClient *client.Client
}

// SwarmTasks defines interfaces with the required methods
type SwarmTasks interface {
	GetTask(ctx context.Context, filter filters.Args) ([]swarm.Task, error)
}

// NewSwarmTask returns an instance of SwarmTask
func NewSwarmTask(client *client.Client) *SwarmTask {
	return &SwarmTask{
		client,
	}
}

// GetTask all tasks related to a service running in the cluster
// You will find the available filters on https://docs.docker.com/engine/api/v1.37/#operation/TaskList
func (t *SwarmTask) GetTask(ctx context.Context, filter filters.Args) ([]swarm.Task, error) {
	tasks, err := t.DockerClient.TaskList(ctx, types.TaskListOptions{Filters: filter})

	if err != nil {
		return []swarm.Task{}, err
	}

	return tasks, nil
}
