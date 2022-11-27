package app

import (
	"io/fs"
	"strings"
)

var GeneratorExceptionList = map[string]bool{"errors.js": true}

func (app *App) InitProject() error {
	err := fs.WalkDir(app.Templates, ".", func(path string, object fs.DirEntry, err error) error {
		if !strings.Contains(path, app.Config.ServiceDir) {
			objectPath := strings.Replace(path, app.AppTemplateDir, ".", 1)
			if object.IsDir() {
				if objectPath != "." {
					if err := app.CreateDir(objectPath); err != nil {
						return err
					}
				}
				return nil
			}

			// erros.go is created during service generation
			if !strings.Contains(objectPath, "errors.js") {
				if err := app.TemplateToFile(path, objectPath, app.AppModel); err != nil {
					return err
				}
			}
		}
		return nil
	})
	return err
}
