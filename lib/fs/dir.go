package fs

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
)

func EnsureEmptyDir(name string) error {
	err := DeleteDirIfExists(name)
	if err != nil {
		return err
	}
	return CreateDir(name)
}

func DeleteDir(name string) error {
	f, e := os.Stat(name)
	if e != nil {
		return e
	}
	return deleteFileDirectory(name, f)
}

func deleteFileDirectory(name string, f os.FileInfo) error {
	if f.IsDir() {
		return os.RemoveAll(name)
	} else {
		return fmt.Errorf("%s not a directory", name)
	}
}

func DeleteDirIfExists(name string) error {
	f, e := os.Stat(name)
	if e != nil {
		return nil
	}
	return deleteFileDirectory(name, f)
}

func CreateDir(name string) error {
	f, e := os.Stat(name)
	if e != nil {
		return os.MkdirAll(name, os.ModePerm)
	}
	return fmt.Errorf("%s already exists", f.Name())
}

func CreateDirIfNotExists(name string) error {
	_, e := os.Stat(name)
	if e != nil {
		return os.MkdirAll(name, os.ModePerm)
	}
	return nil
}

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

func ReadBytes(f string) ([]byte, error) {
	rc, err := os.ReadFile(f)
	if err != nil {
		return nil, err
	}
	return rc, nil
}

func CloseSafely(fs io.Closer) {
	_ = fs.Close()
}
