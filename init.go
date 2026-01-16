package bjson

import (
	"fmt"
	"os"
	"reflect"
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

func Marshal(v interface{}, isPretty bool, targets ...string) ([]byte, error) {
	bj, err := NewBJSON(v)
	if err != nil {
		return nil, err
	}

	return bj.Marshal(isPretty, targets...)
}

func Unmarshal(data any, v any, targets ...string) error {
	bj, err := NewBJSON(data)
	if err != nil {
		return err
	}

	return bj.Unmarshal(&v, targets...)
}

func UnmarshalAndUnwarp[T any](data any, targets ...string) (*T, error) {
	var t T

	rv := reflect.ValueOf(t)
	if rv.Kind() == reflect.Pointer {
		return nil, fmt.Errorf("T must be a non pointer object. retrieved T type: %T", t)
	}

	if err := Unmarshal(data, &t, targets...); err != nil {
		return nil, err
	}

	return &t, nil
}

func MarshalWrite(path string, v interface{}, isPretty bool, targets ...string) error {
	data, err := Marshal(v, isPretty, targets...)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, os.ModePerm)
}

func UnmarshalRead(path string, v interface{}, targets ...string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return Unmarshal(data, v, targets...)
}
