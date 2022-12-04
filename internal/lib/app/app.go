package app

import (
	"embed"
	"errors"
	"fmt"
	"net/http"

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
	Templates      embed.FS
	ServiceName    string
}

type AppModel struct {
	AppName        string
	AppDescription string
}

type ServiceSpecs struct {
	ServiceName string
	Method      string
	Spec        *openapi3.Operation
}

func (a *App) Execute(command string) error {
	var err error
	switch command {
	case "init":
		fmt.Printf("Building %s project... \n", a.AppModel.AppName)
		err = a.InitProject()
	case "services":
		fmt.Printf("Building %s... \n", command)
		err = a.CreateService()

	case "http":
		fmt.Printf("Building %s... \n", command)
		err = a.CreateHttpAdapter()
	case "storage":
		fmt.Printf("Building %s... \n", command)
		err = a.CreateStorageAdapter()
	default:
		return fmt.Errorf(constants.UndefinedCommandMsg, command)
	}

	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Printf(constants.ProjectBuiltMsg, a.AppModel.AppName)
	return nil
}

func (a *App) GetServiceSpecs(service string) ([]*ServiceSpecs, error) {
	pathsLen := len(a.Spec.Paths)

	if service != nil {
		pathsLen = 1
	}

	specs := make(map[string][]*ServiceSpecs)

	for path, doc := range a.Spec.Paths {
		specList := map[string]*openapi3.Operation{
			http.MethodGet:    doc.Get,
			http.MethodPost:   doc.Post,
			http.MethodPatch:  doc.Patch,
			http.MethodPut:    doc.Put,
			http.MethodDelete: doc.Delete,
		}

		serviceName := ""
		for method, spec := range specList {
			if spec != nil {
				if len(spec.Tags) == 0 {
					return errors.New(constants.TagMissingErr)
				}

				if serviceName != "" {
					if serviceName != spec.Tags[0] {
						err := fmt.Sprintf(constants.TagDifferentErr)
						return errors.New(constants.TagDifferentErr, spec.Tags[0], path)
					}
				}

			}
		}
	}
}
