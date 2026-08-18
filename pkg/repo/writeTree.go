package repo

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	u "github.com/devxdh/gitingo/pkg/utils"
)

func WriteTree(targetDir string) (string, error) {
	if exists, err := u.DotGitExists(targetDir); err != nil || !exists {
		return "", fmt.Errorf("fatal: not a git repository (or any of the parent directories): .git")
	}

	return writeTreeRecursive(targetDir, targetDir)
}

func writeTreeRecursive(repoRoot string, currentDir string) (string, error) {
	dirEntries, err := os.ReadDir(currentDir)
	if err != nil {
		return "", fmt.Errorf("failed to read directory %s: %w", currentDir, err)
	}

	var entries []TreeEntry

	for _, entry := range dirEntries {
		name := entry.Name()
		if name == ".git" {
			continue
		}

		fullPath := filepath.Join(currentDir, name)
		info, err := entry.Info()
		if err != nil {
			return "", fmt.Errorf("failed to get info for %s: %w", fullPath, err)
		}

		if entry.IsDir() {
			subTreeHash, err := writeTreeRecursive(repoRoot, fullPath)
			if err != nil {
				return "", err
			}

			entries = append(entries, TreeEntry{
				Mode: "40000",
				Name: name,
				Hash: subTreeHash,
			})
		} else {
			blobHash, err := HashObject(repoRoot, fullPath, true)
			if err != nil {
				return "", err
			}

			mode := "100644"
			if info.Mode()&0o111 != 0 {
				mode = "100755"
			}

			entries = append(entries, TreeEntry{
				Mode: mode,
				Name: name,
				Hash: blobHash,
			})
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name < entries[j].Name
	})

	var treePayload []byte
	for _, entry := range entries {
		rawSha, err := hex.DecodeString(entry.Hash)
		if err != nil {
			return "", fmt.Errorf("invalid Hash %s: %w", entry.Hash, err)
		}

		entryHeader := fmt.Sprintf("%s %s\x00", entry.Mode, entry.Name)
		treePayload = append(treePayload, []byte(entryHeader)...)
		treePayload = append(treePayload, rawSha...)
	}

	header := fmt.Sprintf("tree %d\x00", len(treePayload))
	fullPayload := append([]byte(header), treePayload...)

	hasher := sha1.New()
	hasher.Write(fullPayload)
	treeHash := fmt.Sprintf("%x", hasher.Sum(nil))

	if err := writeObject(repoRoot, treeHash, fullPayload); err != nil {
		return "", err
	}

	return treeHash, nil
}
