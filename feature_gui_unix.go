// +build dragonfly freebsd linux nacl netbsd openbsd solaris, !darwin

// Build for all unix platforms, except OSX (darwin)
package main

import (
    "github.com/fsouza/go-dockerclient"
    "errors"
    "os"
)

func enableGui() (map[docker.Port][]docker.PortBinding, []string, []string, error) {
    portBindings := map[docker.Port][]docker.PortBinding{}
    envVariables := []string{}
    binds := make([]string, 0)

    display := os.Getenv("DISPLAY")
    if display == "" {
        return portBindings, envVariables, binds, errors.New("$DISPLAY is undefined")
    } else {
        envVariables = append(envVariables, "DISPLAY=unix" + display)
        binds = append(binds, "/tmp/.X11-unix" + ":" + "/tmp/.X11-unix")
        return portBindings, envVariables, binds, nil
    }
}
