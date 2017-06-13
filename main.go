package main

import (
    "os"
    "sort"
    "github.com/codegangsta/cli"
    log "github.com/Sirupsen/logrus"
    "reflect"
    "github.com/matthieudelaro/nut/persist"
    Utils "github.com/matthieudelaro/nut/utils"
    Config "github.com/matthieudelaro/nut/config"
)

func main() {
    log.SetLevel(log.ErrorLevel)
    // log.SetLevel(log.DebugLevel)

    // get the folder where nut has been called
    pwd, err := os.Getwd()
    if err != nil {
        log.Fatal("find path error: ", err)
    }

    context, err := Utils.NewContext(pwd, pwd)
    if err != nil {
        log.Fatal("Error with context: ", err)
    }

    log.Debug("main first context: ", context)
    // try to parse the folder's configuration to add macros
    var macros []cli.Command
    project, projectContext, err := Config.FindProject(context)
    log.Debug("main second context: ", projectContext)
    projectMacros := map[string]Config.Macro{}

    initFlag := false
    logsFlag := false
    cleanFlag := false
    execFlag := ""
    inheriteConfigMacroFlag := ""
    gitHubFlag := ""
    macroFlag := ""
    useDockerCLIFlag := DOCKERCLI_FLAG_DEFAULT_VALUE

    if err == nil {
        // Macros are stored in a random order.
        // But we want to display them in the same order everytime.
        // So sort the names of the macros.
        projectMacros = Config.GetMacros(project)
        macroNamesOrdered := make([]string, 0, len(projectMacros))
        for key, _ := range projectMacros {
            macroNamesOrdered = append(macroNamesOrdered, key)
        }
        sort.Strings(macroNamesOrdered)


        macros = make([]cli.Command, len(projectMacros))
        log.Debug(len(macroNamesOrdered), " macros")
        for index, name := range macroNamesOrdered {
            macro := projectMacros[name]

            log.Debug("macro ", name, ": ", macro)
            // If the nut file containes a macro which has not any field defined,
            // then the macro will be nil.
            // So check value of macro:
            if macro == reflect.Zero(reflect.TypeOf(macro)).Interface() { // just check whether it macro is nil
                // it seems uselessly complicated, but:
                // if macro == nil { // doesn't work the trivial way one could expect
                // if macro.(*MacroBase) == nil { // panic: interface conversion: main.Macro is *main.MacroV3, not *main.MacroBase
                // Checking for nil doesn't seems like a good solution. TODO: ? require at least a "usage" field for each macro in the nut file?
                log.Warn("Undefined properties of macro " + name + ".")
            } else {
                nameClosure := name
                usage := "macro: "
                if macroUsage := Config.GetUsage(macro); macroUsage == "" {
                    usage += "undefined usage. define one with 'usage' property for this macro."
                } else {
                    usage += macroUsage
                }

                macros[index] = cli.Command{
                    Name:  nameClosure,
                    Usage: usage,
                    Aliases: Config.GetAliases(macro),
                    UsageText: Config.GetUsageText(macro),
                    Description: Config.GetDescription(macro),
                    Action: func(c *cli.Context) error {
                        if logsFlag {
                            log.SetLevel(log.DebugLevel)
                        }
                        execMacro(macro, projectContext, useDockerCLIFlag != DOCKERCLI_FLAG_DEFAULT_VALUE)
                        return nil
                    },
                }
            }
        }
    } else {
        log.Error("Could not parse configuration: " + err.Error())
        macros = []cli.Command{}
    }


    app := cli.NewApp()
    app.Name = "nut"
    app.Version = "0.1.4 dev"   //RELEASE_BUILD_PLACEHOLDER//
    app.Usage = "the development environment, containerized"
    // app.EnableBashCompletion = true
    app.Flags = []cli.Flag {
        cli.BoolFlag{
            Name:  "clean",
            Usage: "clean all data stored in .nut",
            Destination: &cleanFlag,
        },
        cli.BoolFlag{
            Name:        "init",
            Usage:       "initialize a nut project",
            Destination: &initFlag,
        },
        cli.StringFlag{
            Name:  "github",
            Usage: "Use with --init: provide a GitHub repository to initialize Nut.",
            Destination: &gitHubFlag,
        },
        cli.BoolFlag{
            Name:        "logs",
            Usage:       "Use with --exec: display log messages. Useful for contributors and to report an issue",
            Destination: &logsFlag,
        },
        cli.StringFlag{
            Name:  "exec",
            Usage: "execute a command in a container.",
            Destination: &execFlag,
        },
        cli.StringFlag{
            Name:  "macro",
            Usage: "Name of the macro to execute. Use with --logs.",
            Destination: &macroFlag,
        },
        cli.StringFlag{
            Name:  "extend-macro",
            Usage: "Use with --exec: name of the macro from which to inherite the configuration",
            Destination: &inheriteConfigMacroFlag,
        },
    }
    if DOCKERCLI_FLAG_DEFAULT_VALUE == false {
        app.Flags = append(app.Flags, cli.BoolFlag{
            Name:  "dockercli",
            Usage: "Use Docker CLI, instead of using Docker API to reach the host directly.",
            Destination: &useDockerCLIFlag,
        })
    } else {
        app.Flags = append(app.Flags, cli.BoolFlag{
            Name:  "dockerapi",
            Usage: "Use Docker API to reach the host directly, instead of using Docker CLI.",
            Destination: &useDockerCLIFlag,
        })
    }

    defaultAction := app.Action
    app.Action = func(c *cli.Context) error {
        if logsFlag {
            log.SetLevel(log.DebugLevel)
        }

        if cleanFlag {
            if project != nil {
                persist.CleanStoreFromProject(projectContext.GetRootDirectory())
            }
        } else if initFlag {
            log.Debug("main context for init: ", context)
            initSubcommand(c, context, gitHubFlag)
        } else if execFlag != "" {
            if project != nil {
                if inheriteConfigMacroFlag == "" {
                    execInContainer([]string{execFlag}, project, projectContext, useDockerCLIFlag != DOCKERCLI_FLAG_DEFAULT_VALUE)
                } else if macro, ok := projectMacros[inheriteConfigMacroFlag]; ok && project != nil {
                    execInContainer([]string{execFlag}, macro, projectContext, useDockerCLIFlag != DOCKERCLI_FLAG_DEFAULT_VALUE)
                } else {
                    log.Error("Undefined macro " + macroFlag)
                    cli.ShowAppHelp(c)
                    return cli.NewExitError("", 42)
                }
            } else {
                log.Error("Could not find nut configuration.")
                cli.ShowAppHelp(c)
                return cli.NewExitError("", 42)
            }
        } else if macroFlag != "" {
            if macro, ok := projectMacros[macroFlag]; ok && project != nil {
                execMacro(macro, projectContext, useDockerCLIFlag != DOCKERCLI_FLAG_DEFAULT_VALUE)
            } else {
                log.Error("Undefined macro " + macroFlag)
                cli.ShowAppHelp(c)
                return cli.NewExitError("", 42)
            }
        } else {
            cli.HandleAction(defaultAction, c)
        }
        return nil
    }

    app.Commands = macros
    app.CommandNotFound = func(c *cli.Context, macroName string) {
        log.Error("Undefined macro " + macroName)
        cli.ShowAppHelp(c)
        os.Exit(42)
    }

    app.Run(os.Args)
}
