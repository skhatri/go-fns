package converters

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

type TestStruct struct {
	Name  string `json:"name" yaml:"name"`
	Value int    `json:"value" yaml:"value"`
}

// Custom type that implements Marshaler interface but returns an error
type errorMarshaler struct{}

func (e errorMarshaler) MarshalYAML() (interface{}, error) {
	return nil, fmt.Errorf("marshal error")
}

func (e errorMarshaler) MarshalJSON() ([]byte, error) {
	return nil, fmt.Errorf("marshal error")
}

func TestUnmarshalFile(t *testing.T) {
	t.Run("unmarshal yaml file", func(t *testing.T) {
		tempDir := t.TempDir()
		testFile := filepath.Join(tempDir, "test.yaml")
		content := []byte("name: test\nvalue: 42")

		if err := os.WriteFile(testFile, content, 0644); err != nil {
			t.Fatal(err)
		}

		var result TestStruct
		err := UnmarshalFile(testFile, &result)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if result.Name != "test" || result.Value != 42 {
			t.Errorf("Expected {test 42}, got {%s %d}", result.Name, result.Value)
		}
	})

	t.Run("error on non-existent file", func(t *testing.T) {
		var result TestStruct
		err := UnmarshalFile("nonexistent.yaml", &result)
		if err == nil {
			t.Error("Expected error for non-existent file")
		}
	})

	t.Run("error on invalid yaml", func(t *testing.T) {
		tempDir := t.TempDir()
		testFile := filepath.Join(tempDir, "invalid.yaml")
		content := []byte("name: test\nvalue: invalid")

		if err := os.WriteFile(testFile, content, 0644); err != nil {
			t.Fatal(err)
		}

		var result TestStruct
		err := UnmarshalFile(testFile, &result)
		if err == nil {
			t.Error("Expected error for invalid yaml")
		}
	})
}

func TestUnmarshalJsonFile(t *testing.T) {
	t.Run("unmarshal json file", func(t *testing.T) {
		tempDir := t.TempDir()
		testFile := filepath.Join(tempDir, "test.json")
		content := []byte(`{"name": "test", "value": 42}`)

		if err := os.WriteFile(testFile, content, 0644); err != nil {
			t.Fatal(err)
		}

		var result TestStruct
		err := UnmarshalJsonFile(testFile, &result)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if result.Name != "test" || result.Value != 42 {
			t.Errorf("Expected {test 42}, got {%s %d}", result.Name, result.Value)
		}
	})

	t.Run("error on non-existent file", func(t *testing.T) {
		var result TestStruct
		err := UnmarshalJsonFile("nonexistent.json", &result)
		if err == nil {
			t.Error("Expected error for non-existent file")
		}
	})

	t.Run("error on invalid json", func(t *testing.T) {
		tempDir := t.TempDir()
		testFile := filepath.Join(tempDir, "invalid.json")
		content := []byte(`{"name": "test", "value": "invalid"}`)

		if err := os.WriteFile(testFile, content, 0644); err != nil {
			t.Fatal(err)
		}

		var result TestStruct
		err := UnmarshalJsonFile(testFile, &result)
		if err == nil {
			t.Error("Expected error for invalid json")
		}
	})
}

func TestUnmarshalJson(t *testing.T) {
	t.Run("unmarshal valid json", func(t *testing.T) {
		content := []byte(`{"name": "test", "value": 42}`)
		var result TestStruct
		err := UnmarshalJson(content, &result)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if result.Name != "test" || result.Value != 42 {
			t.Errorf("Expected {test 42}, got {%s %d}", result.Name, result.Value)
		}
	})

	t.Run("error on invalid json", func(t *testing.T) {
		content := []byte(`{"name": "test", "value": "invalid"}`)
		var result TestStruct
		err := UnmarshalJson(content, &result)
		if err == nil {
			t.Error("Expected error for invalid json")
		}
	})
}

func TestUnmarshalYaml(t *testing.T) {
	t.Run("unmarshal valid yaml", func(t *testing.T) {
		content := []byte("name: test\nvalue: 42")
		var result TestStruct
		err := UnmarshalYaml(content, &result)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if result.Name != "test" || result.Value != 42 {
			t.Errorf("Expected {test 42}, got {%s %d}", result.Name, result.Value)
		}
	})

	t.Run("error on invalid yaml", func(t *testing.T) {
		content := []byte("name: test\nvalue: invalid")
		var result TestStruct
		err := UnmarshalYaml(content, &result)
		if err == nil {
			t.Error("Expected error for invalid yaml")
		}
	})
}

