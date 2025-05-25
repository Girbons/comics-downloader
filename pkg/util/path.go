package util

import (
	"fmt"
	"os"
	"path/filepath"
)

// createPath create folders given the path.
func createPath(path string) (string, error) {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return path, err
	}

	dir, err := filepath.Abs(path)
	if err != nil {
		return dir, err
	}

	return dir, err
}

// PathSetup creates the folders where the comic will be saved.
// when `createDefaultPath` is false the comic is stored without prepending
// the default folder path `comics/source/name/[comic.format]`.
func PathSetup(createDefaultPath bool, outputFolder, source, name string) (string, error) {
	path := fmt.Sprintf("%s/comics/%s/%s/", outputFolder, source, name)

	if !createDefaultPath {
		path = fmt.Sprintf("%s/", outputFolder)
	}

	return createPath(path)
}

// ImagesPathSetup creates the folders for the images to be saved.
// when `createDefaultPath` is false the images are stored without prepending
// the default folder path `comics/source/name/[comic.format]`.
func ImagesPathSetup(createDefaultPath bool, outputFolder, source, name, issueFolderName, issueNumber string) (string, error) {
	path := fmt.Sprintf("%s/comics/%s/%s/images-%s/", outputFolder, source, name, issueNumber)

	if !createDefaultPath {
		path = fmt.Sprintf("%s/%s%s", outputFolder, issueFolderName, issueNumber)
	}

	return createPath(path)
}

// CurrentDir returns the path where the executable was called.
func CurrentDir() (string, error) {
	exePath, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return exePath, nil
}

// DirectoryOrFileDoesNotExist check if a directory/file exist.
func DirectoryOrFileDoesNotExist(filePath string) bool {
	_, err := os.Stat(filePath)

	return os.IsNotExist(err)
}

// GetPathToFile returns the path where the file should be saved.
func GetPathToFile(dir, name, issueNumber, format string, issueNumberOnly bool) string {
	if issueNumberOnly {
		return fmt.Sprintf("%s/%s.%s", dir, issueNumber, format)
	}
	return fmt.Sprintf("%s/%s-%s.%s", dir, name, issueNumber, format)
}
