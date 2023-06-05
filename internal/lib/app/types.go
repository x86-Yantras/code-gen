package app

import (
	"fmt"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/x86-Yantras/code-gen/config"
	"github.com/x86-Yantras/code-gen/internal/adapters/templates"
	"github.com/x86-Yantras/code-gen/internal/constants"
)

var typeCreated = false

type TypeDependencies struct {
	Config    *config.Config
	Operation string
	Spec      *openapi3.Operation
	Templater templates.TemplatesIface
}

type ServiceType struct {
	spec        *openapi3.Operation
	fileExt     string
	testFileExt string
	templater   templates.TemplatesIface
	templateDir string
	serviceDir  string
	serviceName string
	operation   string
}

type TypeModel struct {
	Operation      string
	Validation     map[string]string
	ServicePayload map[string]interface{}
}

func NewServiceTypes(deps *TypeDependencies) (*ServiceType, error) {
	return &ServiceType{
		spec:        deps.Spec,
		fileExt:     deps.Config.FileExt,
		testFileExt: deps.Config.TestFileExt,
		serviceDir:  deps.Config.ServiceDir,
		templateDir: fmt.Sprintf("%s/%s", constants.TemplatesDir, deps.Config.Language),
		operation:   deps.Operation,
		templater:   deps.Templater,
		serviceName: deps.Spec.Tags[0],
	}, nil
}

func (t *ServiceType) Build() error {

	payload := t.BuildPayload()
	validation := t.BuildValidation()

	servicePath := fmt.Sprintf("%s/%s/%s/%s", t.templateDir, t.serviceDir, constants.ServiceDirPlaceholder, constants.Service)
	serviceTypesTemplatePath := fmt.Sprintf("%s/%s%s%s", servicePath, constants.TypesFile, t.fileExt, constants.TemplateExtension)
	serviceTypesValidatorTemplatePath := fmt.Sprintf("%s/%s%s%s", servicePath, constants.TypesValidatorFile, t.fileExt, constants.TemplateExtension)
	typesFile := fmt.Sprintf("%s/%s/%s/%s%s", t.serviceDir, t.serviceName, constants.Service, constants.TypesFile, t.fileExt)

	// Create types import first
	if !typeCreated {
		serviceTypesImportTemplatePath := fmt.Sprintf("%s/%s%s%s", servicePath, constants.TypesImportFile, t.fileExt, constants.TemplateExtension)
		err := t.templater.Create(&templates.FileCreateParams{
			FileName:     typesFile,
			TemplatePath: serviceTypesImportTemplatePath,
			Overwrite:    true,
		})

		if err != nil {
			return err
		}
		typeCreated = true
	}

	// Create types first
	err := t.templater.Create(&templates.FileCreateParams{
		FileName:     typesFile,
		TemplatePath: serviceTypesTemplatePath,
		Data: &TypeModel{
			Operation:      t.spec.OperationID,
			ServicePayload: payload,
		},
		Append: true,
	})

	if err != nil {
		return err
	}

	// append validators then
	err = t.templater.Create(&templates.FileCreateParams{
		FileName:     typesFile,
		TemplatePath: serviceTypesValidatorTemplatePath,
		Data: &TypeModel{
			Operation:      t.spec.OperationID,
			Validation:     validation,
			ServicePayload: payload,
		},
		Append: true,
	})

	if err != nil {
		return err
	}

	return nil
}

func (t *ServiceType) extractParametes(pararms openapi3.Parameters, payload map[string]interface{}) map[string]interface{} {
	for _, param := range pararms {
		payload[param.Value.Name] = t.getDefaultType(param.Value.Schema.Value.Type)
	}
	return payload
}

func (t *ServiceType) extractBody(body *openapi3.RequestBodyRef, payload map[string]interface{}) map[string]interface{} {
	if body != nil {
		props := body.Value.Content[constants.ContentJson].Schema.Value.Properties
		for key, prop := range props {
			payload[key] = t.getDefaultType(prop.Value.Type)
		}
	}

	return payload
}

func (t *ServiceType) BuildPayload() map[string]interface{} {
	spec := t.spec

	payload := map[string]interface{}{}

	t.extractParametes(spec.Parameters, payload)

	switch t.operation {
	case http.MethodGet:
		// default limit and offset
		payload[constants.PayloadLimit] = 50
		payload[constants.PayloadOffset] = 0
	default:
		t.extractBody(spec.RequestBody, payload)
	}
	return payload
}

func (t *ServiceType) BuildValidation() map[string]string {
	spec := t.spec
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
	return validationMap
}

func (t *ServiceType) getDefaultType(inType string) interface{} {
	switch inType {
	case "string":
		return "\"\""
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
