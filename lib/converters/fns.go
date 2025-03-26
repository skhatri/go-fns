// Package converters provides utility functions for converting between different data formats
// and handling file-based data operations.
package converters

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

// UnmarshalFile reads a file and unmarshals its contents into the provided data structure.
// The file format is determined by its extension (.json or .yaml).
// Returns an error if the file cannot be read or unmarshaled.
func UnmarshalFile(file string, t interface{}) error {
	content, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("file: [%s], error: [%v]", file, err)
	}
	err = UnmarshalYaml(content, t)
	if err != nil {
		return fmt.Errorf("file: [%s], error: [%v]", file, err)
	}
	return nil
}

// UnmarshalJsonFile reads a JSON file and unmarshals its contents into the provided data structure.
// Returns an error if the file cannot be read or unmarshaled as JSON.
func UnmarshalJsonFile(file string, t interface{}) error {
	content, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("file: [%s], error: [%v]", file, err)
	}
	err = UnmarshalJson(content, t)
	if err != nil {
		return fmt.Errorf("file: [%s], error: [%v]", file, err)
	}
	return nil
}

// UnmarshalJson unmarshals a JSON byte slice into the provided data structure.
// Returns an error if the data cannot be unmarshaled as JSON.
func UnmarshalJson(content []byte, t interface{}) error {
	err := json.Unmarshal(content, t)
	if err != nil {
		return fmt.Errorf("error unmarshalling to %T, with error %v", t, err)
	}
	return nil
}

// UnmarshalYaml unmarshals a YAML byte slice into the provided data structure.
// Returns an error if the data cannot be unmarshaled as YAML.
func UnmarshalYaml(content []byte, t interface{}) error {
	err := yaml.Unmarshal(content, t)
	if err != nil {
		return fmt.Errorf("error unmarshalling to %T, with error %v", t, err)
	}
	return nil
}

// MarshalToYamlFile marshals the provided data structure to YAML and writes it to a file.
// Returns an error if the data cannot be marshaled or if the file cannot be written.
func MarshalToYamlFile(t interface{}, path string) error {
	data, marshalErr := yaml.Marshal(t)
	if marshalErr != nil {
		return marshalErr
	}
	return os.WriteFile(path, data, os.ModePerm)
}

// MarshalToJsonFile marshals the provided data structure to JSON and writes it to a file.
// Returns an error if the data cannot be marshaled or if the file cannot be written.
func MarshalToJsonFile(t interface{}, path string) error {
	data, err := marshalJson(t, false)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, os.ModePerm)
}

// MarshalToJsonPrettyFile marshals the provided data structure to pretty-printed JSON
// and writes it to a file. Returns an error if the data cannot be marshaled or if the file
// cannot be written.
func MarshalToJsonPrettyFile(t interface{}, path string) error {
	data, e := marshalJson(t, true)
	if e != nil {
		return e
	}
	return os.WriteFile(path, data, os.ModePerm)
}

// MarshalToJson marshals the provided data structure to JSON.
// Returns the JSON byte slice and any error that occurred during marshaling.
func MarshalToJson(item interface{}) string {
	b, err := marshalJson(item, true)
	if err == nil {
		return string(b)
	}
	return ""
}

// marshalJson is an internal function that marshals data to JSON with optional pretty printing.
// The pretty parameter determines whether the output should be formatted with indentation.
func marshalJson(item interface{}, pretty bool) ([]byte, error) {
	bff := bytes.Buffer{}
	encoder := json.NewEncoder(&bff)
	if pretty {
		encoder.SetIndent("", "  ")
	}
	err := encoder.Encode(item)
	if err == nil {
		return bff.Bytes(), nil
	}
	return nil, err
}

// MarshalToJsonPretty marshals the provided data structure to pretty-printed JSON.
// Returns the formatted JSON byte slice and any error that occurred during marshaling.
func MarshalToJsonPretty(item interface{}) string {
	data, err := marshalJson(item, true)
	if err == nil {
		return string(data)
	}
	return ""
}

// ReadTo reads data from a reader and unmarshals it into the provided data structure.
// The format is determined by the provided format string ("json" or "yaml").
// Returns an error if the data cannot be read or unmarshaled.
func ReadTo(src io.Reader, t interface{}) error {
	bb := bytes.Buffer{}
	_, err := bb.ReadFrom(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(bb.Bytes(), t)
}
