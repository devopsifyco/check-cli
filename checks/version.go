package checks

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"sort"
	"strings"
	"time"
	"gopkg.in/yaml.v3"
)

const (
	apiEndpoint = "https://api.opsify.dev/checks"
)

type VersionResponseShort struct {
	Name    string     `json:"name"`
	Version string     `json:"version"`
	EOL     *time.Time `json:"eol_date"`
}

type VersionResponseFull struct {
	ID                    int        `json:"id"`
	Name                  string     `json:"name"`
	Version               string     `json:"version"`
	Vendor                *string    `json:"vendor"`
	ReleaseDate           time.Time  `json:"release_date"`
	ActiveSupportEndDate  *time.Time `json:"active_support_end_date"`
	SecuritySupportEndDate *time.Time `json:"security_support_end_date"`
	EOL                   *time.Time `json:"eol_date"`
}

// VersionHistory represents a single version entry
type VersionHistory struct {
	ProductName           string     `json:"product_name"`
	Version              string     `json:"version"`
	ReleaseDate          *time.Time `json:"release_date"`
	ActiveSupportEndDate *time.Time `json:"active_support_end_date"`
	SecuritySupportEndDate *time.Time `json:"security_support_end_date"`
	EOL                  *time.Time `json:"eol_date"`
	ID                   int        `json:"id"`
	Vendor               string     `json:"vendor"`
}

// VersionHistoryList represents a list of version histories
type VersionHistoryList []VersionHistory

// CVEResponse represents a CVE entry for a product version
// Updated to match the new flat structure
// Sample:
// {
//     "cve_id": "CVE-2026-0002",
//     "state": "PUBLISHED",
//     "published_date": "2026-02-19T16:55:30.675Z",
//     "score": null,
//     "title": "Heap buffer overflow ...",
//     "references": [ ... ]
// }
type CVEResponse struct {
	CVEID         string    `json:"cve_id"`
	State         string    `json:"state"`
	PublishedDate string    `json:"published_date"`
	Score         *float64  `json:"score"`
	Title         string    `json:"title"`
	References    []string  `json:"references"`
}

// CVEProduct represents an affected product in a CVE
// Used for AffectedProducts field
// You may need to adjust the fields based on actual API response
// Example assumes ProductName and Version
type CVEProduct struct {
	ProductName string
	Version     string
}

// getSystemDateFormat returns the date format based on the OS and locale settings
func getSystemDateFormat() string {
	switch runtime.GOOS {
	case "windows":
		return getWindowsDateFormat()
	case "darwin", "linux":
		return getUnixDateFormat()
	default:
		return "2006-01-02" // Default ISO format
	}
}

// formatTime helper function to safely format nullable time values using system date format
func formatTime(t *time.Time) string {
	if t == nil {
		return "N/A"
	}
	return t.Format(getSystemDateFormat())
}

// formatTimeForOutput helper function to format time for JSON/YAML output
func formatTimeForOutput(t *time.Time) string {
	if t == nil {
		return ""
	}
	// Always use ISO format for structured output
	return t.Format("2006-01-02")
}

