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
	EscapeElement(targets ...string) error
	UnescapeElement(targets ...string) error
	Copy() BJSON

	String() string
	Len() int
}

func NewBJSON(data interface{}) (BJSON, error) {
	switch d := data.(type) {
	case string:
		data = []byte(d)
	}

	val, err := deepCopy(data)
	if err != nil {
		return nil, err
	}

	return &bjson{value: val}, nil
}

func NewBJSONFromFile(path string) (BJSON, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading file at path '%s': %w", path, err)
	}

	return NewBJSON(data)
}
