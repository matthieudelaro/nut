package main

import (
	"errors"
	log "github.com/Sirupsen/logrus"
    "fmt"
	"gopkg.in/yaml.v2"
    "path/filepath"
    "github.com/matthieudelaro/nut/persist"
)

////// README
// Nut must remain backward compatible with respect to nut configuration files.
// It must also be easy to add new features in nut files, without issues of
// performance to parse the file, and without applying modifications to all
// older versions to insure backward compatibility.
//
// The chosen solution is that new features should be defined in the interface,
// and default behavior/feature should be defined in base class (to be
// accessible) to all syntaxes version. Then override this default behavior in
// the new syntax version.
//
//  Interfaces Names| Base Class Names    | Version 2         | Version 3         | ...
//  ----------------+---------------------+-------------------+-------------------+----
//  Macro           | MacroBase           | MacroV2           | MacroV3           | ...
//  MountingPoint   | MountingPointBase   | MountingPointV2   | MountingPointV3   | ...
//  BaseEnvironment | BaseEnvironmentBase | BaseEnvironmentV2 | BaseEnvironmentV3 | ...
//  Project         | ProjectBase         | ProjectV2         | ProjectV3         | ...
//
// For each new syntax, copy/paste a project_vXXX.go file (ie the latest one)
// and change implementation to suit new requirements.
// For any new feature, create a virtual method in the interface Project
// (or its components Macro, MountingPoint, BaseEnvironment, etc) and
// implement a method in ProjectBase class (or its components MacroBase, etc)
// to define a default behavior for older versions.
//
// Old syntax versions should never be modified. Only the interface, the
// base class, and the new syntax should be updated.
//
// Note: When working with structs and embedding, everything is STATICALLY LINKED. All references are resolved at compile time.
// See https://github.com/luciotato/golang-notes/blob/master/OOP.md for more details.


////// Interfaces
/// Define what a project (and its components) should be able to do
type Macro interface {
    getUsage() string
    getActions() []string
    getAliases() []string
    getUsageText() string
    getDescription() string
    // override project settings
    setParentProject(project Project)
    getParentProject() Project
    getDockerImageName() (string, error)
    getWorkingDir() string
    getMountingPoints() map[string]MountingPoint
    getEnvironmentVariables() map[string]string
    getPorts() []string
    getEnableGui() bool
    getPrivileged() bool
    getEnableNvidiaDevices() bool
    getSecurityOpts() []string
}

type MountingPoint interface {
    fullHostPath() (string, error)
    fullContainerPath() (string, error)
}

type BaseEnvironment interface {
    getDockerImageName() (string, error)
    getFilePath() string
    setParentBase(parentBase BaseEnvironment) error
    getGitHub() string
    setGitHub(repositoryName string)
    setFilePath(filePath string)
}

type BooleanFeature interface {
    getEnable() bool
}

type Project interface {
    getSyntaxVersion() string
    getName() string
    getBaseEnv() BaseEnvironment
    getWorkingDir() string
    getMountingPoints() map[string]MountingPoint
    getMacros() map[string]Macro
    getEnvironmentVariables() map[string]string
    getPorts() []string
    getEnableGui() bool
    getEnableNvidiaDevices() bool
    toYAML() string
    fromYAML(bytes []byte) error
    getParentProject() Project
    setParentProject(project Project) error
    getPrivileged() bool
    getSecurityOpts() []string
}

