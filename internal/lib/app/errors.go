package app

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/x86-Yantras/code-gen/internal/adapters/templates"
	"github.com/x86-Yantras/code-gen/internal/constants"
)

type Errors struct {
	Name           string
	HttpStatusCode int `json:"httpStatusCode"`
}

type ErrCodes struct {
	Name string
	Code string
}

func (app *App) BuildErrors(errs interface{}) error {

	// Build errors dir
	errDir := fmt.Sprintf("%s/%s", app.Config.LibDir, constants.Errors)

	if err := app.CreateDir(errDir); err != nil {
		return err
	}

	errors := []*Errors{}

	if err := json.Unmarshal(errs.(json.RawMessage), &errors); err != nil {
		return err
	}

	// Build custom errors
	for _, err := range errors {
		errFile := fmt.Sprintf("%s/%s/%s%s", app.Config.LibDir, constants.Errors, app.ToLowerFirst(err.Name), app.Config.FileExt)
		tplPath := fmt.Sprintf("%s/%s/%s/%s%s%s", app.AppTemplateDir, app.Config.LibDir, constants.Errors, constants.Errors, app.Config.FileExt, constants.TemplateExtension)

		err := app.Templater.Create(&templates.FileCreateParams{
			FileName:     errFile,
			TemplatePath: tplPath,
			Data:         err,
			Overwrite:    true,
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func (app *App) BuildErrorCodes(errCodes interface{}) error {

	errorCodes := []*ErrCodes{}

	if err := json.Unmarshal(errCodes.(json.RawMessage), &errorCodes); err != nil {
		return err
	}

	// Build custom errors
	errFile := fmt.Sprintf("%s/%s/%s%s", app.Config.LibDir, constants.Errors, constants.ErrCodes, app.Config.FileExt)
	tplPath := fmt.Sprintf("%s/%s/%s/%s%s%s", app.AppTemplateDir, app.Config.LibDir, constants.Errors, constants.ErrCodes, app.Config.FileExt, constants.TemplateExtension)

	fmt.Println("building err codes")
	err := app.Templater.Create(&templates.FileCreateParams{
		FileName:     errFile,
		TemplatePath: tplPath,
		Data:         errorCodes,
		Overwrite:    true,
	})

	return err
}

func (app *App) ToLowerFirst(input string) string {
	return strings.ToLower(string(input[0])) + input[1:]
}
