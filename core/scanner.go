package core

import (
	"os"
	"path/filepath"
	"strings"
)

var ignoredFolders = map[string]bool{
	".git":         true,
	"node_modules": true,
	"vendor":       true,
	"dist":         true,
	"build":        true,
}

var allowedExtensions = map[string]bool{
	".go":    true,
	".java":  true,
	".py":    true,
	".rb":    true,
	".php":   true,
	".rs":    true,
	".c":     true,
	".cpp":   true,
	".cs":    true,
	".js":    true,
	".ts":    true,
	".jsx":   true,
	".tsx":   true,
	".html":  true,
	".css":   true,
	".scss":  true,
	".swift": true,
	".kt":    true,
}

func isReadme(name string) bool {
	return strings.ToLower(name) == "readme.md"
}

func ScanFiles(root string) ([]string, error) {
	var files []string

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			if ignoredFolders[d.Name()] {
				return filepath.SkipDir
			}
			return nil
		}

		name := d.Name()
		ext := strings.ToLower(filepath.Ext(name))

		if !allowedExtensions[ext] && !isReadme(name) {
			return nil
		}

		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		files = append(files, rel)
		return nil
	})

	return files, err
}
