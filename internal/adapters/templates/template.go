package templates

import (
	"embed"
	"fmt"
	"os"
	"text/template"
)

type Templates struct {
	Tmpls embed.FS
}
type Files struct {
	FilePath     string
	TemplatePath string
	Overwrite    bool
}

type FileCreateParams struct {
	FileName     string
	TemplatePath string
	Data         interface{}
	Overwrite    bool
	Append       bool
}

func (t *Templates) Create(params *FileCreateParams) error {

	tpl, err := template.ParseFS(t.Tmpls, params.TemplatePath)

	if err != nil {
		return err
	}

	_, err = os.Stat(params.FileName)

	if os.Getenv("DEV_MODE") == "true" && !params.Append {
		params.Overwrite = true
	}

	NotExists := os.IsNotExist(err)

	if NotExists || params.Append || params.Overwrite {

		if params.Overwrite && !params.Append {
			if !NotExists {
				if err := os.Remove(params.FileName); err != nil {
					return err
				}
			}
		}

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

func (t *Templates) CreateMany(data interface{}, files ...*FileCreateParams) error {
	for _, file := range files {
		if err := t.Create(&FileCreateParams{
			FileName:     file.FileName,
			TemplatePath: file.TemplatePath,
			Data:         data,
			Overwrite:    file.Overwrite,
		}); err != nil {
			return err
		}
	}
	return nil
}

func (f *Templates) CreateDir(name string) error {
	if _, err := os.Stat(name); os.IsNotExist(err) {
		err := os.MkdirAll(name, 0700)

		if err != nil {
			return err
		}
		return nil
	}
	fmt.Printf("Skipping: Dir %s already exists\n", name)
	return nil
}
