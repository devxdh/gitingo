// Package utils contains utility functions which servers other packages
package utils

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func DirExists(dirPath string) (bool, error) {
	info, err := os.Stat(dirPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}

		return false, err
	}

	if !info.IsDir() {
		return false, fmt.Errorf("path %s is a file not a folder", dirPath)
	}
	return true, nil
}

func DotGitExists(dirPath string) (bool, error) {
	dotGitPath := filepath.Join(dirPath, ".git")

	exists, err := DirExists(dotGitPath)
	if err != nil || !exists {
		return false, fmt.Errorf("fatal: not a git repository: %w", err)
	}

	return exists, nil
}
