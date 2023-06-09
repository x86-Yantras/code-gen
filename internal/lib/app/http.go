package app

import (
	"errors"
	"fmt"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/x86-Yantras/code-gen/config"
	"github.com/x86-Yantras/code-gen/internal/adapters/templates"
	"github.com/x86-Yantras/code-gen/internal/constants"
)

// only one route per service,
// append each route for adapter.Build calls to this array
var routes *RoutesModel

type HandlerModel struct {
	HandlerName            string
	OperationID            string
	ServiceResponsePayload interface{}
	HttpStatusCode         string
	Path                   string
	Method                 string
}

// ServiceAdapter stores data required to build api adapters
type ServiceAdapter struct {
	spec        *openapi3.Operation
	fileExt     string
	testFileExt string
	templater   templates.TemplatesIface
	templateDir string
	adapterDir  string
	serviceDir  string
	serviceName string
	operation   string
	adapterType string
	path        string
	routesFile  string
}
type RoutesModel struct {
	HandlerImports []string
	Routes         []*RouteHandler
}
type RouteHandler struct {
	Method  string
	Handler string
	Path    string
}

type AdapterDependencies struct {
	Config      *config.Config
	Operation   string
	Spec        *openapi3.Operation
	Templater   templates.TemplatesIface
	AdapterType string
	Path        string
}

func NewServiceAdapter(deps *AdapterDependencies) (*ServiceAdapter, error) {
	spec := deps.Spec

	fmt.Println(deps.Config.AdapterDir)
	fmt.Println(fmt.Sprintf(deps.Config.AdapterDir, spec.Tags[0]))

	adapter := &ServiceAdapter{
		spec:        spec,
		templater:   deps.Templater,
		fileExt:     deps.Config.FileExt,
		testFileExt: deps.Config.TestFileExt,
		// adapterDir:  deps.Config.AdapterDir,
		adapterDir:  fmt.Sprintf(deps.Config.AdapterDir, spec.Tags[0]),
		templateDir: fmt.Sprintf("%s/%s", constants.TemplatesDir, deps.Config.Language),
		operation:   deps.Operation,
		serviceName: spec.Tags[0],
		adapterType: deps.AdapterType,
		path:        deps.Path,
		serviceDir:  deps.Config.ServiceDir,
		routesFile:  deps.Config.RoutesFile,
	}

	if routes == nil {
		routes = &RoutesModel{
			HandlerImports: make([]string, 0),
			Routes:         make([]*RouteHandler, 0),
		}
	}

	return adapter, nil
}

func (a *ServiceAdapter) Build() error {
	err := a.BuildHttpHandler()

	if err != nil {
		return err
	}

	a.AppendRoute()
	return nil
}

func (a *ServiceAdapter) BuildHttpHandler() error {
	// Take first tag as handler name
	spec := a.spec

	if err := a.BuildAdapterDir(); err != nil {
		return err
	}

	statusCode := "200"
	// find the success statusCode
	for resp := range spec.Responses {
		if strings.Contains(resp, "20") || strings.Contains(resp, "30") || strings.Contains(resp, "10") {
			statusCode = resp
		}
	}

	handler := &HandlerModel{
		OperationID:            spec.OperationID,
		HandlerName:            spec.OperationID,
		ServiceResponsePayload: nil,
		Path:                   a.formatPath(a.path),
		HttpStatusCode:         statusCode,
		Method:                 strings.ToLower(a.operation),
	}

	if err := a.BuildHandlerFile(handler); err != nil {
		return err
	}

	return nil
}

func (h *HandlerModel) Validate() error {
	if h.OperationID == "" || h.HandlerName == "" {
		return errors.New("invalid spec file, operationId missing")
	}
	return nil
}

func (a *ServiceAdapter) BuildAdapterDir() error {
	// create adapters dir

	if err := a.templater.CreateDir(a.adapterDir); err != nil {
		return err
	}

	adapterPath := fmt.Sprintf("%s/%s", a.adapterDir, a.adapterType)
	if err := a.templater.CreateDir(strings.ToLower(adapterPath)); err != nil {
		return err
	}

	return nil
}

func (a *ServiceAdapter) BuildHandlerFile(handler *HandlerModel) error {
	operationID := a.spec.OperationID
	if err := handler.Validate(); err != nil {
		return err
	}

	// BuildHandler
	handlerFile := fmt.Sprintf("%s/%s/%s", a.adapterDir, constants.APIHTTPAdapter, operationID)
	handlerFilePath := fmt.Sprintf("%s%s", handlerFile, a.fileExt)
	handlerTestPath := fmt.Sprintf("%s%s", handlerFile, a.testFileExt)

	handlerTemplate := fmt.Sprintf("%s/%s/%s/adapters/api/http/%s", a.templateDir, a.serviceDir, constants.ServiceDirPlaceholder, constants.HandlerPlaceHolder)
	handlerTemplatePath := fmt.Sprintf("%s%s%s", handlerTemplate, a.fileExt, constants.TemplateExtension)
	handlerTestTemplatePath := fmt.Sprintf("%s%s%s", handlerTemplate, a.testFileExt, constants.TemplateExtension)

	files := []*templates.FileCreateParams{
		{FileName: handlerFilePath, TemplatePath: handlerTemplatePath},
		{FileName: handlerTestPath, TemplatePath: handlerTestTemplatePath},
	}

	return a.templater.CreateMany(handler, files...)
}

func (a *ServiceAdapter) formatPath(path string) string {
	formattedPath := strings.Replace(path, "{", ":", 1)
	formattedPath = strings.Replace(formattedPath, "}", "", 1)
	return formattedPath
}

func (a *ServiceAdapter) AppendRoute() {
	handler := &RouteHandler{
		Method:  a.operation,
		Handler: a.spec.OperationID,
		Path:    a.path,
	}

	routes.Routes = append(routes.Routes, handler)
	routes.HandlerImports = append(routes.HandlerImports, handler.Handler)
}

func (a *ServiceAdapter) BuildRoutes() error {

	routesFile := fmt.Sprintf("%s/%s/%s", a.adapterDir, constants.APIHTTPAdapter, a.routesFile)
	routerTemplatePath := fmt.Sprintf("%s/%s/%s/adapters/api/http/%s%s", a.templateDir, a.serviceDir, constants.ServiceDirPlaceholder, a.routesFile, constants.TemplateExtension)

	err := a.templater.Create(&templates.FileCreateParams{
		FileName:     routesFile,
		TemplatePath: routerTemplatePath,
		Data:         routes,
		Overwrite:    true,
	})

	return err
}
