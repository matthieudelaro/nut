package main

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"
    "path/filepath"
)

type MacroV2 struct {
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
        func (self *MacroV2) getUsage() string {
            return self.Usage
        }
        func (self *MacroV2) getActions() []string {
            return self.Actions
        }
        func (self *MacroV2) getAliases() []string {
            return self.Aliases
        }
        func (self *MacroV2) getUsageText() string {
            return self.UsageText
        }
        func (self *MacroV2) getDescription() string {
            return self.Description
        }

type MountingPointV2 struct {
	MountingPointBase `yaml:"inheritedValues,inline"`

    host string `yaml:host_path`
    container string `yaml:container_path`
}
        func (self *MountingPointV2) fullHostPath() (string, error) {
            absolutePath, err := filepath.Abs(self.host)
            return absolutePath, err
        }
        func (self *MountingPointV2) fullContainerPath() (string, error) {
            absolutePath, err := filepath.Abs(self.container)
            return absolutePath, err
        }

type BaseEnvironmentV2 struct {
	BaseEnvironmentBase `yaml:"inheritedValues,inline"`

    DockerImage string `yaml:"docker_image,omitempty"`
}
        func (self *BaseEnvironmentV2) getDockerImageName() (string, error) {
            if self.DockerImage != "" {
                return self.DockerImage, nil
            } else {
                return "", errors.New("Docker image has not been specified.")
            }
        }

type ProjectV2 struct {
	ProjectBase `yaml:"inheritedValues,inline"`

    SyntaxVersion string `yaml:"syntax_version"`
    ProjectName string `yaml:"project_name"`
    Base BaseEnvironmentV2 `yaml:"based_on"`
    WorkingDir string `yaml:"container_working_directory,omitempty"`
    Mount map[string][]string `yaml:"mount,omitempty"`
    Macros map[string]*MacroV2 `yaml:"macros,omitempty"`
}
        func (self *ProjectV2) getSyntaxVersion() string {
            return self.SyntaxVersion
        }
        func (self *ProjectV2) getName() string {
            return self.ProjectName
        }
        func (self *ProjectV2) getBaseEnv() BaseEnvironment {
            return &self.Base
        }
        func (self *ProjectV2) getWorkingDir() string {
            return self.WorkingDir
        }
        func (self *ProjectV2) getMountingPoints() map[string]MountingPoint {
            mountingPoints := make(map[string]MountingPoint)
            for name, data := range(self.Mount) {
                mountingPoints[name] = &MountingPointV2{
                    host: data[0],
                    container: data[1],
                }
            }
            return mountingPoints
        }
        func (self *ProjectV2) getMacros() map[string]Macro {
        	macros := make(map[string]Macro) //, 0, len(self.Macros))
        	for name, data := range self.Macros {
        		macros[name] = data
                macros[name].setParentProject(self)
        	}
        	return macros
        }
        func (self *ProjectV2) toYAML() string {
            d, err := yaml.Marshal(&self)
            if err != nil {
                log.Fatalf("error: %v", err)
            }
            return string(d)
        }
	    func (self *ProjectV2) fromYAML(bytes []byte) error {
	    	err := yaml.Unmarshal(bytes, self)
			if err == nil {
				if self.getSyntaxVersion() == "2" {
					return nil
				} else {
					return errors.New("Unexpected version number.")
				}
			} else {
				return err
			}
	    }

func NewProjectV2() *ProjectV2 {
    project := &ProjectV2 {
        SyntaxVersion: "2",
    }
    return project
}
