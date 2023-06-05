package app

import (
	"errors"
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/x86-Yantras/code-gen/internal/adapters/templates"
	"github.com/x86-Yantras/code-gen/internal/constants"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type StorageModel struct {
	EntitySchema map[string]interface{}
	Entity       string
	EntityObj    string
	DBOperation  string
}

type StorageIndexModel struct {
	Imports     map[string]bool
	StorageType string
}

func (s *StorageModel) Validate() error {
	if s.DBOperation == "" {
		return errors.New("invalid spec file, operationId missing")
	}

	if s.EntitySchema == nil {
		return errors.New("invalid spec file, attributes missing")
	}
	return nil
}

func (app *App) CreateStorageAdapter() error {

	entities := map[string]map[string]bool{}

	for _, doc := range app.Spec.Paths {

		specList := []*openapi3.Operation{
			doc.Get,
			doc.Post,
			doc.Patch,
			doc.Put,
			doc.Delete,
		}

		subEntity := map[string]bool{}

		for _, spec := range specList {
			if spec != nil {
				serviceName := spec.Tags[0]

				if app.ServiceName != "" && app.ServiceName != serviceName {
					continue
				}

				if entities[serviceName] != nil {
					subEntity = entities[serviceName]
				}

				err := app.BuildStorageAdapter(spec, &app.Spec.Components)
				if err != nil {
					return err
				}

				subEntity[serviceName] = true
				entities[serviceName] = subEntity
			}
		}
	}
	err := app.BuildStorageIndex(entities)
	return err
}

func (app *App) BuildStorageAdapter(spec *openapi3.Operation, components *openapi3.Components) error {
	// Take first tag as service name
	serviceName := spec.Tags[0]
	titleCase := cases.Title(language.English)

	if err := app.BuildAdapterDir(serviceName, constants.Storage); err != nil {
		return err
	}

	adapterDir := fmt.Sprintf(app.Config.AdapterDir, serviceName)
	adapterTypeDir := fmt.Sprintf("%s/%s/%s", adapterDir, constants.Storage, constants.DefaultStorage)

	if err := app.Templater.CreateDir(adapterTypeDir); err != nil {
		return err
	}

	entities := spec.Tags

	for _, entity := range entities {
		schema := components.Schemas[titleCase.String(entity)]

		if schema == nil {
			fmt.Printf("schema not found for entity %s\n", entity)
			continue
		}

		payload := app.BuildSchema(schema.Value.Properties)

		storage := &StorageModel{
			Entity:       titleCase.String(serviceName),
			EntityObj:    serviceName,
			DBOperation:  spec.OperationID,
			EntitySchema: payload,
		}

		if err := app.BuildStorageFiles(serviceName, spec.OperationID, storage); err != nil {
			return err
		}
	}

	return nil
}

func (app *App) BuildSchema(spec openapi3.Schemas) map[string]interface{} {
	payload := map[string]interface{}{}

	for key, param := range spec {
		payload[key] = app.getTypeString(param.Value.Type)
	}

	return payload
}

func (app *App) getTypeString(inType string) interface{} {
	switch inType {
	case "string":
		return "String"
	case "number", "integer":
		return "Number"
	case "boolean":
		return "Boolean"
	case "array":
		return "[]"
	case "object":
		return "{}"
	default:
		return "{}"
	}
}

func (app *App) BuildStorageFiles(serviceName, operation string, storage *StorageModel) error {
	if err := storage.Validate(); err != nil {
		return err
	}

	adapterDir := fmt.Sprintf(app.Config.AdapterDir, serviceName)
	adapterTemplateDir := fmt.Sprintf(app.Config.AdapterDir, constants.ServiceDirPlaceholder)

	storageEntityFile := fmt.Sprintf("%s/%s/%s/%s%s", adapterDir, constants.Storage, constants.DefaultStorage, serviceName, app.Config.FileExt)

	storageEntityTemplate := fmt.Sprintf("%s/%s/%s/%s/%s%s%s", app.AppTemplateDir, adapterTemplateDir, constants.Storage, constants.DefaultStorage, constants.StorageEntityPlaceholder, app.Config.FileExt, constants.TemplateExtension)

	storageEntityOpTemplate := fmt.Sprintf("%s/%s/%s/%s/%s%s%s", app.AppTemplateDir, adapterTemplateDir, constants.Storage, constants.DefaultStorage, constants.StorageEntityOpPlaceholder, app.Config.FileExt, constants.TemplateExtension)

	files := []*templates.FileCreateParams{
		{FileName: storageEntityFile, TemplatePath: storageEntityTemplate},
	}

	err := app.Templater.CreateMany(storage, files...)

	if err != nil {
		return err
	}

	if err := app.Templater.Create(&templates.FileCreateParams{
		FileName:     storageEntityFile,
		TemplatePath: storageEntityOpTemplate,
		Append:       true,
		Data:         storage,
	}); err != nil {
		return err
	}

	return nil

}

func (app *App) BuildStorageIndex(storageImports map[string]map[string]bool) error {

	if len(storageImports) == 0 {
		return nil
	}

	for serviceName, storage := range storageImports {
		adapterDir := fmt.Sprintf(app.Config.AdapterDir, serviceName)
		adapterTemplateDir := fmt.Sprintf(app.Config.AdapterDir, constants.ServiceDirPlaceholder)

		storageIndexFile := fmt.Sprintf("%s/%s/%s", adapterDir, constants.Storage, app.Config.IndexFile)

		storageIndexTemplate := fmt.Sprintf("%s/%s/%s/%s%s", app.AppTemplateDir, adapterTemplateDir, constants.Storage, app.Config.IndexFile, constants.TemplateExtension)

		fmt.Println("overwriting storage index")
		err := app.Templater.Create(&templates.FileCreateParams{
			FileName:     storageIndexFile,
			TemplatePath: storageIndexTemplate,
			Data: &StorageIndexModel{
				StorageType: constants.DefaultStorage,
				Imports:     storage,
			},
			Overwrite: true,
		})

		return err
	}

	return nil
}
