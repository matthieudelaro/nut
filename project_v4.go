package main

import (
    "errors"
    log "github.com/Sirupsen/logrus"
    "gopkg.in/yaml.v2"
    "path/filepath"
)

type MacroV4 struct {
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
}
        func (self *MacroV4) getUsage() string {
            return self.Usage
        }
        func (self *MacroV4) getActions() []string {
            return self.Actions
        }
        func (self *MacroV4) getAliases() []string {
            return self.Aliases
        }
        func (self *MacroV4) getUsageText() string {
            return self.UsageText
        }
        func (self *MacroV4) getDescription() string {
            return self.Description
        }

type MountingPointV4 struct {
    MountingPointBase `yaml:"inheritedValues,inline"`

    host string `yaml:host_path`
    container string `yaml:container_path`
}
        func (self *MountingPointV4) fullHostPath() (string, error) {
            absolutePath, err := filepath.Abs(self.host)
            return absolutePath, err
        }
        func (self *MountingPointV4) fullContainerPath() (string, error) {
            absolutePath, err := filepath.Abs(self.container)
            return absolutePath, err
        }

type BaseEnvironmentV4 struct {
    BaseEnvironmentBase `yaml:"inheritedValues,inline"`

    DockerImage string `yaml:"docker_image,omitempty"`
    // FileURL string `yaml:"nut_file_url,omitempty"`
    FilePath string `yaml:"nut_file_path,omitempty"`
    GitHub string `yaml:"github,omitempty"`
    parentBase BaseEnvironment
}
        func (self *BaseEnvironmentV4) getDockerImageName() (string, error) {
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
        func (self *BaseEnvironmentV4) getFilePath() string{
            return self.FilePath
        }
        func (self *BaseEnvironmentV4) getGitHub() string{
            return self.GitHub
        }
        func (self *BaseEnvironmentV4) setGitHub(repositoryName string) {
            self.GitHub = repositoryName
        }
        func (self *BaseEnvironmentV4) setFilePath(filePath string) {
            self.FilePath = filePath
        }
        func (self *BaseEnvironmentV4) setParentBase(parentBase BaseEnvironment) error {
            self.parentBase = parentBase
            return nil
        }


// Introduce the notion of inheritance between nut files
type ProjectV4 struct {
    ProjectBase `yaml:"inheritedValues,inline"`

    SyntaxVersion string `yaml:"syntax_version"`
    ProjectName string `yaml:"project_name"`
    Base BaseEnvironmentV4 `yaml:"based_on"`
    WorkingDir string `yaml:"container_working_directory,omitempty"`
    Mount map[string][]string `yaml:"mount,omitempty"`
    Macros map[string]*MacroV4 `yaml:"macros,omitempty"`
    EnableGUI string `yaml:"enable_gui,omitempty"`
    Privileged string `yaml:"privileged,omitempty"`
    parentProject Project
    cacheMountingPoints map[string]MountingPoint
    cacheMacros map[string]Macro
}
        func (self *ProjectV4) getSyntaxVersion() string {
            return self.SyntaxVersion
        }
        func (self *ProjectV4) getName() string {
            if self.ProjectName == "" && self.parentProject != nil {
                return self.parentProject.getName()
            } else {
                return self.ProjectName
            }
        }
        func (self *ProjectV4) getBaseEnv() BaseEnvironment {
            return &self.Base
        }
        func (self *ProjectV4) getWorkingDir() string {
            if self.WorkingDir == "" && self.parentProject != nil {
                return self.parentProject.getWorkingDir()
            } else {
                return self.WorkingDir
            }
        }
        // Build cacheMountingPoints if it is nil.
        // Returns cache miss.
        func (self *ProjectV4) __cacheMountingPoints() bool {
            if self.cacheMountingPoints == nil {
                self.cacheMountingPoints = make(map[string]MountingPoint)
                for name, data := range(self.Mount) {
                    self.cacheMountingPoints[name] = &MountingPointV4{
                        host: data[0],
                        container: data[1],
                    }
                }
                return true
            } else {
                return false
            }

        }

        func (self *ProjectV4) getMountingPoints() map[string]MountingPoint {
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
        func (self *ProjectV4) getMacros() map[string]Macro {
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
        func (self *ProjectV4) toYAML() string {
            d, err := yaml.Marshal(&self)
            if err != nil {
                log.Fatalf("error: %v", err)
            }
            return string(d)
        }
        func (self *ProjectV4) fromYAML(bytes []byte) error {
            err := yaml.Unmarshal(bytes, self)
            if err == nil {
                if self.getSyntaxVersion() == "4" {
                    return nil
                } else {
                    return errors.New("Unexpected version number.")
                }
            } else {
                return err
            }
        }
        func (self *ProjectV4) getEnableGui() bool {
            if self.EnableGUI == "" && self.parentProject != nil {
                return self.parentProject.getEnableGui()
            } else {
                if self.EnableGUI == "true" {
                    return true
                } else {
                    return false
                }
            }
            return false
        }
        func (self *ProjectV4) getParentProject() Project {
            return self.parentProject
        }
        func (self *ProjectV4) setParentProject(project Project) error {
            self.parentProject = project
            self.cacheMacros = nil
            self.cacheMountingPoints = nil
            self.Base.setParentBase(project.getBaseEnv())
            return nil
        }
        func (self *ProjectV4) getPrivileged() bool {
            if self.Privileged == "" && self.parentProject != nil {
                return self.parentProject.getPrivileged()
            } else {
                if self.Privileged == "true" {
                    return true
                } else {
                    return false
                }
            }
            return false
        }

func NewProjectV4() *ProjectV4 {
    project := &ProjectV4 {
        SyntaxVersion: "4",
        Base: BaseEnvironmentV4{},
        Mount: make(map[string][]string),
        Macros: make(map[string]*MacroV4),
        parentProject: nil,
    }
    return project
}
