package converters

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"gopkg.in/yaml.v3"
)

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

func UnmarshalJson(content []byte, t interface{}) error {
	err := json.Unmarshal(content, t)
	if err != nil {
		return fmt.Errorf("error unmarshalling to %T, with error %v", t, err)
	}
	return nil
}

func UnmarshalYaml(content []byte, t interface{}) error {
	err := yaml.Unmarshal(content, t)
	if err != nil {
		return fmt.Errorf("error unmarshalling to %T, with error %v", t, err)
	}
	return nil
}

func MarshalToYamlFile(t interface{}, path string) error {
	data, marshalErr := yaml.Marshal(t)
	if marshalErr != nil {
		return marshalErr
	}
	return os.WriteFile(path, data, os.ModePerm)
}

func MarshalToJsonFile(t interface{}, path string) error {
	data, err := marshalJson(t, false)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, os.ModePerm)
}

func MarshalToJsonPrettyFile(t interface{}, path string) error {
	data, e := marshalJson(t, true)
	if e != nil {
		return e
	}
	return os.WriteFile(path, data, os.ModePerm)
}

func MarshalToJson(item interface{}) string {
	b, err := marshalJson(item, true)
	if err == nil {
		return string(b)
	}
	return ""
}

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

func MarshalToJsonPretty(item interface{}) string {
	data, err := marshalJson(item, true)
	if err == nil {
		return string(data)
	}
	return ""
}

func ReadTo(src io.Reader, t interface{}) error {
	bb := bytes.Buffer{}
	_, err := bb.ReadFrom(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(bb.Bytes(), t)
}
