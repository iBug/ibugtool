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

func Client() *client.Client {
	c, err := GetDockerClient()
	if err != nil {
		panic(err)
	}
	return c
}

func ContainerProcesses(ctx context.Context, containerID string) ([]int, error) {
	container, err := Client().ContainerTop(ctx, containerID, nil)
	if err != nil {
		return nil, err
	}
	pids := make([]int, 0, len(container.Processes))
	for _, proc := range container.Processes {
		pid, err := strconv.ParseInt(proc[1], 10, 64)
		if err != nil {
			return nil, err
		}
		pids = append(pids, int(pid))
	}
	return pids, nil
}
