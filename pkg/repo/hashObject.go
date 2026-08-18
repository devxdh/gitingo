package repo

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"os"
	"path/filepath"

	u "github.com/devxdh/gitingo/pkg/utils"
)

func HashObject(targetDir string, filePath string, write bool) (string, error) {
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	header := fmt.Sprintf("blob %d\x00", len(fileBytes))
	payload := append([]byte(header), fileBytes...)

	hasher := sha1.New()
	if _, err = hasher.Write(payload); err != nil {
		return "", fmt.Errorf("failed to hash payload: %w", err)
	}
	hash := fmt.Sprintf("%x", hasher.Sum(nil))

	if !write {
		return hash, nil
	}

	if exists, err := u.DotGitExists(targetDir); err != nil || !exists {
		return "", err
	}

	if err := writeObject(targetDir, hash, payload); err != nil {
		return "", err
	}

	return hash, nil
}

func writeObject(targetDir string, hash string, payload []byte) error {
	dirName := hash[:2]
	fileName := hash[2:]

	objectDir := filepath.Join(targetDir, ".git", "objects", dirName)
	objectPath := filepath.Join(objectDir, fileName)

	if _, err := os.Stat(objectPath); err == nil {
		return nil
	}

	if err := os.MkdirAll(objectDir, 0o755); err != nil {
		return fmt.Errorf("failed to create object directory %s: %w", objectDir, err)
	}

	var buf bytes.Buffer
	zw := zlib.NewWriter(&buf)
	if _, err := zw.Write(payload); err != nil {
		return fmt.Errorf("failed to compress payload: %w", err)
	}

	if err := zw.Close(); err != nil {
		return fmt.Errorf("failed to close zlib writer: %w", err)
	}

	if err := os.WriteFile(objectPath, buf.Bytes(), 0o444); err != nil {
		return fmt.Errorf("failed to write object file %s: %w", objectPath, err)
	}

	return nil
}
