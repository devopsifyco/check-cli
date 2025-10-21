package dependencies

import (
	"encoding/xml"
	"fmt"
	"os"
	"strings"
)

// MavenProject represents the root element of a pom.xml file.
type MavenProject struct {
	XMLName      xml.Name          `xml:"project"`
	Parent       *Parent           `xml:"parent"`
	GroupID      string           `xml:"groupId"`
	ArtifactID   string           `xml:"artifactId"`
	Version      string           `xml:"version"`
	Properties   Properties        `xml:"properties"`
	Dependencies []MavenDependency `xml:"dependencies>dependency"`
	DependencyManagement *DependencyManagement `xml:"dependencyManagement"`
}

// Parent represents the parent element in a pom.xml file
type Parent struct {
	GroupID    string `xml:"groupId"`
	ArtifactID string `xml:"artifactId"`
	Version    string `xml:"version"`
}

// DependencyManagement represents the dependencyManagement section
type DependencyManagement struct {
	Dependencies []MavenDependency `xml:"dependencies>dependency"`
}

// Properties represents Maven properties
type Properties struct {
	List []Property `xml:",any"`
}

// Property represents a single Maven property
type Property struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

// MavenDependency represents a single dependency in a pom.xml file.
type MavenDependency struct {
	GroupID    string `xml:"groupId"`
	ArtifactID string `xml:"artifactId"`
	Version    string `xml:"version"`
}

// resolveVersion resolves a version string that might contain property references
func resolveVersion(version string, project *MavenProject) string {
	// If the version doesn't contain property reference, return as is
	if !strings.Contains(version, "${") {
		return version
	}

	// Extract property name from ${property.name} format
	propName := strings.TrimPrefix(strings.TrimSuffix(version, "}"), "${")

	// Look for the property in the properties section
	for _, prop := range project.Properties.List {
		if prop.XMLName.Local == propName {
			return prop.Value
		}
	}

	// If property not found, return original version
	fmt.Printf("Warning: Property %s not found, using raw version string\n", propName)
	return version
}

// findVersionInDependencyManagement looks for a dependency version in the dependencyManagement section
func findVersionInDependencyManagement(groupID, artifactID string, project *MavenProject) string {
	if project.DependencyManagement == nil {
		return ""
	}

	for _, dep := range project.DependencyManagement.Dependencies {
		if dep.GroupID == groupID && dep.ArtifactID == artifactID {
			return resolveVersion(dep.Version, project)
		}
	}

	return ""
}

// findVersionFromParent looks for a dependency version based on parent groupId
func findVersionFromParent(dep MavenDependency, project *MavenProject) string {
	// Check if the dependency belongs to the parent's group
	if project.Parent != nil && strings.HasPrefix(dep.GroupID, project.Parent.GroupID) {
		return project.Parent.Version
	}

	// Check if the dependency belongs to the project's group
	if project.GroupID != "" && strings.HasPrefix(dep.GroupID, project.GroupID) {
		return project.Version
	}

	return ""
}

// parsePomXML parses a pom.xml file and extracts dependencies.
func parsePomXML(filePath string) ([]Dependency, error) {
	xmlFile, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening pom.xml %s: %w", filePath, err)
	}
	defer xmlFile.Close()

	var project MavenProject
	if err := xml.NewDecoder(xmlFile).Decode(&project); err != nil {
		return nil, fmt.Errorf("error decoding pom.xml %s: %w", filePath, err)
	}

	deps := make([]Dependency, 0, len(project.Dependencies))
	for _, mvnDep := range project.Dependencies {
		// Try to resolve version in this order:
		// 1. Direct version
		// 2. Version from dependencyManagement
		// 3. Version from parent/project groupId
		version := resolveVersion(mvnDep.Version, &project)
		if version == "" {
			version = findVersionInDependencyManagement(mvnDep.GroupID, mvnDep.ArtifactID, &project)
		}
		if version == "" {
			version = findVersionFromParent(mvnDep, &project)
		}
		
		if version == "" {
			// Skip dependencies without a version after all resolution attempts
			fmt.Printf("Skipping dependency %s:%s with missing version in %s\n", mvnDep.GroupID, mvnDep.ArtifactID, filePath)
			continue
		}
		
		deps = append(deps, Dependency{
			Name:    fmt.Sprintf("%s", mvnDep.ArtifactID),
			Version: version,
			Manager: "maven",
		})
	}

	return deps, nil
} 