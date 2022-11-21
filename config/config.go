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
}

func New(configType string) *Config {
	switch configType {
	case "node":
		return NodeConfig()
	}
	return &Config{}
}
