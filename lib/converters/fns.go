package converters

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"os"
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

func MarshalToJson(item interface{}) string {
	bff := bytes.Buffer{}
	err := json.NewEncoder(&bff).Encode(item)
	if err == nil {
		return bff.String()
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
