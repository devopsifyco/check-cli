package checks

import (
	"encoding/json"
	"fmt"
	"os"
	"gopkg.in/yaml.v3"
)

// CheckResult represents the result of any check operation
type CheckResult interface {
	Print(outputFormat string)
}

// CheckCommand represents a check command that can be executed
type CheckCommand interface {
	Execute(args []string) (CheckResult, error)
}

// BaseCheckCommand provides common functionality for all check commands
type BaseCheckCommand struct {
	Name        string
	Description string
	Usage       string
	Args        int
}

// NewBaseCheckCommand creates a new base check command
func NewBaseCheckCommand(name, description, usage string, args int) *BaseCheckCommand {
	return &BaseCheckCommand{
		Name:        name,
		Description: description,
		Usage:       usage,
		Args:        args,
	}
}

// PrintResult is a helper function to print check results
func PrintResult(result CheckResult, outputFormat string) {
	if result == nil {
		return
	}
	result.Print(outputFormat)
}

// PrintError is a helper function to print errors consistently
func PrintError(err error) {
	fmt.Printf("Error: %v\n", err)
	os.Exit(1)
}

// PrintJSON is a helper function to print any result as JSON
func PrintJSON(v interface{}) {
	jsonData, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		PrintError(err)
	}
	fmt.Println(string(jsonData))
}

// PrintYAML is a helper function to print any result as YAML
func PrintYAML(v interface{}) {
	yamlData, err := yaml.Marshal(v)
	if err != nil {
		PrintError(err)
	}
	fmt.Println(string(yamlData))
} 