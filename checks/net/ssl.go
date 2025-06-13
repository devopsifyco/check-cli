package net

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"
	"gopkg.in/yaml.v3"
	"github.com/devopsifyco/check-cli/checks"
)

// SSLCertResult implements CheckResult interface for SSL certificate checks
type SSLCertResult struct {
	Domain        string     `json:"domain"`
	Valid         bool       `json:"valid"`
	Issuer        string     `json:"issuer"`
	Subject       string     `json:"subject"`
	NotBefore     time.Time  `json:"not_before"`
	NotAfter      time.Time  `json:"not_after"`
	ExpiresInDays int        `json:"expires_in_days"`
	SerialNumber  string     `json:"serial_number"`
	Error         string     `json:"error,omitempty"`
	SANs          []string   `json:"sans,omitempty"` // Subject Alternative Names
}

// Print implements CheckResult interface
func (r *SSLCertResult) Print(outputFormat string) {
	switch outputFormat {
	case "json":
		// For JSON output, create a struct with formatted dates
		output := struct {
			Domain        string   `json:"domain"`
			Valid         bool     `json:"valid"`
			Issuer        string   `json:"issuer"`
			Subject       string   `json:"subject"`
			NotBefore     string   `json:"not_before"`
			NotAfter      string   `json:"not_after"`
			ExpiresInDays int      `json:"expires_in_days"`
			SerialNumber  string   `json:"serial_number"`
			Error         string   `json:"error,omitempty"`
			SANs          []string `json:"sans,omitempty"`
		}{
			Domain:        r.Domain,
			Valid:         r.Valid,
			Issuer:        r.Issuer,
			Subject:       r.Subject,
			NotBefore:     checks.FormatTimeForOutput(&r.NotBefore),
			NotAfter:      checks.FormatTimeForOutput(&r.NotAfter),
			ExpiresInDays: r.ExpiresInDays,
			SerialNumber:  r.SerialNumber,
			Error:         r.Error,
			SANs:          r.SANs,
		}
		jsonData, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			fmt.Printf("Error formatting JSON: %v\n", err)
			return
		}
		fmt.Println(string(jsonData))
	case "yaml":
		// For YAML output, use the same formatted date structure
		output := struct {
			Domain        string   `yaml:"domain"`
			Valid         bool     `yaml:"valid"`
			Issuer        string   `yaml:"issuer"`
			Subject       string   `yaml:"subject"`
			NotBefore     string   `yaml:"not_before"`
			NotAfter      string   `yaml:"not_after"`
			ExpiresInDays int      `yaml:"expires_in_days"`
			SerialNumber  string   `yaml:"serial_number"`
			Error         string   `yaml:"error,omitempty"`
			SANs          []string `yaml:"sans,omitempty"`
		}{
			Domain:        r.Domain,
			Valid:         r.Valid,
			Issuer:        r.Issuer,
			Subject:       r.Subject,
			NotBefore:     checks.FormatTimeForOutput(&r.NotBefore),
			NotAfter:      checks.FormatTimeForOutput(&r.NotAfter),
			ExpiresInDays: r.ExpiresInDays,
			SerialNumber:  r.SerialNumber,
			Error:         r.Error,
			SANs:          r.SANs,
		}
		yamlData, err := yaml.Marshal(output)
		if err != nil {
			fmt.Printf("Error formatting YAML: %v\n", err)
			return
		}
		fmt.Println(string(yamlData))
	default:
		fmt.Printf("Domain: %s\n", r.Domain)
		fmt.Printf("Valid: %v\n", r.Valid)
		fmt.Printf("Issuer: %s\n", r.Issuer)
		fmt.Printf("Subject: %s\n", r.Subject)
		fmt.Printf("Valid From: %s\n", checks.FormatTime(&r.NotBefore))
		fmt.Printf("Valid Until: %s\n", checks.FormatTime(&r.NotAfter))
		fmt.Printf("Expires In: %d days\n", r.ExpiresInDays)
		fmt.Printf("Serial Number: %s\n", r.SerialNumber)
		if len(r.SANs) > 0 {
			fmt.Println("Subject Alternative Names:")
			for _, san := range r.SANs {
				fmt.Printf("  - %s\n", san)
			}
		}
		if r.Error != "" {
			fmt.Printf("Error: %s\n", r.Error)
		}
	}
}

// SSLCheckCommand implements the CheckCommand interface for SSL certificate checks
type SSLCheckCommand struct {
	*checks.BaseCheckCommand
	showChain bool
}

// NewSSLCheckCommand creates a new SSL certificate check command
func NewSSLCheckCommand() *SSLCheckCommand {
	return &SSLCheckCommand{
		BaseCheckCommand: checks.NewBaseCheckCommand(
			"ssl",
			"Check SSL certificate for a domain",
			"ssl [domain]",
			1,
		),
	}
}

