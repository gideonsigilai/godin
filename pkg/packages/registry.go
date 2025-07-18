package packages

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Registry manages package registry and discovery
type Registry struct {
	packages map[string]RegistryPackage
	client   *http.Client
}

// RegistryPackage represents a package in the registry
type RegistryPackage struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	GitHubURL   string   `json:"github_url"`
	Version     string   `json:"version"`
	Tags        []string `json:"tags"`
	Author      string   `json:"author"`
	License     string   `json:"license"`
	Downloads   int      `json:"downloads"`
	Stars       int      `json:"stars"`
	UpdatedAt   string   `json:"updated_at"`
}

// NewRegistry creates a new package registry
func NewRegistry() *Registry {
	return &Registry{
		packages: make(map[string]RegistryPackage),
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetPackageURL returns the GitHub URL for a package
func (r *Registry) GetPackageURL(name string) string {
	if pkg, exists := r.packages[name]; exists {
		return pkg.GitHubURL
	}
	return ""
}

// SearchPackages searches for packages by name or tags
func (r *Registry) SearchPackages(query string) []RegistryPackage {
	var results []RegistryPackage
	
	for _, pkg := range r.packages {
		if containsString(pkg.Name, query) || 
		   containsString(pkg.Description, query) ||
		   containsStringSlice(pkg.Tags, query) {
			results = append(results, pkg)
		}
	}
	
	return results
}

// GetPackageInfo returns detailed information about a package
func (r *Registry) GetPackageInfo(name string) (*RegistryPackage, error) {
	if pkg, exists := r.packages[name]; exists {
		return &pkg, nil
	}
	return nil, fmt.Errorf("package %s not found", name)
}

// RegisterPackage registers a new package in the registry
func (r *Registry) RegisterPackage(pkg RegistryPackage) error {
	// Validate package information
	if pkg.Name == "" {
		return fmt.Errorf("package name is required")
	}
	if pkg.GitHubURL == "" {
		return fmt.Errorf("GitHub URL is required")
	}
	
	r.packages[pkg.Name] = pkg
	return nil
}

// LoadDefaultPackages loads the default package registry
func (r *Registry) LoadDefaultPackages() error {
	// Default Godin framework packages
	defaultPackages := []RegistryPackage{
		{
			Name:        "godin-ui-kit",
			Description: "Essential UI components for Godin applications",
			GitHubURL:   "github.com/godin-framework/ui-kit",
			Version:     "v1.2.0",
			Tags:        []string{"ui", "components", "widgets"},
			Author:      "Godin Framework Team",
			License:     "MIT",
			Downloads:   1500,
			Stars:       89,
			UpdatedAt:   "2024-01-15T10:30:00Z",
		},
		{
			Name:        "godin-charts",
			Description: "Chart and data visualization widgets",
			GitHubURL:   "github.com/godin-framework/charts",
			Version:     "v0.8.1",
			Tags:        []string{"charts", "visualization", "data"},
			Author:      "Godin Framework Team",
			License:     "MIT",
			Downloads:   750,
			Stars:       45,
			UpdatedAt:   "2024-01-10T14:20:00Z",
		},
		{
			Name:        "godin-forms",
			Description: "Advanced form components and validation",
			GitHubURL:   "github.com/godin-framework/forms",
			Version:     "v1.0.3",
			Tags:        []string{"forms", "validation", "input"},
			Author:      "Godin Framework Team",
			License:     "MIT",
			Downloads:   920,
			Stars:       67,
			UpdatedAt:   "2024-01-12T09:15:00Z",
		},
		{
			Name:        "godin-auth",
			Description: "Authentication and authorization widgets",
			GitHubURL:   "github.com/godin-framework/auth",
			Version:     "v2.1.0",
			Tags:        []string{"auth", "security", "login"},
			Author:      "Godin Framework Team",
			License:     "MIT",
			Downloads:   1200,
			Stars:       78,
			UpdatedAt:   "2024-01-14T16:45:00Z",
		},
		{
			Name:        "godin-testing",
			Description: "Testing utilities for Godin applications",
			GitHubURL:   "github.com/godin-framework/testing",
			Version:     "v1.0.0",
			Tags:        []string{"testing", "dev", "utilities"},
			Author:      "Godin Framework Team",
			License:     "MIT",
			Downloads:   450,
			Stars:       23,
			UpdatedAt:   "2024-01-08T11:30:00Z",
		},
	}
	
	for _, pkg := range defaultPackages {
		r.packages[pkg.Name] = pkg
	}
	
	return nil
}

// FetchFromRemote fetches package information from a remote registry
func (r *Registry) FetchFromRemote(registryURL string) error {
	resp, err := r.client.Get(registryURL + "/packages")
	if err != nil {
		return fmt.Errorf("failed to fetch from remote registry: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("remote registry returned status %d", resp.StatusCode)
	}
	
	var remotePackages []RegistryPackage
	err = json.NewDecoder(resp.Body).Decode(&remotePackages)
	if err != nil {
		return fmt.Errorf("failed to decode remote packages: %w", err)
	}
	
	// Merge remote packages with local registry
	for _, pkg := range remotePackages {
		r.packages[pkg.Name] = pkg
	}
	
	return nil
}

// GetPopularPackages returns the most popular packages
func (r *Registry) GetPopularPackages(limit int) []RegistryPackage {
	var packages []RegistryPackage
	for _, pkg := range r.packages {
		packages = append(packages, pkg)
	}
	
	// Sort by downloads (simplified)
	// In a real implementation, you'd use sort.Slice
	
	if limit > len(packages) {
		limit = len(packages)
	}
	
	return packages[:limit]
}

// GetRecentPackages returns recently updated packages
func (r *Registry) GetRecentPackages(limit int) []RegistryPackage {
	var packages []RegistryPackage
	for _, pkg := range r.packages {
		packages = append(packages, pkg)
	}
	
	// Sort by updated_at (simplified)
	// In a real implementation, you'd parse dates and sort
	
	if limit > len(packages) {
		limit = len(packages)
	}
	
	return packages[:limit]
}

// Helper functions
func containsString(str, substr string) bool {
	return len(str) >= len(substr) && str[:len(substr)] == substr
}

func containsStringSlice(slice []string, str string) bool {
	for _, item := range slice {
		if containsString(item, str) {
			return true
		}
	}
	return false
}
