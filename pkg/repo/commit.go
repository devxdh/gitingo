package repo

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	u "github.com/devxdh/gitingo/pkg/utils"
)

func Commit(targetDir string, message string) (string, error) {
	if exists, err := u.DotGitExists(targetDir); err != nil || !exists {
		return "", fmt.Errorf("fatal: not a git repository (or any of the parent directories): .git")
	}

	treeHash, err := WriteTree(targetDir)
	if err != nil {
		return "", fmt.Errorf("failed to write tree: %w", err)
	}

	headPath := filepath.Join(targetDir, ".git", "HEAD")
	headContent, err := os.ReadFile(headPath)
	if err != nil {
		return "", fmt.Errorf("failed to read HEAD: %w", err)
	}

	headStr := strings.TrimSpace(string(headContent))
	if !strings.HasPrefix(headStr, "ref: ") {
		return "", fmt.Errorf("detached HEAD or invalid HEAD format: %s", headStr)
	}

	refRelativePath := strings.TrimPrefix(headStr, "ref: ")
	refFullPath := filepath.Join(targetDir, ".git", refRelativePath)

	var parentHash string
	if parentBytes, err := os.ReadFile(refFullPath); err == nil {
		parentHash = strings.TrimSpace(string(parentBytes))
	}

	commitHash, err := CommitTree(targetDir, treeHash, parentHash, message)
	if err != nil {
		return "", fmt.Errorf("failed to commit tree: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(refFullPath), 0o755); err != nil {
		return "", fmt.Errorf("failed to create refs directory: %w", err)
	}

	if err := os.WriteFile(refFullPath, []byte(commitHash+"\n"), 0o644); err != nil {
		return "", fmt.Errorf("failed to update ref %s: %w", refRelativePath, err)
	}

	branchName := filepath.Base(refRelativePath)
	isRoot := parentHash == ""
	if isRoot {
		fmt.Printf("[%s (root-commit) %s] %s\n", branchName, commitHash[:7], strings.TrimSpace(message))
	} else {
		fmt.Printf("[%s %s] %s\n", branchName, commitHash[:7], strings.TrimSpace(message))
	}

	return commitHash, nil
}
