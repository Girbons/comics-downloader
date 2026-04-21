package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathSetup(t *testing.T) {
	result, err := OutputPathSetup(true, filepath.Dir(os.Args[0]), "example-source", "comic-name")

	assert.Nil(t, err)
	assert.Contains(t, result, "example-source")
	assert.Contains(t, result, "comic-name")
}

func TestGenerateFileName(t *testing.T) {
	result := GetPathToFile("path/to/something", "comic-name", "invalid_character", "pdf", false)
	assert.Equal(t, "path/to/something/comic-name - invalid_character.pdf", result)
	result = GetPathToFile("path/to/something", "comic-name", "invalid_character", "pdf", true)
	assert.Equal(t, "path/to/something/invalid_character.pdf", result)
}

func TestDirectoryOrFileDoesNotExist(t *testing.T) {

	path, _ := ImagesPathSetup(true, filepath.Dir(os.Args[0]), "source", "name", "issue-", "issueNumber")
	defer os.RemoveAll(path)

	result := DirectoryOrFileDoesNotExist(path)

	assert.False(t, result)
}

func TestTrimNameLength(t *testing.T) {
	shortName := "short-name"
	exactLengthName := strings.Repeat("a", NameLength)
	longName := strings.Repeat("b", NameLength+10)

	assert.Equal(t, shortName, TrimNameLength(shortName))
	assert.Equal(t, exactLengthName, TrimNameLength(exactLengthName))
	assert.Equal(t, strings.Repeat("b", NameLength), TrimNameLength(longName))
}

func TestPathSetupTrimsComicName(t *testing.T) {
	outputFolder := t.TempDir()
	longComicName := strings.Repeat("comic", 30)

	result, err := OutputPathSetup(true, outputFolder, "example-source", longComicName)

	assert.Nil(t, err)
	assert.Contains(t, result, filepath.Join("comics", "example-source"))
	assert.Contains(t, result, strings.Repeat("comic", 20))
	assert.NotContains(t, result, longComicName)
}

func TestImagesPathSetupTrimsIssueNumber(t *testing.T) {
	outputFolder := t.TempDir()
	longIssueNumber := strings.Repeat("issue", 30)

	result, err := ImagesPathSetup(true, outputFolder, "source", "name", "issue-", longIssueNumber)

	assert.Nil(t, err)
	assert.Contains(t, result, "images-")
	assert.Contains(t, result, fmt.Sprintf("images-%s", strings.Repeat("issue", 20)))
	assert.NotContains(t, result, longIssueNumber)
}

func TestImagesPathSetupTrimsCustomFolderName(t *testing.T) {
	outputFolder := t.TempDir()
	longIssueFolderName := strings.Repeat("folder", 20)
	longIssueNumber := strings.Repeat("number", 20)

	result, err := ImagesPathSetup(false, outputFolder, "source", "name", longIssueFolderName, longIssueNumber)

	assert.Nil(t, err)
	assert.Equal(t, NameLength, len(filepath.Base(result)))
	assert.NotContains(t, result, longIssueFolderName+longIssueNumber)
}

func TestGetPathToFileTrimsNameAndIssueNumber(t *testing.T) {
	dir := "path/to/something"
	longName := strings.Repeat("n", NameLength+20)
	longIssueNumber := strings.Repeat("i", NameLength+15)

	result := GetPathToFile(dir, longName, longIssueNumber, "pdf", false)
	expected := fmt.Sprintf("%s/%s - %s.pdf", dir, strings.Repeat("n", NameLength), strings.Repeat("i", NameLength))
	assert.Equal(t, expected, result)

	resultIssueOnly := GetPathToFile(dir, longName, longIssueNumber, "pdf", true)
	expectedIssueOnly := fmt.Sprintf("%s/%s.pdf", dir, strings.Repeat("i", NameLength))
	assert.Equal(t, expectedIssueOnly, resultIssueOnly)
}
