package bjson

import (
	"fmt"
	"os"
)

type jsonElement struct {
	value interface{}
}

type JSONElement interface {
	AddElement(value interface{}, targets ...string) error
	GetElement(targets ...string) (JSONElement, error)
	SetElement(value interface{}, targets ...string) error
	RemoveElement(targets ...string) error

	Marshal(isPretty bool, targets ...string) ([]byte, error)
	MarshalWrite(path string, isPretty bool, targets ...string) error
	EscapeElement(targets ...string) error
	UnescapeElement(targets ...string) error
	Copy() JSONElement

	String() string
	Len() int
}

func NewJSONElement(data interface{}) (JSONElement, error) {
	switch d := data.(type) {
	case string:
		data = []byte(d)
	case *jsonElement:
		return NewJSONElement(d.value)
	}

	val, err := deepCopy(data)
	if err != nil {
		return nil, err
	}

	return &jsonElement{value: val}, nil
}

func NewJSONElementFromFile(path string) (JSONElement, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading file at path '%s': %w", path, err)
	}

	return NewJSONElement(data)
}
