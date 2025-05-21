package dependencies

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// parseGoMod parses a go.mod file and extracts dependencies.
func parseGoMod(filePath string) ([]Dependency, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening go.mod %s: %w", filePath, err)
	}
	defer file.Close()

	var deps []Dependency
	scanner := bufio.NewScanner(file)
	
	inRequireBlock := false
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		
		// Skip empty lines and module/go version declarations
		if line == "" || strings.HasPrefix(line, "module") || strings.HasPrefix(line, "go ") {
			continue
		}

		// Handle require blocks
		if line == "require (" {
			inRequireBlock = true
			continue
		}
		if line == ")" {
			inRequireBlock = false
			continue
		}

		// Handle single-line require statements
		if strings.HasPrefix(line, "require ") {
			dep := parseDependencyLine(strings.TrimPrefix(line, "require "))
			if dep != nil {
				deps = append(deps, *dep)
			}
			continue
		}

		// Handle dependencies inside require block
		if inRequireBlock && line != "" {
			dep := parseDependencyLine(line)
			if dep != nil {
				deps = append(deps, *dep)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading go.mod %s: %w", filePath, err)
	}

	return deps, nil
}

// parseDependencyLine parses a single dependency line from go.mod
func parseDependencyLine(line string) *Dependency {
	// Skip comments and empty lines
	line = strings.TrimSpace(line)
	if line == "" || strings.HasPrefix(line, "//") {
		return nil
	}

	// Split the line into parts
	parts := strings.Fields(line)
	if len(parts) < 2 {
		return nil
	}

	name := parts[0]
	version := parts[1]
	
	// Remove any trailing comments including "// indirect"
	if strings.Contains(version, "//") {
		version = strings.Fields(version)[0]
	}

	return &Dependency{
		Name:    name,
		Version: version,
		Manager: "go",
	}
} 