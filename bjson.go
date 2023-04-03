/*
Package bjson provides a simple and flexible way to manipulate JSON data in Go.

It allows users to perform various operations on JSON data, such as setting a root JSON element, escaping and unescaping JSON elements, and removing JSON elements.

Features:

  - Read JSON data from a file, string, or byte slice
  - Unmarshal and Marshal JSON data
  - Set a root JSON element for marshaling
  - Escape and Unescape JSON elements
  - Remove JSON elements from the object
  - Write marshaled JSON data to a file
*/
package bjson

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

type bjson struct {
	value           map[string]interface{}
	rootJSONElement []string
	escapedElements map[string]bool
}

/*
BJSON is an interface that provides methods for manipulating JSON data.

Methods:

  - UnmarshalJSON: Unmarshals the provided JSON data into the BJSON object.
  - MarshalJSON: Marshals the BJSON object into a JSON string.
  - MarshalJSONPretty: Marshals the BJSON object into a formatted JSON string.
  - WriteMarshalJSON: Marshals the BJSON object into a JSON string and writes it to a local file.
  - SetMarshalRootJSONElement: Sets the root JSON element for marshaling. The provided targetElement
    is a slice of strings representing the JSON path to the root element.
  - ResetMarshalRootJSONElement: Resets the root JSON element to nil, which causes the entire
    BJSON object to be marshaled.
  - RemoveElement: Removes the JSON element at the provided targetElement. The targetElement
    is a slice of strings representing the JSON path to the element.
    Returns an error if the element is not found.
  - EscapeJSONElement: Escapes the JSON element at the provided targetElement by marshaling
    it into a JSON string. The targetElement is a slice of strings representing the JSON path
    to the element. Returns an error if the element is not found or is already escaped or not a
    valid JSON object or array.
  - UnescapeJSONElement: Unescapes the JSON element at the provided targetElement by unmarshaling
    it from a JSON string into a JSON object or array. The targetElement is a slice of strings
    representing the JSON path to the element. Returns an error if the element is not found or
    is not escaped or is not a valid JSON object or array.
*/
type BJSON interface {
	// UnmarshalJSON unmarshals the provided JSON data into the BJSON object.
	UnmarshalJSON(data []byte) error

	// MarshalJSON marshals the BJSON object into a JSON string.
	MarshalJSON() ([]byte, error)

	// MarshalJSONPretty marshals the BJSON object into a formatted JSON string.
	MarshalJSONPretty() ([]byte, error)

	// WriteMarshalJSON marshals the BJSON object into a JSON string and writes it to a local file.
	WriteMarshalJSON(path string, isPretty bool) error

	// SetMarshalRootJSONElement sets the root JSON element for marshaling. The provided targetElement
	// is a slice of strings representing the JSON path to the root element.
	SetMarshalRootJSONElement(targetElement []string) error

	// ResetMarshalRootJSONElement resets the root JSON element to nil, which causes the entire
	// BJSON object to be marshaled.
	ResetMarshalRootJSONElement()

	// RemoveElement removes the JSON element at the provided targetElement. The targetElement
	// is a slice of strings representing the JSON path to the element.
	// Returns an error if the element is not found.
	RemoveElement(targetElement []string) error

	// EscapeJSONElement escapes the JSON element at the provided targetElement by marshaling
	// it into a JSON string. The targetElement is a slice of strings representing the JSON path
	// to the element. Returns an error if the element is not found or is already escaped or not a
	// valid JSON object or array.
	EscapeJSONElement(targetElement []string) error

	// UnescapeJSONElement unescapes the JSON element at the provided targetElement by unmarshaling
	// it from a JSON string into a JSON object or array. The targetElement is a slice of strings
	// representing the JSON path to the element. Returns an error if the element is not found or
	// is not escaped or is not a valid JSON object or array.
	UnescapeJSONElement(targetElement []string) error
}

