package bjson

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

type bjson struct {
	value interface{}
}

type BJSON interface {
	AddElement(value interface{}, targets ...string) error
	GetElement(targets ...string) (BJSON, error)
	SetElement(value interface{}, targets ...string) error
	RemoveElement(targets ...string) error

	Marshal(isPretty bool, targets ...string) ([]byte, error)
	MarshalWrite(path string, isPretty bool, targets ...string) error
	Unmarshal(v any, targets ...string) error

	EscapeElement(targets ...string) error
	UnescapeElement(targets ...string) error

	Len() int
	Copy() (BJSON, error)
	String() string
}

func NewBJSON(data interface{}) (BJSON, error) {
	dataString, ok := data.(string)
	if ok {
		data = []byte(dataString)
	}

	bjValue, err := deepCopy(data)
	if err != nil {
		return nil, err
	}

	return &bjson{value: bjValue}, nil
}

func NewBJSONFromFile(path string) (BJSON, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading file at path '%s': %w", path, err)
	}

	return NewBJSON(data)
}

func MarshalWrite(path string, v interface{}, isPretty bool) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	if isPretty {
		buff := bytes.NewBuffer(nil)
		_ = json.Indent(buff, data, "", "\t")
		data = buff.Bytes()
	}

	return os.WriteFile(path, data, os.ModePerm)
}

func UnmarshalRead(path string, v interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, v)
}
