package checks

import (
	"fmt"
	"os"

	"github.com/opsify/check/checks/dependencies"
)

// DepsCheckCommand handles the execution of the dependencies check.
type DepsCheckCommand struct {
	*BaseCheckCommand
	outputFormat string
}

// CleanDependency represents a dependency without empty fields
type CleanDependency struct {
	Name    string   `json:"name" yaml:"name"`
	Version string   `json:"version" yaml:"version"`
	Manager string   `json:"manager" yaml:"manager"`
	Tags    []string `json:"tags,omitempty" yaml:"tags,omitempty"`
}

// DepsResult holds the outcome of the dependencies check.
type DepsResult struct {
	Directory    string                  `json:"directory,omitempty"`
	Dependencies []dependencies.Dependency `json:"-" yaml:"-"` // Hide the original dependencies
	Error        string                  `json:"error,omitempty"`
}

// cleanDependencies converts dependencies to clean format
func (r *DepsResult) cleanDependencies() []CleanDependency {
	clean := make([]CleanDependency, len(r.Dependencies))
	for i, dep := range r.Dependencies {
		clean[i] = CleanDependency{
			Name:    dep.Name,
			Version: dep.Version,
			Manager: dep.Manager,
		}
		// Only include tags if they exist
		if len(dep.Tags) > 0 {
			clean[i].Tags = dep.Tags
		}
	}
	return clean
}

// Print implements CheckResult interface
func (r *DepsResult) Print(outputFormat string) {
	switch outputFormat {
	case "json", "yaml":
		// Create clean output structure for structured formats
		cleanResult := struct {
			Dependencies []CleanDependency `json:"dependencies" yaml:"dependencies"`
		}{
			Dependencies: r.cleanDependencies(),
		}
		if outputFormat == "json" {
			PrintJSON(cleanResult)
		} else {
			PrintYAML(cleanResult)
		}
	default:
		fmt.Printf("--- Dependencies for %s ---\n", r.Directory)
		if r.Error != "" {
			fmt.Printf("Error: %s\n", r.Error)
			return
		}
		if len(r.Dependencies) == 0 {
			fmt.Println("No dependencies found. Supported files:")
			fmt.Println("  - pom.xml (Maven)")
			fmt.Println("  - package.json (Node.js)")
			fmt.Println("  - project.json (DotNet)")
			fmt.Println("  - requirements.txt (Python)")
			fmt.Println("  - go.mod (Go)")
		} else {
			for _, dep := range r.Dependencies {
				fmt.Printf("- %s (%s) [%s]\n", dep.Name, dep.Version, dep.Manager)
			}
		}
	}
}

// NewDepsCheckCommand creates a new command for checking dependencies.
func NewDepsCheckCommand(outputFormat string) *DepsCheckCommand {
	return &DepsCheckCommand{
		BaseCheckCommand: NewBaseCheckCommand(
			"deps",
			"Check project dependencies",
			"deps [path]",
			0, // 0 required args as path is optional
		),
		outputFormat: outputFormat,
	}
}

// Execute implements the CheckCommand interface
func (c *DepsCheckCommand) Execute(args []string) (CheckResult, error) {
	var targetPath string
	if len(args) > 0 {
		targetPath = args[0]
	} else {
		// Default to current directory if no argument is provided
		wd, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("failed to get current working directory: %w", err)
		}
		targetPath = wd
	}

	deps, err := dependencies.ReadDependencies(targetPath, c.outputFormat)
	result := &DepsResult{
		Dependencies: deps,
	}
	
	if c.outputFormat != "json" && c.outputFormat != "yaml" {
		result.Directory = targetPath
	}

	if err != nil {
		result.Error = err.Error()
		return result, err
	}

	return result, nil
} 