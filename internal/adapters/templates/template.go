package templates

import (
	"fmt"
	"os"
	"text/template"
)

type Templates struct{}
type Files struct {
	FilePath     string
	TemplatePath string
}

type FileCreateParams struct {
	FileName     string
	TemplatePath string
	Data         interface{}
	Overwrite    bool
	Append       bool
}

func (t *Templates) Create(params *FileCreateParams) error {
	tpl, err := template.ParseFiles(params.TemplatePath)

	if err != nil {
		return err
	}

	_, err = os.Stat(params.FileName)

	if os.IsNotExist(err) || params.Append || params.Overwrite {
		fileFlags := os.O_WRONLY | os.O_CREATE

		if params.Append {
			fileFlags = fileFlags | os.O_APPEND
			fmt.Printf("appending to file %s\n", params.FileName)
		}

		file, err := os.OpenFile(params.FileName, fileFlags, 0644)

		if err != nil {
			return err
		}

		defer file.Close()

		if err := tpl.Execute(file, params.Data); err != nil {
			return err
		}
		return nil
	}

	fmt.Printf("Skipping: file %s already exists\n", params.FileName)
	return nil
}

func (t *Templates) CreateMany(data interface{}, files ...*Files) error {
	for _, file := range files {
		if err := t.Create(&FileCreateParams{
			FileName:     file.FilePath,
			TemplatePath: file.TemplatePath,
			Data:         data,
		}); err != nil {
			return err
		}
	}
	return nil
}
