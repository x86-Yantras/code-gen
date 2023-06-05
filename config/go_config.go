package config

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

func GoConfig() *Config {

	goPath := os.Getenv("GOPATH")

	if goPath == "" {
		panic("GOPATH not found in the env")
	}
	workSpace := fmt.Sprintf("%s/%s", goPath, "src")

	pwd, err := os.Getwd()

	if err != nil {
		panic(err)
	}
	projectPath := strings.Split(pwd, workSpace)[1]
	projectPath = strings.Replace(projectPath, "/", "", 1)

	serviceDir := "internal/services"
	servicePath := fmt.Sprintf("%s/%s", projectPath, serviceDir)

	libDir := "pkg"
	libPath := fmt.Sprintf("%s/%s", projectPath, libDir)
	adapterDir := "internal/services/%s/adapters"
	adapterPath := fmt.Sprintf("%s/%s", projectPath, "internal/services/%s/adapters")

	return &Config{
		ServiceDir:     serviceDir,
		LibDir:         libDir,
		AdapterDir:     adapterDir,
		PackageManager: "go.mod",
		ReadmeFile:     "Readme.md",
		FileExt:        ".go",
		TestFileExt:    "_test.go",
		RoutesFile:     "routes.go",
		IndexFile:      "main.go",
		SchemaFile:     "types.go",
		ProjectPath:    projectPath,
		RuntimeVersion: runtime.Version(),
		Language:       "go",
		LibPath:        libPath,
		ServicePath:    servicePath,
		AdapterPath:    adapterPath,
	}
}
