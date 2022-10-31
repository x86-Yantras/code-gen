package app

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/x86-Yantras/code-gen/config"
	"github.com/x86-Yantras/code-gen/internal/adapters/filesys"
	"github.com/x86-Yantras/code-gen/internal/adapters/templates"
	"github.com/x86-Yantras/code-gen/internal/constants"
)

type App struct {
	filesys.FsIface
	*AppModel
	Templater      templates.TemplatesIface
	Spec           *openapi3.T
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
		fmt.Printf("Building %s project... \n", a.AppModel.AppName)
		err = a.InitProject()
	case "services":
		fmt.Printf("Building services... \n")
		err = a.CreateService()

		// build cases for adapters

	default:
		return fmt.Errorf(constants.UndefinedCommandMsg, command)
	}
	fmt.Println()
	fmt.Printf(constants.ProjectBuiltMsg, a.AppModel.AppName)
	return err
}
