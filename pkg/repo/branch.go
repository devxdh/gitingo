package repo

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	u "github.com/devxdh/gitingo/pkg/utils"
)

type BranchInfo struct {
	Name     string
	IsActive bool
}

func ListBranches(targetDir string) ([]BranchInfo, error) {
	if exists, err := u.DotGitExists(targetDir); err != nil || !exists {
		return nil, fmt.Errorf("fatal: not a git repository (or any of the parent directories): .git")
	}

	headContent, err := os.ReadFile(filepath.Join(targetDir, ".git", "HEAD"))
	if err != nil {
		return nil, fmt.Errorf("failed to read HEAD: %w", err)
	}

	headStr := strings.TrimSpace(string(headContent))
	activeBranch := ""
	if strings.HasPrefix(headStr, "ref: refs/heads/") {
		activeBranch = strings.TrimPrefix(headStr, "ref: refs/heads/")
	}

	headsDir := filepath.Join(targetDir, ".git", "refs", "heads")
	entries, err := os.ReadDir(headsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read refs/heads: %w", err)
	}

	var branches []BranchInfo
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		branches = append(branches, BranchInfo{
			Name:     name,
			IsActive: name == activeBranch,
		})
	}

	return branches, nil
}

func CreateBranch(targetDir string, branchName string) error {
	if exists, err := u.DotGitExists(targetDir); err != nil || !exists {
		return fmt.Errorf("fatal: not a git repository (or any of the parent directories): .git")
	}

	if branchName == "" || strings.Contains(branchName, " ") || strings.Contains(branchName, "..") {
		return fmt.Errorf("fatal: '%s' is not a valid branch name", branchName)
	}

	currentCommit, err := GetHeadCommitHash(targetDir)
	if err != nil {
		return fmt.Errorf("fatal: not a valid object name: 'HEAD'")
	}

	branchPath := filepath.Join(targetDir, ".git", "refs", "heads", branchName)

	if _, err := os.Stat(branchPath); err == nil {
		return fmt.Errorf("fatal: a branch named '%s' already exists", branchName)
	}

	if err := os.MkdirAll(filepath.Dir(branchPath), 0o755); err != nil {
		return fmt.Errorf("failed to create directory for ref: %w", err)
	}

	if err := os.WriteFile(branchPath, []byte(currentCommit+"\n"), 0o644); err != nil {
		return fmt.Errorf("failed to write branch ref: %w", err)
	}

	return nil
}

func DeleteBranch(targetDir string, branchName string) error {
	if exists, err := u.DotGitExists(targetDir); err != nil || !exists {
		return fmt.Errorf("fatal: not a git repository (or any of the parent directories): .git")
	}

	headContent, err := os.ReadFile(filepath.Join(targetDir, ".git", "HEAD"))
	if err != nil {
		return fmt.Errorf("failed to read HEAD: %w", err)
	}

	headStr := strings.TrimSpace(string(headContent))
	if headStr == fmt.Sprintf("ref: refs/heads/%s", branchName) {
		return fmt.Errorf("fatal: cannot delete branch '%s' checked out at '%s'", branchName, targetDir)
	}

	branchPath := filepath.Join(targetDir, ".git", "refs", "heads", branchName)
	if _, err := os.Stat(branchPath); os.IsNotExist(err) {
		return fmt.Errorf("error: branch '%s' not found", branchName)
	}

	if err := os.Remove(branchPath); err != nil {
		return fmt.Errorf("failed to delete branch ref: %w", err)
	}

	fmt.Printf("Deleted branch %s\n", branchName)
	return nil
}
