package repo

import (
	"bytes"
	"encoding/hex"
	"fmt"
)

type TreeEntry struct {
	Mode string
	Name string
	Hash string
}

func ParseTree(content []byte) ([]TreeEntry, error) {
	var entries []TreeEntry

	for len(content) > 0 {
		spaceIdx := bytes.IndexByte(content, ' ')
		if spaceIdx == -1 {
			return nil, fmt.Errorf("corrupted tree: missing space after mode")
		}
		mode := string(content[:spaceIdx])

		nullIdx := bytes.IndexByte(content[spaceIdx+1:], 0)
		if nullIdx == -1 {
			return nil, fmt.Errorf("corrputed tree: missing null byte after file name")
		}

		absNullIdx := spaceIdx + 1 + nullIdx
		name := string(content[spaceIdx+1 : absNullIdx])

		shaStart := absNullIdx + 1
		shaEnd := shaStart + 20
		if len(content) < shaEnd {
			return nil, fmt.Errorf("corrupted tree: unexpected EOF reading SHA")
		}

		rawSha := content[shaStart:shaEnd]
		hashHex := hex.EncodeToString(rawSha)

		entries = append(entries, TreeEntry{
			Mode: mode,
			Name: name,
			Hash: hashHex,
		})

		content = content[shaEnd:]
	}

	return entries, nil
}

func LsTree(targetDir string, treeHash string) ([]TreeEntry, error) {
	objType, content, err := ReadObject(targetDir, treeHash)
	if err != nil {
		return nil, err
	}

	if objType != "tree" {
		return nil, fmt.Errorf("fatal: not a tree object: %s (is %s)", treeHash, objType)
	}

	return ParseTree(content)
}
