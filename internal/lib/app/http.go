package app

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/x86-Yantras/code-gen/internal/constants"
)

type HandlerModel struct {
	HandlerName            string
	ServiceName            string
	ServiceResponsePayload interface{}
}

type RoutesModel struct {
	HandlerImports []string
	Routes         []*RouteHandler
}
type RouteHandler struct {
	Method  string
	Handler string
}

func (app *App) CreateHttpAdapter() error {
	for _, doc := range app.Spec.Paths {

		specList := map[string]*openapi3.Operation{
			http.MethodGet:    doc.Get,
			http.MethodPost:   doc.Post,
			http.MethodPatch:  doc.Patch,
			http.MethodPut:    doc.Put,
			http.MethodDelete: doc.Delete,
		}

		handlers := map[string]string{}
		var serviceName string

		for method, spec := range specList {
			if spec != nil {
				err := app.BuildHttpHandler(spec, method)
				if err != nil {
					return err
				}
				handlers[method] = spec.OperationID
				if serviceName == "" {
					serviceName = spec.Tags[0]
				}
			}
		}
		app.BuildRoutes(handlers, serviceName)
	}
	return nil
}

func (app *App) BuildHttpHandler(spec *openapi3.Operation, operationType string) error {
	// Take first tag as handler name
	serviceName := spec.Tags[0]

	if err := app.BuildAdapterDir(serviceName, constants.APIHTTPAdapter); err != nil {
		return err
	}

	handler := &HandlerModel{
		ServiceName:            spec.OperationID,
		HandlerName:            spec.OperationID,
		ServiceResponsePayload: nil,
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

	if err := app.Templater.Create(handlerFilePath, handlerTemplatePath, handler); err != nil {
		return err
	}

	if err := app.Templater.Create(handlerTestPath, handlerTestTemplatePath, handler); err != nil {
		return err
	}

	return nil
}

func (app *App) BuildRoutes(handlers map[string]string, serviceName string) error {
	handlerImports := []string{}
	routes := make([]*RouteHandler, 0, len(handlers))

	fmt.Println("running here")
	adapterDir := fmt.Sprintf(app.Config.AdapterDir, serviceName)

	for method, handler := range handlers {
		handlerImports = append(handlerImports, handler)
		routes = append(routes, &RouteHandler{
			Method:  method,
			Handler: handler,
		})
	}
	routesFile := fmt.Sprintf("%s/%s/%s", adapterDir, constants.APIHTTPAdapter, app.Config.RoutesFile)
	routerTemplatePath := fmt.Sprintf("%s/%s/%s/adapters/api/http/%s%s", app.AppTemplateDir, app.Config.ServiceDir, constants.ServiceDirPlaceholder, app.Config.RoutesFile, constants.TemplateExtension)

	fmt.Printf("%+v", handlers)
	fmt.Printf("%+v", routes)

	err := app.Templater.Create(routesFile, routerTemplatePath, &RoutesModel{
		HandlerImports: handlerImports,
		Routes:         routes,
	}, true)

	if err != nil {
		return err
	}

	return nil
}