// NewBJSONFromByte creates a new BJSON object by unmarshaling the provided JSON data.
// Returns an error if the data is not valid JSON.
func NewBJSONFromByte(data []byte) (BJSON, error) {
	bj := newBJSON()
	if err := bj.UnmarshalJSON(data); err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON from byte data: %w", err)
	}
	return bj, nil
}

// NewBJSONFromString creates a new BJSON object by unmarshaling the provided JSON data in string form.
// Returns an error if the data is not valid JSON.
func NewBJSONFromString(data string) (BJSON, error) {
	bj := newBJSON()
	if err := bj.UnmarshalJSON([]byte(data)); err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON from string data: %w", err)
	}
	return bj, nil
}

// NewBJSONFromFile creates a new BJSON object by reading the JSON data from the specified file path.
// Returns an error if the file cannot be read or the data is not valid JSON.
func NewBJSONFromFile(path string) (BJSON, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading file at path '%s': %w", path, err)
	}

	return NewBJSONFromByte(data)
}

func newBJSON() *bjson {
	return &bjson{
		value:           make(map[string]interface{}),
		rootJSONElement: nil,
		escapedElements: make(map[string]bool),
	}
}

func (bj *bjson) UnmarshalJSON(data []byte) error {
	*bj = *newBJSON()
	err := json.Unmarshal(data, &bj.value)
	if err != nil {
		return err
	}

	bj.populateEscapedElements("", bj.value)
	return nil
}

func (bj *bjson) MarshalJSON() ([]byte, error) {
	if bj.rootJSONElement == nil || len(bj.rootJSONElement) == 0 {
		return json.Marshal(bj.value)
	}

	element, err := bj.getElement(bj.rootJSONElement)
	if err != nil {
		return nil, err
	}

	return json.Marshal(element)
}

func (bj *bjson) MarshalJSONPretty() ([]byte, error) {
	data, err := bj.MarshalJSON()
	if err != nil {
		return nil, err
	}

	// format the encoded json data
	buff := bytes.NewBuffer(nil)
	if err = json.Indent(buff, data, "", "\t"); err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}

