package templates

type TemplatesIface interface {
	Create(*FileCreateParams) error
	CreateMany(service interface{}, files ...*FileCreateParams) error
}
