package packages

import (
	"fmt"
	"os"
	"path/filepath"
	"gopkg.in/yaml.v3"
)

// PackageManager manages package installation and loading
type PackageManager struct {
	packages map[string]*Package
	registry *Registry
}

// NewPackageManager creates a new package manager
func NewPackageManager() *PackageManager {
	return &PackageManager{
		packages: make(map[string]*Package),
		registry: NewRegistry(),
	}
}

// Package represents a Godin package
type Package struct {
	Name         string            `yaml:"name"`
	Version      string            `yaml:"version"`
	Description  string            `yaml:"description"`
	GitHubURL    string            `yaml:"github"`
	Dependencies map[string]string `yaml:"dependencies"`
	Widgets      map[string]Widget `yaml:"-"`
	Styles       map[string]string `yaml:"-"`
	Path         string            `yaml:"-"`
}

// PackageConfig represents the package.yaml configuration
type PackageConfig struct {
	Name         string                       `yaml:"name"`
	Version      string                       `yaml:"version"`
	Description  string                       `yaml:"description"`
	Dependencies map[string]PackageDependency `yaml:"dependencies"`
	DevDependencies map[string]PackageDependency `yaml:"dev_dependencies"`
	Scripts      map[string]string            `yaml:"scripts"`
	Config       PackageAppConfig             `yaml:"config"`
}

// PackageDependency represents a package dependency
type PackageDependency struct {
	GitHub  string `yaml:"github"`
	Version string `yaml:"version"`
}

// PackageAppConfig represents application configuration
type PackageAppConfig struct {
	Server struct {
		Port string `yaml:"port"`
		Host string `yaml:"host"`
	} `yaml:"server"`
	WebSocket struct {
		Enabled bool   `yaml:"enabled"`
		Path    string `yaml:"path"`
	} `yaml:"websocket"`
	Static struct {
		Dir   string `yaml:"dir"`
		Cache bool   `yaml:"cache"`
	} `yaml:"static"`
}

// LoadPackageConfig loads package.yaml configuration
func (pm *PackageManager) LoadPackageConfig(path string) (*PackageConfig, error) {
	configPath := filepath.Join(path, "package.yaml")
	
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read package.yaml: %w", err)
	}
	
	var config PackageConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse package.yaml: %w", err)
	}
	
	return &config, nil
}

// SavePackageConfig saves package.yaml configuration
func (pm *PackageManager) SavePackageConfig(path string, config *PackageConfig) error {
	configPath := filepath.Join(path, "package.yaml")
	
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal package.yaml: %w", err)
	}
	
	err = os.WriteFile(configPath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write package.yaml: %w", err)
	}
	
	return nil
}

// InstallPackage installs a package from GitHub
func (pm *PackageManager) InstallPackage(githubURL, version string) (*Package, error) {
	// Download package from GitHub
	packagePath, err := pm.downloadPackage(githubURL, version)
	if err != nil {
		return nil, fmt.Errorf("failed to download package: %w", err)
	}
	
	// Load package
	pkg, err := pm.LoadPackage(packagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load package: %w", err)
	}
	
	// Install dependencies
	err = pm.installDependencies(pkg)
	if err != nil {
		return nil, fmt.Errorf("failed to install dependencies: %w", err)
	}
	
	pm.packages[pkg.Name] = pkg
	return pkg, nil
}

// LoadPackage loads a package from a directory
func (pm *PackageManager) LoadPackage(path string) (*Package, error) {
	configPath := filepath.Join(path, "package.yaml")
	
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read package.yaml: %w", err)
	}
	
	var pkg Package
	err = yaml.Unmarshal(data, &pkg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse package.yaml: %w", err)
	}
	
	pkg.Path = path
	
	// Load widgets and styles
	err = pm.loadPackageAssets(&pkg)
	if err != nil {
		return nil, fmt.Errorf("failed to load package assets: %w", err)
	}
	
	return &pkg, nil
}

// downloadPackage downloads a package from GitHub
func (pm *PackageManager) downloadPackage(githubURL, version string) (string, error) {
	// TODO: Implement GitHub package downloading
	// This would:
	// 1. Clone or download the repository
	// 2. Checkout the specified version/tag
	// 3. Place it in the packages/ directory
	
	packageName := extractPackageName(githubURL)
	packagePath := filepath.Join("packages", packageName)
	
	// For now, just create the directory
	err := os.MkdirAll(packagePath, 0755)
	if err != nil {
		return "", err
	}
	
	return packagePath, nil
}

// installDependencies installs package dependencies
func (pm *PackageManager) installDependencies(pkg *Package) error {
	for depName, depVersion := range pkg.Dependencies {
		if _, exists := pm.packages[depName]; !exists {
			// Resolve dependency GitHub URL from registry
			githubURL := pm.registry.GetPackageURL(depName)
			if githubURL == "" {
				return fmt.Errorf("package %s not found in registry", depName)
			}
			
			_, err := pm.InstallPackage(githubURL, depVersion)
			if err != nil {
				return fmt.Errorf("failed to install dependency %s: %w", depName, err)
			}
		}
	}
	return nil
}

// loadPackageAssets loads widgets and styles from package
func (pm *PackageManager) loadPackageAssets(pkg *Package) error {
	// TODO: Implement asset loading
	// This would:
	// 1. Scan widgets/ directory for Go files
	// 2. Parse and register widgets
	// 3. Load CSS files from styles/ directory
	// 4. Load static assets
	
	pkg.Widgets = make(map[string]Widget)
	pkg.Styles = make(map[string]string)
	
	return nil
}

// GetPackage retrieves a loaded package by name
func (pm *PackageManager) GetPackage(name string) (*Package, bool) {
	pkg, exists := pm.packages[name]
	return pkg, exists
}

// ListPackages returns all loaded packages
func (pm *PackageManager) ListPackages() map[string]*Package {
	return pm.packages
}

// extractPackageName extracts package name from GitHub URL
func extractPackageName(githubURL string) string {
	// Simple implementation - extract from URL
	// e.g., "github.com/user/package" -> "package"
	parts := filepath.Base(githubURL)
	return parts
}

// Widget interface placeholder
type Widget interface {
	Render() string
}
