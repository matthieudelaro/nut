package main

import (
    "errors"
    log "github.com/Sirupsen/logrus"
    "gopkg.in/yaml.v2"
    "path/filepath"
)

type MacroV5 struct {
    MacroBase `yaml:"inheritedValues,inline"`

    // A short description of the usage of this macro
    Usage string `yaml:"usage,omitempty"`
    // The commands to execute when this macro is invoked
    Actions []string `yaml:"actions,omitempty"`
    // A list of aliases for the macro
    Aliases []string `yaml:"aliases,omitempty"`
    // Custom text to show on USAGE section of help
    UsageText string `yaml:"usage_for_help_section,omitempty"`
    // A longer explanation of how the macro works
    Description string `yaml:"description,omitempty"`
    DockerImage string `yaml:"docker_image,omitempty"`
    WorkingDir string `yaml:"container_working_directory,omitempty"`
    Mount map[string][]string `yaml:"mount,omitempty"`
    EnvironmentVariables map[string]string `yaml:"environment,omitempty"`
    Ports []string `yaml:"ports,omitempty"`
    EnableGUI string `yaml:"enable_gui,omitempty"`
    EnableNvidiaDevices string `yaml:"enable_nvidia_devices,omitempty"`
    Privileged string `yaml:"privileged,omitempty"`
    SecurityOpts []string `yaml:"security_opts,omitempty"`
}
        func (self *MacroV5) getUsage() string {
            return self.Usage
        }
        func (self *MacroV5) getActions() []string {
            return self.Actions
        }
        func (self *MacroV5) getAliases() []string {
            return self.Aliases
        }
        func (self *MacroV5) getUsageText() string {
            return self.UsageText
        }
        func (self *MacroV5) getDescription() string {
            return self.Description
        }
        func (self *MacroV5) getDockerImageName() (string, error) {
            if self.DockerImage != "" { 
                return self.DockerImage, nil 
            }
            return self.project.getBaseEnv().getDockerImageName()
        }
        func (self *MacroV5) getWorkingDir() string {
            if self.WorkingDir != "" { 
                return self.WorkingDir 
            }
            return self.project.getWorkingDir()
        }
        func (self *MacroV5) getMountingPoints() map[string]MountingPoint {
            if len(self.Mount) > 0 {
                points := make(map[string]MountingPoint)
                for name, data := range(self.Mount) {
                    points[name] = &MountingPointV5{
                        host: data[0],
                        container: data[1],
                    }
                }
                return points
            }
            return self.project.getMountingPoints()
        }
        func (self *MacroV5) getEnvironmentVariables() map[string]string {
            if len(self.EnvironmentVariables) > 0 {
                return self.EnvironmentVariables
            }
            return self.project.getEnvironmentVariables()
        }
        func (self *MacroV5) getPorts() []string {
            if len(self.Ports) > 0 {
                return self.Ports
            }
            return self.project.getPorts()
        }
        func (self *MacroV5) getEnableGui() bool {
            if self.EnableGUI != "" {
                return self.EnableGUI == "true"
            }
            return self.project.getEnableGui()
        }
        func (self *MacroV5) getEnableNvidiaDevices() bool {
            if self.EnableNvidiaDevices != "" {
                return self.EnableNvidiaDevices == "true"
            }
            return self.project.getEnableNvidiaDevices()
        }
        func (self *MacroV5) getPrivileged() bool {
            if self.Privileged != "" {
                return self.Privileged == "true"
            }
            return self.project.getPrivileged()
        }
        func (self *MacroV5) getSecurityOpts() []string {
            if len(self.SecurityOpts) > 0 {
                return self.SecurityOpts
            }
            return self.project.getSecurityOpts()
        }

type MountingPointV5 struct {
    MountingPointBase `yaml:"inheritedValues,inline"`

    host string `yaml:host_path`
    container string `yaml:container_path`
}
        func (self *MountingPointV5) fullHostPath() (string, error) {
            absolutePath, err := filepath.Abs(self.host)
            return absolutePath, err
        }
        func (self *MountingPointV5) fullContainerPath() (string, error) {
            absolutePath, err := filepath.Abs(self.container)
            return absolutePath, err
        }

