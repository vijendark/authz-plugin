package authz

import (
    "strings"

    "github.com/Sirupsen/logrus"
    "github.com/docker/docker/pkg/authorization"
)

// AuthorizeDeleteAction checks a delete command against a list of swarm infrastructure to prevent
// accidental deletion of these containers.
func AuthorizeDeleteAction(uri string) *authorization.Response {
    target := extractDeleteTarget(uri)

    for _, name := range CarinaInfrastructureContainers {
        uuid, err := GetActiveContainerUUID(name)
        if err != nil {
            logrus.Error(err.Error())
            return actionDenied("Error: Internal docker authorization error")
        }

        // Allow if container is not running
        if uuid == "" {
            continue
        }

        if targetMatches(target, name, uuid) {
            return actionDenied("Error: Attempted to delete infrastructure container")
        }
    }

    return nil
}

// targetMatches checks if the target string matches the given name or UUID (including partial match).
func targetMatches(target, name, uuid string) bool {
    if target == name || strings.HasPrefix(uuid, target) {
        return true
    }

    return false
}

// extractDeleteTarget parses a container delete url and returns the delete target.
func extractDeleteTarget(uri string) string {
    //Example URL: /v1.23/containers/big_minsky?force=1
    parts := strings.Split(uri, "/")
    target := strings.Split(parts[len(parts)-1], "?")

    return target[0]
}

