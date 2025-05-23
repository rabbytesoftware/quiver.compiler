package packageconfig

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"time"
)

// GenerateBuildNumber generates a random build number and adds it to the config
func (c *PackageConfig) GenerateBuildNumber() {
	// Generate a random number between 10000 and 99999
	max := big.NewInt(90000)
	randInt, _ := rand.Int(rand.Reader, max)
	randNum := randInt.Int64() + 10000
	
	// Generate timestamp part (last 6 digits of current Unix timestamp)
	timestamp := time.Now().Unix() % 1000000
	
	// Combine them for a unique number
	buildNum := randNum*1000000 + timestamp
	
	// Convert to string
	c.BuildNumber = fmt.Sprintf("%d", buildNum)
}

// SavePackageConfig saves the package configuration to the given file path
func (c *PackageConfig) SavePackageConfig(filePath string) error {
	// Convert to JSON with indentation
	data, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	
	// Write to file
	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	
	return nil
}