package app

import (
	"fmt"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/x86-Yantras/code-gen/config"
	"github.com/x86-Yantras/code-gen/internal/adapters/templates"
	"github.com/x86-Yantras/code-gen/internal/constants"
)

// ServiceInternal stores data required to build business logic layer
type ServiceInternal struct {
	spec        *openapi3.Operation
	fileExt     string
	testFileExt string
	templater   templates.TemplatesIface
	templateDir string
	serviceDir  string
	serviceName string
	operation   string
}

// Dependencies is used to inject dependencies when initializing ServiceInternal
type InternalDependencies struct {
	Config    *config.Config
	Operation string
	Spec      *openapi3.Operation
	Templater templates.TemplatesIface
}

// SericeModel is the struct used to build service internal files from template files
type ServiceModel struct {
	ServiceDir         string
	ServiceName        string
	ServiceDescription string
	Operation          string
	OperationMethod    string
	Payload            map[string]interface{}
	Validations        map[string]string
}

func NewServiceInternal(deps *InternalDependencies) (*ServiceInternal, error) {
	spec := deps.Spec

	internal := &ServiceInternal{
		spec:        spec,
		templater:   deps.Templater,
		fileExt:     deps.Config.FileExt,
		testFileExt: deps.Config.TestFileExt,
		serviceDir:  deps.Config.ServiceDir,
		templateDir: fmt.Sprintf("%s/%s", constants.TemplatesDir, deps.Config.Language),
		operation:   deps.Operation,
		serviceName: spec.Tags[0],
	}

	return internal, nil
}

func (s *ServiceInternal) Build() error {

	err := s.BuildServiceDir()

	if err != nil {
		return err
	}

	err = s.BuildServiceFiles()

	if err != nil {
		return err
	}

	return nil
}

func (s *ServiceInternal) BuildServiceDir() error {
	// create services dir
	serviceDir := s.serviceDir

	if err := s.templater.CreateDir(serviceDir); err != nil {
		return err
	}

	servicePath := fmt.Sprintf("%s/%s/%s", serviceDir, s.serviceName, constants.Service)
	if err := s.templater.CreateDir(strings.ToLower(servicePath)); err != nil {
		return err
	}

	return nil
}

func (s *ServiceInternal) BuildServiceFiles() error {
	serviceDir := s.serviceDir
	serviceFile := fmt.Sprintf("%s/%s/%s/%s", serviceDir, s.serviceName, constants.Service, s.spec.OperationID)

	serviceFilePath := fmt.Sprintf("%s%s", serviceFile, s.fileExt)
	serviceTestPath := fmt.Sprintf("%s%s", serviceFile, s.testFileExt)

	servicePath := fmt.Sprintf("%s/%s/%s/%s", s.templateDir, serviceDir, constants.ServiceDirPlaceholder, constants.Service)

	serviceTemplate := fmt.Sprintf("%s/%s", servicePath, constants.ServiceFilePlaceholder)

	serviceTemplatePath := fmt.Sprintf("%s%s%s", serviceTemplate, s.fileExt, constants.TemplateExtension)
	serviceTestTemplatePath := fmt.Sprintf("%s%s%s", serviceTemplate, s.testFileExt, constants.TemplateExtension)

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

	return s.templater.CreateMany(&ServiceModel{
		ServiceDir:         s.serviceDir,
		ServiceName:        s.spec.Tags[0],
		ServiceDescription: s.spec.Description,
		OperationMethod:    s.operation,
		Operation:          s.spec.OperationID,
	}, files...)
}
