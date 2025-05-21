package dependencies

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// NuGetConfig represents the structure of a NuGet.config file
type NuGetConfig struct {
	XMLName         xml.Name `xml:"configuration"`
	PackageSources  struct {
		Clear bool `xml:"clear"`
		Sources []struct {
			Key   string `xml:"key,attr"`
			Value string `xml:"value,attr"`
		} `xml:"add"`
	} `xml:"packageSources"`
	PackageSourceCredentials struct {
		Sources []struct {
			Name     string `xml:",chardata"`
			Username struct {
				Value string `xml:"value,attr"`
			} `xml:"add"`
			Password struct {
				Value string `xml:"value,attr"`
			} `xml:"add"`
		} `xml:",any"`
	} `xml:"packageSourceCredentials"`
	PackageSourceMapping struct {
		Sources []struct {
			Key      string `xml:"key,attr"`
			Packages []struct {
				Pattern string `xml:"pattern,attr"`
			} `xml:"package"`
		} `xml:"packageSource"`
	} `xml:"packageSourceMapping"`
}

// parseNuGetConfig parses a NuGet.config file and extracts package sources
func parseNuGetConfig(filePath string) ([]Dependency, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading NuGet.config %s: %w", filePath, err)
	}

	var config NuGetConfig
	if err := xml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error parsing NuGet.config %s: %w", filePath, err)
	}

	deps := make([]Dependency, 0)

	// Add package sources as dependencies
	for _, source := range config.PackageSources.Sources {
		deps = append(deps, Dependency{
			Name:    source.Key,
			Version: source.Value,
			Manager: "nuget-source",
			Tags:    []string{"package-source"},
		})
	}

	return deps, nil
}

// ProjectJSON represents the structure of a project.json file
type ProjectJSON struct {
	Dependencies map[string]string `json:"dependencies"`
	Frameworks   map[string]struct {
		Dependencies map[string]struct {
			Version string `json:"version"`
		} `json:"dependencies"`
	} `json:"frameworks"`
}

// CSProject represents the structure of a .csproj file
type CSProject struct {
	XMLName xml.Name `xml:"Project"`
	ItemGroups []struct {
		PackageReferences []struct {
			Include string `xml:"Include,attr"`
			Version string `xml:"Version,attr"`
		} `xml:"PackageReference"`
	} `xml:"ItemGroup"`
}

// parseProjectJSONFile parses a project.json file and extracts dependencies
func parseProjectJSONFile(filePath string) ([]Dependency, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading project.json %s: %w", filePath, err)
	}

	var proj ProjectJSON
	if err := json.Unmarshal(data, &proj); err != nil {
		return nil, fmt.Errorf("error parsing project.json %s: %w", filePath, err)
	}

	deps := make([]Dependency, 0)

	// Process top-level dependencies
	for name, version := range proj.Dependencies {
		deps = append(deps, Dependency{
			Name:    name,
			Version: version,
			Manager: "nuget",
		})
	}

	// Process framework-specific dependencies
	for _, framework := range proj.Frameworks {
		for name, dep := range framework.Dependencies {
			deps = append(deps, Dependency{
				Name:    name,
				Version: dep.Version,
				Manager: "nuget",
			})
		}
	}

	return deps, nil
}

// parseCSProj parses a .csproj file and extracts dependencies
func parseCSProj(filePath string) ([]Dependency, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading .csproj %s: %w", filePath, err)
	}

	var proj CSProject
	if err := xml.Unmarshal(data, &proj); err != nil {
		return nil, fmt.Errorf("error parsing .csproj %s: %w", filePath, err)
	}

	deps := make([]Dependency, 0)

	// Process all PackageReference items
	for _, group := range proj.ItemGroups {
		for _, pkg := range group.PackageReferences {
			deps = append(deps, Dependency{
				Name:    pkg.Include,
				Version: pkg.Version,
				Manager: "nuget",
			})
		}
	}

	return deps, nil
}

// parseProjectJSON parses a project.json or .csproj file and extracts dependencies.
func parseProjectJSON(filePath string) ([]Dependency, error) {
	ext := strings.ToLower(filepath.Ext(filePath))
	base := strings.ToLower(filepath.Base(filePath))
	
	switch {
	case ext == ".json":
		return parseProjectJSONFile(filePath)
	case ext == ".csproj":
		return parseCSProj(filePath)
	case base == "nuget.config":
		return parseNuGetConfig(filePath)
	default:
		// Try to detect file type by reading first few bytes
		data, err := os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("error reading file %s: %w", filePath, err)
		}

		// Check if it's XML (csproj or nuget.config)
		if len(data) > 0 && data[0] == '<' {
			// Try to determine if it's a NuGet.config file
			if strings.Contains(string(data), "<configuration") && strings.Contains(string(data), "packageSources") {
				return parseNuGetConfig(filePath)
			}
			return parseCSProj(filePath)
		}

		// Check if it's JSON (project.json)
		if len(data) > 0 && (data[0] == '{' || data[0] == '[') {
			return parseProjectJSONFile(filePath)
		}

		return nil, fmt.Errorf("unsupported project file format: %s", filePath)
	}
}

// PackagesConfig represents the structure of a packages.config file
type PackagesConfig struct {
	XMLName  xml.Name         `xml:"packages"`
	Packages []PackageConfig `xml:"package"`
}

// PackageConfig represents a single package in packages.config
type PackageConfig struct {
	ID                   string `xml:"id,attr"`
	Version              string `xml:"version,attr"`
	TargetFramework      string `xml:"targetFramework,attr"`
	DevelopmentDependency bool   `xml:"developmentDependency,attr"`
}

// parsePackagesConfig parses a packages.config file and extracts dependencies
func parsePackagesConfig(filePath string) ([]Dependency, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading packages.config %s: %w", filePath, err)
	}

	var config PackagesConfig
	if err := xml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error parsing packages.config %s: %w", filePath, err)
	}

	deps := make([]Dependency, 0, len(config.Packages))
	for _, pkg := range config.Packages {
		tags := []string{fmt.Sprintf("framework:%s", pkg.TargetFramework)}
		if pkg.DevelopmentDependency {
			tags = append(tags, "development")
		}
		
		deps = append(deps, Dependency{
			Name:    pkg.ID,
			Version: pkg.Version,
			Manager: "nuget",
			Tags:    tags,
		})
	}

	return deps, nil
} 