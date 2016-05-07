package main

import (
    "github.com/codegangsta/cli"
)

type MacroExec struct {
	MacroBase
    actions []string
}
        func (self *MacroExec) setActions(actions []string) {
            self.actions = actions
        }
        func (self *MacroExec) getActions() []string {
            return self.actions
        }

func exec(project Project, c *cli.Context, command string) {
    macro := &MacroExec{
        actions: []string{command},
    }
    // macro.setActions()
    execMacro(macro)
}
