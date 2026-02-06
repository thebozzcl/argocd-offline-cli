package preview

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLocalRepos(t *testing.T) {
	mappings := map[string]string{
		"https://github.com/org/repo.git": "/path/to/repo",
	}
	repos := NewLocalRepos(mappings)
	assert.NotNil(t, repos)
	assert.Equal(t, mappings, repos.repoMappings)
}

func TestNormalizeRepoURL(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"https://github.com/org/repo.git", "https://github.com/org/repo"},
		{"https://github.com/org/repo", "https://github.com/org/repo"},
		{"https://github.com/org/repo/", "https://github.com/org/repo"},
		{"https://github.com/org/repo.git/", "https://github.com/org/repo"},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			result := normalizeRepoURL(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestGetLocalPath(t *testing.T) {
	mappings := map[string]string{
		"https://github.com/org/repo.git": "/path/to/repo",
	}
	repos := NewLocalRepos(mappings)

	tests := []struct {
		repoURL     string
		expected    string
		expectError bool
	}{
		{"https://github.com/org/repo.git", "/path/to/repo", false},
		{"https://github.com/org/repo", "/path/to/repo", false},
		{"https://github.com/org/other", "", true},
	}

	for _, tc := range tests {
		t.Run(tc.repoURL, func(t *testing.T) {
			result, err := repos.getLocalPath(tc.repoURL)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

func TestGetDirectories(t *testing.T) {
	tempDir := t.TempDir()

	// Create test directory structure
	dirs := []string{
		"app1",
		"app2",
		"app2/subdir",
		"config",
	}
	for _, dir := range dirs {
		err := os.MkdirAll(filepath.Join(tempDir, dir), 0755)
		require.NoError(t, err)
	}

	// Create some files (these should not be included)
	files := []string{
		"README.md",
		"app1/values.yaml",
	}
	for _, file := range files {
		err := os.WriteFile(filepath.Join(tempDir, file), []byte("test"), 0644)
		require.NoError(t, err)
	}

	mappings := map[string]string{
		"https://github.com/org/repo.git": tempDir,
	}
	repos := NewLocalRepos(mappings)

	directories, err := repos.GetDirectories(context.Background(), "https://github.com/org/repo.git", "main", "", false, false)
	require.NoError(t, err)

	assert.Contains(t, directories, "app1")
	assert.Contains(t, directories, "app2")
	assert.Contains(t, directories, "app2/subdir")
	assert.Contains(t, directories, "config")
	assert.NotContains(t, directories, "README.md")
}

func TestGetFiles(t *testing.T) {
	tempDir := t.TempDir()

	// Create test directory structure with files
	err := os.MkdirAll(filepath.Join(tempDir, "apps"), 0755)
	require.NoError(t, err)

	files := map[string]string{
		"apps/app1.yaml":   "name: app1",
		"apps/app2.yaml":   "name: app2",
		"apps/config.json": `{"key": "value"}`,
		"README.md":        "# README",
	}
	for file, content := range files {
		err := os.WriteFile(filepath.Join(tempDir, file), []byte(content), 0644)
		require.NoError(t, err)
	}

	mappings := map[string]string{
		"https://github.com/org/repo.git": tempDir,
	}
	repos := NewLocalRepos(mappings)

	t.Run("match yaml files in apps dir", func(t *testing.T) {
		result, err := repos.GetFiles(context.Background(), "https://github.com/org/repo.git", "main", "", "apps/*.yaml", false, false)
		require.NoError(t, err)

		assert.Len(t, result, 2)
		assert.Contains(t, result, "apps/app1.yaml")
		assert.Contains(t, result, "apps/app2.yaml")
		assert.Equal(t, []byte("name: app1"), result["apps/app1.yaml"])
	})

	t.Run("match all files in apps dir", func(t *testing.T) {
		result, err := repos.GetFiles(context.Background(), "https://github.com/org/repo.git", "main", "", "apps/*", false, false)
		require.NoError(t, err)

		assert.Len(t, result, 3)
	})

	t.Run("match using double star", func(t *testing.T) {
		result, err := repos.GetFiles(context.Background(), "https://github.com/org/repo.git", "main", "", "**/*.yaml", false, false)
		require.NoError(t, err)

		assert.Len(t, result, 2)
	})
}

func TestGetFilesNoMatch(t *testing.T) {
	tempDir := t.TempDir()

	mappings := map[string]string{
		"https://github.com/org/repo.git": tempDir,
	}
	repos := NewLocalRepos(mappings)

	result, err := repos.GetFiles(context.Background(), "https://github.com/org/repo.git", "main", "", "*.yaml", false, false)
	require.NoError(t, err)
	assert.Empty(t, result)
}

func TestGetDirectoriesUnknownRepo(t *testing.T) {
	repos := NewLocalRepos(map[string]string{})

	_, err := repos.GetDirectories(context.Background(), "https://github.com/org/unknown.git", "main", "", false, false)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no local mapping found")
}