// Custom time parsing for API dates
func (v *VersionResponseFull) UnmarshalJSON(data []byte) error {
	type Alias VersionResponseFull
	aux := &struct {
		ReleaseDate          *string `json:"release_date"`
		ActiveSupportEndDate *string `json:"active_support_end_date"`
		SecuritySupportEndDate *string `json:"security_support_end_date"`
		EOL                  *string `json:"eol_date"`
		*Alias
	}{
		Alias: (*Alias)(v),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Parse dates with support for null values
	layouts := []string{"2006-01-02", "2006-01-02T15:04:05"}

	// Parse ReleaseDate (now optional)
	if aux.ReleaseDate != nil {
		var parseError error
		for _, layout := range layouts {
			if t, err := time.Parse(layout, *aux.ReleaseDate); err == nil {
				v.ReleaseDate = t
				parseError = nil
				break
			} else {
				parseError = err
			}
		}
		if parseError != nil {
			return fmt.Errorf("failed to parse ReleaseDate: %v", parseError)
		}
	} else {
		// Set a zero time if release date is null
		v.ReleaseDate = time.Time{}
	}

	// Parse nullable dates
	if aux.ActiveSupportEndDate != nil {
		for _, layout := range layouts {
			if t, err := time.Parse(layout, *aux.ActiveSupportEndDate); err == nil {
				v.ActiveSupportEndDate = &t
				break
			}
		}
	}
	if aux.SecuritySupportEndDate != nil {
		for _, layout := range layouts {
			if t, err := time.Parse(layout, *aux.SecuritySupportEndDate); err == nil {
				v.SecuritySupportEndDate = &t
				break
			}
		}
	}
	if aux.EOL != nil {
		for _, layout := range layouts {
			if t, err := time.Parse(layout, *aux.EOL); err == nil {
				v.EOL = &t
				break
			}
		}
	}

	return nil
}

// Custom time parsing for version history items
func (v *VersionHistory) UnmarshalJSON(data []byte) error {
	type Alias VersionHistory
	aux := &struct {
		ReleaseDate         *string `json:"release_date"`
		ActiveSupportEndDate *string `json:"active_support_end_date"`
		SecuritySupportEndDate *string `json:"security_support_end_date"`
		EOL                  *string `json:"eol_date"`
		*Alias
	}{
		Alias: (*Alias)(v),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Parse nullable dates
	layouts := []string{"2006-01-02", "2006-01-02T15:04:05"}
	if aux.ReleaseDate != nil {
		for _, layout := range layouts {
			if t, err := time.Parse(layout, *aux.ReleaseDate); err == nil {
				v.ReleaseDate = &t
				break
			}
		}
	}
	if aux.ActiveSupportEndDate != nil {
		for _, layout := range layouts {
			if t, err := time.Parse(layout, *aux.ActiveSupportEndDate); err == nil {
				v.ActiveSupportEndDate = &t
				break
			}
		}
	}
	if aux.SecuritySupportEndDate != nil {
		for _, layout := range layouts {
			if t, err := time.Parse(layout, *aux.SecuritySupportEndDate); err == nil {
				v.SecuritySupportEndDate = &t
				break
			}
		}
	}
	if aux.EOL != nil {
		for _, layout := range layouts {
			if t, err := time.Parse(layout, *aux.EOL); err == nil {
				v.EOL = &t
				break
			}
		}
	}

	return nil
}

// APIClient represents the API client configuration
type APIClient struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}

// NewAPIClient creates a new API client instance
func NewAPIClient(apiKey string) *APIClient {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // Skip certificate verification
		},
	}

	return &APIClient{
		BaseURL: apiEndpoint,
		APIKey:  apiKey,
		HTTPClient: &http.Client{
			Transport: tr,
		},
	}
}

// VersionService handles version-related operations
type VersionService struct {
	client *APIClient
}

// NewVersionService creates a new version service
func NewVersionService(client *APIClient) *VersionService {
	return &VersionService{
		client: client,
	}
}

// GetVersion retrieves version information for a component
func (s *VersionService) GetVersion(component string) (*VersionResponseFull, error) {
	url := fmt.Sprintf("%s/release/%s/latest?apikey=%s", s.client.BaseURL, component, s.client.APIKey)
	result, err := makeRequest[VersionResponseFull](s.client.HTTPClient, url)
	if err != nil {
		return nil, err
	}
	// Check if result is empty (no version found)
	if result.Version == "" {
		return nil, fmt.Errorf("version not found for component: %s", component)
	}
	// Ensure component name is set
	if result.Name == "" {
		result.Name = component
	}
	return result, nil
}

// GetSpecificVersion retrieves information for a specific version of a component
func (s *VersionService) GetSpecificVersion(component string, version string) (*VersionResponseFull, error) {
	url := fmt.Sprintf("%s/release/%s/%s?apikey=%s", s.client.BaseURL, component, version, s.client.APIKey)
	result, err := makeRequest[VersionResponseFull](s.client.HTTPClient, url)
	if err != nil {
		return nil, err
	}
	// Check if result is empty (no version found)
	if result.Version == "" {
		return nil, fmt.Errorf("version %s not found for component: %s", version, component)
	}
	// Ensure component name is set
	if result.Name == "" {
		result.Name = component
	}
	return result, nil
}

// GetVersions retrieves version history for a component
func (s *VersionService) GetVersions(component string) (VersionHistoryList, error) {
	url := fmt.Sprintf("%s/release/%s?apikey=%s", s.client.BaseURL, component, s.client.APIKey)
	result, err := makeRequest[[]VersionHistory](s.client.HTTPClient, url)
	if err != nil {
		return nil, err
	}
	// Check if result is empty (no versions found)
	if result == nil || len(*result) == 0 {
		return nil, fmt.Errorf("no versions found for component: %s", component)
	}
	return *result, nil
}

