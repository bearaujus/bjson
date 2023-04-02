package bjson

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type bjson struct {
	value           map[string]interface{}
	rootJSONElement []string
	escapedElements map[string]bool
}

type BJSON interface {
	/*
		UnmarshalJSON unmarshals the provided JSON data into the BJSON object.
	*/
	UnmarshalJSON(data []byte) error

	/*
		MarshalJSON marshals the BJSON object into a JSON string.
	*/
	MarshalJSON() ([]byte, error)

	/*
		MarshalJSONPretty marshals the BJSON object into a formatted JSON string.
	*/
	MarshalJSONPretty() ([]byte, error)

	/*
		SetMarshalRootJSONElement sets the root JSON element for Marshaling. The provided targetElement
		is a slice of strings representing the JSON path to the root element.
	*/
	SetMarshalRootJSONElement(targetElement []string) error

	/*
		ResetMarshalRootJSONElement resets the root JSON element to nil, which causes the entire
		BJSON object to be marshaled.
	*/
	ResetMarshalRootJSONElement()

	/*
		RemoveElement removes the JSON element at the provided targetElement. The targetElement
		is a slice of strings representing the JSON path to the element.
		Returns an error if the element is not found.
	*/
	RemoveElement(targetElement []string) error

	/*
		EscapeJSONElement escapes the JSON element at the provided targetElement by marshaling
		it into a JSON string. The targetElement is a slice of strings representing the JSON path
		to the element. Returns an error if the element is not found or is already escaped.
	*/
	EscapeJSONElement(targetElement []string) error

	/*
		UnescapeJSONElement unescapes the JSON element at the provided targetElement by unmarshaling
		it from a JSON string into a JSON object or array. The targetElement is a slice of strings
		representing the JSON path to the element. Returns an error if the element is not found or
		is not escaped or is not a valid JSON object or array.
	*/
	UnescapeJSONElement(targetElement []string) error
}

/*
NewBJSON returns a new instance of BJSON, unmarshaling the provided byte slice of JSON data into it.

Parameters:
  - data: a byte slice containing valid JSON data.

Returns:
  - BJSON: a pointer to a new instance of BJSON, containing the unmarshalled JSON data.
  - error: if an error occurred during the unmarshalling process, this will contain a descriptive error message.

Example:

	bjsonData := []byte(`{"foo": "bar", "num": 42, "nested": {"a": [1, 2, 3], "b": true}}`)
	bj, err := NewBJSON(bjsonData)
	if err != nil {
		log.Fatal(err)
	}
	// bj now contains the unmarshalled JSON data
*/
func NewBJSON(data []byte) (BJSON, error) {
	bj := newBJSON()
	if err := bj.UnmarshalJSON(data); err != nil {
		return nil, err
	}
	return bj, nil
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

func (bj *bjson) SetMarshalRootJSONElement(targetElement []string) error {
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
		return errors.New("targetElement is nil or empty")
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
	if isElementEquals(bj.rootJSONElement, targetElement) {
		return fmt.Errorf("cannot escape the root JSON element: %v", strings.Join(targetElement, "."))
	}

	elementKey := strings.Join(targetElement, ".")
	if bj.escapedElements[elementKey] {
		return fmt.Errorf("element is already escaped: %v", elementKey)
	}

	element, err := bj.getElement(targetElement)
	if err != nil {
		return err
	}

	escaped, err := json.Marshal(element)
	if err != nil {
		return fmt.Errorf("element is not a valid JSON: %v", strings.Join(targetElement, "."))
	}

	bj.escapedElements[elementKey] = true // Mark the element as escaped
	return bj.setElement(targetElement, string(escaped))
}

func (bj *bjson) UnescapeJSONElement(targetElement []string) error {
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
			var v interface{}
			err := json.Unmarshal([]byte(s), &v)
			if err == nil {
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
		if !exists {
			return fmt.Errorf("element not found: %v", strings.Join(targetElement, "."))
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
