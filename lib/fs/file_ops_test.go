package fs

import (
	"os"
	"path/filepath"
	"testing"
)

func stringSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
func TestListFiles(t *testing.T) {
	tempDir := t.TempDir()
	createTestFiles(tempDir, []string{"a.txt", "b.txt", "c.yaml", "subdir/d.txt"})

	t.Run("ListFiles with valid extension", func(t *testing.T) {
		ext := ".txt"
		expected := []string{filepath.Join(tempDir, "a.txt"),
			filepath.Join(tempDir, "b.txt"),
			filepath.Join(tempDir, "subdir/d.txt"),
		}
		result := ListFiles(tempDir, ext)

		if !stringSlicesEqual(result, expected) {
			t.Errorf("Expected %v, but got %v", expected, result)
		}
	})

	t.Run("ListFiles with invalid extension", func(t *testing.T) {
		ext := ".pdf"
		expected := []string{}
		result := ListFiles(tempDir, ext)

		if !stringSlicesEqual(result, expected) {
			t.Errorf("Expected %v, but got %v", expected, result)
		}
	})
}

func createTestFiles(dir string, files []string) {
	for _, file := range files {
		filePath := filepath.Join(dir, file)
		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			panic(err)
		}
		file, err := os.Create(filePath)
		if err != nil {
			panic(err)
		}
		file.Close()
	}
}
