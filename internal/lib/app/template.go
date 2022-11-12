package app

import (
	"strings"

	"github.com/x86-Yantras/code-gen/internal/constants"
)

func (app *App) TemplateToFile(templatePath, filepath string, data interface{}) error {

	fileName := strings.Replace(filepath, constants.TemplateExtension, "", 1)
	err := app.Templater.Create(fileName, templatePath, data)

	if err != nil {
		return err
	}

	return nil
}
