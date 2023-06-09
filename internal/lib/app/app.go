package app

import (
	"embed"
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/x86-Yantras/code-gen/config"
	"github.com/x86-Yantras/code-gen/internal/adapters/templates"
	"github.com/x86-Yantras/code-gen/internal/constants"
)

type App struct {
	templates.TemplatesIface
	*AppModel
	Templater      templates.TemplatesIface
	Spec           *openapi3.T
	Config         *config.Config
	AppTemplateDir string
	Templates      embed.FS
	ServiceName    string
}

type AppModel struct {
	AppName        string
	AppDescription string
	ProjectPath    string
	LibDir         string
	RuntimeVersion string
	LibPath        string
}

func (a *App) Execute(command string) error {
	core, err := NewCore(&Dependencies{
		Spec:      a.Spec,
		Service:   a.ServiceName,
		Config:    a.Config,
		Templater: a.Templater,
	})

	if err != nil {
		return fmt.Errorf("%s, %s", "error initalizing core", err.Error())
	}

	switch command {
	case "init":
		fmt.Printf("Building %s project... \n", a.AppModel.AppName)
		err = a.InitProject()
	case "services":
		fmt.Printf("Building %s... \n", command)
		err = core.CreateService()
	case "http":
		fmt.Printf("Building %s... \n", command)
		err = core.CreateHttpAdapter()
	case "storage":
		fmt.Printf("Building %s... \n", command)
		err = a.CreateStorageAdapter()
	default:
		err = fmt.Errorf(constants.UndefinedCommandMsg, command)
	}

	if err != nil {
		return err
	}

	err = core.CreateErrors()
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Printf(constants.ProjectBuiltMsg, a.AppModel.AppName)
	return nil
}