// GetCVEs retrieves CVE information for a specific version of a component
func (s *VersionService) GetCVEs(component string, version string, limit *int) ([]CVEResponse, error) {
	url := fmt.Sprintf("%s/release/%s/%s/cves", s.client.BaseURL, component, version)
	
	// Set default limit to 100 if not provided or if 0
	//limitValue := 100
	//if limit != nil && *limit > 0 {
	//	limitValue = *limit
	//}
	//url = fmt.Sprintf("%s&limit=%d", url, limitValue)
	
	url = fmt.Sprintf("%s?apikey=%s", url, s.client.APIKey)	

	result, err := makeRequest[[]CVEResponse](s.client.HTTPClient, url)
	if err != nil {
		return nil, err
	}
	return *result, nil
}

// ErrorResponse represents an error response from the API
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// makeRequest is a generic function to make HTTP requests and parse responses
func makeRequest[T any](client *http.Client, url string) (*T, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.Unmarshal(body, &errResp); err == nil && (errResp.Error != "" || errResp.Message != "") {
			// If we can parse the error response and it has content
			return nil, fmt.Errorf("API error (status %d): %s - %s", 
				resp.StatusCode, 
				errResp.Error, 
				errResp.Message)
		}
		// If we can't parse the error response or it's empty, include the raw body in the error
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Check for empty response
	if len(body) == 0 {
		return nil, fmt.Errorf("empty response received from API")
	}

	// Try to parse the response
	var result T
	if err := json.Unmarshal(body, &result); err != nil {
		// If parsing fails, try to determine if it's a known error format
		var errResp ErrorResponse
		if jsonErr := json.Unmarshal(body, &errResp); jsonErr == nil && (errResp.Error != "" || errResp.Message != "") {
			return nil, fmt.Errorf("API error: %s - %s", 
				errResp.Error, 
				errResp.Message)
		}
		// If it's not a known error format, include both the parsing error and the raw response
		return nil, fmt.Errorf("error parsing response: %w (raw response: %s)", err, string(body))
	}

	// For pointer types, check if the result is empty/nil
	if ptr, ok := any(&result).(interface{ IsZero() bool }); ok && ptr.IsZero() {
		return nil, fmt.Errorf("API returned empty result")
	}

	return &result, nil
}

// VersionPrinter handles formatting and displaying version information
type VersionPrinter struct {
	jsonOutput bool
	fullOutput bool
}

// NewVersionPrinter creates a new version printer
func NewVersionPrinter(jsonOutput, fullOutput bool) *VersionPrinter {
	return &VersionPrinter{
		jsonOutput: jsonOutput,
		fullOutput: fullOutput,
	}
}

// PrintVersion prints version information
func (p *VersionPrinter) PrintVersion(result *VersionResponseFull) {
	if p.jsonOutput {
		var output interface{}
		if p.fullOutput {
			output = result
		} else {
			output = &VersionResponseShort{
				Name:    result.Name,
				Version: result.Version,
				EOL:     result.EOL,
			}
		}
		jsonData, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			fmt.Printf("Error formatting JSON: %v\n", err)
			return
		}
		fmt.Println(string(jsonData))
	} else {
		if p.fullOutput {
			p.printFullVersion(result)
		} else {
			// Print short version
			fmt.Printf("Name: %s\n", result.Name)
			fmt.Printf("Version: %s\n", result.Version)
			if result.EOL != nil {
				fmt.Printf("EOL Date: %s\n", formatTime(result.EOL))
			}
		}
	}
}

// PrintFullVersion prints detailed version information
func (p *VersionPrinter) printFullVersion(result *VersionResponseFull) {
	fmt.Printf("Name: %s\n", result.Name)
	fmt.Printf("Version: %s\n", result.Version)
	if result.Vendor != nil {
		fmt.Printf("Vendor: %s\n", *result.Vendor)
	}
	if !result.ReleaseDate.IsZero() {
		fmt.Printf("Release Date: %s\n", formatTime(&result.ReleaseDate))
	}
	if result.ActiveSupportEndDate != nil {
		fmt.Printf("Active Support End: %s\n", formatTime(result.ActiveSupportEndDate))
	}
	if result.SecuritySupportEndDate != nil {
		fmt.Printf("Security Support End: %s\n", formatTime(result.SecuritySupportEndDate))
	}
	if result.EOL != nil {
		fmt.Printf("EOL Date: %s\n", formatTime(result.EOL))
	}
}

