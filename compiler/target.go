package compiler

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"

	entrypoint "github.com/rabbytesoftware/quiver.compiler/shared/base/entrypoint"
	pc "github.com/rabbytesoftware/quiver.compiler/shared/base/package-config"
)

// Compiler handles compiling Go projects for multiple platforms
type Compiler struct {
	SourceDir     string
	OutputFile    string
	TempDir       string
	Targets       []entrypoint.Target
	FastMode      bool
	PackageConfig *pc.PackageConfig
}

// NewCompiler creates a new compiler instance
func NewCompiler(sourceDir, outputFile string, fastMode bool) *Compiler {
	var targets []entrypoint.Target
	
	if fastMode {
		// In fast mode, only compile for the current platform
		var name string
		if runtime.GOOS == "windows" {
			name = "win-" + runtime.GOARCH + ".exe"
		} else {
			name = runtime.GOOS + "-" + runtime.GOARCH
		}
		
		targets = []entrypoint.Target{
			{OS: runtime.GOOS, Arch: runtime.GOARCH, Name: name},
		}
	} else {
		targets = entrypoint.DefaultTargets
	}
	
	return &Compiler{
		SourceDir:  sourceDir,
		OutputFile: outputFile,
		Targets:    targets,
		TempDir:    filepath.Join(os.TempDir(), "watcher-compiler"),
		FastMode:   fastMode,
	}
}

// Run executes the compilation and packaging process
func (c *Compiler) Run() error {
	// Load and validate package configuration
	packageConfig, err := LoadPackageConfig(c.SourceDir)
	if err != nil {
		return fmt.Errorf("failed to load package configuration: %w", err)
	}
	
	// Store the package config
	c.PackageConfig = packageConfig
	
	// Generate build number
	packageConfig.GenerateBuildNumber()
	fmt.Printf("Generated build number: %s\n", packageConfig.BuildNumber)
	
	// Create temp directory
	err = os.MkdirAll(c.TempDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(c.TempDir)
	
	// Find main.go in source directory
	mainFile := ""
	err = filepath.Walk(c.SourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Base(path) == "main.go" {
			mainFile = path
			return filepath.SkipAll
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error searching for main.go: %w", err)
	}
	
	if mainFile == "" {
		return fmt.Errorf("main.go not found in source directory")
	}
	
	// Compile for all targets
	for _, target := range c.Targets {
		outputPath := filepath.Join(c.TempDir, target.Name)
		err = c.compileForTarget(mainFile, outputPath, target)
		if err != nil {
			return fmt.Errorf("failed to compile for %s/%s: %w", target.OS, target.Arch, err)
		}
	}
	
	// Copy public folder if it exists
	publicDir := filepath.Join(c.SourceDir, "public")
	if _, err := os.Stat(publicDir); !os.IsNotExist(err) {
		destPublicDir := filepath.Join(c.TempDir, "public")
		err = copyDir(publicDir, destPublicDir)
		if err != nil {
			return fmt.Errorf("failed to copy public directory: %w", err)
		}
	}
	
	// Save the updated package config to the temp directory
	configPath := filepath.Join(c.TempDir, "package.json")
	err = c.PackageConfig.SavePackageConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to save updated package config: %w", err)
	}
	
	// Archive everything
	err = createWatcherArchive(c.TempDir, c.OutputFile)
	if err != nil {
		return fmt.Errorf("failed to create .watcher archive: %w", err)
	}
	
	return nil
}

// compileForTarget compiles the Go application for a specific target
func (c *Compiler) compileForTarget(mainFile, outputPath string, target entrypoint.Target) error {
	fmt.Printf("Compiling for %s/%s...\n", target.OS, target.Arch)
	
	// Create a build script
	scriptContent, err := c.createBuildScript(mainFile, outputPath, target)
	if err != nil {
		return fmt.Errorf("failed to create build script: %w", err)
	}
	
	// Write the script to a temporary file
	scriptExt := ".sh"
	if runtime.GOOS == "windows" {
		scriptExt = ".bat"
	}
	
	scriptPath := filepath.Join(c.TempDir, "build"+scriptExt)
	err = os.WriteFile(scriptPath, []byte(scriptContent), 0755)
	if err != nil {
		return fmt.Errorf("failed to write build script: %w", err)
	}
	
	// Execute the build script
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", scriptPath)
	} else {
		cmd = exec.Command("/bin/sh", scriptPath)
	}
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("compilation failed: %w\n%s", err, output)
	}
	
	return nil
}

// createBuildScript generates a shell/batch script to build the project
func (c *Compiler) createBuildScript(mainFile, outputPath string, target entrypoint.Target) (string, error) {
	// Get the directory of the main file
	mainDir := filepath.Dir(mainFile)
	
	// Template for the build script
	var scriptTemplate string
	if runtime.GOOS == "windows" {
		scriptTemplate = `@echo off
cd "{{.MainDir}}"
set GOOS={{.Target.OS}}
set GOARCH={{.Target.Arch}}
go build -o "{{.OutputPath}}" .
`
	} else {
		scriptTemplate = `#!/bin/sh
cd "{{.MainDir}}"
export GOOS={{.Target.OS}}
export GOARCH={{.Target.Arch}}
go build -o "{{.OutputPath}}" .
`
	}
	
	// Create template data
	data := struct {
		MainDir    string
		OutputPath string
		Target     entrypoint.Target
	}{
		MainDir:    mainDir,
		OutputPath: outputPath,
		Target:     target,
	}
	
	// Parse and execute the template
	tmpl, err := template.New("build").Parse(scriptTemplate)
	if err != nil {
		return "", err
	}
	
	var scriptBuilder strings.Builder
	err = tmpl.Execute(&scriptBuilder, data)
	if err != nil {
		return "", err
	}
	
	return scriptBuilder.String(), nil
}
