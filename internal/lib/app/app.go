package app

import (
	"fmt"

	"github.com/x86-Yantras/code-gen/config"
	"github.com/x86-Yantras/code-gen/internal/adapters/filesys"
	"github.com/x86-Yantras/code-gen/internal/adapters/templates"
	"github.com/x86-Yantras/code-gen/internal/constants"
)

type App struct {
	filesys.FsIface
	*AppModel
	Templater      templates.TemplatesIface
	Config         *config.Config
	AppTemplateDir string
}

type AppModel struct {
	AppName        string
	AppDescription string
}

func (a *App) Execute(command string) error {
	var err error
	switch command {
	case "init":
		fmt.Printf("Building %s project \n", a.AppModel.AppName)
		err = a.InitProject()

		// build cases for service and adapters
	default:
		fmt.Printf(constants.UndefinedCommandMsg, command)
	}
	fmt.Printf(constants.ProjectBuiltMsg, a.AppModel.AppName)
	return err
}
