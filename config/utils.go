package config

import (
    "gopkg.in/yaml.v2"
    "errors"
    "fmt"
    "io/ioutil"
    "path/filepath"
    log "github.com/Sirupsen/logrus"
    Utils "github.com/matthieudelaro/nut/utils"
    Persist "github.com/matthieudelaro/nut/persist"
)

const NutFileName = "nut.yml"
const NutOverrideFileName = "nut.override.yml"
const NutProjectFolderKey = "NUT_PROJECT_FOLDER"

/// Returns the boolean value, and whether there is or not a value
func TruthyString(s string) (bool, bool) {
    if s == "" {
        return false, false // undefined value, default is false
    } else {
        return (s == "true"), true // defined value
    }
}

// Compares a map of Volume, and a given new Volume.
// Returns the first conflict element from the map, or nil if
// there wasn't any conflict.
func CheckConflict(context Utils.Context, key string, newPoint Volume, mountingPoints map[string]Volume) Volume {
    h, errh := newPoint.getFullHostPath(context)
    c, errc := newPoint.getFullContainerPath(context)

    for key2, mountingPoint2 := range mountingPoints {
        // log.Debug("child point ", key)
        h2, errh2 := mountingPoint2.getFullHostPath(context)
        c2, errc2 := mountingPoint2.getFullContainerPath(context)
        if key2 == key ||
           h == h2 ||
           c == c2 ||
           (newPoint.getVolumeName() != "" && newPoint.getVolumeName() == mountingPoint2.getVolumeName()){
           // || errh != nil || errc != nil || errh2 != nil || errc2 != nil {
            log.Debug("conflic between mounting points ", key, " and ", key2)
            if errh != nil || errc != nil || errh2 != nil || errc2 != nil {
                log.Debug("warning while checking conflic between volumes ", key, " and ", key2)
            }
            return mountingPoint2
       }
    }
    return nil
}

func GetSyntaxes() []Project {
    return []Project{
        NewProjectV8(nil),
        NewProjectV7(nil),
        NewProjectV6(nil),
        NewProjectV5(nil),
        // NewProjectV4(), // represented by V5
        // NewProjectV3(), // represented by V5
        // NewProjectV2(), // represented by V5
    }
}

func ToYAML(in interface{}) (out []byte, err error) {
    return yaml.Marshal(in)
}

func ProjectFromYAML(bytes []byte) (Project, error) {
    logs := ""
    for _, syntax := range GetSyntaxes() {
        version := syntax.getSyntaxVersion()
        if err := yaml.Unmarshal(bytes, syntax); err == nil {
            parsedVersion := syntax.getSyntaxVersion()
            if parsedVersion == version {
                return syntax, nil
            } else if version == "5" && (
                parsedVersion == "4" ||
                parsedVersion == "3" ||
                parsedVersion == "2") { // support legacy syntaxes
                    return syntax, nil
            } else {
                logs += "\nVersion " + version + ": " + "syntax_version does not match."
            }
        } else {
            logs += "\nVersion " + version + ": " + err.Error()
        }
    }
    return nil, errors.New("Doesn't match any known syntax: " + logs)
}

/// Load project from given Github name
/// If not found on the file system, download it.
/// Returns the name of the file where
/// it has been saved from Github, and an error.
func DownloadFromGithub(name string, store Persist.Store) (string, error) {
    githubFile := filepath.Join(Persist.EnvironmentsFolder, name, NutFileName)
    _, fullPath, err := Persist.ReadFile(store, githubFile)
    if err != nil {
        fmt.Println("File from GitHub (" + name + ") not available yet. Downloading...")
        fullPath, err = Persist.StoreFile(store,
            githubFile,
            []byte{0})
        if err != nil {
            return "", errors.New(
                "Could not prepare destination for file from GitHub: " + err.Error())
        }
        err = Utils.Wget("https://raw.githubusercontent.com/" + name + "/master/nut.yml",
            fullPath)
        if err != nil {
            return "", errors.New("Could not download from GitHub: " + err.Error())
        }
        fmt.Println("File from GitHub (" + name + ") downloaded.")
    }
    return fullPath, nil
}

