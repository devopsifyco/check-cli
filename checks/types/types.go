package types

import (
	"crypto/tls"
	"net/http"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

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
		BaseURL: "https://api.opsify.dev/checks",
		APIKey:  apiKey,
		HTTPClient: &http.Client{
			Transport: tr,
		},
	}
}

// VersionService handles version-related operations
type VersionService struct {
	Client *APIClient
}

// NewVersionService creates a new version service
func NewVersionService(client *APIClient) *VersionService {
	return &VersionService{
		Client: client,
	}
}

// Add GetCVEs method to VersionService
func (vs *VersionService) GetCVEs(product, version string, vendor *string) ([]CVEResponse, error) {
	url := fmt.Sprintf("%s/cve?product_name=%s&version=%s", vs.Client.BaseURL, product, version)
	if vendor != nil {
		url += "&vendor=" + *vendor
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+vs.Client.APIKey)

	resp, err := vs.Client.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var cves []CVEResponse
	if err := json.Unmarshal(body, &cves); err != nil {
		return nil, err
	}
	return cves, nil
} 