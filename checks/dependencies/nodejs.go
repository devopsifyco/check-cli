package dependencies

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// NodePackage represents the structure of a package.json file
type NodePackage struct {
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

// NodePackageLock represents the structure of a package-lock.json file
type NodePackageLock struct {
	LockfileVersion int `json:"lockfileVersion"`
	Dependencies map[string]struct {
		Version   string `json:"version"`
		Resolved  string `json:"resolved,omitempty"`
		Integrity string `json:"integrity,omitempty"`
		Dev       bool   `json:"dev,omitempty"`
		Requires  map[string]string `json:"requires,omitempty"`
	} `json:"dependencies"`
}

// parsePackageJSON parses a package.json file and extracts dependencies.
func parsePackageJSON(filePath string) ([]Dependency, error) {
	// First try to read package-lock.json if it exists
	lockFile := filepath.Join(filepath.Dir(filePath), "package-lock.json")
	if _, err := os.Stat(lockFile); err == nil {
		// package-lock.json exists, use it as primary source
		return parsePackageLock(lockFile)
	}

	// Fallback to package.json if package-lock.json doesn't exist
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading package.json %s: %w", filePath, err)
	}

	var pkg NodePackage
	if err := json.Unmarshal(data, &pkg); err != nil {
		return nil, fmt.Errorf("error parsing package.json %s: %w", filePath, err)
	}

	// Calculate total capacity for the dependencies slice
	totalDeps := len(pkg.Dependencies) + len(pkg.DevDependencies)
	deps := make([]Dependency, 0, totalDeps)

	// Process regular dependencies
	for name, version := range pkg.Dependencies {
		deps = append(deps, Dependency{
			Name:    name,
			Version: version,
			Manager: "npm",
		})
	}

	// Process dev dependencies
	for name, version := range pkg.DevDependencies {
		deps = append(deps, Dependency{
			Name:    name,
			Version: version,
			Manager: "npm",
			Tags:    []string{"dev"},
		})
	}

	return deps, nil
}

// parsePackageLock parses a package-lock.json file and returns the list of dependencies
func parsePackageLock(filePath string) ([]Dependency, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading package-lock.json %s: %w", filePath, err)
	}

	var lock NodePackageLock
	if err := json.Unmarshal(data, &lock); err != nil {
		return nil, fmt.Errorf("error parsing package-lock.json %s: %w", filePath, err)
	}

	// Create a slice to hold all dependencies
	deps := make([]Dependency, 0, len(lock.Dependencies))

	// Process dependencies from the dependencies section
	for name, dep := range lock.Dependencies {
		depInfo := Dependency{
			Name:    name,
			Version: dep.Version,
			Manager: "npm",
		}
		if dep.Dev {
			depInfo.Tags = []string{"dev"}
		}
		deps = append(deps, depInfo)
	}

	return deps, nil
} 