package compiler

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	pc "github.com/rabbytesoftware/quiver.compiler/shared/base/package-config"
)

func LoadPackageConfig(sourceDir string) (*pc.PackageConfig, error) {
	// Find any .json file in the source directory
	var configPath string
	
	entries, err := os.ReadDir(sourceDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read source directory: %w", err)
	}
	
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".json" {
			configPath = filepath.Join(sourceDir, entry.Name())
			break
		}
	}
	
	if configPath == "" {
		return nil, fmt.Errorf("no JSON configuration file found in source directory")
	}
	
	// Read the config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	
	// Parse the JSON
	var config pc.PackageConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}
	
	// Validate required fields
	if config.URL == "" {
		return nil, fmt.Errorf("missing required field 'url' in config file")
	}
	if config.Name == "" {
		return nil, fmt.Errorf("missing required field 'name' in config file")
	}
	if config.Version == "" {
		return nil, fmt.Errorf("missing required field 'version' in config file")
	}
	if len(config.Maintainers) == 0 {
		return nil, fmt.Errorf("missing required field 'maintainers' in config file")
	}
	
	return &config, nil
}