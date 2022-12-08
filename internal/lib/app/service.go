package app

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/x86-Yantras/code-gen/internal/adapters/templates"
	"github.com/x86-Yantras/code-gen/internal/constants"
)

type ServiceModel struct {
	ServiceName        string
	ServiceDescription string
	Operation          string
	ServicePayload     interface{}
	Validation         map[string]string
}

func (app *App) CreateService() error {
	// Build Custom Service errors
	if app.Spec.Components.ExtensionProps.Extensions != nil && app.Spec.Components.ExtensionProps.Extensions[constants.Errors] != nil {
		err := app.BuildErrors(app.Spec.Components.ExtensionProps.Extensions[constants.Errors])
		if err != nil {
			return err
		}
	}

	// Build Service error codes
	if app.Spec.Components.ExtensionProps.Extensions != nil && app.Spec.Components.ExtensionProps.Extensions[constants.ErrCodes] != nil {
		err := app.BuildErrorCodes(app.Spec.Components.ExtensionProps.Extensions[constants.ErrCodes])
		if err != nil {
			return err
		}
	}

	serviceModels := map[string][]*ServiceModel{}
	servicesList := map[string]bool{}

	// Build Service for each path
	for _, doc := range app.Spec.Paths {
		specList := map[string]*openapi3.Operation{
			http.MethodGet:    doc.Get,
			http.MethodPost:   doc.Post,
			http.MethodPatch:  doc.Patch,
			http.MethodPut:    doc.Put,
			http.MethodDelete: doc.Delete,
		}

		models := []*ServiceModel{}

		for method, spec := range specList {
			if spec != nil {
				serviceName := spec.Tags[0]

				if app.ServiceName != "" && app.ServiceName != serviceName {
					continue
				}
				service, err := app.BuildService(spec, method)

				if err != nil {
					return err
				}

				if serviceModels[serviceName] != nil {
					models = serviceModels[serviceName]
				}

				models = append(models, service)

				servicesList[serviceName] = true
				serviceModels[serviceName] = models
			}
		}
	}

	err := app.BuildServiceTypes(serviceModels, servicesList)

	if err != nil {
		return err
	}

	return nil
}

func (app *App) BuildService(spec *openapi3.Operation, operationType string) (*ServiceModel, error) {
	// Take first tag as service name
	serviceName := spec.Tags[0]

	if err := app.BuildServiceDir(serviceName); err != nil {
		return nil, err
	}

	service := &ServiceModel{
		Operation:          spec.OperationID,
		ServiceDescription: spec.Description,
		ServiceName:        serviceName,
	}
	payload := app.BuildPayload(spec, operationType)

	requiredValidation := []string{}

	if spec.Parameters != nil {
		for _, params := range spec.Parameters {
			if params.Value.Required {
				requiredValidation = append(requiredValidation, params.Value.Name)
			}
		}
	}

	if spec.RequestBody != nil {
		requiredValidation = append(requiredValidation, spec.RequestBody.Value.Content[constants.ContentJson].Schema.Value.Required...)
	}

	validationMap := map[string]string{}

	for _, attrib := range requiredValidation {
		validationMap[attrib] = attrib
	}

	service.ServicePayload = payload
	service.Validation = validationMap

	if err := app.BuildServiceFiles(serviceName, spec.OperationID, service); err != nil {
		return nil, err
	}

	return service, nil
}

func (s *ServiceModel) Validate() error {
	if s.ServiceName == "" {
		return errors.New("invalid spec file, operationId missing")
	}
	return nil
}

func (app *App) BuildPayload(spec *openapi3.Operation, operationType string) map[string]interface{} {
	var payload map[string]interface{}

	switch operationType {
	case http.MethodGet:
		payload = map[string]interface{}{}
		// default limit and offset
		payload[constants.PayloadLimit] = 50
		payload[constants.PayloadOffset] = 0
		app.extractParametes(spec.Parameters, payload)
	case http.MethodPost:
		payload = map[string]interface{}{}
		app.extractParametes(spec.Parameters, payload)
		app.extractBody(spec.RequestBody, payload)
	case http.MethodPatch:
		payload = map[string]interface{}{}
		app.extractParametes(spec.Parameters, payload)
		app.extractBody(spec.RequestBody, payload)
	case http.MethodPut:
		payload = map[string]interface{}{}
		app.extractParametes(spec.Parameters, payload)
		app.extractBody(spec.RequestBody, payload)
	case http.MethodDelete:
		payload = map[string]interface{}{}
		app.extractParametes(spec.Parameters, payload)
		app.extractBody(spec.RequestBody, payload)
	default:
		payload = map[string]interface{}{}
		app.extractParametes(spec.Parameters, payload)
		app.extractBody(spec.RequestBody, payload)
	}
	return payload
}

