package preview

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
)

// LocalRepos implements the services.Repos interface for local filesystem directories.
// It maps remote Git repository URLs to local directory paths.
type LocalRepos struct {
	repoMappings map[string]string
}

// NewLocalRepos creates a new LocalRepos instance with the given URL-to-path mappings.
func NewLocalRepos(mappings map[string]string) *LocalRepos {
	return &LocalRepos{
		repoMappings: mappings,
	}
}

// normalizeRepoURL normalizes a repository URL for comparison by removing
// trailing slashes and .git suffix.
func normalizeRepoURL(url string) string {
	url = strings.TrimSuffix(url, "/")
	url = strings.TrimSuffix(url, ".git")
	return url
}

// getLocalPath returns the local path for a given repository URL.
func (r *LocalRepos) getLocalPath(repoURL string) (string, error) {
	normalizedURL := normalizeRepoURL(repoURL)

	// Try exact match first
	if localPath, ok := r.repoMappings[repoURL]; ok {
		return localPath, nil
	}

	// Try normalized URL match
	for url, localPath := range r.repoMappings {
		if normalizeRepoURL(url) == normalizedURL {
			return localPath, nil
		}
	}

	return "", fmt.Errorf("no local mapping found for repository URL: %s", repoURL)
}

// GetFiles returns content of files (not directories) within the target repo
// that match the given pattern.
func (r *LocalRepos) GetFiles(ctx context.Context, repoURL, revision, project, pattern string, noRevisionCache, verifyCommit bool) (map[string][]byte, error) {
	localPath, err := r.getLocalPath(repoURL)
	if err != nil {
		return nil, err
	}

	result := make(map[string][]byte)

	// Use doublestar for glob matching (supports **)
	matches, err := doublestar.Glob(os.DirFS(localPath), pattern)
	if err != nil {
		return nil, fmt.Errorf("error matching pattern %s: %w", pattern, err)
	}

	for _, match := range matches {
		fullPath := filepath.Join(localPath, match)
		info, err := os.Stat(fullPath)
		if err != nil {
			return nil, fmt.Errorf("error getting file info for %s: %w", fullPath, err)
		}

		// Skip directories
		if info.IsDir() {
			continue
		}

		content, err := os.ReadFile(fullPath)
		if err != nil {
			return nil, fmt.Errorf("error reading file %s: %w", fullPath, err)
		}

		result[match] = content
	}

	return result, nil
}

// GetDirectories returns a list of directories (not files) within the target repo.
func (r *LocalRepos) GetDirectories(ctx context.Context, repoURL, revision, project string, noRevisionCache, verifyCommit bool) ([]string, error) {
	localPath, err := r.getLocalPath(repoURL)
	if err != nil {
		return nil, err
	}

	var directories []string

	err = filepath.WalkDir(localPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip the root directory itself
		if path == localPath {
			return nil
		}

		// Only include directories
		if d.IsDir() {
			// Get relative path from the repo root
			relPath, err := filepath.Rel(localPath, path)
			if err != nil {
				return err
			}
			// Normalize path separators for consistency
			relPath = filepath.ToSlash(relPath)
			directories = append(directories, relPath)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking directory %s: %w", localPath, err)
	}

	return directories, nil
}
