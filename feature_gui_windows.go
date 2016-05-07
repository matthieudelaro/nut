// Build for Windows
package main

import (
    "github.com/fsouza/go-dockerclient"
    "errors"
)

func enableGui() (map[docker.Port][]docker.PortBinding, []string, []string, error) {
    portBindings := map[docker.Port][]docker.PortBinding{}
    envVariables := []string{}
    binds := make([]string, 0)

    return portBindings, envVariables, binds, errors.New("Could not enable GUI: it has not been implemented for linux yet.")
}
