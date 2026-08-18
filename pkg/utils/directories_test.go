package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDirExists(t *testing.T) {
	tempDir := t.TempDir()

	// Test existing directory
	exists, err := DirExists(tempDir)
	if err != nil || !exists {
		t.Fatalf("expected directory to exist, got exists=%v, err=%v", exists, err)
	}

	// Test non-existing directory
	nonExistent := filepath.Join(tempDir, "does-not-exist")
	exists, err = DirExists(nonExistent)
	if err != nil || exists {
		t.Fatalf("expected non-existing directory to return false, got exists=%v, err=%v", exists, err)
	}

	// Test path that is a file, not a directory
	filePath := filepath.Join(tempDir, "file.txt")
	if err := os.WriteFile(filePath, []byte("hello"), 0o644); err != nil {
		t.Fatalf("failed to create dummy file: %v", err)
	}

	exists, err = DirExists(filePath)
	if err == nil {
		t.Fatalf("expected error for file path, got nil")
	}
	if exists {
		t.Fatalf("expected exists=false for file path")
	}
}

func TestDotGitExists(t *testing.T) {
	tempDir := t.TempDir()

	// Without .git
	exists, err := DotGitExists(tempDir)
	if exists || err == nil {
		t.Fatalf("expected DotGitExists to fail when .git is missing")
	}

	// Create .git directory
	gitDir := filepath.Join(tempDir, ".git")
	if err := os.Mkdir(gitDir, 0o755); err != nil {
		t.Fatalf("failed to create .git dir: %v", err)
	}

	exists, err = DotGitExists(tempDir)
	if err != nil || !exists {
		t.Fatalf("expected DotGitExists to succeed when .git exists, got exists=%v, err=%v", exists, err)
	}
}
