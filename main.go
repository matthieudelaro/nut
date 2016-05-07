package main

import (
	"os"
    "sort"
	"github.com/codegangsta/cli"
	log "github.com/Sirupsen/logrus"
    "reflect"
    "github.com/matthieudelaro/nut/persist"
)

func main() {
    log.SetLevel(log.ErrorLevel)
	// log.SetLevel(log.DebugLevel)

	// try to parse the folder's configuration to add macros
	var macros []cli.Command
	project, err := LoadProject()
	if err == nil {
        // Macros are stored in a random order.
        // But we want to display them in the same order everytime.
        // So sort the names of the macros.
        projectMacros := project.getMacros()
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
                if macro.getUsage() == "" {
                    usage += "undefined usage. define one with 'usage' property for this macro."
                } else {
                    usage += macro.getUsage()
                }

    			macros[index] = cli.Command{
    				Name:  nameClosure,
    				Usage: usage,
                    Aliases: macro.getAliases(),
                    UsageText: macro.getUsageText(),
                    Description: macro.getDescription(),
    				Action: func(c *cli.Context) {
    					execMacro(macro)
    				},
    			}
            }
		}
	} else {
        log.Error("Could not parse nut.yml: " + err.Error())
		macros = []cli.Command{}
	}

    initFlag := false
    statusFlag := false
    logsFlag := false
    cleanFlag := false
    execFlag := ""
    gitHubFlag := ""

    app := cli.NewApp()
	app.Name = "nut"
	app.Version = "0.1.0 dev"
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
        cli.BoolFlag{
            Name:  "status",
            Usage: "gives status of the dev env",
            Destination: &statusFlag,
        },
    }
    defaultAction := app.Action
    app.Action = func(c *cli.Context) {
        if logsFlag {
            log.SetLevel(log.DebugLevel)
        }
        if cleanFlag {
            persist.CleanStoreFromProject(".") // TODO: do better than "."
        }
        if statusFlag {
            status()
            return
        }
        if initFlag {
            initSubcommand(c, gitHubFlag)
            return
        }
        if execFlag != "" {
            if project != nil {
                exec(project, c, execFlag)
            } else {
                log.Error("Could not find nut configuration.")
            }
            return
        }
        defaultAction(c)
    }

	app.Commands = macros

    app.Run(os.Args)
}
