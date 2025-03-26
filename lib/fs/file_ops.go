// Package fs provides utilities for file system operations, including directory management,
// file reading, and zip file handling.
package fs

import (
	"os"
	"path/filepath"
	"strings"
)

// ParsePasswordEntry parses a password entry from a file.
// The expected format is "username:password" per line.
// Returns the username and password, and any error that occurred during parsing.
func ParsePasswordEntry(passSource string) (string, error) {
	passData := ""
	if strings.HasPrefix(passSource, "file:") {
		passSource = strings.Replace(passSource, "file:", "", 1)
		passBytes, passErr := os.ReadFile(passSource)
		if passErr != nil {
			return "", passErr
		}
		passData = string(passBytes)
	} else {
		passData = passSource
	}
	return passData, nil
}

// ListFiles returns a list of files in the specified directory.
// The pattern parameter can be used to filter files using a glob pattern.
// Returns a slice of file names and any error that occurred during listing.
func ListFiles(rootPath string, ext string) []string {
	files := make([]string, 0)
	walkerFn := func(path string, fi os.FileInfo, err error) error {
		if err == nil && !fi.IsDir() {
			if filepath.Ext(path) == ext && strings.HasSuffix(path, ext) {
				files = append(files, path)
			}
		}
		return err
	}
	filepath.Walk(rootPath, walkerFn)
	return files
}
