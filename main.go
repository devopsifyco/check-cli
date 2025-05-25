package main

import (
	"fmt"
	"os"

	"github.com/devopsifyco/check-cli/checks"
	"github.com/spf13/cobra"
)

var (
	outputFormat string
	apiKey     string
)

func main() {
	// Create root command
	rootCmd := &cobra.Command{
		Use:   "check",
		Short: "DevOpsify Check Tool for various system checks",
		Long:  "DevOpsify Check Tool for performing various system checks including version, OS, speed, and SSL certificate checks.",
	}

	// Add global flags
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "", "Output format (json, yaml)")
	rootCmd.PersistentFlags().StringVar(&apiKey, "apikey", "", "API key for version checks")

	// Create command registry
	registry := make(map[string]checks.CheckCommand)

	// Register commands with their names
	commands := map[string]checks.CheckCommand{
		"version":   checks.NewVersionCheckCommand(apiKey, "", false, false, false, false),
		"os":        checks.NewOSCheckCommand(),
		"speed":     checks.NewSpeedCheckCommand(),
		"ssl":       checks.NewSSLCheckCommand(),
		"deps":      checks.NewDepsCheckCommand(outputFormat, false),
	}

	// Add commands to root
	for name, checkCmd := range commands {
		// Store command in registry
		registry[name] = checkCmd

		// Create a new cobra command
		cmdName := name // Create a copy of the name for the closure
		cobraCmd := &cobra.Command{
			Use:   name + " [args]",
			Short: name + " check",
			Run: func(cobraCmd *cobra.Command, args []string) {
				// Get the check command from registry
				cmd := registry[cmdName]
				if cmd == nil {
					fmt.Printf("Error: Command %s not found\n", cmdName)
					os.Exit(1)
				}

				// For version command, create a new instance with the current flags
				if cmdName == "version" {
					// Get command-specific flags
					client, _ := cobraCmd.Flags().GetBool("client")
					full, _ := cobraCmd.Flags().GetBool("full")
					history, _ := cobraCmd.Flags().GetBool("history")
					output, _ := cobraCmd.Flags().GetString("output")
					cve, _ := cobraCmd.Flags().GetBool("cve")

					// Use command-specific output format if provided, otherwise use global output format
					if output != "" {
						outputFormat = output
					}

					cmd = checks.NewVersionCheckCommand(apiKey, outputFormat, full, history, client, cve)
					registry[cmdName] = cmd
				}

				// For deps command, create a new instance with the current flags
				if cmdName == "deps" {
					cve, _ := cobraCmd.Flags().GetBool("cve")
					cmd = checks.NewDepsCheckCommand(outputFormat, cve)
					registry[cmdName] = cmd
				}

				// Execute the check command
				result, err := cmd.Execute(args)
				if err != nil {
					fmt.Printf("Error: %v\n", err)
					os.Exit(1)
				}

				// Print the result
				result.Print(outputFormat)
			},
		}

		// Add version-specific flags
		if name == "version" {
			cobraCmd.Flags().BoolP("client", "c", false, "Use local client instead of remote API")
			cobraCmd.Flags().BoolP("full", "f", false, "Show full version information")
			cobraCmd.Flags().BoolP("history", "H", false, "Show version history")
			cobraCmd.Flags().StringP("output", "o", "", "Output format (json, yaml)")
			cobraCmd.Flags().Bool("cve", false, "Include CVEs in the response")
		}

		// Add deps-specific flags
		if name == "deps" {
			cobraCmd.Flags().Bool("cve", false, "Include CVEs in the response")
		}

		// Add the command to root
		rootCmd.AddCommand(cobraCmd)
	}

	// Execute root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
} 