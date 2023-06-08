package app

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/x86-Yantras/code-gen/config"
	"github.com/x86-Yantras/code-gen/internal/adapters/templates"
	"github.com/x86-Yantras/code-gen/mocks"
)

const errNotNil = "err not nil"
const coreNil = "core is nil"

func setupCore(t *testing.T, deps *Dependencies) (*Core, *mocks.TemplatesIface) {
	templaterMock := mocks.NewTemplatesIface(t)
	core, err := NewCore(&Dependencies{
		Spec:      deps.Spec,
		Service:   deps.Service,
		Config:    deps.Config,
		Templater: templaterMock,
	})

	assert.NotNil(t, core, coreNil)
	assert.Nil(t, err, errNotNil)

	return core, templaterMock
}

func TestCore(t *testing.T) {
	// cfg := config.NodeConfig()
	serviceName := "it"

	t.Run("should initialize core", func(t *testing.T) {
		setupCore(t, &Dependencies{
			Spec:      &openapi3.T{},
			Service:   serviceName,
			Config:    config.NodeConfig(),
			Templater: &templates.Templates{},
		})
	})

	t.Run("should return err spec is nil", func(t *testing.T) {
		core, err := NewCore(&Dependencies{
			Spec:      nil,
			Service:   serviceName,
			Config:    config.NodeConfig(),
			Templater: &templates.Templates{},
		})

		assert.Nil(t, core)
		assert.NotNil(t, err)
	})

	t.Run("should build services", func(t *testing.T) {
		operationIDs := []string{"GetTest", "PutTest", "PatchTest", "DeleteTest", "PostTest"}
		// operationMethods := []string{"GET", "PUT", "PATCH", "DELERE", "POST"}

		core, templater := setupCore(t, &Dependencies{
			Spec: &openapi3.T{
				Paths: openapi3.Paths{
					"/test": &openapi3.PathItem{
						Get: &openapi3.Operation{
							Tags:        []string{serviceName},
							OperationID: operationIDs[0],
						},
						Patch: &openapi3.Operation{
							Tags:        []string{serviceName},
							OperationID: operationIDs[1],
						},
						Put: &openapi3.Operation{
							Tags:        []string{serviceName},
							OperationID: operationIDs[2],
						},
						Delete: &openapi3.Operation{
							Tags:        []string{serviceName},
							OperationID: operationIDs[3],
						},
						Post: &openapi3.Operation{
							Tags:        []string{serviceName},
							OperationID: operationIDs[4],
						},
					},
				},
			},
			Service:   serviceName,
			Config:    config.NodeConfig(),
			Templater: &templates.Templates{},
		})

		templater.On("CreateDir", "src/services/it/service").Return(nil).Times(5)
		templater.On("CreateDir", "src/services").Return(nil).Times(5)

		// mock anything of type because order is not preserved
		templater.On("CreateMany", mock.AnythingOfType("*app.ServiceModel"), mock.AnythingOfType("*templates.FileCreateParams"), mock.AnythingOfType("*templates.FileCreateParams")).Times(5).Return(nil)

		// Type Create
		templater.On("Create", &templates.FileCreateParams{
			FileName:     "src/services/it/service/types.js",
			TemplatePath: "templates/node/src/services/servicename/service/types_import.js.tmpl",
			Overwrite:    true,
		}).Return(nil).Times(1)

		// Type schema append
		templater.On("Create", mock.AnythingOfType("*templates.FileCreateParams")).Return(nil).Times(5)

		// Type validation append
		templater.On("Create", mock.AnythingOfType("*templates.FileCreateParams")).Return(nil).Times(5)

		err := core.CreateService()

		assert.Nil(t, err, errNotNil)
		templater.AssertExpectations(t)
	})
}