/////// Base classes
/// Placeholder for basic behaviors, and default behaviors for retro-compatibility.
type MacroBase struct {
    project Project
}
        func (self *MacroBase) getUsage() string {
            return ""
        }
        func (self *MacroBase) getActions() []string {
            return []string{}
        }
        func (self *MacroBase) getAliases() []string {
            return []string{}
        }
        func (self *MacroBase) getUsageText() string {
            return ""
        }
        func (self *MacroBase) getDescription() string {
            return ""
        }
        func (self *MacroBase) setParentProject(project Project) {
            self.project = project
        }
        func (self *MacroBase) getParentProject() Project {
            return self.project
        }
        func (self *MacroBase) getDockerImageName() (string, error) {
            return self.project.getBaseEnv().getDockerImageName()
        }
        func (self *MacroBase) getWorkingDir() string {
            return self.project.getWorkingDir()
        }
        func (self *MacroBase) getMountingPoints() map[string]MountingPoint {
            return self.project.getMountingPoints()
        }
        func (self *MacroBase) getEnvironmentVariables() map[string]string {
            return self.project.getEnvironmentVariables()
        }
        func (self *MacroBase) getPorts() []string {
            return self.project.getPorts()
        }
        func (self *MacroBase) getEnableGui() bool {
            return self.project.getEnableGui()
        }
        func (self *MacroBase) getPrivileged() bool {
            return self.project.getPrivileged()
        }
        func (self *MacroBase) getSecurityOpts() []string {
        	return self.project.getSecurityOpts()
        }
        func (self *MacroBase) getEnableNvidiaDevices() bool {
            return self.project.getEnableNvidiaDevices()
        }

type MountingPointBase struct {
}
        func (self *MountingPointBase) fullHostPath() (string, error) {
            return "", errors.New("MountingPointBase.fullHostPath() must be overloaded.")
        }
        func (self *MountingPointBase) fullContainerPath() (string, error) {
            return "", errors.New("MountingPointBase.fullContainerPath() must be overloaded.")
        }

type BaseEnvironmentBase struct {
}
        func (self *BaseEnvironmentBase) getDockerImageName() (string, error) {
            return "", errors.New("BaseEnvironmentBase.getDockerImageName() must be overloaded.")
        }
        func (self *BaseEnvironmentBase) getFilePath() string{
            return ""
        }
        func (self *BaseEnvironmentBase) setParentBase(parentBase BaseEnvironment) error {
            return errors.New("This version of configuration cannot inherite configuration.")
        }
        func (self *BaseEnvironmentBase) getGitHub() string{
            return ""
        }
        func (self *BaseEnvironmentBase) setGitHub(repositoryName string) {
            log.Error("This version of configuration cannot inherite ",
                "configuration from GitHub ", repositoryName)
        }

        func (self *BaseEnvironmentBase) setFilePath(filePath string) {
            log.Error("This version of configuration cannot inherite ",
                "configuration from file ", filePath)
        }

type ProjectBase struct {
}
        func (self *ProjectBase) getSyntaxVersion() string {
            return "1"
        }
        func (self *ProjectBase) getName() string {
        	return "nut"
        }
        func (self *ProjectBase) getBaseEnv() BaseEnvironment {
            return &BaseEnvironmentBase{}
        }
        func (self *ProjectBase) getWorkingDir() string {
        	str, _ := filepath.Abs(".")
            return str
        }
        func (self *ProjectBase) getMountingPoints() map[string]MountingPoint {
            return make(map[string]MountingPoint)
        }
        func (self *ProjectBase) getMacros() map[string]Macro {
            return make(map[string]Macro)
        }
        func (self *ProjectBase) getEnvironmentVariables() map[string]string {
            return make(map[string]string)
        }
        func (self *ProjectBase) getPorts() []string {
        	return []string{}
        }
        func (self *ProjectBase) toYAML() string {
            d, err := yaml.Marshal(&self)
            if err != nil {
                log.Fatalf("error: %v", err)
            }
            return string(d)
        }
	    func (self *ProjectBase) fromYAML(bytes []byte) error {
	    	return errors.New("ProjectBase.fromYAML() must be overloaded.")
	    }
        func (self *ProjectBase) getEnableGui() bool {
            return false
        }
        func (self *ProjectBase) getEnableNvidiaDevices() bool {
            return false
        }
        func (self *ProjectBase) getParentProject() Project {
            return nil
        }
        func (self *ProjectBase) setParentProject(project Project) error {
            return errors.New("This version of configuration cannot inherite configuration.")
        }
        func (self *ProjectBase) getPrivileged() bool {
        	return false
        }
        func (self *ProjectBase) getSecurityOpts() []string {
        	return []string{}
        }


