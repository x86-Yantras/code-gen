package main

import (
	"context"
	"fmt"
	"os"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/x86-Yantras/code-gen/config"
	"github.com/x86-Yantras/code-gen/internal/adapters/filesys"
	"github.com/x86-Yantras/code-gen/internal/adapters/templates"
	"github.com/x86-Yantras/code-gen/internal/constants"
	"github.com/x86-Yantras/code-gen/internal/lib/app"
)

func main() {

	args := os.Args

	if len(args) < 4 {
		fmt.Println(`Usage:
		code-gen [api-specs.yaml][language][command]

		command list:
		init: init project`)
		os.Exit(0)
	}

	specFile := args[1]
	appLang := args[2]
	command := args[3]

	ctx := context.Background()
	loader := &openapi3.Loader{
		Context:               ctx,
		IsExternalRefsAllowed: true,
	}

	doc, err := loader.LoadFromFile(specFile)

	if err != nil {
		panic(err)
	}

	config, err := config.New(appLang)

	if err != nil {
		panic(err)
	}

	templater := &templates.Templates{}
	appModel := &app.AppModel{
		AppName:        doc.Info.Title,
		AppDescription: doc.Info.Description,
	}
	app := &app.App{
		&filesys.Fs{},
		appModel,
		templater,
		config,
		fmt.Sprintf("%s/%s", constants.TemplatesDir, appLang),
	}

	err = app.Execute(command)

	if err != nil {
		panic(err)
	}

	// for _, path := range doc.Paths {
	// 	fmt.Printf("%+v", path)
	// }
}
