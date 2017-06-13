package main

import (
    log "github.com/Sirupsen/logrus"
    "io/ioutil"
    "github.com/codegangsta/cli"
    "fmt"
    "path/filepath"
    Config "github.com/matthieudelaro/nut/config"
    Utils "github.com/matthieudelaro/nut/utils"
    Persist "github.com/matthieudelaro/nut/persist"
    containerFilepath "github.com/matthieudelaro/nut/container/filepath"
)

// create a nut.yml at the current path
func initSubcommand(c *cli.Context, context Utils.Context, gitHubFlag string) {
    log.Debug("initSubcommand context: ", context)
    name := filepath.Base(context.GetUserDirectory())

    var project Config.Project
    if gitHubFlag == "" {
        defaultProject := Config.NewProjectV6(nil)
        defaultProject.ProjectName = name
        defaultProject.DockerImage = "matthieudelaro/golang:1.7-cross"
        defaultProject.Macros["build"] = &Config.MacroV6{
            Usage: "build the project in the container",
            Actions: []string{"go build -o nut"},
            ConfigV6: *Config.NewConfigV6(nil),
        }
        defaultProject.Macros["run"] = &Config.MacroV6{
            Usage: "run the project in the container",
            Actions: []string{"./nut"},
            ConfigV6: *Config.NewConfigV6(nil),
        }
        project = defaultProject
    } else {
        if store, err := Persist.InitStore(context.GetUserDirectory()); err == nil {
            log.Debug("Store initialized: ", store)
            if filename, err := Config.DownloadFromGithub(gitHubFlag, store); err == nil {
                log.Debug("Downloaded from Github to ", filename)
                if parentProject, err := Config.LoadProjectFromFile(filename); err == nil {
                    log.Debug("Parsed from Github: ", parentProject)
                    if err = Config.ResolveDependencies(parentProject, store, filename); err == nil {
                        log.Debug("Resolved dependencies from Github: ", filename)
                        finalProject := Config.NewProjectV6(parentProject)
                        finalProject.ProjectName = name
                        finalProject.Base.GitHub = gitHubFlag

                        // modified project depending on parent configuration
                        // 1 - mount folder "." if not already mounted by parent configuration
                        mountingPointName := "main"
                        hostDir := "."
                        containerDir := containerFilepath.Join("/nut", name)
                        mountingPoint := &Config.VolumeV6{
                            Host: hostDir,
                            Container: containerDir,
                        }
                        if Config.CheckConflict(context, mountingPointName, mountingPoint, Config.GetVolumes(finalProject, context)) == nil {
                            finalProject.ConfigV6.Mount[mountingPointName] = []string{hostDir, containerDir}
                        }
                        // 2 - set working directory to "." if not specified otherwise
                        if Config.GetWorkingDir(finalProject) == "" {
                            finalProject.ConfigV6.WorkingDir = containerDir
                        }

                        project = finalProject
                    }
                } else {
                    log.Error("Error while parsing from GitHub: ", err)
                }
            } else {
                log.Error("Error while loading from GitHub: ", err)
            }
        } else {
            log.Error("Error while loading nut files: ", err)
        }
    }

    if project == nil {
        return
    }
    if data, err := Config.ToYAML(project); err != nil {
        log.Error("Could not generate default project configuration:", err)
    } else {
        fmt.Println("Project configuration:")
        fmt.Println("")
        dataString := string(data)
        fmt.Println(dataString)
        // check is nut.yml exists at the current path
        nutfileName := filepath.Join(context.GetUserDirectory(), Config.NutFileName)  // TODO: pick this name from a well documented and centralized list of legit nut file names
        if exists, err := Utils.FileExists(nutfileName); exists && err == nil {
            log.Error("Could not save new Nut project because a nut.yml file already exists")
            return
        } else {
            err := ioutil.WriteFile(nutfileName, data, 0644) // TODO: discuss this permission level
            if err != nil {
                log.Error(err)
            } else {
                fmt.Println("Project configuration saved in ", nutfileName) // TODO: make sure that bug Project configuration saved in %!(EXTRA string=./nut.yml) is fixed.
            }
        }
    }
}