type BaseEnvironmentV5 struct {
    BaseEnvironmentBase `yaml:"inheritedValues,inline"`

    DockerImage string `yaml:"docker_image,omitempty"`
    // FileURL string `yaml:"nut_file_url,omitempty"`
    FilePath string `yaml:"nut_file_path,omitempty"`
    GitHub string `yaml:"github,omitempty"`
    parentBase BaseEnvironment
}
        func (self *BaseEnvironmentV5) getDockerImageName() (string, error) {
            if self.DockerImage != "" {
                log.Debug("self.DockerImage is not ''")
                return self.DockerImage, nil
            } else if self.parentBase != nil {
                log.Debug("self.parentBase != nil")
                if name, err := self.parentBase.getDockerImageName(); err == nil {
                    log.Debug("err == nil, name=", name)
                    return name, nil
                }
                log.Debug("err != nil")
            }
            log.Debug("returning error")
            return "", errors.New("Docker image has not been specified.")
        }
        func (self *BaseEnvironmentV5) getFilePath() string{
            return self.FilePath
        }
        func (self *BaseEnvironmentV5) getGitHub() string{
            return self.GitHub
        }
        func (self *BaseEnvironmentV5) setGitHub(repositoryName string) {
            self.GitHub = repositoryName
        }
        func (self *BaseEnvironmentV5) setFilePath(filePath string) {
            self.FilePath = filePath
        }
        func (self *BaseEnvironmentV5) setParentBase(parentBase BaseEnvironment) error {
            self.parentBase = parentBase
            return nil
        }

