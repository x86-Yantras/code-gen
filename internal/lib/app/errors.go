package app

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/x86-Yantras/code-gen/config"
	"github.com/x86-Yantras/code-gen/internal/adapters/templates"
	"github.com/x86-Yantras/code-gen/internal/constants"
)

type Errors struct {
	errorDir    string
	fileExt     string
	templateDir string
	templater   templates.TemplatesIface
	Errors      []*Error
	ErrCodes    []*ErrCode
}

// Error stores data to build custom errors
type Error struct {
	Name           string
	HttpStatusCode int `json:"httpStatusCode"`
}

// ErrCode stores data to build custom error codes
type ErrCode struct {
	Name string
	Code string
}

type ErrDependencies struct {
	Config    *config.Config
	Spec      *openapi3.T
	Templater templates.TemplatesIface
}

func NewErrors(deps *ErrDependencies) (*Errors, error) {
	extensions := deps.Spec.Components.ExtensionProps.Extensions
	errors := []*Error{}
	errCodes := []*ErrCode{}

	if extensions != nil && extensions[constants.Errors] != nil {
		errs := extensions[constants.Errors]

		if err := json.Unmarshal(errs.(json.RawMessage), &errors); err != nil {
			return nil, err
		}
	}

	if extensions != nil && extensions[constants.ErrCodes] != nil {
		errs := extensions[constants.ErrCodes]

		if err := json.Unmarshal(errs.(json.RawMessage), &errCodes); err != nil {
			return nil, err
		}
	}

	return &Errors{
		templater:   deps.Templater,
		Errors:      errors,
		ErrCodes:    errCodes,
		fileExt:     deps.Config.FileExt,
		errorDir:    deps.Config.LibDir,
		templateDir: fmt.Sprintf("%s/%s", constants.TemplatesDir, deps.Config.Language),
	}, nil
}

func (e *Errors) BuildErrors() error {

	// Build errors dir
	errDir := fmt.Sprintf("%s/%s", e.errorDir, constants.Errors)

	if err := e.templater.CreateDir(errDir); err != nil {
		return err
	}

	// Build custom errors
	fmt.Println("building errors")
	for _, err := range e.Errors {
		errFile := fmt.Sprintf("%s/%s/%s%s", e.errorDir, constants.Errors, ToLowerFirst(err.Name), e.fileExt)
		tplPath := fmt.Sprintf("%s/%s/%s/%s%s%s", e.templateDir, e.errorDir, constants.Errors, constants.Errors, e.fileExt, constants.TemplateExtension)

		err := e.templater.Create(&templates.FileCreateParams{
			FileName:     errFile,
			TemplatePath: tplPath,
			Data:         err,
			Overwrite:    true,
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func (e *Errors) BuildErrorCodes() error {

	// Build custom errors
	errFile := fmt.Sprintf("%s/%s/%s%s", e.errorDir, constants.Errors, constants.ErrCodes, e.fileExt)
	tplPath := fmt.Sprintf("%s/%s/%s/%s%s%s", e.templateDir, e.errorDir, constants.Errors, constants.ErrCodes, e.fileExt, constants.TemplateExtension)

	fmt.Println("building err codes")
	err := e.templater.Create(&templates.FileCreateParams{
		FileName:     errFile,
		TemplatePath: tplPath,
		Data:         e.ErrCodes,
		Overwrite:    true,
	})

	return err
}

func ToLowerFirst(input string) string {
	return strings.ToLower(string(input[0])) + input[1:]
}

func (app *App) ToLowerFirst(input string) string {
	return strings.ToLower(string(input[0])) + input[1:]
}
