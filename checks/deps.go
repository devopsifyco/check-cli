package checks

import (
	"fmt"
	"os"

	"github.com/devopsifyco/check-cli/checks/dependencies"
	"github.com/devopsifyco/check-cli/checks/utilities/output"
)

// DepsCheckCommand handles the execution of the dependencies check.
type DepsCheckCommand struct {
	*BaseCheckCommand
	outputFormat string
	cve          bool // Add cve flag
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
	CVEs         map[string][]CVEResponse `json:"cves,omitempty" yaml:"cves,omitempty"` // Map of dep name@version to CVEs
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
	case "json":
		// Output only the array of dependencies (with CVEs as before, if needed)
		output.PrintJSON(r.cleanDependencies())
	case "yaml":
		// --- Custom YAML output to match requested format, no top-level 'dependencies' key ---
		type YAMLCVE struct {
			CVEID         string    `yaml:"cveid"`
			State         string    `yaml:"state"`
			PublishedDate string    `yaml:"publisheddate"`
			Score         *float64  `yaml:"score"`
			Title         string    `yaml:"title"`
			References    []string  `yaml:"references"`
		}
		type YAMLDependency struct {
			Name     string     `yaml:"name"`
			Version  string     `yaml:"version"`
			Manager  string     `yaml:"manager"`
			Tags     []string   `yaml:"tags,omitempty"`
			CVEs     []YAMLCVE  `yaml:"cves,omitempty"`
		}
	
		yamlDeps := make([]YAMLDependency, 0, len(r.Dependencies))
		for _, dep := range r.Dependencies {
			cves := []YAMLCVE{}
			if r.CVEs != nil {
				key := dep.Name + "@" + dep.Version
				if depCVEs, ok := r.CVEs[key]; ok {
					for _, cve := range depCVEs {
						cves = append(cves, YAMLCVE{
							CVEID:         cve.CVEID,
							State:         cve.State,
							PublishedDate: cve.PublishedDate,
							Score:         cve.Score,
							Title:         cve.Title,
							References:    cve.References,
						})
					}
				}
			}
			depYaml := YAMLDependency{
				Name:    dep.Name,
				Version: dep.Version,
				Manager: dep.Manager,
				Tags:    dep.Tags,
				CVEs:    cves,
			}
			yamlDeps = append(yamlDeps, depYaml)
		}
		output.PrintYAML(yamlDeps)
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
				key := dep.Name + "@" + dep.Version
				if r.CVEs != nil && len(r.CVEs[key]) > 0 {
					fmt.Println("  CVEs:")
					for _, cve := range r.CVEs[key] {
						fmt.Printf("    - %s: %s\n", cve.CVEID, cve.Title)
					}
				}
			}
		}
	}
}

// NewDepsCheckCommand creates a new command for checking dependencies.
func NewDepsCheckCommand(outputFormat string, cve bool) *DepsCheckCommand {
	return &DepsCheckCommand{
		BaseCheckCommand: NewBaseCheckCommand(
			"deps",
			"Check project dependencies",
			"deps [path]",
			0, // 0 required args as path is optional
		),
		outputFormat: outputFormat,
		cve:          cve,
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

	// If cve flag is set, fetch CVEs for each dependency
	if c.cve {
		result.CVEs = make(map[string][]CVEResponse)
		apiKey := os.Getenv("CHECK_API_KEY") // Optionally get from env, or make configurable
		if apiKey == "" {
			apiKey = "SPK1HgBWcxO5EmLsCSP6aIRNhX6wXMYa" // fallback demo key
		}
		apiClient := NewAPIClient(apiKey)
		versionService := NewVersionService(apiClient)
		for _, dep := range deps {
			if dep.Name == "" || dep.Version == "" {
				continue
			}
			cves, err := versionService.GetCVEs(dep.Name, dep.Version, nil)
			if err == nil && len(cves) > 0 {
				key := dep.Name + "@" + dep.Version
				result.CVEs[key] = cves
			}
		}
	}

	return result, nil
} 