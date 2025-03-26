// Package fs provides utilities for file system operations, including directory management,
// file reading, and zip file handling.
package fs

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
)

// EnsureEmptyDir creates an empty directory at the specified path.
// If the directory already exists, it will be emptied.
// Returns an error if the directory cannot be created or emptied.
func EnsureEmptyDir(name string) error {
	err := DeleteDirIfExists(name)
	if err != nil {
		return err
	}
	return CreateDir(name)
}

// DeleteDir removes a directory and all its contents.
// Returns an error if the directory cannot be deleted.
func DeleteDir(name string) error {
	f, e := os.Stat(name)
	if e != nil {
		return e
	}
	return deleteFileDirectory(name, f)
}

// deleteFileDirectory is an internal function that recursively deletes a file or directory.
// It handles both files and directories, ensuring proper cleanup of all contents.
func deleteFileDirectory(name string, f os.FileInfo) error {
	if f.IsDir() {
		return os.RemoveAll(name)
	} else {
		return fmt.Errorf("%s not a directory", name)
	}
}

// DeleteDirIfExists removes a directory and all its contents if it exists.
// If the directory does not exist, no error is returned.
// Returns an error if the directory exists but cannot be deleted.
func DeleteDirIfExists(name string) error {
	f, e := os.Stat(name)
	if e != nil {
		return nil
	}
	return deleteFileDirectory(name, f)
}

// CreateDir creates a new directory at the specified path.
// Returns an error if the directory cannot be created.
func CreateDir(name string) error {
	f, e := os.Stat(name)
	if e != nil {
		return os.MkdirAll(name, os.ModePerm)
	}
	return fmt.Errorf("%s already exists", f.Name())
}

// CreateDirIfNotExists creates a new directory at the specified path if it does not exist.
// If the directory already exists, no error is returned.
// Returns an error if the directory needs to be created but cannot be.
func CreateDirIfNotExists(name string) error {
	_, e := os.Stat(name)
	if e != nil {
		return os.MkdirAll(name, os.ModePerm)
	}
	return nil
}

// ReadZipEntry reads the contents of a specific entry from a zip file.
// Returns the contents of the zip entry and any error that occurred during reading.
func ReadZipEntry(f *zip.File) ([]byte, error) {
	rc, err := f.Open()
	if err != nil {
		return nil, err
	}
	defer CloseSafely(rc)

	content, err := io.ReadAll(rc)
	if err != nil {
		return nil, err
	}
	return content, nil
}

// ReadBytes reads all bytes from a reader and returns them as a byte slice.
// Returns the byte slice and any error that occurred during reading.
func ReadBytes(f string) ([]byte, error) {
	rc, err := os.ReadFile(f)
	if err != nil {
		return nil, err
	}
	return rc, nil
}

// CloseSafely safely closes an io.Closer, ignoring any errors that occur during closing.
// This is useful when you want to ensure a resource is closed but don't need to handle
// any errors that might occur during the closing process.
func CloseSafely(fs io.Closer) {
	_ = fs.Close()
}