func (app *App) BuildServiceDir(serviceName string) error {
	// create services dir
	serviceDir := app.Config.ServiceDir

	if err := app.CreateDir(serviceDir); err != nil {
		return err
	}

	servicePath := fmt.Sprintf("%s/%s/%s", serviceDir, serviceName, constants.Service)
	if err := app.CreateDir(strings.ToLower(servicePath)); err != nil {
		return err
	}

	return nil
}

func (app *App) BuildServiceFiles(serviceName, operation string, service *ServiceModel) error {
	serviceDir := app.Config.ServiceDir
	serviceFile := fmt.Sprintf("%s/%s/%s/%s", serviceDir, serviceName, constants.Service, operation)

	serviceFilePath := fmt.Sprintf("%s%s", serviceFile, app.Config.FileExt)
	serviceTestPath := fmt.Sprintf("%s%s", serviceFile, app.Config.TestFileExt)

	servicePath := fmt.Sprintf("%s/%s/%s/%s", app.AppTemplateDir, serviceDir, constants.ServiceDirPlaceholder, constants.Service)

	serviceTemplate := fmt.Sprintf("%s/%s", servicePath, constants.ServiceFilePlaceholder)

	serviceTemplatePath := fmt.Sprintf("%s%s%s", serviceTemplate, app.Config.FileExt, constants.TemplateExtension)
	serviceTestTemplatePath := fmt.Sprintf("%s%s%s", serviceTemplate, app.Config.TestFileExt, constants.TemplateExtension)

	files := []*templates.FileCreateParams{
		{
			FileName:     serviceFilePath,
			TemplatePath: serviceTemplatePath,
		},
		{
			FileName:     serviceTestPath,
			TemplatePath: serviceTestTemplatePath,
		},
	}

	return app.Templater.CreateMany(service, files...)
}

func (app *App) BuildServiceTypes(services map[string][]*ServiceModel, servicesList map[string]bool) error {

	servicePath := fmt.Sprintf("%s/%s/%s/%s", app.AppTemplateDir, app.Config.ServiceDir, constants.ServiceDirPlaceholder, constants.Service)
	serviceTypesTemplatePath := fmt.Sprintf("%s/%s%s%s", servicePath, constants.TypesFile, app.Config.FileExt, constants.TemplateExtension)
	serviceTypesValidatorTemplatePath := fmt.Sprintf("%s/%s%s%s", servicePath, constants.TypesValidatorFile, app.Config.FileExt, constants.TemplateExtension)

	// Create types first
	for serviceName := range servicesList {
		typesFile := fmt.Sprintf("%s/%s/%s/%s%s", app.Config.ServiceDir, serviceName, constants.Service, constants.TypesFile, app.Config.FileExt)
		err := app.Templater.Create(&templates.FileCreateParams{
			FileName:     typesFile,
			TemplatePath: serviceTypesTemplatePath,
			Data:         services[serviceName],
			Overwrite:    true,
		})

		if err != nil {
			return err
		}
	}

	// append validators then
	for serviceName := range servicesList {
		typesFile := fmt.Sprintf("%s/%s/%s/%s%s", app.Config.ServiceDir, serviceName, constants.Service, constants.TypesFile, app.Config.FileExt)
		err := app.Templater.Create(&templates.FileCreateParams{
			FileName:     typesFile,
			TemplatePath: serviceTypesValidatorTemplatePath,
			Data:         services[serviceName],
			Append:       true,
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func (app *App) extractParametes(pararms openapi3.Parameters, payloadObject map[string]interface{}) {
	for _, param := range pararms {
		payloadObject[param.Value.Name] = app.getDefaultType(param.Value.Schema.Value.Type)
	}
}

func (app *App) extractBody(body *openapi3.RequestBodyRef, payloadObject map[string]interface{}) {
	if body != nil {
		props := body.Value.Content[constants.ContentJson].Schema.Value.Properties
		for key, prop := range props {
			payloadObject[key] = app.getDefaultType(prop.Value.Type)
		}
	}
}

func (app *App) getDefaultType(inType string) interface{} {
	switch inType {
	case "string":
		return ""
	case "number":
		return 0
	case "boolean":
		return false
	case "array":
		return "[]"
	case "object":
		return "{}"
	default:
		return ""
	}
}
