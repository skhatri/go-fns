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

func TestParsePasswordEntry(t *testing.T) {
	t.Run("direct password string", func(t *testing.T) {
		input := "mypassword123"
		result, err := ParsePasswordEntry(input)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if result != input {
			t.Errorf("Expected %s, got %s", input, result)
		}
	})

	t.Run("password from file", func(t *testing.T) {
		// Create a temporary file with password
		tempDir := t.TempDir()
		passFile := filepath.Join(tempDir, "password.txt")
		password := "secretpassword123"
		if err := os.WriteFile(passFile, []byte(password), 0600); err != nil {
			t.Fatal(err)
		}

		result, err := ParsePasswordEntry("file:" + passFile)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if result != password {
			t.Errorf("Expected %s, got %s", password, result)
		}
	})

	t.Run("non-existent file", func(t *testing.T) {
		_, err := ParsePasswordEntry("file:/nonexistent/path/password.txt")
		if err == nil {
			t.Error("Expected error for non-existent file, got nil")
		}
	})
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

	t.Run("ListFiles with non-existent directory", func(t *testing.T) {
		nonExistentDir := filepath.Join(tempDir, "nonexistent")
		result := ListFiles(nonExistentDir, ".txt")
		if len(result) != 0 {
			t.Errorf("Expected empty slice for non-existent directory, got %v", result)
		}
	})

	t.Run("ListFiles with inaccessible directory", func(t *testing.T) {
		if os.Getuid() == 0 {
			t.Skip("Skipping test when running as root")
		}

		restrictedDir := filepath.Join(tempDir, "restricted")
		if err := os.MkdirAll(restrictedDir, 0); err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(restrictedDir)

		result := ListFiles(restrictedDir, ".txt")
		if len(result) != 0 {
			t.Errorf("Expected empty slice for inaccessible directory, got %v", result)
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
