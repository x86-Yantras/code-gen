package main

import (
	"context"
	"embed"
	"fmt"
	"os"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/x86-Yantras/code-gen/config"
	"github.com/x86-Yantras/code-gen/internal/adapters/filesys"
	"github.com/x86-Yantras/code-gen/internal/adapters/templates"
	"github.com/x86-Yantras/code-gen/internal/constants"
	"github.com/x86-Yantras/code-gen/internal/lib/app"
)

//go:embed templates
var template embed.FS

func main() {

	args := os.Args

	if len(args) < 4 {
		fmt.Println(`Usage:
		code-gen [api-specs.yaml][language][command][service]

		command list:
		init: init project
		services: generate services
		http: generates http layer
		storage: generates storage layer
		
		service(optional): name of the service, should match the first tag value in spec file`)
		os.Exit(0)
	}

	specFile := args[1]
	appLang := args[2]
	command := args[3]
	service := ""

	if len(args) == 5 {
		service = args[4]
	}

	ctx := context.Background()
	loader := &openapi3.Loader{
		Context:               ctx,
		IsExternalRefsAllowed: true,
	}

	doc, err := loader.LoadFromFile(specFile)

	if err != nil {
		panic(err)
	}

	config := config.New(appLang)

	if err != nil {
		panic(err)
	}

	templater := &templates.Templates{
		template,
	}
	appModel := &app.AppModel{
		AppName:        doc.Info.Title,
		AppDescription: doc.Info.Description,
		ProjectPath:    config.ProjectPath,
		LibDir:         config.LibDir,
	}
	app := &app.App{
		&filesys.Fs{},
		appModel,
		templater,
		doc,
		config,
		fmt.Sprintf("%s/%s", constants.TemplatesDir, appLang),
		template,
		service,
	}

	err = app.Execute(command)

	if err != nil {
		panic(err)
	}
}
