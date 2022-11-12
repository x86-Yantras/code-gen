package templates

import (
	"fmt"
	"os"
	"text/template"
)

type Templates struct{}

func (t *Templates) Create(fileName, templatePath string, data interface{}, overwrite ...bool) error {
	tpl, err := template.ParseFiles(templatePath)

	if err != nil {
		return err
	}
	_, err = os.Stat(fileName)

	if os.IsNotExist(err) || (len(overwrite) > 0 && overwrite[0]) {
		file, err := os.Create(fileName)

		if err != nil {
			return err
		}

		defer file.Close()

		if err := tpl.Execute(file, data); err != nil {
			return err
		}
		return nil
	}

	fmt.Printf("Skipping: file %s already exists\n", fileName)
	return nil
}
