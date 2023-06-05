package app

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/x86-Yantras/code-gen/config"
	"github.com/x86-Yantras/code-gen/internal/adapters/templates"
	"github.com/x86-Yantras/code-gen/internal/constants"
)

// Core stores all the data required to build project parts
type Core struct {
	AppName        string
	AppDescription string
	services       []*Service
	errors         *Errors
}

// Service stores data required to build a single service
type Service struct {
	internals       []*ServiceInternal
	types           []*ServiceType
	apiAdapters     []*ApiHandler
	storageAdapters []*Storage
}

// ApiHandler stores data required to build api handlers
type ApiHandler struct {
	Name       string
	StatusCode string
	Method     string
	Path       string
}

// Storage stores data required to build storage adapter
type Storage struct {
	EntitySchema map[string]interface{}
	Entity       string
	EntityObj    string
	Operation    string
	StorageType  string
}

type Dependencies struct {
	Spec      *openapi3.T
	Service   string
	Config    *config.Config
	Templater templates.TemplatesIface
}

// NewCore takes Dependencies to build and returns an instance of core
func NewCore(deps *Dependencies) (*Core, error) {
	if deps.Spec == nil {
		return nil, errors.New("empty spec")
	}
	service := deps.Service
	spec := deps.Spec
	config := deps.Config
	templater := deps.Templater

	if spec == nil {
		return nil, errors.New(constants.ErrorEmptySpec)
	}

	servicesMap := map[string]*Service{}

	for _, doc := range spec.Paths {
		specList := map[string]*openapi3.Operation{
			http.MethodGet:    doc.Get,
			http.MethodPost:   doc.Post,
			http.MethodPatch:  doc.Patch,
			http.MethodPut:    doc.Put,
			http.MethodDelete: doc.Delete,
		}

		for method, operation := range specList {
			if operation != nil {
				serviceName := operation.Tags[0]

				// skip services(tags) that do no match service input, if service input
				if service != "" && service != serviceName {
					continue
				}

				serviceInternal, err := NewServiceInternal(&InternalDependencies{
					Config:    config,
					Operation: method,
					Spec:      operation,
					Templater: templater,
				})

				if err != nil {
					return nil, err
				}

				serviceTypes, err := NewServiceTypes(&TypeDependencies{
					Config:    config,
					Operation: method,
					Spec:      operation,
					Templater: templater,
				})

				if err != nil {
					return nil, err
				}

				if servicesMap[serviceName] == nil {
					servicesMap[serviceName] = &Service{
						internals:       make([]*ServiceInternal, 0, 5),
						apiAdapters:     make([]*ApiHandler, 0),
						types:           make([]*ServiceType, 0),
						storageAdapters: make([]*Storage, 0),
					}
				}

				servicesMap[serviceName].internals = append(servicesMap[serviceName].internals, serviceInternal)
				servicesMap[serviceName].types = append(servicesMap[serviceName].types, serviceTypes)
			}
		}

	}

	services := make([]*Service, 0, len(servicesMap))
	for _, service := range servicesMap {
		services = append(services, service)
	}

	errors, err := NewErrors(&ErrDependencies{
		Config:    config,
		Spec:      spec,
		Templater: templater,
	})

	if err != nil {
		return nil, err
	}

	core := &Core{
		services: services,
		errors:   errors,
	}

	if spec.Info != nil {
		core.AppName = spec.Info.Title
		core.AppDescription = spec.Info.Description
	}

	return core, nil
}

func (c *Core) CreateService() error {
	if len(c.services) == 0 || c.services == nil {
		return errors.New("no services to build")
	}

	for _, service := range c.services {
		if len(service.internals) == 0 || service.internals == nil {
			fmt.Println("no internals")
			break
		}

		for _, internal := range service.internals {
			err := internal.Build()
			if err != nil {
				return err
			}
		}

		for _, types := range service.types {
			err := types.Build()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *Core) CreateErrors() error {
	if c.errors == nil {
		return errors.New("no errors to build")
	}

	if err := c.errors.BuildErrors(); err != nil {
		return err
	}

	if err := c.errors.BuildErrorCodes(); err != nil {
		return err
	}

	return nil
}
