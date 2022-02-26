package config

import(
    "path/filepath"
    Utils "github.com/matthieudelaro/nut/utils"
    containerFilepath "github.com/matthieudelaro/nut/container/filepath"

    log "github.com/Sirupsen/logrus"
    "errors"
    "strings"

    "os"
)

type VolumeV8 struct {
    VolumeBase `yaml:"inheritedValues,inline"`

    VolumeName string `yaml:"volume_name,omitempty"`
    Host string `yaml:"host_path,omitempty"`
    Container string `yaml:"container_path"`
    Options string `yaml:"options,omitempty"`
}
        func (self *VolumeV8) getVolumeName() string {
            return self.VolumeName
        }
        func (self *VolumeV8) getOptions() string {
            return self.Options
        }
        // TODO: move the content of the two next function to an appropriate place
        // Same for V6 and V5
        func (self *VolumeV8) getFullHostPath(context Utils.Context) (string, error) {
            if self.Host == "" {
                return "", errors.New("Undefined host path")
            } else {
                // log.Debug("getFullHostPath value ", self.Host)
                res := filepath.Clean(self.Host)
                // log.Debug("getFullHostPath clean ", res)
                if !filepath.IsAbs(res) {
                    // log.Debug("getFullHostPath is not Abs")
                    res = filepath.Join(context.GetRootDirectory(), res)
                }
                // log.Debug("getFullHostPath value ", res)
                if strings.Contains(res, `:\`) { // path on windows. Eg: C:\\Users\
                    // log.Debug("getFullHostPath windows ", `:\`)
                    parts := strings.Split(res, `:\`)
                    parts[0] = strings.ToLower(parts[0]) // drive letter should be lower case
                    res = "//" + parts[0] + "/" + filepath.ToSlash(parts[1])
                }

                log.Debug("getFullHostPath res ", res)
                return res, nil
            }
        }
        func (self *VolumeV8) getFullContainerPath(context Utils.Context) (string, error) {
            if self.Container == "" {
                return "", errors.New("Undefined container path")
            } else {
                // log.Debug("getFullContainerPath value ", self.Container)
                clean := containerFilepath.ToSlash(containerFilepath.Clean(self.Container))
                // log.Debug("getFullContainerPath clean ", clean)
                if containerFilepath.IsAbs(clean) {
                    // log.Debug("getFullContainerPath isAbs")
                    return clean, nil
                } else {
                    // log.Debug("getFullContainerPath is not Abs")
                    log.Debug("getFullContainerPath return ", containerFilepath.Join(context.GetRootDirectory(), clean))
                    return containerFilepath.Join(context.GetRootDirectory(), clean), nil
                }
            }
        }


type DeviceV8 struct {
    DeviceBase `yaml:"inheritedValues,inline"`

    Host string `yaml:"host_path"`
    Container string `yaml:"container_path"`
    Options string `yaml:"options,omitempty"`
}
        func (self *DeviceV8) getHostPath() string {
            return self.Host
        }
        func (self *DeviceV8) getContainerPath() string {
            return self.Container
        }
        func (self *DeviceV8) getOptions() string {
            return self.Options
        }


type BaseEnvironmentV8 struct {
    BaseEnvironmentBase `yaml:"inheritedValues,inline"`

    FilePath string `yaml:"nut_file_path,omitempty"`
    GitHub string `yaml:"github,omitempty"`
}
        func (self *BaseEnvironmentV8) getFilePath() string{
            return self.FilePath
        }
        func (self *BaseEnvironmentV8) getGitHub() string{
            return self.GitHub
        }

type ConfigV8 struct {
    ConfigBase `yaml:"inheritedValues,inline"`

    DockerImage string `yaml:"docker_image,omitempty"`
    Volumes map[string]*VolumeV8 `yaml:"volumes,omitempty"`
    WorkingDir string `yaml:"container_working_directory,omitempty"`
    EnvironmentVariables map[string]string `yaml:"environment,omitempty"`
    Ports []string `yaml:"ports,omitempty"`
    EnableGUI string `yaml:"enable_gui,omitempty"`
    EnableNvidiaDevices string `yaml:"enable_nvidia_devices,omitempty"`
    Privileged string `yaml:"privileged,omitempty"`
    SecurityOpts []string `yaml:"security_opts,omitempty"`
    Detached string `yaml:"detached,omitempty"`
    UTSMode string `yaml:"uts,omitempty"`
    NetworkMode string `yaml:"net,omitempty"`
    Devices map[string]*DeviceV8 `yaml:"devices,omitempty"`
    EnableCurrentUser string `yaml:"enable_current_user,omitempty"`
    WorkInProjectFolderAs string `yaml:"work_in_project_folder_as,omitempty"`
    parent Config
}
        func (self *ConfigV8) getDockerImage() string {
            return self.DockerImage
        }
        func (self *ConfigV8) getNetworkMode() string {
            return self.NetworkMode
        }
        func (self *ConfigV8) getUTSMode() string {
            return self.UTSMode
        }
        func (self *ConfigV8) getParent() Config {
            return self.parent
        }
        func (self *ConfigV8) getWorkingDir() string {
            return self.WorkingDir
        }
        func (self *ConfigV8) getVolumes() map[string]Volume {
            cacheVolumes := make(map[string]Volume)
            for name, data := range(self.Volumes) {
                cacheVolumes[name] = data
            }
            return cacheVolumes
        }
        func (self *ConfigV8) getEnvironmentVariables() map[string]string {
            return self.EnvironmentVariables
        }
        func (self *ConfigV8) getDevices() map[string]Device {
            cacheVolumes := make(map[string]Device)
            for name, data := range(self.Devices) {
                cacheVolumes[name] = data
            }
            return cacheVolumes
        }
        func (self *ConfigV8) getPorts() []string {
            return self.Ports
        }
        func (self *ConfigV8) getEnableGui() (bool, bool) {
            return TruthyString(self.EnableGUI)
        }
        func (self *ConfigV8) getEnableNvidiaDevices() (bool, bool) {
            return TruthyString(self.EnableNvidiaDevices)
        }
        func (self *ConfigV8) getPrivileged() (bool, bool) {
            return TruthyString(self.Privileged)
        }
        func (self *ConfigV8) getDetached() (bool, bool) {
            return TruthyString(self.Detached)
        }
        func (self *ConfigV8) getEnableCurrentUser() (bool, bool) {
            return TruthyString(self.EnableCurrentUser)
        }
        func (self *ConfigV8) getSecurityOpts() []string {
            return self.SecurityOpts
        }
        func (self *ConfigV8) getWorkInProjectFolderAs() string {
            return self.WorkInProjectFolderAs
        }

type ProjectV8 struct {
    SyntaxVersion string `yaml:"syntax_version"`
    ProjectName string `yaml:"project_name"`
    Base BaseEnvironmentV8 `yaml:"based_on,omitempty"`
    Macros map[string]*MacroV8 `yaml:"macros,omitempty"`
    parent Project

    ProjectBase `yaml:"inheritedValues,inline"`
    ConfigV8 `yaml:"inheritedValues,inline"`
}
        func (self *ProjectV8) getSyntaxVersion() string {
            return self.SyntaxVersion
        }
        func (self *ProjectV8) getProjectName() string {
            return self.ProjectName
        }
        func (self *ProjectV8) getBaseEnv() BaseEnvironment {
            return &self.Base
        }
        func (self *ProjectV8) getMacros() map[string]Macro {
            // make the list of macros
            cacheMacros := make(map[string]Macro)
            for name, data := range self.Macros {
                data.parent = self
                cacheMacros[name] = data
            }
            return cacheMacros
        }
        // func (self *ProjectV8) createMacro(usage string, commands []string) Macro {
        //     return &MacroV8 {
        //         ConfigV8: *NewConfigV8(self,),
        //         Usage: usage,
        //         Actions: commands,
        //     }
        // }
        func (self *ProjectV8) getParent() Config {
            return self.parent
        }
        func (self *ProjectV8) getParentProject() Project {
            return self.parent
        }
        func (self *ProjectV8) setParentProject(project Project) {
            self.parent = project
        }

type MacroV8 struct {
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

    MacroBase `yaml:"inheritedValues,inline"`
    ConfigV8 `yaml:"inheritedValues,inline"`
}
        func (self *MacroV8) setParentProject(project Project) {
            self.ConfigV8.parent = project
        }
        func (self *MacroV8) getUsage() string {
            return self.Usage
        }
        func (self *MacroV8) getActions() []string {
            for i := range self.Actions {
                self.Actions[i] = os.ExpandEnv(self.Actions[i])
            }
            return self.Actions
        }
        func (self *MacroV8) getAliases() []string {
            return self.Aliases
        }
        func (self *MacroV8) getUsageText() string {
            return self.UsageText
        }
        func (self *MacroV8) getDescription() string {
            return self.Description
        }


func NewConfigV8(parent Config) *ConfigV8 {
    return &ConfigV8{
        Volumes: make(map[string]*VolumeV8),
        Devices: make(map[string]*DeviceV8),
        EnvironmentVariables: map[string]string{},
        parent: parent,
    }
}

func NewProjectV8(parent Project) *ProjectV8 {
    project := &ProjectV8 {
        SyntaxVersion: "8",
        Macros: make(map[string]*MacroV8),
        ConfigV8: *NewConfigV8(nil),
        parent: parent,
    }
    return project
}
