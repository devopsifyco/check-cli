package dependencies

// Dependency represents a single project dependency.
type Dependency struct {
	Name    string   `json:"name"`
	Version string   `json:"version"`
	Manager string   `json:"manager"` // e.g., "maven", "npm", "pip", "go", "nuget"
	Tags    []string `json:"tags,omitempty"`
} 