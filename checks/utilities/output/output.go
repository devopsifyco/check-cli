package output

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"strings"
)

func PrintJSON(v interface{}) {
	data, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(data))
}

func PrintYAML(v interface{}) {
	data, _ := yaml.Marshal(v)
	fmt.Println(string(data))
}

// PrintTable prints a table with custom column widths and alignment.
// header: slice of column names
// rows: slice of rows, each a slice of strings
// colWidths: width for each column
// rightAlign: true for right-align, false for left-align per column
func PrintTable(header []string, rows [][]string, colWidths []int, rightAlign []bool) {
	if len(header) != len(colWidths) || len(header) != len(rightAlign) {
		fmt.Println("PrintTable: header, colWidths, and rightAlign must have the same length")
		return
	}
	// Helper to format a cell
	formatCell := func(s string, width int, right bool) string {
		ellipsis := "..."
		if len(s) > width {
			if width > len(ellipsis) {
				s = s[:width-len(ellipsis)] + ellipsis
			} else if width > 0 {
				s = s[:width]
			} else {
				s = ""
			}
		}
		pad := width - len(s)
		if pad < 0 {
			pad = 0
		}
		if right {
			return fmt.Sprintf("%*s", width, s)
		}
		return fmt.Sprintf("%-*s", width, s)
	}
	// Print header
	for i, h := range header {
		fmt.Print(formatCell(h, colWidths[i], rightAlign[i]))
		if i < len(header)-1 {
			fmt.Print("  ")
		}
	}
	fmt.Println()
	// Print separator
	for i, w := range colWidths {
		fmt.Print(strings.Repeat("-", w))
		if i < len(colWidths)-1 {
			fmt.Print("  ")
		}
	}
	fmt.Println()
	// Print rows
	for _, row := range rows {
		for i, cell := range row {
			cellStr := cell
			if i < len(colWidths) {
				cellStr = formatCell(cell, colWidths[i], rightAlign[i])
			}
			fmt.Print(cellStr)
			if i < len(row)-1 {
				fmt.Print("  ")
			}
		}
		fmt.Println()
	}
} 