package authz

import (
    "github.com/Sirupsen/logrus"
    "github.com/docker/engine-api/types"
    "github.com/docker/engine-api/types/filters"
    "golang.org/x/net/context"
)

const SwarmManagerContainerName string = "swarm-manager"
const SwarmAgentContainerName string = "swarm-agent"
const SwarmDataContainerName string = "swarm-data"

var InfrastructureContainers = []string{
    SwarmManagerContainerName,
    SwarmAgentContainerName,
    SwarmDataContainerName,
}

// GetActiveContainerUUID attempts to retrieve the container UUID matching the given name via the Docker API
func GetActiveContainerUUID(name string) (uuid string, err error) {
    cli := NewApiClient()

    args := filters.NewArgs()
    args.Add("name", name)
    args.Add("status", "running")
    options := types.ContainerListOptions{All: true, Filter: args}

    containers, err := cli.ContainerList(context.Background(), options)
    if err != nil {
        return "", err
    }

    if len(containers) != 1 {
        logrus.Warnf("%s not found in running containers.", name)
        return "", nil
    }

    return containers[0].ID, nil
}

