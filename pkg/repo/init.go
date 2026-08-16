// Package repo contains all the core Version Control System operations.
// For example := initializing the repository, making a commit etc etc
package repo

import (
	"fmt"
	"os"
	"path/filepath"

	u "github.com/devxdh/gitingo/pkg/utils"
)

func Init(targetDir string) error {
	cleanedTarget := filepath.Clean(targetDir)

	exists, _ := u.DotGitExists(cleanedTarget)
	if exists {
		fmt.Printf("%s is already a repository\n", cleanedTarget)
		return nil
	}

	dirs := []string{
		filepath.Join(cleanedTarget, ".git", "objects"),
		filepath.Join(cleanedTarget, ".git", "refs", "heads"),
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0o755); err != nil {
			return fmt.Errorf("failed to create git directories %s in %s: %v", d, cleanedTarget, err)
		}
	}

	headPath := filepath.Join(cleanedTarget, ".git", "HEAD")
	err := os.WriteFile(headPath, []byte("ref: refs/heads/main\n"), 0o644)
	if err != nil {
		return fmt.Errorf("failed to write HEAD into %s: %v", headPath, err)
	}

	fmt.Printf("Initialized Gitingo repository in %s\n", filepath.Join(cleanedTarget, ".git"))
	return nil
}