// VersionResult implements CheckResult interface for version checks
type VersionResult struct {
	Name              string     `json:"name,omitempty"`
	Version           string     `json:"version,omitempty"`
	ReleaseDate       *time.Time `json:"release_date,omitempty"`
	ActiveSupportEnd  *time.Time `json:"active_support_end,omitempty"`
	SecuritySupportEnd *time.Time `json:"security_support_end,omitempty"`
	EOLDate           *time.Time `json:"eol_date,omitempty"`
	ID                int        `json:"id,omitempty"`
	CreatedAt         *time.Time `json:"created_at,omitempty"`
	Error             string     `json:"error,omitempty"`
	Message           string     `json:"message,omitempty"`
}

// UnmarshalJSON implements custom JSON unmarshaling for VersionResult
func (v *VersionResult) UnmarshalJSON(data []byte) error {
	type Alias VersionResult
	aux := &struct {
		ReleaseDate       interface{} `json:"release_date"`
		ActiveSupportEnd  interface{} `json:"active_support_end"`
		SecuritySupportEnd interface{} `json:"security_support_end"`
		EOLDate          interface{} `json:"eol_date"`
		CreatedAt        interface{} `json:"created_at"`
		*Alias
	}{
		Alias: (*Alias)(v),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Helper function to parse dates with multiple formats
	parseDate := func(field interface{}, fieldName string) (*time.Time, error) {
		if field == nil {
			return nil, nil
		}

		dateStr, ok := field.(string)
		if !ok {
			return nil, fmt.Errorf("invalid date format for %s", fieldName)
		}

		formats := []string{
			"2006-01-02",
			"2006-01-02T15:04:05Z",
			"2006-01-02T15:04:05",
		}

		for _, format := range formats {
			if t, err := time.Parse(format, dateStr); err == nil {
				return &t, nil
			}
		}
		return nil, fmt.Errorf("could not parse %s date: %s", fieldName, dateStr)
	}

	var err error
	if v.ReleaseDate, err = parseDate(aux.ReleaseDate, "release_date"); err != nil {
		return err
	}
	if v.ActiveSupportEnd, err = parseDate(aux.ActiveSupportEnd, "active_support_end"); err != nil {
		return err
	}
	if v.SecuritySupportEnd, err = parseDate(aux.SecuritySupportEnd, "security_support_end"); err != nil {
		return err
	}
	if v.EOLDate, err = parseDate(aux.EOLDate, "eol_date"); err != nil {
		return err
	}
	if v.CreatedAt, err = parseDate(aux.CreatedAt, "created_at"); err != nil {
		return err
	}

	return nil
}

// Helper struct for JSON/YAML output
type structuredOutput struct {
	Name               string `json:"name,omitempty" yaml:"name,omitempty"`
	Version            string `json:"version,omitempty" yaml:"version,omitempty"`
	Vendor             string `json:"vendor,omitempty" yaml:"vendor,omitempty"`
	ReleaseDate        string `json:"release_date,omitempty" yaml:"release_date,omitempty"`
	ActiveSupportEnd   string `json:"active_support_end,omitempty" yaml:"active_support_end,omitempty"`
	SecuritySupportEnd string `json:"security_support_end,omitempty" yaml:"security_support_end,omitempty"`
	EOLDate            string `json:"eol_date,omitempty" yaml:"eol_date,omitempty"`
	ID                 int    `json:"id,omitempty" yaml:"id,omitempty"`
	CreatedAt          string `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	UpdatedAt          string `json:"updated_at,omitempty" yaml:"updated_at,omitempty"`
	Error              string `json:"error,omitempty" yaml:"error,omitempty"`
	Message            string `json:"message,omitempty" yaml:"message,omitempty"`
}

// Print implements CheckResult interface
func (r *VersionResult) Print(outputFormat string) {
	// If there's an error, only print error information
	if r.Error != "" || r.Message != "" {
		errResult := structuredOutput{
			Error:   r.Error,
			Message: r.Message,
		}
		
		switch outputFormat {
		case "json":
			jsonData, err := json.MarshalIndent(errResult, "", "  ")
			if err != nil {
				fmt.Printf("Error formatting JSON: %v\n", err)
				return
			}
			fmt.Println(string(jsonData))
		case "yaml":
			yamlData, err := yaml.Marshal(errResult)
			if err != nil {
				fmt.Printf("Error formatting YAML: %v\n", err)
				return
			}
			fmt.Println(string(yamlData))
		default:
			if r.Error != "" {
				fmt.Printf("Error: %s\n", r.Error)
			}
			if r.Message != "" {
				fmt.Printf("Message: %s\n", r.Message)
			}
		}
		return
	}

	// Format for structured output
	output := structuredOutput{
		Name:               r.Name,
		Version:            r.Version,
		ReleaseDate:        formatTimeForOutput(r.ReleaseDate),
		ActiveSupportEnd:   formatTimeForOutput(r.ActiveSupportEnd),
		SecuritySupportEnd: formatTimeForOutput(r.SecuritySupportEnd),
		EOLDate:           formatTimeForOutput(r.EOLDate),
		//ID:                r.ID,
		CreatedAt:         formatTimeForOutput(r.CreatedAt),
	}

	switch outputFormat {
	case "json":
		jsonData, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			fmt.Printf("Error formatting JSON: %v\n", err)
			return
		}
		fmt.Println(string(jsonData))
	case "yaml":
		yamlData, err := yaml.Marshal(output)
		if err != nil {
			fmt.Printf("Error formatting YAML: %v\n", err)
			return
		}
		fmt.Println(string(yamlData))
	default:
		fmt.Printf("Name: %s\n", r.Name)
		fmt.Printf("Version: %s\n", r.Version)
		fmt.Printf("Release Date: %s\n", formatTime(r.ReleaseDate))
		fmt.Printf("Active Support End: %s\n", formatTime(r.ActiveSupportEnd))
		fmt.Printf("Security Support End: %s\n", formatTime(r.SecuritySupportEnd))
		fmt.Printf("EOL Date: %s\n", formatTime(r.EOLDate))		
		if r.CreatedAt != nil {
			fmt.Printf("Created At: %s\n", formatTime(r.CreatedAt))
		}
	}
}

// String implements Stringer interface for VersionResult
func (r *VersionResult) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "Name: %s\n", r.Name)
	fmt.Fprintf(&b, "Version: %s\n", r.Version)
	fmt.Fprintf(&b, "Release Date: %s\n", formatTime(r.ReleaseDate))
	fmt.Fprintf(&b, "Active Support End: %s\n", formatTime(r.ActiveSupportEnd))
	fmt.Fprintf(&b, "Security Support End: %s\n", formatTime(r.SecuritySupportEnd))
	fmt.Fprintf(&b, "EOL Date: %s\n", formatTime(r.EOLDate))
	if r.Error != "" {
		fmt.Fprintf(&b, "Error: %s\n", r.Error)
	}
	if r.Message != "" {
		fmt.Fprintf(&b, "Message: %s\n", r.Message)
	}
	return b.String()
}

// VersionCheckCommand implements the CheckCommand interface for version checks
type VersionCheckCommand struct {
	*BaseCheckCommand
	Name         string
	Version      string
	Service      *VersionService
	apiKey       string
	outputFormat string
	fullOutput   bool
	history      bool
	client       bool
	CVE          bool
}

// NewVersionCheckCommand creates a new version check command
func NewVersionCheckCommand(apiKey string, outputFormat string, fullOutput bool, history bool, client bool, cve bool) *VersionCheckCommand {
	// Use demo API key if none provided
	if apiKey == "" {
		apiKey = "SPK1HgBWcxO5EmLsCSP6aIRNhX6wXMYa"
	}

	// Create API client and version service
	apiClient := NewAPIClient(apiKey)
	versionService := NewVersionService(apiClient)

	return &VersionCheckCommand{
		BaseCheckCommand: NewBaseCheckCommand(
			"version",
			"Check version information for a component",
			"version [component] [version]",
			1,
		),
		Service:      versionService,
		apiKey:       apiKey,
		outputFormat: outputFormat,
		fullOutput:   fullOutput,
		history:      history,
		client:       client,
		CVE:         cve,
	}
}

// sortCVEs sorts CVEs by score (higher first) as the primary key, and published_date (most recent first) as the secondary key
func sortCVEs(cves []CVEResponse) {
	sort.Slice(cves, func(i, j int) bool {
		// Compare by score (higher first)
		var scoreI, scoreJ float64
		if cves[i].Score != nil {
			scoreI = *cves[i].Score
		}
		if cves[j].Score != nil {
			scoreJ = *cves[j].Score
		}
		if scoreI != scoreJ {
			return scoreI > scoreJ
		}

		// Compare by published date (most recent first)
		parseTime := func(s string) time.Time {
			t, err := time.Parse(time.RFC3339, s)
			if err != nil {
				return time.Time{}
			}
			return t
		}
		timeI := parseTime(cves[i].PublishedDate)
		timeJ := parseTime(cves[j].PublishedDate)
		if !timeI.Equal(timeJ) {
			return timeI.After(timeJ)
		}

		// No Priority/Severity fields, so skip those
		return false
	})
}

// CombinedResult represents a combined version and CVE result
type CombinedResult struct {
	Version    interface{}   `json:"version"`
	CVEs       []CVEResponse `json:"cves,omitempty"`
	fullOutput bool         // Add fullOutput field
}

// NewCombinedResult creates a new CombinedResult with the specified output format
func NewCombinedResult(version interface{}, cves []CVEResponse, fullOutput bool) *CombinedResult {
	return &CombinedResult{
		Version:    version,
		CVEs:       cves,
		fullOutput: fullOutput,
	}
}

// Print implements CheckResult interface
func (cr *CombinedResult) Print(outputFormat string) {
	switch outputFormat {
	case "json":
		jsonData, err := json.MarshalIndent(cr, "", "  ")
		if err != nil {
			fmt.Printf("Error formatting JSON: %v\n", err)
			return
		}
		fmt.Println(string(jsonData))
	case "yaml":
		yamlData, err := yaml.Marshal(cr)
		if err != nil {
			fmt.Printf("Error formatting YAML: %v\n", err)
			return
		}
		fmt.Println(string(yamlData))
	default:
		// Print version information
		if short, ok := cr.Version.(*VersionResponseShort); ok {
			fmt.Printf("Name: %s\n", short.Name)
			fmt.Printf("Version: %s\n", short.Version)
			fmt.Printf("EOL Date: %s\n", formatTime(short.EOL))
		} else if full, ok := cr.Version.(*VersionResponseFull); ok {
			// Name & Vendor on one line, aligned
			if full.Vendor != nil {
				fmt.Printf("Name: %-20s          Vendor: %s\n", full.Name, *full.Vendor)
			} else {
				fmt.Printf("Name: %s\n", full.Name)
			}
			// Version & Release Date on one line, aligned
			if !full.ReleaseDate.IsZero() {
				fmt.Printf("Version: %-17s          Release Date: %s\n", full.Version, formatTime(&full.ReleaseDate))
			} else {
				fmt.Printf("Version: %s\n", full.Version)
			}
			if full.ActiveSupportEndDate != nil {
				fmt.Printf("Active Support End: %s\n", formatTime(full.ActiveSupportEndDate))
			}
			if full.SecuritySupportEndDate != nil {
				fmt.Printf("Security Support End: %s\n", formatTime(full.SecuritySupportEndDate))
			}
			if full.EOL != nil {
				fmt.Printf("EOL Date: %s\n", formatTime(full.EOL))
			}
		}

		// Print CVE information
		if len(cr.CVEs) > 0 {
			// Use deps.go style table
			sep := "  " + strings.Repeat("-", 18) + "  " + strings.Repeat("-", 10) + "  " + strings.Repeat("-", 6) + "  " + strings.Repeat("-", 50)
			headFmt := "  %-18s  %-10s  %-6s  %s\n"
			fmt.Println(sep)
			fmt.Printf(headFmt, "CVE ID", "Published", "Score", "Title")
			fmt.Println(sep)
			for _, cve := range cr.CVEs {
				cveID := cve.CVEID
				if len(cveID) > 18 {
					cveID = cveID[:18]
				}
				published := cve.PublishedDate
				if len(published) >= 10 {
					published = published[:10]
				}
				score := ""
				if cve.Score != nil {
					score = fmt.Sprintf("%.1f", *cve.Score)
				}
				title := cve.Title
				if len(title) > 50 {
					title = title[:47] + "..."
				}
				fmt.Printf(headFmt, cveID, published, score, title)
			}
			fmt.Println(sep)
		} else {
			fmt.Println("No CVEs found")
		}

	}
}

// Execute runs the version check command
func (c *VersionCheckCommand) Execute(args []string) (CheckResult, error) {
	// Parse arguments
	if len(args) > 0 {
		c.Name = args[0]
	}
	if len(args) > 1 {
		c.Version = args[1]
	}

	// Validate input parameters
	if c.Name == "" {
		return nil, fmt.Errorf("component name is required")
	}

	// Handle --history flag
	if c.history {
		versions, err := c.Service.GetVersions(c.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to get version history for component %q: %w", c.Name, err)
		}
		return versions, nil
	}

	// Get version information first
	var version *VersionResponseFull
	var err error
	
	if c.Version == "" {
		// Get latest version
		version, err = c.Service.GetLatestVersion(c.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to get latest version for component %q: %w", c.Name, err)
		}
	} else {
		// Get specific version
		version, err = c.Service.GetSpecificVersion(c.Name, c.Version)
		if err != nil {
			return nil, fmt.Errorf("failed to get version %q for component %q: %w", c.Version, c.Name, err)
		}
	}

	// If CVE flag is set, get CVEs and return combined result
	if c.CVE {
		cves, err := c.Service.GetCVEs(c.Name, version.Version, nil) // always pass nil for limit
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to get CVEs: %v\n", err)
			cves = []CVEResponse{} 
		}

		sortCVEs(cves)

		return NewCombinedResult(version, cves, c.fullOutput), nil
	}

	// Return short or full version based on fullOutput flag
	if !c.fullOutput {
		return &VersionResponseShort{
			Name:    version.Name,
			Version: version.Version,
			EOL:     version.EOL,
		}, nil
	}

	return version, nil
}

// CheckVersion is the main function to check version information
func CheckVersion(component string, version string, apiKey string, outputFormat string, fullOutput bool, history bool, client bool, cve bool) {
	cmd := NewVersionCheckCommand(apiKey, outputFormat, fullOutput, history, client, cve)
	args := []string{component}
	if version != "" {
		args = append(args, version)
	}
	
	result, err := cmd.Execute(args)
	if err != nil {
		PrintError(err)
		return
	}
	result.Print(outputFormat)
}

// Print implements CheckResult interface
func (r *VersionResponseShort) Print(outputFormat string) {
	// Check if this is a combined result with CVEs
	if combined, ok := interface{}(r).(interface{
		GetVersion() interface{}
		GetCVEs() []CVEResponse
	}); ok {
		version := combined.GetVersion().(*VersionResponseShort)
		cves := combined.GetCVEs()

		output := struct {
			structuredOutput
			CVEs []CVEResponse `json:"cves,omitempty"`
		}{
			structuredOutput: structuredOutput{
				Name:    version.Name,
				Version: version.Version,
				EOLDate: formatTimeForOutput(version.EOL),
			},
			CVEs: cves,
		}

		switch outputFormat {
		case "json":
			jsonData, err := json.MarshalIndent(output, "", "  ")
			if err != nil {
				fmt.Printf("Error formatting JSON: %v\n", err)
				return
			}
			fmt.Println(string(jsonData))
		case "yaml":
			yamlData, err := yaml.Marshal(output)
			if err != nil {
				fmt.Printf("Error formatting YAML: %v\n", err)
				return
			}
			fmt.Println(string(yamlData))
		default:
			fmt.Printf("Name: %s\n", version.Name)
			fmt.Printf("Version: %s\n", version.Version)
			fmt.Printf("EOL Date: %s\n", formatTime(version.EOL))
			
			if len(cves) > 0 {
				fmt.Println("\nCVEs:")
				for _, cve := range cves {
					fmt.Printf("\nCVE ID: %s\n", cve.CVEID)
					fmt.Printf("Title: %s\n", cve.Title)
					fmt.Printf("State: %s\n", cve.State)
					fmt.Printf("Published Date: %s\n", cve.PublishedDate)
					if cve.Score != nil {
						fmt.Printf("Score: %.2f\n", *cve.Score)
					} else {
						fmt.Printf("Score: null\n")
					}
					fmt.Println("---")
				}
			}
		}
		return
	}

	// If not a combined result, print normal version information
	output := structuredOutput{
		Name:    r.Name,
		Version: r.Version,
		EOLDate: formatTimeForOutput(r.EOL),
	}

	switch outputFormat {
	case "json":
		jsonData, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			fmt.Printf("Error formatting JSON: %v\n", err)
			return
		}
		fmt.Println(string(jsonData))
	case "yaml":
		yamlData, err := yaml.Marshal(output)
		if err != nil {
			fmt.Printf("Error formatting YAML: %v\n", err)
			return
		}
		fmt.Println(string(yamlData))
	default:
		fmt.Printf("Name: %s\n", r.Name)
		fmt.Printf("Version: %s\n", r.Version)
		fmt.Printf("EOL Date: %s\n", formatTime(r.EOL))
	}
}

// Print implements CheckResult interface
func (r VersionHistoryList) Print(outputFormat string) {
	switch outputFormat {
	case "json", "yaml":
		versions := make([]structuredOutput, len(r))
		for i, version := range r {
			versions[i] = structuredOutput{
				Name:               version.ProductName,
				Version:            version.Version,
				Vendor:             version.Vendor,
				ReleaseDate:        formatTimeForOutput(version.ReleaseDate),
				EOLDate:            formatTimeForOutput(version.EOL),
			}
		}

		var data []byte
		var err error
		if outputFormat == "json" {
			data, err = json.MarshalIndent(versions, "", "  ")
		} else {
			data, err = yaml.Marshal(versions)
		}

		if err != nil {
			fmt.Printf("Error formatting %s: %v\n", outputFormat, err)
			return
		}
		fmt.Println(string(data))
	default:
		// Print as table
		headFmt := "%-12s %-12s %-18s %-20s %-12s\n"
		rowFmt := "%-12s %-12s %-18s %-20s %-12s\n"
		fmt.Printf(headFmt, "Version", "Release Date", "Active Support End", "Security Support End", "EOL Date")
		fmt.Printf("%s\n", strings.Repeat("-", 74))
		for _, version := range r {
			fmt.Printf(rowFmt,
				version.Version,
				formatTime(version.ReleaseDate),
				formatTime(version.ActiveSupportEndDate),
				formatTime(version.SecuritySupportEndDate),
				formatTime(version.EOL),
			)
		}
	}
}

// Helper function to print errors in the requested format
func printError(err string, outputFormat string) {
	errResult := struct {
		Error string `json:"error" yaml:"error"`
	}{
		Error: err,
	}

	switch outputFormat {
	case "json":
		jsonData, err := json.MarshalIndent(errResult, "", "  ")
		if err != nil {
			fmt.Printf("Error formatting JSON: %v\n", err)
			return
		}
		fmt.Println(string(jsonData))
	case "yaml":
		yamlData, err := yaml.Marshal(errResult)
		if err != nil {
			fmt.Printf("Error formatting YAML: %v\n", err)
			return
		}
		fmt.Println(string(yamlData))
	default:
		fmt.Printf("Error: %s\n", err)
	}
}

// getWindowsDateFormat returns the date format for Windows
func getWindowsDateFormat() string {
	// TODO: Implement Windows-specific date format detection
	// For now, return a default format
	return "2006-01-02"
}

// getUnixDateFormat returns the date format for Unix-like systems
func getUnixDateFormat() string {
	// TODO: Implement Unix-specific date format detection
	// For now, return a default format
	return "2006-01-02"
}

// checkLocalVersion checks the local version of a component
func checkLocalVersion(component string, outputFormat string, fullOutput bool) (CheckResult, error) {
	// For now, we'll implement a basic version check that returns the CLI version
	version := "0.1.0" // This should be replaced with actual version detection
	
	result := &VersionResponseFull{
		Name:        component,
		Version:     version,
		ReleaseDate: time.Now(), // This should be replaced with actual release date
	}
	
	return result, nil
}

// GetLatestVersion retrieves the latest version information for a component
func (s *VersionService) GetLatestVersion(component string) (*VersionResponseFull, error) {
	return s.GetVersion(component)
}

// Print implements CheckResult interface for VersionResponseFull
func (r *VersionResponseFull) Print(outputFormat string) {
	output := structuredOutput{
		Name:               r.Name,
		Version:            r.Version,
		ReleaseDate:        formatTimeForOutput(&r.ReleaseDate),
		ActiveSupportEnd:   formatTimeForOutput(r.ActiveSupportEndDate),
		SecuritySupportEnd: formatTimeForOutput(r.SecuritySupportEndDate),
		EOLDate:           formatTimeForOutput(r.EOL),
		//ID:                r.ID,
	}
	if r.Vendor != nil {
		output.Vendor = *r.Vendor
	}

	switch outputFormat {
	case "json":
		jsonData, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			fmt.Printf("Error formatting JSON: %v\n", err)
			return
		}
		fmt.Println(string(jsonData))
	case "yaml":
		yamlData, err := yaml.Marshal(output)
		if err != nil {
			fmt.Printf("Error formatting YAML: %v\n", err)
			return
		}
		fmt.Println(string(yamlData))
	default:
		fmt.Printf("Name: %s\n", r.Name)
		fmt.Printf("Version: %s\n", r.Version)
		if r.Vendor != nil {
			fmt.Printf("Vendor: %s\n", *r.Vendor)
		}
		fmt.Printf("Release Date: %s\n", formatTime(&r.ReleaseDate))
		fmt.Printf("Active Support End: %s\n", formatTime(r.ActiveSupportEndDate))
		fmt.Printf("Security Support End: %s\n", formatTime(r.SecuritySupportEndDate))
		fmt.Printf("EOL Date: %s\n", formatTime(r.EOL))
		//fmt.Printf("ID: %d\n", r.ID)
	}
} 