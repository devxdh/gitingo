package repo

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	u "github.com/devxdh/gitingo/pkg/utils"
)

func Checkout(targetDir string, target string, createBranch bool) error {
	if exists, err := u.DotGitExists(targetDir); err != nil || !exists {
		return fmt.Errorf("fatal: not a git repository (or any of the parent directories): .git")
	}

	if createBranch {
		if err := CreateBranch(targetDir, target); err != nil {
			return err
		}
	}

	var commitHash string
	isBranch := false
	branchRefPath := filepath.Join(targetDir, ".git", "refs", "heads", target)

	if refBytes, err := os.ReadFile(branchRefPath); err == nil {
		commitHash = strings.TrimSpace(string(refBytes))
		isBranch = true
	} else if len(target) == 40 {
		commitHash = target
	} else {
		return fmt.Errorf("error: pathspec '%s' did not match any file(s) known to git", target)
	}

	objType, commitContent, err := ReadObject(targetDir, commitHash)
	if err != nil {
		return fmt.Errorf("failed to read commit %s: %w", commitHash, err)
	}
	if objType != "commit" {
		return fmt.Errorf("object %s is not a commit", commitHash)
	}

	commit, err := ParseCommit(commitHash, commitContent)
	if err != nil {
		return fmt.Errorf("failed to parse commit %s: %w", commitHash, err)
	}

	if err := restoreTree(targetDir, targetDir, commit.Tree); err != nil {
		return fmt.Errorf("failed to restore tree: %w", err)
	}

	headPath := filepath.Join(targetDir, ".git", "HEAD")
	if isBranch {
		if err := os.WriteFile(headPath, []byte(fmt.Sprintf("ref: refs/heads/%s\n", target)), 0o644); err != nil {
			return fmt.Errorf("failed to update HEAD: %w", err)
		}
		if createBranch {
			fmt.Printf("Switched to a new branch '%s'\n", target)
		} else {
			fmt.Printf("Switched to branch '%s'\n", target)
		}
	} else {
		if err := os.WriteFile(headPath, []byte(commitHash+"\n"), 0o644); err != nil {
			return fmt.Errorf("failed to update HEAD: %w", err)
		}
		fmt.Printf("Note: switching to '%s' (detached HEAD)\n", commitHash[:7])
	}

	return nil
}

func restoreTree(repoRoot string, currentDir string, treeHash string) error {
	entries, err := LsTree(repoRoot, treeHash)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		destPath := filepath.Join(currentDir, entry.Name)

		if entry.Mode == "040000" || entry.Mode == "40000" {
			if err := os.MkdirAll(destPath, 0o755); err != nil {
				return err
			}
			if err := restoreTree(repoRoot, destPath, entry.Hash); err != nil {
				return err
			}
		} else {
			_, blobContent, err := ReadObject(repoRoot, entry.Hash)
			if err != nil {
				return err
			}

			perm := os.FileMode(0o644)
			if entry.Mode == "100755" {
				perm = 0o755
			}

			if err := os.WriteFile(destPath, blobContent, perm); err != nil {
				return err
			}
		}
	}

	return nil
}
