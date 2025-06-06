package code

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hhatto/gocloc"
	"github.com/devopsifyco/check-cli/checks"
	"github.com/devopsifyco/check-cli/checks/utilities/output"
)

// LocResult holds the outcome of the LOC (lines of code) check.
// It contains a list of languages with their code statistics and a total summary.
type LocResult struct {
	Languages []gocloc.ClocLanguage `json:"languages" yaml:"languages"`
	Total     gocloc.ClocLanguage   `json:"total" yaml:"total"`
	Error     string                `json:"error,omitempty" yaml:"error,omitempty"`
}

// Print outputs the LOC result in the specified format (table, json, or yaml).
func (r *LocResult) Print(outputFormat string) {
	switch outputFormat {
	case "json":
		output.PrintJSON(r)
	case "yaml":
		output.PrintYAML(r)
	default:
		if r.Error != "" {
			fmt.Printf("Error: %s\n", r.Error)
			return
		}
		headers := []string{"Language", "Files", "Blank", "Comment", "Code"}
		colWidths := []int{15, 8, 8, 10, 10}
		rightAlign := []bool{false, true, true, true, true}
		rows := make([][]string, 0, len(r.Languages))
		for _, lang := range r.Languages {
			rows = append(rows, []string{
				lang.Name,
				fmt.Sprintf("%d", lang.FilesCount),
				fmt.Sprintf("%d", lang.Blanks),
				fmt.Sprintf("%d", lang.Comments),
				fmt.Sprintf("%d", lang.Code),
			})
		}
		output.PrintTable(headers, rows, colWidths, rightAlign)
		fmt.Printf("TOTAL: Files=%d Blank=%d Comment=%d Code=%d\n",
			r.Total.FilesCount, r.Total.Blanks, r.Total.Comments, r.Total.Code)
	}
}

// LocCheckCommand implements checks.CheckCommand for counting lines of code (LOC).
type LocCheckCommand struct {
	*checks.BaseCheckCommand
	outputFormat string
}

// NewLocCheckCommand creates a new command for counting lines of code (LOC).
// outputFormat specifies the output format ("json", "yaml", or default table).
func NewLocCheckCommand(outputFormat string) *LocCheckCommand {
	return &LocCheckCommand{
		BaseCheckCommand: checks.NewBaseCheckCommand(
			"loc",
			"Count lines of code in a directory or file",
			"loc [path]",
			0,
		),
		outputFormat: outputFormat,
	}
}

// Execute runs the LOC check for the given arguments (path).
// If no path is provided, it defaults to the current directory.
func (c *LocCheckCommand) Execute(args []string) (checks.CheckResult, error) {
	var target string
	if len(args) > 0 && args[0] != "" {
		target = args[0]
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			return &LocResult{Error: err.Error()}, err
		}
		target = cwd
	}
	// Expand to absolute path
	absPath, err := filepath.Abs(target)
	if err != nil {
		return &LocResult{Error: err.Error()}, err
	}
	// Prepare gocloc
	langs := gocloc.NewDefinedLanguages()
	opts := gocloc.NewClocOptions()
	processor := gocloc.NewProcessor(langs, opts)
	result, err := processor.Analyze([]string{absPath})
	if err != nil {
		return &LocResult{Error: err.Error()}, err
	}
	// Convert map to sorted slice
	languages := make([]gocloc.ClocLanguage, 0, len(result.Languages))
	var totalFiles int32
	for _, lang := range result.Languages {
		filesCount := int32(len(lang.Files))
		totalFiles += filesCount
		languages = append(languages, gocloc.ClocLanguage{
			Name:       lang.Name,
			FilesCount: filesCount,
			Code:       lang.Code,
			Comments:   lang.Comments,
			Blanks:     lang.Blanks,
		})
	}
	// Sort by code lines descending
	for i := 0; i < len(languages); i++ {
		for j := i + 1; j < len(languages); j++ {
			if languages[j].Code > languages[i].Code {
				languages[i], languages[j] = languages[j], languages[i]
			}
		}
	}
	return &LocResult{
		Languages: languages,
		Total: gocloc.ClocLanguage{
			Name:       result.Total.Name,
			FilesCount: totalFiles,
			Code:       result.Total.Code,
			Comments:   result.Total.Comments,
			Blanks:     result.Total.Blanks,
		},
	}, nil
} 