/// Parse Project from file. Given filename must be absolute.
func LoadProjectFromFile(fullPath string) (Project, error) {
    if bytes, err := ioutil.ReadFile(fullPath); err == nil {
        return ProjectFromYAML(bytes)
    } else {
        return nil, err
    }
}

/// Resolve dependencies of the given project
/// (e.g. by loading other projects from files or downloading them
/// from Github)
func ResolveDependencies(project Project, store Persist.Store, nutFilePath string) (err error) {
    baseEnv := GetBaseEnv(project)
    parentFilePath := baseEnv.getFilePath()
    parentGitHub := baseEnv.getGitHub()
    var parent Project
    if parentFilePath != "" && parentGitHub != "" {
        err = errors.New("Cannot inherite both from GitHub" +
            " and from a file.")
    } else if parentFilePath != "" {
        parentFilePath = filepath.Join(filepath.Dir(nutFilePath), parentFilePath)
        if parent, err = LoadProjectFromFile(parentFilePath); err == nil {
            if err = ResolveDependencies(parent, store, parentFilePath); err == nil {
                SetParentProject(project, parent)
            }
        }
    } else if parentGitHub != "" {
        var githubFilePath string
        if githubFilePath, err = DownloadFromGithub(parentGitHub, store); err == nil {
            if parent, err = LoadProjectFromFile(githubFilePath); err == nil {
                if err = ResolveDependencies(parent, store, githubFilePath); err == nil {
                    SetParentProject(project, parent)
                }
            }
        }
    }
    return
}


/// Look for a file from which to parse configuration (nut.yml in current
/// directory). Parse the file, and returns an updated context (root directory)
/// TODO: look for nut.yml file in parent folders
func FindProject(context Utils.Context) (Project, Utils.Context, error) {
    // var parentProject Project
    var project Project
    var err error
    var newContext Utils.Context
    var store Persist.Store

    foundDirectory := context.GetUserDirectory()
    // fullpath := filepath.Join(foundDirectory, NutFileName)
    fullpath := filepath.Join(foundDirectory, NutOverrideFileName)
    previousFoundDirectory := ""
    var exists bool

    searchHigher := true
    found := false
    for searchHigher == true {
        if exists, err = Utils.FileExists(fullpath); exists && err == nil {
            found = true
            searchHigher = false
        } else {
            fullpath = filepath.Join(foundDirectory, NutFileName)

            if exists, err = Utils.FileExists(fullpath); exists && err == nil {
                found = true
                searchHigher = false
            } else {
                previousFoundDirectory = foundDirectory
                foundDirectory = filepath.Dir(foundDirectory)
                // fullpath = filepath.Join(foundDirectory, NutFileName)
                fullpath = filepath.Join(foundDirectory, NutOverrideFileName)

                if foundDirectory == previousFoundDirectory {
                    searchHigher = false
                }
            }
        }
    }

    if found {
        if project, err = LoadProjectFromFile(fullpath); err == nil {
            log.Debug("Parsed from file: ", fullpath)
            if newContext, err = Utils.NewContext(foundDirectory, context.GetUserDirectory()); err == nil {
                log.Debug("Context updated: ", newContext)
                if store, err = Persist.InitStore(newContext.GetRootDirectory()); err == nil {
                    log.Debug("Store initialized: ", store)
                    if err = ResolveDependencies(project, store, fullpath); err == nil {
                        log.Debug("Resolved dependencies from file: ", fullpath)
                        return project, newContext, err
                    }
                }
            }
        }
    } else {
        err = errors.New("Could not find '" + NutFileName + "', neither in current directory nor in its parents.")
    }
    return nil, nil, err
}
