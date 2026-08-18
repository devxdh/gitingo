package repo

import (
	"crypto/sha1"
	"fmt"
	"strings"
	"time"

	u "github.com/devxdh/gitingo/pkg/utils"
)

func CommitTree(targetDir string, treeHash string, parentHash string, message string) (string, error) {
	if exists, err := u.DotGitExists(targetDir); err != nil || !exists {
		return "", fmt.Errorf("fatal: not a git repository (or any of the parent directories): .git")
	}

	if len(treeHash) != 40 {
		return "", fmt.Errorf("fatal: not a valid object name %s", treeHash)
	}

	now := time.Now()
	timestamp := now.Unix()
	_, offsetSeconds := now.Zone()
	offsetHours := offsetSeconds / 3600
	offsetMinutes := (offsetSeconds % 3600) / 60
	timezone := fmt.Sprintf("%+03d%02d", offsetHours, offsetMinutes)

	authorName := "user"
	authorEmail := "user@example.com"
	authorLine := fmt.Sprintf("%s <%s> %d %s", authorName, authorEmail, timestamp, timezone)

	var builder strings.Builder
	fmt.Fprintf(&builder, "tree %s\n", treeHash)

	if parentHash != "" {
		if len(parentHash) != 40 {
			return "", fmt.Errorf("fatal: not a valid parent object name %s", parentHash)
		}
		fmt.Fprintf(&builder, "parent %s\n", parentHash)
	}
	fmt.Fprintf(&builder, "author %s\n", authorLine)
	fmt.Fprintf(&builder, "committer %s\n", authorLine)
	fmt.Fprintf(&builder, "\n")
	fmt.Fprintf(&builder, "%s", strings.TrimSpace(message))
	fmt.Fprintf(&builder, "\n")

	payload := []byte(builder.String())

	header := fmt.Sprintf("commit %d\x00", len(payload))
	fullPayload := append([]byte(header), payload...)

	hasher := sha1.New()
	hasher.Write(fullPayload)
	commitHash := fmt.Sprintf("%x", hasher.Sum(nil))

	if err := writeObject(targetDir, commitHash, fullPayload); err != nil {
		return "", err
	}

	return commitHash, nil
}
