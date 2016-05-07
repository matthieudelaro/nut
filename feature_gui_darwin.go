// Builds for OSX
package main

import (
    // "github.com/Sirupsen/logrus"
    "github.com/fsouza/go-dockerclient"
    // "bytes"
    "os"
    // "net"
    // "io"
    // log "github.com/Sirupsen/logrus"
    "errors"
    // "time"
    // "strings" // words := strings.Fields(someString)
)

func enableGui() (map[docker.Port][]docker.PortBinding, []string, []string, error) {
    guiPortNumber := "6000"
    portBindings := map[docker.Port][]docker.PortBinding{}
    envVariables := []string{}
    binds := make([]string, 0)

    displayVariable := os.Getenv("DISPLAY")
    if displayVariable == "" {
        return portBindings, envVariables, binds, errors.New(
            "DISPLAY variable is empty. (Did you install XQuartz properly?)")
    } else {
        ip, err := bridgeUnixSocketToPort(displayVariable, guiPortNumber)
        if err != nil {
            return portBindings, envVariables, binds,
                   errors.New("Network bridging failed: " + err.Error())
        } else {
            // log.Debug("bridged to ip:", string(ip))
            envVariables = append(envVariables, "DISPLAY=" + ip + ":0")
            return portBindings, envVariables, binds, nil
        }
    }
}
