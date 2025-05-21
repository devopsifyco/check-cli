package os

import (
	"github.com/spf13/cobra"
)

// RegisterCommand registers the OS check command
func RegisterCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "os",
		Short: "Check operating system information",
		Long: `Gather and display detailed information about the operating system, including:
- Basic OS information (name, architecture, CPU count, etc.)
- Memory usage
- Disk usage
- CPU information and usage
- Network interfaces and statistics
- Process information`,
		RunE: func(cmd *cobra.Command, args []string) error {
			format, _ := cmd.Flags().GetString("format")
			
			osCmd := NewOSCheckCommand()
			osCmd.Format = format
			
			return osCmd.Execute()
		},
	}

	// Add flags
	cmd.Flags().StringP("format", "f", "json", "Output format (json)")

	return cmd
} 