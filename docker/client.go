package docker

import (
	"github.com/docker/docker/client"
)

var dockerAPIVersion = "v1.37"

// NewDockerClient creates a docker client connection
func NewDockerClient(host string, defaultHeaders map[string]string) (*client.Client, error) {
	if host == "" {
		host = "unix:///var/run/docker.sock"
	}

	client, err := client.NewClient(host, dockerAPIVersion, nil, defaultHeaders)
	if err != nil {
		return nil, err
	}

	return client, nil
}