// Execute implements the CheckCommand interface
func (c *SSLCheckCommand) Execute(args []string) (checks.CheckResult, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("domain is required")
	}
	domain := args[0]

	// Add default port if not specified
	if _, _, err := net.SplitHostPort(domain); err != nil {
		domain = net.JoinHostPort(domain, "443")
	}

	// Create a connection with timeout
	conn, err := tls.DialWithDialer(
		&net.Dialer{Timeout: 10 * time.Second},
		"tcp",
		domain,
		&tls.Config{
			InsecureSkipVerify: false, // Enable certificate verification
		},
	)
	if err != nil {
		return &SSLCertResult{
			Domain: domain,
			Valid:  false,
			Error:  err.Error(),
		}, nil
	}
	defer conn.Close()

	// Get the certificate chain
	state := conn.ConnectionState()
	if len(state.PeerCertificates) == 0 {
		return &SSLCertResult{
			Domain: domain,
			Valid:  false,
			Error:  "no certificates found",
		}, nil
	}

	// Get the leaf certificate
	cert := state.PeerCertificates[0]
	now := time.Now()
	expiresIn := int(cert.NotAfter.Sub(now).Hours() / 24)

	result := &SSLCertResult{
		Domain:        domain,
		Valid:         now.After(cert.NotBefore) && now.Before(cert.NotAfter),
		Issuer:        cert.Issuer.String(),
		Subject:       cert.Subject.String(),
		NotBefore:     cert.NotBefore,
		NotAfter:      cert.NotAfter,
		ExpiresInDays: expiresIn,
		SerialNumber:  cert.SerialNumber.String(),
	}

	// Extract SANs from the certificate
	for _, dnsName := range cert.DNSNames {
		result.SANs = append(result.SANs, dnsName)
	}

	return result, nil
}

// CheckSSLCert checks the SSL certificate for a given domain
func CheckSSLCert(domain string, jsonOutput bool) {
	result, err := checkSSLCert(domain)
	if err != nil {
		fmt.Printf("Error checking SSL certificate: %v\n", err)
		os.Exit(1)
	}

	if jsonOutput {
		jsonData, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(jsonData))
	} else {
		printSSLCertResult(result)
	}
}

// checkSSLCert performs the actual SSL certificate check
func checkSSLCert(domain string) (*SSLCertResult, error) {
	result := &SSLCertResult{
		Domain: domain,
	}

	// Add default port if not specified
	if _, _, err := net.SplitHostPort(domain); err != nil {
		domain = net.JoinHostPort(domain, "443")
	}

	// Create a connection with timeout
	conn, err := tls.DialWithDialer(
		&net.Dialer{Timeout: 10 * time.Second},
		"tcp",
		domain,
		&tls.Config{
			InsecureSkipVerify: false, // Enable certificate verification
		},
	)
	if err != nil {
		result.Error = err.Error()
		return result, err
	}
	defer conn.Close()

	// Get the certificate chain
	state := conn.ConnectionState()
	if len(state.PeerCertificates) == 0 {
		result.Error = "no certificates found"
		return result, fmt.Errorf("no certificates found")
	}

	// Get the leaf certificate
	cert := state.PeerCertificates[0]

	// Fill in the result details
	result.Issuer = cert.Issuer.String()
	result.Subject = cert.Subject.String()
	result.NotBefore = cert.NotBefore
	result.NotAfter = cert.NotAfter
	result.ExpiresInDays = int(cert.NotAfter.Sub(time.Now()).Hours() / 24)
	result.SerialNumber = cert.SerialNumber.String()

	// Extract SANs from the certificate
	for _, dnsName := range cert.DNSNames {
		result.SANs = append(result.SANs, dnsName)
	}

	// Check if the certificate is valid based on current time
	now := time.Now()
	if now.Before(cert.NotBefore) || now.After(cert.NotAfter) {
		result.Valid = false
		result.Error = "certificate is not valid for current date"
	} else {
		result.Valid = true
	}

	return result, nil
}

// printSSLCertResult prints the SSL certificate result in a human-readable format
func printSSLCertResult(result *SSLCertResult) {
	fmt.Printf("SSL Certificate Check for %s:\n", result.Domain)
	fmt.Printf("Status: %s\n", getStatusString(result.Valid))
	fmt.Printf("Issuer: %s\n", result.Issuer)
	fmt.Printf("Subject: %s\n", result.Subject)
	fmt.Printf("Valid From: %s\n", checks.FormatTime(&result.NotBefore))
	fmt.Printf("Valid Until: %s\n", checks.FormatTime(&result.NotAfter))
	fmt.Printf("Expires In: %d days\n", result.ExpiresInDays)
	fmt.Printf("Serial Number: %s\n", result.SerialNumber)
	if len(result.SANs) > 0 {
		fmt.Println("Subject Alternative Names:")
		for _, san := range result.SANs {
			fmt.Printf("  - %s\n", san)
		}
	}
	if result.Error != "" {
		fmt.Printf("Error: %s\n", result.Error)
	}
}

// getStatusString returns a human-readable status string
func getStatusString(valid bool) string {
	if valid {
		return "Valid"
	}
	return "Invalid"
} 