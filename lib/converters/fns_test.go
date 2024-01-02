package converters

import (
	"github.com/skhatri/go-fns/lib/fs"
	"os"
	"path/filepath"
	"testing"
)

func TestUnmarshalYaml(t *testing.T) {
	t.Run("UnmarshalYaml with valid content", func(t *testing.T) {
		content := []byte("key: value")

		var result map[string]string
		err := UnmarshalYaml(content, &result)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		expected := map[string]string{"key": "value"}
		if !mapStringEqual(result, expected) {
			t.Errorf("Expected %v, but got %v", expected, result)
		}
	})

	t.Run("UnmarshalYaml with invalid content", func(t *testing.T) {
		content := []byte("invalid: yaml: content")

		var result map[string]string
		err := UnmarshalYaml(content, &result)

		if err == nil {
			t.Error("Expected an error, but got nil")
		}
	})
}
func mapStringEqual(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for key, value := range a {
		if bValue, ok := b[key]; !ok || value != bValue {
			return false
		}
	}
	return true
}

func writeContent(name string, content string) {
	f, _ := os.OpenFile(name, os.O_RDWR, 0644)
	defer fs.CloseSafely(f)
	_, err := f.Write([]byte(content))
	if err != nil {
		panic(err)
	}
}

func TestUnmarshalYamlFile(t *testing.T) {
	tempDir := t.TempDir()
	createTestFiles(tempDir, []string{"valid/d.yaml", "invalid/b.yaml"})
	writeContent(filepath.Join(tempDir, "valid/d.yaml"), "abc: \"1234\"")
	writeContent(filepath.Join(tempDir, "invalid/b.yaml"), "stuff")

	t.Run("Unmarshall Valid Files", func(t *testing.T) {
		m := make(map[string]string)
		fileName := filepath.Join(tempDir, "valid/d.yaml")
		err := UnmarshalFile(fileName, m)
		if err != nil {
			t.Errorf("error unmarshalling file %s", fileName)
		}
		if m["abc"] != "1234" {
			t.Errorf("expected content is not the same for unmarshalled file")
		}
	})

	t.Run("Unmarshall Invalid Files", func(t *testing.T) {
		for _, name := range []string{
			"invalid/b.yaml",
			"invalid/nonexist.yaml",
		} {
			m := make(map[string]string)
			fileName := filepath.Join(tempDir, name)
			err := UnmarshalFile(fileName, m)
			if err == nil {
				t.Errorf("expected error when unmarshalling invalid file %s", fileName)
			}
			if len(m) != 0 {
				t.Errorf("expected map to be empty")
			}
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
