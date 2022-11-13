package config

func NodeConfig() *Config {
	return &Config{
		ServiceDir:     "src/services",
		AdapterDir:     "src/services/%s/adapters",
		PackageManager: "package.json",
		ReadmeFile:     "Readme.md",
		FileExt:        ".js",
		TestFileExt:    ".spec.js",
		RoutesFile:     "routes.js",
		IndexFile:      "index.js",
		SchemaFile:     "schema.js",
	}
}
