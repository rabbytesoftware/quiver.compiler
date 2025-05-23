package compiler

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// createWatcherArchive creates a .watcher archive from the source directory
func createWatcherArchive(sourceDir, outputPath string) error {
	// Create output file
	out, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer out.Close()
	
	// Create gzip writer
	gw := gzip.NewWriter(out)
	defer gw.Close()
	
	// Create tar writer
	tw := tar.NewWriter(gw)
	defer tw.Close()
	
	// Walk through the source directory and add files to the archive
	return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		// Get relative path
		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}
		
		// Skip the root directory
		if relPath == "." {
			return nil
		}
		
		// Create header
		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return fmt.Errorf("failed to create header: %w", err)
		}
		header.Name = relPath
		
		// Write header
		if err := tw.WriteHeader(header); err != nil {
			return fmt.Errorf("failed to write header: %w", err)
		}
		
		// If it's a file, write content
		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("failed to open file: %w", err)
			}
			defer file.Close()
			
			if _, err := io.Copy(tw, file); err != nil {
				return fmt.Errorf("failed to copy file content: %w", err)
			}
		}
		
		return nil
	})
}

// copyDir recursively copies a directory structure
func copyDir(src, dst string) error {
	// Get file info
	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	
	// Create destination directory
	if err = os.MkdirAll(dst, info.Mode()); err != nil {
		return err
	}
	
	// Read directory entries
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())
		
		if entry.IsDir() {
			// Recurse into subdirectories
			if err = copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			// Copy files
			if err = copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}
	
	return nil
}

// copyFile copies a single file
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	
	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	
	// Make sure executable permissions are preserved
	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	return os.Chmod(dst, info.Mode())
}