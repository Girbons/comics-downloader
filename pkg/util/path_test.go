package util

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathSetup(t *testing.T) {
	result, err := PathSetup(true, filepath.Dir(os.Args[0]), "example-source", "comic-name")

	assert.Nil(t, err)
	assert.Contains(t, result, "example-source")
	assert.Contains(t, result, "comic-name")
}

func TestGenerateFileName(t *testing.T) {
	result := GetPathToFile("path/to/something", "comic-name", "invalid_character", "pdf", false)
	assert.Equal(t, "path/to/something/comic-name-invalid_character.pdf", result)
	result = GetPathToFile("path/to/something", "comic-name", "invalid_character", "pdf", true)
	assert.Equal(t, "path/to/something/invalid_character.pdf", result)
}

func TestDirectoryOrFileDoesNotExist(t *testing.T) {

	path, _ := ImagesPathSetup(true, filepath.Dir(os.Args[0]), "source", "name", "issue-", "issueNumber")
	defer os.RemoveAll(path)

	result := DirectoryOrFileDoesNotExist(path)

	assert.False(t, result)
}
