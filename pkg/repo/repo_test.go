package repo

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRepoWorkflow(t *testing.T) {
	tempDir := t.TempDir()

	// 1. Test Init
	if err := Init(tempDir); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	headContent, err := os.ReadFile(filepath.Join(tempDir, ".git", "HEAD"))
	if err != nil {
		t.Fatalf("Failed to read HEAD: %v", err)
	}
	if string(headContent) != "ref: refs/heads/main\n" {
		t.Fatalf("Unexpected HEAD content: %s", string(headContent))
	}

	// Re-init should succeed gracefully
	if err := Init(tempDir); err != nil {
		t.Fatalf("Re-init failed: %v", err)
	}

	// 2. HashObject & Write Blob
	sampleFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(sampleFile, []byte("hello gitingo\n"), 0o644); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	hash, err := HashObject(tempDir, sampleFile, true)
	if err != nil {
		t.Fatalf("HashObject failed: %v", err)
	}
	if hash == "" {
		t.Fatalf("Expected non-empty hash")
	}

	// Read Object back
	objType, content, err := ReadObject(tempDir, hash)
	if err != nil {
		t.Fatalf("ReadObject failed: %v", err)
	}
	if objType != "blob" {
		t.Fatalf("Expected objType 'blob', got '%s'", objType)
	}
	if string(content) != "hello gitingo\n" {
		t.Fatalf("Expected blob content 'hello gitingo\\n', got '%s'", string(content))
	}

	// 3. WriteTree
	treeHash, err := WriteTree(tempDir)
	if err != nil {
		t.Fatalf("WriteTree failed: %v", err)
	}
	if treeHash == "" {
		t.Fatalf("Expected non-empty tree hash")
	}

	// 4. Commit
	commitHash, err := Commit(tempDir, "Initial commit")
	if err != nil {
		t.Fatalf("Commit failed: %v", err)
	}
	if commitHash == "" {
		t.Fatalf("Expected non-empty commit hash")
	}

	// 5. Log
	commits, err := Log(tempDir)
	if err != nil {
		t.Fatalf("Log failed: %v", err)
	}
	if len(commits) != 1 {
		t.Fatalf("Expected 1 commit in log, got %d", len(commits))
	}
	if commits[0].Message != "Initial commit" {
		t.Fatalf("Expected commit message 'Initial commit', got '%s'", commits[0].Message)
	}
}