func TestMarshalToYamlFile(t *testing.T) {
	t.Run("marshal to yaml file", func(t *testing.T) {
		tempDir := t.TempDir()
		testFile := filepath.Join(tempDir, "test.yaml")
		data := TestStruct{Name: "test", Value: 42}

		err := MarshalToYamlFile(data, testFile)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		content, err := os.ReadFile(testFile)
		if err != nil {
			t.Fatal(err)
		}

		var result TestStruct
		err = UnmarshalYaml(content, &result)
		if err != nil {
			t.Fatal(err)
		}
		if result.Name != "test" || result.Value != 42 {
			t.Errorf("Expected {test 42}, got {%s %d}", result.Name, result.Value)
		}
	})

	t.Run("error on non-writable directory", func(t *testing.T) {
		if os.Getuid() == 0 {
			t.Skip("Skipping test when running as root")
		}

		tempDir := t.TempDir()
		restrictedDir := filepath.Join(tempDir, "restricted")
		if err := os.MkdirAll(restrictedDir, 0); err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(restrictedDir)

		testFile := filepath.Join(restrictedDir, "test.yaml")
		data := TestStruct{Name: "test", Value: 42}

		err := MarshalToYamlFile(data, testFile)
		if err == nil {
			t.Error("Expected error for non-writable directory")
		}
	})

	t.Run("error on marshal failure", func(t *testing.T) {
		tempDir := t.TempDir()
		testFile := filepath.Join(tempDir, "test.yaml")
		data := errorMarshaler{}

		err := MarshalToYamlFile(data, testFile)
		if err == nil {
			t.Error("Expected error for marshal failure")
		}
	})
}

func TestMarshalToJsonFile(t *testing.T) {
	t.Run("marshal to json file", func(t *testing.T) {
		tempDir := t.TempDir()
		testFile := filepath.Join(tempDir, "test.json")
		data := TestStruct{Name: "test", Value: 42}

		err := MarshalToJsonFile(data, testFile)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		content, err := os.ReadFile(testFile)
		if err != nil {
			t.Fatal(err)
		}

		var result TestStruct
		err = UnmarshalJson(content, &result)
		if err != nil {
			t.Fatal(err)
		}
		if result.Name != "test" || result.Value != 42 {
			t.Errorf("Expected {test 42}, got {%s %d}", result.Name, result.Value)
		}
	})

	t.Run("error on non-writable directory", func(t *testing.T) {
		if os.Getuid() == 0 {
			t.Skip("Skipping test when running as root")
		}

		tempDir := t.TempDir()
		restrictedDir := filepath.Join(tempDir, "restricted")
		if err := os.MkdirAll(restrictedDir, 0); err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(restrictedDir)

		testFile := filepath.Join(restrictedDir, "test.json")
		data := TestStruct{Name: "test", Value: 42}

		err := MarshalToJsonFile(data, testFile)
		if err == nil {
			t.Error("Expected error for non-writable directory")
		}
	})

	t.Run("error on marshal failure", func(t *testing.T) {
		tempDir := t.TempDir()
		testFile := filepath.Join(tempDir, "test.json")
		data := errorMarshaler{}

		err := MarshalToJsonFile(data, testFile)
		if err == nil {
			t.Error("Expected error for marshal failure")
		}
	})
}