func (bj *bjson) WriteMarshalJSON(path string, isPretty bool) error {
	var data []byte
	var err error
	if isPretty {
		data, err = bj.MarshalJSONPretty()
	} else {
		data, err = bj.MarshalJSON()
	}
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func (bj *bjson) SetMarshalRootJSONElement(targetElement []string) error {
	if targetElement == nil || len(targetElement) == 0 {
		return fmt.Errorf("targetElement is nil or empty")
	}

	element, err := bj.getElement(targetElement)
	if err != nil {
		return err
	}

	if element == nil {
		return fmt.Errorf("element not found: %v", strings.Join(targetElement, "."))
	}

	if _, ok := element.(map[string]interface{}); !ok {
		return fmt.Errorf("element is not a JSON object: %v", strings.Join(targetElement, "."))
	}

	bj.rootJSONElement = targetElement
	return nil
}

func (bj *bjson) ResetMarshalRootJSONElement() {
	bj.rootJSONElement = nil
}

func (bj *bjson) RemoveElement(targetElement []string) error {
	if targetElement == nil || len(targetElement) == 0 {
		return fmt.Errorf("targetElement is nil or empty")
	}

	current := bj.value
	for _, key := range targetElement[:len(targetElement)-1] {
		elem, exists := current[key]
		if !exists {
			return fmt.Errorf("element not found: %v", strings.Join(targetElement, "."))
		}
		if sub, ok := elem.(map[string]interface{}); ok {
			current = sub
		} else {
			return fmt.Errorf("invalid element type: %v", strings.Join(targetElement, "."))
		}
	}

	// Check if the element exists
	if _, exists := current[targetElement[len(targetElement)-1]]; !exists {
		return fmt.Errorf("element not found: %v", strings.Join(targetElement, "."))
	}

	delete(current, targetElement[len(targetElement)-1])
	return nil
}

func (bj *bjson) EscapeJSONElement(targetElement []string) error {
	if targetElement == nil || len(targetElement) == 0 {
		return fmt.Errorf("targetElement is nil or empty")
	}

	if isElementEquals(bj.rootJSONElement, targetElement) {
		return fmt.Errorf("cannot escape the root JSON element: %v", strings.Join(targetElement, "."))
	}

	element, err := bj.getElement(targetElement)
	if err != nil {
		return err
	}

	// Check if element is a valid JSON object or array
	isValidJSON := false
	if _, ok := element.(map[string]interface{}); ok {
		isValidJSON = true
	} else if _, ok := element.([]interface{}); ok {
		isValidJSON = true
	}

	if !isValidJSON {
		return fmt.Errorf("element is not a valid JSON object or array: %v", strings.Join(targetElement, "."))
	}

	escaped, err := json.Marshal(element)
	if err != nil {
		return fmt.Errorf("element is not a valid JSON: %v", strings.Join(targetElement, "."))
	}

	bj.escapedElements[strings.Join(targetElement, ".")] = true // Mark the element as escaped
	return bj.setElement(targetElement, string(escaped))
}

func (bj *bjson) UnescapeJSONElement(targetElement []string) error {
	if targetElement == nil || len(targetElement) == 0 {
		return fmt.Errorf("targetElement is nil or empty")
	}

	elementKey := strings.Join(targetElement, ".")
	if !bj.escapedElements[elementKey] {
		return fmt.Errorf("element is not escaped: %v", elementKey)
	}

	element, err := bj.getElement(targetElement)
	if err != nil {
		return err
	}

	escaped, ok := element.(string)
	if !ok {
		return fmt.Errorf("element is not a string: %v", strings.Join(targetElement, "."))
	}

	var unescaped interface{}
	err = json.Unmarshal([]byte(escaped), &unescaped)
	if err != nil {
		return fmt.Errorf("element is not a valid escaped JSON: %v", strings.Join(targetElement, "."))
	}

	// Check if unescaped is a valid JSON object or array
	if _, ok = unescaped.(map[string]interface{}); !ok {
		if _, ok = unescaped.([]interface{}); !ok {
			return fmt.Errorf("element is not a valid JSON object or array: %v", strings.Join(targetElement, "."))
		}
	}

	bj.escapedElements[elementKey] = false // Mark the element as unescaped
	return bj.setElement(targetElement, unescaped)
}

func (bj *bjson) populateEscapedElements(prefix string, m map[string]interface{}) {
	for key, value := range m {
		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}

		if sub, ok := value.(map[string]interface{}); ok {
			bj.populateEscapedElements(fullKey, sub)
		} else if s, ok := value.(string); ok {
			if strings.HasPrefix(s, "{") || strings.HasPrefix(s, "[") {
				bj.escapedElements[fullKey] = true
			}
		}
	}
}

func (bj *bjson) getElement(targetElement []string) (interface{}, error) {
	if targetElement == nil || len(targetElement) == 0 {
		return nil, errors.New("targetElement is nil or empty")
	}

	current := bj.value
	for _, key := range targetElement[:len(targetElement)-1] {
		elem, exists := current[key]
		if !exists {
			return nil, fmt.Errorf("element not found: %v", strings.Join(targetElement, "."))
		}
		if sub, ok := elem.(map[string]interface{}); ok {
			current = sub
		} else {
			return nil, fmt.Errorf("invalid element type: %v", strings.Join(targetElement, "."))
		}
	}

	return current[targetElement[len(targetElement)-1]], nil
}

func (bj *bjson) setElement(targetElement []string, value interface{}) error {
	if targetElement == nil || len(targetElement) == 0 {
		return errors.New("targetElement is nil or empty")
	}

	current := bj.value
	for _, key := range targetElement[:len(targetElement)-1] {
		elem, exists := current[key]
		if !exists || elem == nil {
			return fmt.Errorf("element not found or nil: %v", strings.Join(targetElement, "."))
		}
		if sub, ok := elem.(map[string]interface{}); ok {
			current = sub
		} else {
			return fmt.Errorf("invalid element type: %v", strings.Join(targetElement, "."))
		}
	}

	current[targetElement[len(targetElement)-1]] = value
	return nil
}

func isElementEquals(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
