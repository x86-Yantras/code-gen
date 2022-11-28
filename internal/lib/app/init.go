package app

import (
	"io/fs"
	"strings"
)

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

			// errors is created during service generation
			if !strings.Contains(objectPath, "err") {
				if err := app.TemplateToFile(path, objectPath, app.AppModel); err != nil {
					return err
				}
			}
		}
		return nil
	})
	return err
}