func loadProjectInheritance(nutFilePath string, project Project) (Project, error) {
    log.Debug("loadProjectInheritance ", nutFilePath, project)
    var err error
    if project == nil && nutFilePath != "" {
        project, err = parseNutFileAtPath(nutFilePath)
        if err != nil {
            return project, err
        }
    } else if project != nil && nutFilePath == "" {
        // nothing to do, we already have the project
    } else {
        return nil, errors.New("Bad usage: should provide either a project or name of a file.")
    }


    // project, err := parseNutFileAtPath(nutFilePath)
    // if err != nil {
        // return project, err
    // } else {
        parentFilePath := project.getBaseEnv().getFilePath()
        parentGitHub := project.getBaseEnv().getGitHub()
        if parentFilePath != "" && parentGitHub != "" {
            return nil, errors.New("Cannot inherite both from GitHub" +
                " and from a file.")
        }
        if parentGitHub != "" {
            // TODO: handle path properly, instead of just "."
            store, err := persist.InitStore(".") // TODO: handle the store in a smart way
            if err != nil {
                return nil, errors.New("Could not init storage: " + err.Error())
            }

            githubFile := filepath.Join(persist.EnvironmentsFolder, parentGitHub, "nut.yml")
            fullPath := filepath.Join(store.GetPath(), githubFile)
            _, err = persist.ReadFile(store, githubFile)
            if err != nil {
                fmt.Println("File from GitHub (" + parentGitHub + ") not available yet. Downloading...")
                fullPath, err = persist.StoreFile(store,
                    githubFile,
                    []byte{0})
                if err != nil {
                    return nil, errors.New(
                        "Could not prepare destination for file from GitHub: " + err.Error())
                }
                err = wget("https://raw.githubusercontent.com/" + parentGitHub + "/master/nut.yml",
                    fullPath)
                if err != nil {
                    return nil, errors.New("Could not download from GitHub: " + err.Error())
                }
                fmt.Println("File from GitHub (" + parentGitHub + ") downloaded.")
            }

            log.Debug("loadProjectInheritance inherite from ", fullPath)
            parent, err := loadProjectInheritance(fullPath, nil)
            if err != nil {
                return nil, errors.New("Could not inherite configuration from " + fullPath + ": " + err.Error())
            } else {
                project.setParentProject(parent)
            }
        }
        if parentFilePath != "" {
            parentFilePath = filepath.Join(filepath.Dir(nutFilePath), parentFilePath)
            log.Debug("loadProjectInheritance inherite from ", parentFilePath)
            parent, err := loadProjectInheritance(parentFilePath, nil)
            if err != nil {
                return nil, errors.New("Could not inherite configuration from " + parentFilePath + ": " + err.Error())
            } else {
                project.setParentProject(parent)
            }
        }
    // }
    return project, nil
}

func LoadProjectHierarchy(project Project) error {
    _, err := loadProjectInheritance("", project)
    return err
}

func LoadProject() (Project, error) {
    return loadProjectInheritance("nut.yml", nil)
}

func ProjectFromBytes(bytes []byte) (Project, error) {
	var err error
	for _, syntax := range getSyntaxes() {
		version := syntax.getSyntaxVersion()
		err = syntax.fromYAML(bytes)
		if err == nil {
			return syntax, nil
		} else {
			log.Debug("Not syntax ", version, ": ", err)
		}
	}
	return nil, err
}

// Compares a map of MountingPoint, and a given new MountingPoint.
// Returns the first conflict element from the map, or nil if
// there wasn't any conflict.
func CheckConflict(key string, newPoint MountingPoint, mountingPoints map[string]MountingPoint) MountingPoint {
    h, errh := newPoint.fullHostPath()
    c, errc := newPoint.fullContainerPath()

    for key2, mountingPoint2 := range mountingPoints {
        log.Debug("child point ", key)
        h2, errh2 := mountingPoint2.fullHostPath()
        c2, errc2 := mountingPoint2.fullContainerPath()
        if key2 == key ||
           h == h2 ||
           c == c2 ||
           errh != nil || errc != nil || errh2 != nil || errc2 != nil {
            log.Debug("conflic between mounting points ", key, " and ", key2)
            return mountingPoint2
       }
    }
    return nil
}

func getSyntaxes() []Project {
	return []Project{
        NewProjectV5(),
        NewProjectV4(),
		NewProjectV3(),
		NewProjectV2(),
	}
}
