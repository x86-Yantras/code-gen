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

type HandlerModel struct {
	HandlerName            string
	ServiceName            string
	ServiceResponsePayload interface{}
	HttpStatusCode         string
	Path                   string
	Method                 string
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

func (app *App) CreateHttpAdapter() error {
	serviceHandlers := map[string][]*RouteHandler{}

	for path, doc := range app.Spec.Paths {
		specList := map[string]*openapi3.Operation{
			http.MethodGet:    doc.Get,
			http.MethodPost:   doc.Post,
			http.MethodPatch:  doc.Patch,
			http.MethodPut:    doc.Put,
			http.MethodDelete: doc.Delete,
		}

		handlers := []*RouteHandler{}

		for method, spec := range specList {
			if spec != nil {
				if serviceHandlers[spec.Tags[0]] != nil {
					handlers = serviceHandlers[spec.Tags[0]]
				}

				err := app.BuildHttpHandler(spec, method, path, method)
				if err != nil {
					return err
				}

				formattedPath := strings.Replace(path, "{", ":", 1)
				formattedPath = strings.Replace(formattedPath, "}", "", 1)
				handlers = append(handlers, &RouteHandler{
					Method:  method,
					Handler: spec.OperationID,
					Path:    formattedPath,
				})
				serviceHandlers[spec.Tags[0]] = handlers
			}
		}
	}

	app.BuildRoutes(serviceHandlers)
	return nil
}

func (app *App) BuildHttpHandler(spec *openapi3.Operation, operationType, path, method string) error {
	// Take first tag as handler name
	serviceName := spec.Tags[0]

	if err := app.BuildAdapterDir(serviceName, constants.APIHTTPAdapter); err != nil {
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
		ServiceName:            spec.OperationID,
		HandlerName:            spec.OperationID,
		ServiceResponsePayload: nil,
		Path:                   path,
		HttpStatusCode:         statusCode,
		Method:                 strings.ToLower(method),
	}

	if err := app.BuildHandlerFile(serviceName, spec.OperationID, handler); err != nil {
		return err
	}

	return nil
}

func (h *HandlerModel) Validate() error {
	if h.ServiceName == "" || h.HandlerName == "" {
		return errors.New("invalid spec file, operationId missing")
	}
	return nil
}

func (app *App) BuildAdapterDir(serviceName, adapterType string) error {
	// create adapters dir
	adapterDir := fmt.Sprintf(app.Config.AdapterDir, serviceName)

	if err := app.CreateDir(adapterDir); err != nil {
		return err
	}

	adapterPath := fmt.Sprintf("%s/%s", adapterDir, adapterType)
	if err := app.CreateDir(strings.ToLower(adapterPath)); err != nil {
		return err
	}

	return nil
}

func (app *App) BuildHandlerFile(serviceName, operation string, handler *HandlerModel) error {

	if err := handler.Validate(); err != nil {
		return err
	}

	adapterDir := fmt.Sprintf(app.Config.AdapterDir, serviceName)

	// BuildHandler
	handlerFile := fmt.Sprintf("%s/%s/%s", adapterDir, constants.APIHTTPAdapter, operation)
	handlerFilePath := fmt.Sprintf("%s%s", handlerFile, app.Config.FileExt)
	handlerTestPath := fmt.Sprintf("%s%s", handlerFile, app.Config.TestFileExt)

	handlerTemplate := fmt.Sprintf("%s/%s/%s/adapters/api/http/%s", app.AppTemplateDir, app.Config.ServiceDir, constants.ServiceDirPlaceholder, constants.HandlerPlaceHolder)
	handlerTemplatePath := fmt.Sprintf("%s%s%s", handlerTemplate, app.Config.FileExt, constants.TemplateExtension)
	handlerTestTemplatePath := fmt.Sprintf("%s%s%s", handlerTemplate, app.Config.TestFileExt, constants.TemplateExtension)

	files := []*templates.FileCreateParams{
		{FileName: handlerFilePath, TemplatePath: handlerTemplatePath},
		{FileName: handlerTestPath, TemplatePath: handlerTestTemplatePath},
	}

	return app.Templater.CreateMany(handler, files...)
}

func (app *App) BuildRoutes(serviceHandlers map[string][]*RouteHandler) error {

	for serviceName, handlers := range serviceHandlers {
		handlerImports := []string{}
		routes := make([]*RouteHandler, 0, len(handlers))

		adapterDir := fmt.Sprintf(app.Config.AdapterDir, serviceName)

		for _, handler := range handlers {
			handlerImports = append(handlerImports, handler.Handler)
			routes = append(routes, &RouteHandler{
				Method:  strings.ToLower(handler.Method),
				Handler: handler.Handler,
				Path:    handler.Path,
			})
		}
		routesFile := fmt.Sprintf("%s/%s/%s", adapterDir, constants.APIHTTPAdapter, app.Config.RoutesFile)
		routerTemplatePath := fmt.Sprintf("%s/%s/%s/adapters/api/http/%s%s", app.AppTemplateDir, app.Config.ServiceDir, constants.ServiceDirPlaceholder, app.Config.RoutesFile, constants.TemplateExtension)

		err := app.Templater.Create(&templates.FileCreateParams{
			FileName:     routesFile,
			TemplatePath: routerTemplatePath,
			Data: &RoutesModel{
				HandlerImports: handlerImports,
				Routes:         routes,
			},
			Overwrite: true,
		})

		if err != nil {
			return err
		}
	}

	return nil
}
