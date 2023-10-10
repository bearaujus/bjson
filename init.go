package bjson

import (
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
