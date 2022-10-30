package templates

import (
	"fmt"
	"os"
	"text/template"
)

type Templates struct{}

func (t *Templates) Create(fileName, templatePath string, data interface{}) error {
	tpl, err := template.ParseFiles(templatePath)

	if err != nil {
		return err
	}
	if _, err := os.Stat(fileName); os.IsNotExist(err) {

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
