package templates

type TemplatesIface interface {
	Create(templateName, templatePath string, data interface{}, overwrite ...bool) error
}
