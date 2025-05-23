package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/devopsifyco/check-cli/checks"
)

func main() {
	// Define flags
	apiKey := flag.String("apikey", "", "API key for authentication")
	jsonOutput := flag.Bool("json", false, "Output in JSON format")
	fullOutput := flag.Bool("full", false, "Show full version information")
	history := flag.Bool("history", false, "Show version history")
	client := flag.Bool("client", false, "Check local client version")
	cve := flag.Bool("cve", false, "Include CVEs in the response")

	flag.Parse()

	// Get command and arguments
	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("Error: No command specified")
		os.Exit(1)
	}

	command := args[0]
	commandArgs := args[1:]

	// Handle version command
	if command == "version" {
		versionCommand := checks.NewVersionCheckCommand(*apiKey, *jsonOutput, *fullOutput, *history, *client, *cve)
		result, err := versionCommand.Execute(commandArgs)
		if err != nil {
			fmt.Printf("Error executing version command: %v\n", err)
			os.Exit(1)
		}
		if result != nil {
			result.Print(*jsonOutput)
		}
		return
	} else if command == "deps" {
		// Handle dependencies command
		// Example Usage: check deps [--json] [/path/to/project]
		// If no path is provided, it defaults to the current directory.
		depsCommand := checks.NewDepsCheckCommand(*jsonOutput, *cve)
		result, err := depsCommand.Execute(commandArgs) // commandArgs might contain the directory path

		// Print the result regardless of error, as result might contain partial info or error details
		if result != nil {
			result.Print(*jsonOutput)
		}

		// Exit with error code if the execution failed
		if err != nil {
			// Error message is already printed by result.Print or handled within Execute
			// We just need to exit with a non-zero status
			os.Exit(1)
		}
		return
	}

	// Handle other commands...
	fmt.Printf("Error: Unknown command '%s'\n", command)
	os.Exit(1)
} 