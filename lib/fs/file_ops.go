package fs

import (
	"os"
	"path/filepath"
	"strings"
)

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
