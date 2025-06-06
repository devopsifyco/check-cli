package code

import (
	"github.com/devopsifyco/check-cli/checks"
)

// CodeDepsCheckCommand wraps the existing DepsCheckCommand for use as a subcommand under 'code'.
type CodeDepsCheckCommand struct {
	*checks.DepsCheckCommand
}

func NewCodeDepsCheckCommand(outputFormat string, cve bool) *CodeDepsCheckCommand {
	return &CodeDepsCheckCommand{
		DepsCheckCommand: checks.NewDepsCheckCommand(outputFormat, cve),
	}
}

// Execute delegates to the original DepsCheckCommand implementation.
func (c *CodeDepsCheckCommand) Execute(args []string) (checks.CheckResult, error) {
	return c.DepsCheckCommand.Execute(args)
} 