// Introduce environment variables, ports, etc
type ProjectV5 struct {
    ProjectBase `yaml:"inheritedValues,inline"`

    SyntaxVersion string `yaml:"syntax_version"`
    ProjectName string `yaml:"project_name"`
    Base BaseEnvironmentV5 `yaml:"based_on"`
    WorkingDir string `yaml:"container_working_directory,omitempty"`
    Mount map[string][]string `yaml:"mount,omitempty"`
    EnvironmentVariables map[string]string `yaml:"environment,omitempty"`
    Ports []string `yaml:"ports,omitempty"`
    Macros map[string]*MacroV5 `yaml:"macros,omitempty"`
    EnableGUI string `yaml:"enable_gui,omitempty"`
    EnableNvidiaDevices string `yaml:"enable_nvidia_devices,omitempty"`
    Privileged string `yaml:"privileged,omitempty"`
    SecurityOpts []string `yaml:"security_opts,omitempty"`
    parentProject Project
    cacheMountingPoints map[string]MountingPoint
    cacheMacros map[string]Macro
    cacheEnvironmentVariables map[string]string
    cachePorts []string
}
        func (self *ProjectV5) getSyntaxVersion() string {
            return self.SyntaxVersion
        }
        func (self *ProjectV5) getName() string {
            if self.ProjectName == "" && self.parentProject != nil {
                return self.parentProject.getName()
            } else {
                return self.ProjectName
            }
        }
        func (self *ProjectV5) getBaseEnv() BaseEnvironment {
            return &self.Base
        }
        func (self *ProjectV5) getWorkingDir() string {
            if self.WorkingDir == "" && self.parentProject != nil {
                return self.parentProject.getWorkingDir()
            } else {
                return self.WorkingDir
            }
        }
        // Build cacheMountingPoints if it is nil.
        // Returns cache miss.
        func (self *ProjectV5) __cacheMountingPoints() bool {
            if self.cacheMountingPoints == nil {
                self.cacheMountingPoints = make(map[string]MountingPoint)
                for name, data := range(self.Mount) {
                    self.cacheMountingPoints[name] = &MountingPointV5{
                        host: data[0],
                        container: data[1],
                    }
                }
                return true
            } else {
                return false
            }
        }

        func (self *ProjectV5) getMountingPoints() map[string]MountingPoint {
            if self.__cacheMountingPoints() {
                // add the mounting points of the parent, if there is no conflict
                if self.parentProject != nil {
                    log.Debug("there is a parent")
                    for key, mountingPoint := range self.parentProject.getMountingPoints() {
                        log.Debug("parent point ", key)
                        // verify that there is no conflict with current mounting points
                        if CheckConflict(key, mountingPoint, self.cacheMountingPoints) == nil {
                            // add mounting point
                            self.cacheMountingPoints[key] = mountingPoint
                            log.Debug("add mounting point ", key, " ", mountingPoint)
                        }
                    }
                    log.Debug("Done iterating over parent's mouting points")
                }
            } else {
                log.Debug("cacheMountingPoints is ", self.cacheMountingPoints)
            }
            return self.cacheMountingPoints
        }
        func (self *ProjectV5) getMacros() map[string]Macro {
            if self.cacheMacros == nil {
                log.Debug("cacheMacros is nil")
                // make the list of macros
                self.cacheMacros = make(map[string]Macro)
                for name, data := range self.Macros {
                    self.cacheMacros[name] = data
                    self.cacheMacros[name].setParentProject(self)
                }
                // add the macros of the parent, if there is no conflict
                if self.parentProject != nil {
                    for name, macro := range self.parentProject.getMacros() {
                        log.Debug("parent macro ", name)
                        if self.cacheMacros[name] == nil {
                            log.Debug("add it")
                            self.cacheMacros[name] = macro
                            self.cacheMacros[name].setParentProject(self)
                        } else {
                            log.Debug("already exist")
                        }
                    }
                }

            } else {
                log.Debug("cacheMountingPoints is ", self.cacheMountingPoints)
            }
            return self.cacheMacros
        }
        func (self *ProjectV5) toYAML() string {
            d, err := yaml.Marshal(&self)
            if err != nil {
                log.Fatalf("error: %v", err)
            }
            return string(d)
        }
        func (self *ProjectV5) fromYAML(bytes []byte) error {
            err := yaml.Unmarshal(bytes, self)
            if err == nil {
                if self.getSyntaxVersion() == "5" {
                    return nil
                } else {
                    return errors.New("Unexpected version number.")
                }
            } else {
                return err
            }
        }
        func (self *ProjectV5) getEnableGui() bool {
            if self.EnableGUI == "" && self.parentProject != nil {
                return self.parentProject.getEnableGui()
            } else {
                return self.EnableGUI == "true"
            }
            return false
        }
        func (self *ProjectV5) getEnableNvidiaDevices() bool {
            if self.EnableNvidiaDevices == "" && self.parentProject != nil {
                return self.parentProject.getEnableNvidiaDevices()
            } else {
                return self.EnableNvidiaDevices == "true"
            }
            return false
        }
        func (self *ProjectV5) getParentProject() Project {
            return self.parentProject
        }
        func (self *ProjectV5) setParentProject(project Project) error {
            self.parentProject = project
            self.cacheMacros = nil
            self.cacheMountingPoints = nil
            self.cacheEnvironmentVariables = nil
            self.Base.setParentBase(project.getBaseEnv())
            return nil
        }
        func (self *ProjectV5) getPrivileged() bool {
            if self.Privileged == "" && self.parentProject != nil {
                return self.parentProject.getPrivileged()
            } else {
                return self.Privileged == "true"
            }
            return false
        }
        func (self *ProjectV5) getEnvironmentVariables() map[string]string {
            if self.cacheEnvironmentVariables == nil {
                log.Debug("cacheEnvironmentVariables is nil")
                // make the list of items
                self.cacheEnvironmentVariables = make(map[string]string)
                for name, data := range self.EnvironmentVariables {
                    self.cacheEnvironmentVariables[name] = data
                }
                // add the items of the parent, if there is no conflict
                if self.parentProject != nil {
                    for name, data := range self.parentProject.getEnvironmentVariables() {
                        log.Debug("parent env variable ", name)
                        if _, ok := self.cacheEnvironmentVariables[name]; !ok {
                            log.Debug("add it")
                            self.cacheEnvironmentVariables[name] = data
                        } else {
                            log.Debug("already exist")
                        }
                    }
                }

            } else {
                log.Debug("cacheEnvironmentVariables is ", self.cacheEnvironmentVariables)
            }
            return self.cacheEnvironmentVariables
        }
        func (self *ProjectV5) getPorts() []string {
            if self.cachePorts == nil {
                log.Debug("cachePorts is nil")
                // make the list of items
                self.cachePorts = make([]string, len(self.Ports))
                for key, data := range self.Ports {
                    self.cachePorts[key] = data
                }
                // add the items of the parent, if there is no conflict
                if self.parentProject != nil {
                    for _, data := range self.parentProject.getPorts() {
                        log.Debug("parent port ", data)
                        // if _, ok := self.cachePorts[name]; !ok {
                            // TODO: check conflict
                            log.Debug("add it")
                            self.cachePorts = append(self.cachePorts, data)
                        // } else {
                            // log.Debug("already exist")
                        // }
                    }
                }

            } else {
                log.Debug("cachePorts is ", self.cachePorts)
            }
            return self.cachePorts
        }

func NewProjectV5() *ProjectV5 {
    project := &ProjectV5 {
        SyntaxVersion: "5",
        Base: BaseEnvironmentV5{},
        Mount: make(map[string][]string),
        Macros: make(map[string]*MacroV5),
        parentProject: nil,
    }
    return project
}
