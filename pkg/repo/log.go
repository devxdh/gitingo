package repo

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	u "github.com/devxdh/gitingo/pkg/utils"
)

type CommitData struct {
	Hash    string
	Tree    string
	Parent  string
	Author  string
	Date    time.Time
	Message string
}

func ParseCommit(hash string, content []byte) (*CommitData, error) {
	scanner := bufio.NewScanner(bytes.NewReader(content))
	commit := &CommitData{Hash: hash}

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}

		if strings.HasPrefix(line, "tree ") {
			commit.Tree = strings.TrimPrefix(line, "tree ")
		} else if strings.HasPrefix(line, "parent ") {
			commit.Parent = strings.TrimPrefix(line, "parent ")
		} else if strings.HasPrefix(line, "author ") {
			authorRaw := strings.TrimPrefix(line, "author ")
			parts := strings.Split(authorRaw, " ")
			if len(parts) >= 3 {
				tsStr := parts[len(parts)-2]
				commit.Author = strings.Join(parts[:len(parts)-2], " ")

				if sec, err := strconv.ParseInt(tsStr, 10, 64); err == nil {
					commit.Date = time.Unix(sec, 0)
				}
			} else {
				commit.Author = authorRaw
			}
		}
	}

	var msgBuilder strings.Builder
	for scanner.Scan() {
		if msgBuilder.Len() > 0 {
			msgBuilder.WriteString("\n")
		}
		msgBuilder.WriteString(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan commit payload: %w", err)
	}

	commit.Message = strings.TrimSpace(msgBuilder.String())
	return commit, nil
}

func GetHeadCommitHash(targetDir string) (string, error) {
	headPath := filepath.Join(targetDir, ".git", "HEAD")
	headContent, err := os.ReadFile(headPath)
	if err != nil {
		return "", fmt.Errorf("failed to read HEAD: %w", err)
	}

	headStr := strings.TrimSpace(string(headContent))
	if strings.HasPrefix(headStr, "ref: ") {
		refRelative := strings.TrimPrefix(headStr, "ref: ")
		refFullPath := filepath.Join(targetDir, ".git", refRelative)

		refContent, err := os.ReadFile(refFullPath)
		if err != nil {
			if os.IsNotExist(err) {
				return "", fmt.Errorf("fatal: your current branch does not have any commits yet")
			}
			return "", fmt.Errorf("failed to read ref %s: %w", refRelative, err)
		}
		return strings.TrimSpace(string(refContent)), nil
	}

	if len(headStr) == 40 {
		return headStr, nil
	}

	return "", fmt.Errorf("fatal: invalid HEAD format: %s", headStr)
}

func Log(targetDir string) ([]*CommitData, error) {
	if exists, err := u.DotGitExists(targetDir); err != nil || !exists {
		return nil, fmt.Errorf("fatal: not a git repository (or any of the parent directories): .git")
	}

	currentHash, err := GetHeadCommitHash(targetDir)
	if err != nil {
		return nil, err
	}

	var commits []*CommitData

	for currentHash != "" {
		objType, content, err := ReadObject(targetDir, currentHash)
		if err != nil {
			return nil, fmt.Errorf("failed to read commit object %s: %w", currentHash, err)
		}

		if objType != "commit" {
			return nil, fmt.Errorf("fatal: object %s is not a commit (is %s)", currentHash, objType)
		}

		commit, err := ParseCommit(currentHash, content)
		if err != nil {
			return nil, err
		}

		commits = append(commits, commit)
		currentHash = strings.TrimSpace(commit.Parent)
	}

	return commits, nil
}
