package fs

import (
	"archive/zip"
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestEnsureEmptyDir(t *testing.T) {
	t.Run("create new empty dir", func(t *testing.T) {
		tempDir := t.TempDir()
		testDir := filepath.Join(tempDir, "test-dir")

		err := EnsureEmptyDir(testDir)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		// Verify directory exists and is empty
		files, err := os.ReadDir(testDir)
		if err != nil {
			t.Errorf("Failed to read directory: %v", err)
		}
		if len(files) != 0 {
			t.Error("Directory should be empty")
		}
	})

	t.Run("ensure existing dir is emptied", func(t *testing.T) {
		tempDir := t.TempDir()
		testDir := filepath.Join(tempDir, "test-dir")

		// Create directory with a file
		if err := os.MkdirAll(testDir, os.ModePerm); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(testDir, "test.txt"), []byte("test"), 0644); err != nil {
			t.Fatal(err)
		}

		err := EnsureEmptyDir(testDir)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		// Verify directory exists and is empty
		files, err := os.ReadDir(testDir)
		if err != nil {
			t.Errorf("Failed to read directory: %v", err)
		}
		if len(files) != 0 {
			t.Error("Directory should be empty")
		}
	})

	t.Run("error on non-writable parent", func(t *testing.T) {
		if os.Getuid() == 0 {
			t.Skip("Skipping test when running as root")
		}

		tempDir := t.TempDir()
		restrictedDir := filepath.Join(tempDir, "restricted")
		if err := os.MkdirAll(restrictedDir, 0); err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(restrictedDir)

		testDir := filepath.Join(restrictedDir, "test-dir")
		err := EnsureEmptyDir(testDir)
		if err == nil {
			t.Error("Expected error for non-writable parent directory")
		}
	})

	t.Run("error on file exists", func(t *testing.T) {
		tempDir := t.TempDir()
		testFile := filepath.Join(tempDir, "test.txt")

		// Create a file instead of directory
		if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
			t.Fatal(err)
		}

		err := EnsureEmptyDir(testFile)
		if err == nil {
			t.Error("Expected error when path exists as file")
		}
	})
}

func TestDeleteDir(t *testing.T) {
	t.Run("delete existing directory", func(t *testing.T) {
		tempDir := t.TempDir()
		testDir := filepath.Join(tempDir, "test-dir")

		if err := os.MkdirAll(testDir, os.ModePerm); err != nil {
			t.Fatal(err)
		}

		err := DeleteDir(testDir)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if _, err := os.Stat(testDir); !os.IsNotExist(err) {
			t.Error("Directory should not exist")
		}
	})

	t.Run("error on non-existent directory", func(t *testing.T) {
		tempDir := t.TempDir()
		testDir := filepath.Join(tempDir, "nonexistent")

		err := DeleteDir(testDir)
		if err == nil {
			t.Error("Expected error for non-existent directory")
		}
	})

	t.Run("error on file instead of directory", func(t *testing.T) {
		tempDir := t.TempDir()
		testFile := filepath.Join(tempDir, "test.txt")

		if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
			t.Fatal(err)
		}

		err := DeleteDir(testFile)
		if err == nil {
			t.Error("Expected error when trying to delete a file")
		}
	})
}

