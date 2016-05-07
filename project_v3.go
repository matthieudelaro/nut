package main

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"
    "path/filepath"
)

type MacroV3 struct {
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
        func (self *MacroV3) getUsage() string {
            return self.Usage
        }
        func (self *MacroV3) getActions() []string {
            return self.Actions
        }
        func (self *MacroV3) getAliases() []string {
            return self.Aliases
        }
        func (self *MacroV3) getUsageText() string {
            return self.UsageText
        }
        func (self *MacroV3) getDescription() string {
            return self.Description
        }

type MountingPointV3 struct {
	MountingPointBase `yaml:"inheritedValues,inline"`

    host string `yaml:host_path`
    container string `yaml:container_path`
}
        func (self *MountingPointV3) fullHostPath() (string, error) {
            absolutePath, err := filepath.Abs(self.host)
            return absolutePath, err
        }
        func (self *MountingPointV3) fullContainerPath() (string, error) {
            absolutePath, err := filepath.Abs(self.container)
            return absolutePath, err
        }

type BaseEnvironmentV3 struct {
	BaseEnvironmentBase `yaml:"inheritedValues,inline"`

    DockerImage string `yaml:"docker_image,omitempty"`
}
        func (self *BaseEnvironmentV3) getDockerImageName() (string, error) {
            if self.DockerImage != "" {
                return self.DockerImage, nil
            } else {
                return "", errors.New("Docker image has not been specified.")
            }
        }

type ProjectV3 struct {
	ProjectBase `yaml:"inheritedValues,inline"`

    SyntaxVersion string `yaml:"syntax_version"`
    ProjectName string `yaml:"project_name"`
    Base BaseEnvironmentV3 `yaml:"based_on"`
    WorkingDir string `yaml:"container_working_directory,omitempty"`
    Mount map[string][]string `yaml:"mount,omitempty"`
    Macros map[string]*MacroV3 `yaml:"macros,omitempty"`
    EnableGUI bool  `yaml:"enable_gui,omitempty"`
}
        func (self *ProjectV3) getSyntaxVersion() string {
            return self.SyntaxVersion
        }
        func (self *ProjectV3) getName() string {
            return self.ProjectName
        }
        func (self *ProjectV3) getBaseEnv() BaseEnvironment {
            return &self.Base
        }
        func (self *ProjectV3) getWorkingDir() string {
            return self.WorkingDir
        }
        func (self *ProjectV3) getMountingPoints() map[string]MountingPoint {
            mountingPoints := make(map[string]MountingPoint)
            for name, data := range(self.Mount) {
                mountingPoints[name] = &MountingPointV3{
                    host: data[0],
                    container: data[1],
                }
            }
            return mountingPoints
        }
        func (self *ProjectV3) getMacros() map[string]Macro {
        	macros := make(map[string]Macro) //, 0, len(self.Macros))
        	for name, data := range self.Macros {
        		macros[name] = data
                macros[name].setParentProject(self)
        	}
        	return macros
        }
        func (self *ProjectV3) toYAML() string {
            d, err := yaml.Marshal(&self)
            if err != nil {
                log.Fatalf("error: %v", err)
            }
            return string(d)
        }
	    func (self *ProjectV3) fromYAML(bytes []byte) error {
	    	err := yaml.Unmarshal(bytes, self)
			if err == nil {
				if self.getSyntaxVersion() == "3" {
					return nil
				} else {
					return errors.New("Unexpected version number.")
				}
			} else {
				return err
			}
	    }
        func (self *ProjectV3) getEnableGui() bool {
        	return self.EnableGUI
        }

func NewProjectV3() *ProjectV3 {
    project := &ProjectV3 {
        SyntaxVersion: "3",
		Base: BaseEnvironmentV3{},
		Mount: make(map[string][]string),
        Macros: make(map[string]*MacroV3),
    }
    return project
}