func TestMarshalToJsonPrettyFile(t *testing.T) {
	t.Run("marshal to pretty json file", func(t *testing.T) {
		tempDir := t.TempDir()
		testFile := filepath.Join(tempDir, "test.json")
		data := TestStruct{Name: "test", Value: 42}

		err := MarshalToJsonPrettyFile(data, testFile)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		content, err := os.ReadFile(testFile)
		if err != nil {
			t.Fatal(err)
		}

		var result TestStruct
		err = UnmarshalJson(content, &result)
		if err != nil {
			t.Fatal(err)
		}
		if result.Name != "test" || result.Value != 42 {
			t.Errorf("Expected {test 42}, got {%s %d}", result.Name, result.Value)
		}

		// Verify pretty formatting
		if !bytes.Contains(content, []byte("\n  ")) {
			t.Error("Expected pretty formatting with indentation")
		}
	})

	t.Run("error on non-writable directory", func(t *testing.T) {
		if os.Getuid() == 0 {
			t.Skip("Skipping test when running as root")
		}

		tempDir := t.TempDir()
		restrictedDir := filepath.Join(tempDir, "restricted")
		if err := os.MkdirAll(restrictedDir, 0); err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(restrictedDir)

		testFile := filepath.Join(restrictedDir, "test.json")
		data := TestStruct{Name: "test", Value: 42}

		err := MarshalToJsonPrettyFile(data, testFile)
		if err == nil {
			t.Error("Expected error for non-writable directory")
		}
	})

	t.Run("error on marshal failure", func(t *testing.T) {
		tempDir := t.TempDir()
		testFile := filepath.Join(tempDir, "test.json")
		data := errorMarshaler{}

		err := MarshalToJsonPrettyFile(data, testFile)
		if err == nil {
			t.Error("Expected error for marshal failure")
		}
	})
}

func TestMarshalToJson(t *testing.T) {
	t.Run("marshal to json string", func(t *testing.T) {
		data := TestStruct{Name: "test", Value: 42}
		result := MarshalToJson(data)
		if result == "" {
			t.Error("Expected non-empty json string")
		}

		var unmarshaled TestStruct
		err := UnmarshalJson([]byte(result), &unmarshaled)
		if err != nil {
			t.Fatal(err)
		}
		if unmarshaled.Name != "test" || unmarshaled.Value != 42 {
			t.Errorf("Expected {test 42}, got {%s %d}", unmarshaled.Name, unmarshaled.Value)
		}
	})

	t.Run("empty string on marshal error", func(t *testing.T) {
		// Create a channel to force json.Marshal to fail
		data := make(chan int)
		result := MarshalToJson(data)
		if result != "" {
			t.Error("Expected empty string on marshal error")
		}
	})
}

func TestMarshalToJsonPretty(t *testing.T) {
	t.Run("marshal to pretty json string", func(t *testing.T) {
		data := TestStruct{Name: "test", Value: 42}
		result := MarshalToJsonPretty(data)
		if result == "" {
			t.Error("Expected non-empty json string")
		}

		var unmarshaled TestStruct
		err := UnmarshalJson([]byte(result), &unmarshaled)
		if err != nil {
			t.Fatal(err)
		}
		if unmarshaled.Name != "test" || unmarshaled.Value != 42 {
			t.Errorf("Expected {test 42}, got {%s %d}", unmarshaled.Name, unmarshaled.Value)
		}

		// Verify pretty formatting
		if !bytes.Contains([]byte(result), []byte("\n  ")) {
			t.Error("Expected pretty formatting with indentation")
		}
	})

	t.Run("empty string on marshal error", func(t *testing.T) {
		// Create a channel to force json.Marshal to fail
		data := make(chan int)
		result := MarshalToJsonPretty(data)
		if result != "" {
			t.Error("Expected empty string on marshal error")
		}
	})
}

func TestReadTo(t *testing.T) {
	t.Run("read from reader to struct", func(t *testing.T) {
		content := []byte(`{"name": "test", "value": 42}`)
		reader := bytes.NewReader(content)

		var result TestStruct
		err := ReadTo(reader, &result)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if result.Name != "test" || result.Value != 42 {
			t.Errorf("Expected {test 42}, got {%s %d}", result.Name, result.Value)
		}
	})

	t.Run("error on invalid json", func(t *testing.T) {
		content := []byte(`{"name": "test", "value": "invalid"}`)
		reader := bytes.NewReader(content)

		var result TestStruct
		err := ReadTo(reader, &result)
		if err == nil {
			t.Error("Expected error for invalid json")
		}
	})

	t.Run("error on read failure", func(t *testing.T) {
		reader := &errorReader{}
		var result TestStruct
		err := ReadTo(reader, &result)
		if err == nil {
			t.Error("Expected error for read failure")
		}
	})
}

// Helper type for testing ReadTo with read error
type errorReader struct{}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, os.ErrInvalid
}
