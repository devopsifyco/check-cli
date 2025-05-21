package dependencies

import (
	"fmt"
	"os"
	"path/filepath"
)

// ReadDependencies reads dependency information from supported manifest files in a given directory
// or from a specific dependency file.
func ReadDependencies(path string, outputFormat string) ([]Dependency, error) {
	// Check if path is a file or directory
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("error accessing path %s: %w", path, err)
	}

	// Map file names to their respective parser functions and manager type
	parsers := map[string]struct {
		parser  func(string) ([]Dependency, error)
		manager string
	}{
		"pom.xml":          {parsePomXML, "maven"},
		"package.json":     {parsePackageJSON, "npm"},
		"project.json":     {parseProjectJSON, "nuget"},
		"requirements.txt": {parsePythonDependencies, "pip"},
		"go.mod":          {parseGoMod, "go"},
		"packages.config":  {parsePackagesConfig, "nuget"},
	}

	// Add support for file extensions
	extParsers := map[string]struct {
		parser  func(string) ([]Dependency, error)
		manager string
	}{
		".csproj": {parseProjectJSON, "nuget"},
		".toml":   {parsePythonDependencies, "pip"},
	}

	allDependencies := []Dependency{}

	if !fileInfo.IsDir() {
		// If path is a file, try to parse it directly
		fileName := filepath.Base(path)
		ext := filepath.Ext(fileName)

		// First try exact filename match
		if parserInfo, ok := parsers[fileName]; ok {
			deps, err := parserInfo.parser(path)
			if err != nil {
				return nil, fmt.Errorf("error parsing file %s: %w", path, err)
			}
			// Add manager information
			for i := range deps {
				deps[i].Manager = parserInfo.manager
			}
			return deps, nil
		}

		// Then try extension match
		if parserInfo, ok := extParsers[ext]; ok {
			deps, err := parserInfo.parser(path)
			if err != nil {
				return nil, fmt.Errorf("error parsing file %s: %w", path, err)
			}
			// Add manager information
			for i := range deps {
				deps[i].Manager = parserInfo.manager
			}
			return deps, nil
		}

		return nil, fmt.Errorf("unsupported dependency file: %s", fileName)
	}

	// If path is a directory, check for all supported files
	if outputFormat != "json" && outputFormat != "yaml" {
		fmt.Printf("Checking dependencies in directory: %s\n", path)
	}

	// Check for files with exact names
	for fileName, parserInfo := range parsers {
		filePath := filepath.Join(path, fileName)
		if deps, err := tryParseFile(filePath, parserInfo.parser, parserInfo.manager, outputFormat); err == nil {
			allDependencies = append(allDependencies, deps...)
		}
	}

	// Check for files with matching extensions
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("error reading directory %s: %w", path, err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		ext := filepath.Ext(file.Name())
		if parserInfo, ok := extParsers[ext]; ok {
			filePath := filepath.Join(path, file.Name())
			if deps, err := tryParseFile(filePath, parserInfo.parser, parserInfo.manager, outputFormat); err == nil {
				allDependencies = append(allDependencies, deps...)
			}
		}
	}

	return allDependencies, nil
}

// tryParseFile attempts to parse a dependency file
func tryParseFile(filePath string, parser func(string) ([]Dependency, error), manager string, outputFormat string) ([]Dependency, error) {
	// Check if the file exists
	if _, err := os.Stat(filePath); err != nil {
		return nil, err
	}

	// File exists, parse it
	if outputFormat != "json" && outputFormat != "yaml" {
		fmt.Printf("Found %s, parsing...\n", filePath)
	}
	deps, err := parser(filePath)
	if err != nil {
		return nil, fmt.Errorf("error parsing file %s: %w", filePath, err)
	}

	// Add manager information to each dependency
	for i := range deps {
		deps[i].Manager = manager
	}

	return deps, nil
} 