func TestDeleteDirIfExists(t *testing.T) {
	t.Run("delete existing directory", func(t *testing.T) {
		tempDir := t.TempDir()
		testDir := filepath.Join(tempDir, "test-dir")

		if err := os.MkdirAll(testDir, os.ModePerm); err != nil {
			t.Fatal(err)
		}

		err := DeleteDirIfExists(testDir)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if _, err := os.Stat(testDir); !os.IsNotExist(err) {
			t.Error("Directory should not exist")
		}
	})

	t.Run("no error on non-existent directory", func(t *testing.T) {
		tempDir := t.TempDir()
		testDir := filepath.Join(tempDir, "nonexistent")

		err := DeleteDirIfExists(testDir)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("error on file instead of directory", func(t *testing.T) {
		tempDir := t.TempDir()
		testFile := filepath.Join(tempDir, "test.txt")

		if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
			t.Fatal(err)
		}

		err := DeleteDirIfExists(testFile)
		if err == nil {
			t.Error("Expected error when trying to delete a file")
		}
	})
}

func TestCreateDir(t *testing.T) {
	t.Run("create new directory", func(t *testing.T) {
		tempDir := t.TempDir()
		testDir := filepath.Join(tempDir, "test-dir")

		err := CreateDir(testDir)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if fi, err := os.Stat(testDir); err != nil || !fi.IsDir() {
			t.Error("Directory should exist")
		}
	})

	t.Run("error on existing directory", func(t *testing.T) {
		tempDir := t.TempDir()
		testDir := filepath.Join(tempDir, "test-dir")

		if err := os.MkdirAll(testDir, os.ModePerm); err != nil {
			t.Fatal(err)
		}

		err := CreateDir(testDir)
		if err == nil {
			t.Error("Expected error when directory already exists")
		}
	})
}

func TestCreateDirIfNotExists(t *testing.T) {
	t.Run("create new directory", func(t *testing.T) {
		tempDir := t.TempDir()
		testDir := filepath.Join(tempDir, "test-dir")

		err := CreateDirIfNotExists(testDir)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if fi, err := os.Stat(testDir); err != nil || !fi.IsDir() {
			t.Error("Directory should exist")
		}
	})

	t.Run("no error on existing directory", func(t *testing.T) {
		tempDir := t.TempDir()
		testDir := filepath.Join(tempDir, "test-dir")

		if err := os.MkdirAll(testDir, os.ModePerm); err != nil {
			t.Fatal(err)
		}

		err := CreateDirIfNotExists(testDir)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})
}

// ErrorReader is a custom io.Reader that always returns an error
type ErrorReader struct{}

func (er ErrorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("forced read error")
}

func (er ErrorReader) Close() error {
	return nil
}

// MockZipFile is a mock implementation of zip.File for testing
type MockZipFile struct {
	zip.File
	openError error
	reader    io.ReadCloser
}

func (m *MockZipFile) Open() (io.ReadCloser, error) {
	if m.openError != nil {
		return nil, m.openError
	}
	return m.reader, nil
}

func TestReadZipEntry(t *testing.T) {
	t.Run("read valid zip entry", func(t *testing.T) {
		// Create a test zip file in memory
		buf := new(bytes.Buffer)
		w := zip.NewWriter(buf)

		content := []byte("test content")
		f, err := w.Create("test.txt")
		if err != nil {
			t.Fatal(err)
		}
		if _, err := f.Write(content); err != nil {
			t.Fatal(err)
		}
		if err := w.Close(); err != nil {
			t.Fatal(err)
		}

		// Read the zip file
		r, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
		if err != nil {
			t.Fatal(err)
		}

		// Test ReadZipEntry
		data, err := ReadZipEntry(r.File[0])
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if !bytes.Equal(data, content) {
			t.Errorf("Expected content %q, got %q", content, data)
		}
	})

	t.Run("error on corrupted zip entry", func(t *testing.T) {
		// Create a temporary zip file
		tempDir := t.TempDir()
		zipPath := filepath.Join(tempDir, "test.zip")

		// Create a corrupted zip file
		if err := os.WriteFile(zipPath, []byte("corrupted zip data"), 0644); err != nil {
			t.Fatal(err)
		}

		// Try to read the corrupted zip
		r, err := zip.OpenReader(zipPath)
		if err == nil {
			defer r.Close()
			_, err = ReadZipEntry(r.File[0])
			if err == nil {
				t.Error("Expected error for corrupted zip entry")
			}
		}
	})

	t.Run("error on read failure", func(t *testing.T) {
		// Create a temporary zip file
		tempDir := t.TempDir()
		zipPath := filepath.Join(tempDir, "test.zip")

		// Create a zip file with valid content first
		{
			f, err := os.Create(zipPath)
			if err != nil {
				t.Fatal(err)
			}
			w := zip.NewWriter(f)

			file, err := w.Create("test.txt")
			if err != nil {
				t.Fatal(err)
			}
			if _, err := file.Write([]byte("test content")); err != nil {
				t.Fatal(err)
			}
			if err := w.Close(); err != nil {
				t.Fatal(err)
			}
			if err := f.Close(); err != nil {
				t.Fatal(err)
			}
		}

		// Read the zip file
		r, err := zip.OpenReader(zipPath)
		if err != nil {
			t.Fatal(err)
		}
		defer r.Close()

		// Corrupt the underlying file after creating the zip.File
		if err := os.WriteFile(zipPath, []byte("corrupted content"), 0644); err != nil {
			t.Fatal(err)
		}

		// Try to read the content - this should fail since the underlying file is corrupted
		_, err = ReadZipEntry(r.File[0])
		if err == nil {
			t.Error("Expected error when reading corrupted content")
		}
	})

	t.Run("error on zip open failure", func(t *testing.T) {
		// Create a temporary zip file
		tempDir := t.TempDir()
		zipPath := filepath.Join(tempDir, "test.zip")

		// Create a zip file with valid content first
		{
			f, err := os.Create(zipPath)
			if err != nil {
				t.Fatal(err)
			}
			w := zip.NewWriter(f)

			// Create a file with data descriptor enabled
			fh := &zip.FileHeader{
				Name:   "test.txt",
				Method: zip.Deflate,
				Flags:  0x8, // Enable data descriptor
			}
			file, err := w.CreateHeader(fh)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := file.Write([]byte("test content")); err != nil {
				t.Fatal(err)
			}
			if err := w.Close(); err != nil {
				t.Fatal(err)
			}
			if err := f.Close(); err != nil {
				t.Fatal(err)
			}
		}

		// Read the zip file content
		content, err := os.ReadFile(zipPath)
		if err != nil {
			t.Fatal(err)
		}

		// Find and corrupt the data descriptor
		// Data descriptor signature is 'PK\x07\x08'
		for i := 0; i < len(content)-4; i++ {
			if content[i] == 0x50 && content[i+1] == 0x4b && content[i+2] == 0x07 && content[i+3] == 0x08 {
				// Corrupt the CRC-32 value
				content[i+4] = 0xFF
				content[i+5] = 0xFF
				content[i+6] = 0xFF
				content[i+7] = 0xFF
				break
			}
		}

		// Write back the corrupted content
		if err := os.WriteFile(zipPath, content, 0644); err != nil {
			t.Fatal(err)
		}

		// Open the corrupted zip file
		r, err := zip.OpenReader(zipPath)
		if err != nil {
			t.Fatal(err)
		}
		defer r.Close()

		// Try to read the file with corrupted data descriptor
		_, err = ReadZipEntry(r.File[0])
		if err == nil {
			t.Error("Expected error when reading file with corrupted data descriptor")
		}
	})
}

func TestReadBytes(t *testing.T) {
	t.Run("read existing file", func(t *testing.T) {
		tempDir := t.TempDir()
		testFile := filepath.Join(tempDir, "test.txt")
		content := []byte("test content")

		if err := os.WriteFile(testFile, content, 0644); err != nil {
			t.Fatal(err)
		}

		data, err := ReadBytes(testFile)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if !bytes.Equal(data, content) {
			t.Errorf("Expected content %q, got %q", content, data)
		}
	})

	t.Run("error on non-existent file", func(t *testing.T) {
		tempDir := t.TempDir()
		testFile := filepath.Join(tempDir, "nonexistent.txt")

		_, err := ReadBytes(testFile)
		if err == nil {
			t.Error("Expected error for non-existent file")
		}
	})
}

func TestCloseSafely(t *testing.T) {
	t.Run("close without error", func(t *testing.T) {
		// Create a temporary file as a proper io.Closer
		tempFile, err := os.CreateTemp("", "test")
		if err != nil {
			t.Fatal(err)
		}
		CloseSafely(tempFile)
	})

	t.Run("close with error", func(t *testing.T) {
		// Create a closer that always returns an error
		errorCloser := &errorCloser{}
		CloseSafely(errorCloser)
	})
}

// Helper type for testing CloseSafely with error
type errorCloser struct{}

func (e *errorCloser) Close() error {
	return io.ErrClosedPipe
}
