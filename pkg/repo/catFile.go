package repo

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"path/filepath"

	u "github.com/devxdh/gitingo/pkg/utils"
)

func ReadObject(targetDir string, hash string) (string, []byte, error) {
	if len(hash) != 40 {
		return "", nil, fmt.Errorf("fatal: Not a valid object name %s", hash)
	}

	if exists, err := u.DotGitExists(targetDir); err != nil || !exists {
		return "", nil, fmt.Errorf("fatal: not a git repository (or any of the parent directories): .git")
	}

	dirName := hash[:2]
	fileName := hash[2:]
	objectPath := filepath.Join(targetDir, ".git", "objects", dirName, fileName)

	file, err := os.Open(objectPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil, fmt.Errorf("fatal: Not a valid object name %s", hash)
		}
		return "", nil, fmt.Errorf("failed to open object file: %w", err)
	}
	defer file.Close()

	zr, err := zlib.NewReader(file)
	if err != nil {
		return "", nil, fmt.Errorf("failed to initialize zlib reader: %w", err)
	}
	defer zr.Close()

	decompressed, err := io.ReadAll(zr)
	if err != nil {
		return "", nil, fmt.Errorf("failed to read decompressed data: %w", err)
	}

	nullIdx := bytes.IndexByte(decompressed, 0)
	if nullIdx == -1 {
		return "", nil, fmt.Errorf("fatal: corrupted object %s: missing null delimiter", hash)
	}

	header := decompressed[:nullIdx]
	content := decompressed[nullIdx+1:]

	parts := bytes.SplitN(header, []byte(" "), 2)
	if len(parts) < 2 {
		return "", nil, fmt.Errorf("fatal: corrupted object header %s", hash)
	}

	objectType := string(parts[0])

	return objectType, content, nil
}
