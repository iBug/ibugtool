package docker

import (
	"context"
	"strconv"

	"github.com/docker/docker/client"
)

var cli *client.Client

func GetDockerClient() (*client.Client, error) {
	if cli != nil {
		return cli, nil
	}
	c, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return c, err
	}
	cli = c
	return cli, nil
}

func ContainerProcesses(ctx context.Context, containerID string) ([]uint64, error) {
	c, err := GetDockerClient()
	if err != nil {
		return nil, err
	}
	container, err := c.ContainerTop(ctx, containerID, nil)
	if err != nil {
		return nil, err
	}
	pids := make([]uint64, 0, len(container.Processes))
	for _, proc := range container.Processes {
		pid, err := strconv.ParseUint(proc[1], 10, 64)
		if err != nil {
			return nil, err
		}
		pids = append(pids, pid)
	}
	return pids, nil
}
