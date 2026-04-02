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

var ignoredExtensions = map[string]bool{
	".exe": true,
	".bin": true,
	".out": true,
	".png": true,
	".jpg": true,
	".gif": true,
	".mp3": true,
	".mp4": true,
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

		ext := strings.ToLower(filepath.Ext(path))
		if ignoredExtensions[ext] {
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
