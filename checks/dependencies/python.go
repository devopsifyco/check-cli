package dependencies

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

// PyProjectConfig represents the structure of pyproject.toml
type PyProjectConfig struct {
	Project struct {
		Dependencies []string `toml:"dependencies"`
		OptionalDependencies map[string][]string `toml:"optional-dependencies"`
	} `toml:"project"`
}

// parsePythonDependencies parses Python dependency files (requirements.txt or pyproject.toml)
func parsePythonDependencies(filePath string) ([]Dependency, error) {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".txt":
		return parseRequirementsTXT(filePath)
	case ".toml":
		return parsePyProjectTOML(filePath)
	default:
		return nil, fmt.Errorf("unsupported Python dependency file format: %s", ext)
	}
}

// parsePyProjectTOML parses dependencies from pyproject.toml
func parsePyProjectTOML(filePath string) ([]Dependency, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading pyproject.toml %s: %w", filePath, err)
	}

	var config PyProjectConfig
	if err := toml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error parsing pyproject.toml %s: %w", filePath, err)
	}

	var deps []Dependency

	// Parse main dependencies
	for _, dep := range config.Project.Dependencies {
		if parsedDep := parsePythonDependency(dep); parsedDep != nil {
			deps = append(deps, *parsedDep)
		}
	}

	// Parse optional dependencies
	for group, groupDeps := range config.Project.OptionalDependencies {
		for _, dep := range groupDeps {
			if parsedDep := parsePythonDependency(dep); parsedDep != nil {
				// Add the group as a tag to the dependency
				parsedDep.Tags = []string{fmt.Sprintf("optional:%s", group)}
				deps = append(deps, *parsedDep)
			}
		}
	}

	return deps, nil
}

// parseRequirementsTXT parses a requirements.txt file and extracts dependencies.
func parseRequirementsTXT(filePath string) ([]Dependency, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening requirements.txt %s: %w", filePath, err)
	}
	defer file.Close()

	var deps []Dependency
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Handle line continuations
		for strings.HasSuffix(line, "\\") {
			line = strings.TrimSuffix(line, "\\")
			if !scanner.Scan() {
				break
			}
			line += strings.TrimSpace(scanner.Text())
		}

		// Parse the dependency line
		dep := parsePythonDependency(line)
		if dep != nil {
			deps = append(deps, *dep)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading requirements.txt %s: %w", filePath, err)
	}

	return deps, nil
}

// parsePythonDependency parses a single dependency line
func parsePythonDependency(line string) *Dependency {
	// Skip empty lines
	line = strings.TrimSpace(line)
	if line == "" {
		return nil
	}

	// Remove quotes if present
	line = strings.Trim(line, `"'`)

	// Handle different version specifiers
	var name, version string

	// Split on common version specifiers
	for _, sep := range []string{"==", ">=", "<=", "~=", ">", "<", "!="} {
		if parts := strings.SplitN(line, sep, 2); len(parts) == 2 {
			name = strings.TrimSpace(parts[0])
			version = strings.TrimSpace(parts[1])
			// For non-exact versions, keep the operator in the version
			if sep != "==" {
				version = sep + version
			}
			break
		}
	}

	// If no version specifier found, check for egg/wheel format
	if version == "" {
		if parts := strings.Split(line, "@"); len(parts) > 1 {
			// Handle URLs with @ (like git+https://...)
			name = parts[0]
			version = strings.Join(parts[1:], "@")
		} else {
			// Just use the whole line as the name
			name = line
			version = "latest"
		}
	}

	// Remove any extras or environment markers
	if idx := strings.Index(name, "["); idx != -1 {
		name = name[:idx]
	}
	if idx := strings.Index(name, ";"); idx != -1 {
		name = name[:idx]
	}

	// Skip invalid lines
	if name == "" {
		return nil
	}

	// Only create the basic dependency without Tags
	return &Dependency{
		Name:    name,
		Version: version,
		Manager: "pip",
	}
} 