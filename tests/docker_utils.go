package tests

import (
	"fmt"
	"os/exec"
)

// CreateTestService deploys a swarm service
func CreateTestService(name string, labels []string, publish []string, mode string, endpointMode string, image string, environment []string) {
	args := []string{"service", "create", "--name", name}
	for _, v := range labels {
		args = append(args, "-l", v)
	}
	for _, v := range publish {
		args = append(args, "--publish", v)
	}
	for _, v := range environment {
		args = append(args, "-e", v)
	}
	if len(mode) > 0 {
		args = append(args, "--mode", "global")
	}
	if endpointMode != "" {
		args = append(args, "--endpoint-mode", endpointMode)
	}
	args = append(args, image)
	exec.Command("docker", args...).Output()
}

// ScaleTestService scales a swarm service
func ScaleTestService(name string, replicas int) {
	exec.Command("docker", "service", "scale", fmt.Sprintf("%s=%d", name, replicas)).Output()
}

// RemoveTestService remove a swarm service
func RemoveTestService(name string) {
	exec.Command("docker", "service", "rm", name).Output()
}
