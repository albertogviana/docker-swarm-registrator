package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
)

// SwarmNode implements SwarmNodes interface
type SwarmNode struct {
	DockerClient *client.Client
}

// SwarmNodes defines interfaces with the required methods
type SwarmNodes interface {
	GetNodes(ctx context.Context, filter filters.Args) ([]swarm.Node, error)
}

// NewSwarmNode returns an instance of SwarmNode
func NewSwarmNode(client *client.Client) *SwarmNode {
	return &SwarmNode{
		client,
	}
}

// GetNodes all nodes available in the cluster
// You will find more about the filter on https://docs.docker.com/engine/api/v1.37/#operation/NodeList
func (s *SwarmNode) GetNodes(ctx context.Context, filter filters.Args) ([]swarm.Node, error) {
	nodes, err := s.DockerClient.NodeList(ctx, types.NodeListOptions{Filters: filter})

	if err != nil {
		return []swarm.Node{}, err
	}

	return nodes, nil
}

// GetNodeByID get a node by id
func (s *SwarmNode) GetNodeByID(ctx context.Context, id string) (swarm.Node, error) {
	filter := filters.NewArgs()
	filter.Add("id", id)
	node, err := s.GetNodes(ctx, filter)
	if err != nil {
		return swarm.Node{}, err
	}

	return node[0], nil
}
