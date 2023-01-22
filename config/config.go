package config

type Config struct {
	PackageManager string
	ServiceDir     string
	AdapterDir     string
	LibDir         string
	ReadmeFile     string
	FileExt        string
	TestFileExt    string
	RoutesFile     string
	IndexFile      string
	SchemaFile     string
	WorkspacePath  string
	PWD            string
	ProjectPath    string
	RuntimeVersion string
	Language       string
}

func New(configType string) *Config {
	switch configType {
	case "node":
		return NodeConfig()
	case "go":
		return GoConfig()
	}
	return nil
}
