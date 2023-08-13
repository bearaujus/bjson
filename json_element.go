package bjson

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

func (je *jsonElement) AddElement(value interface{}, targets ...string) (err error) {
	return je.updateElement(updateOptionAdd, value, targets...)
}

func (je *jsonElement) GetElement(targets ...string) (JSONElement, error) {
	return je.getElement(targets...)
}

func (je *jsonElement) SetElement(value interface{}, targets ...string) (err error) {
	return je.updateElement(updateOptionSet, value, targets...)
}

func (je *jsonElement) RemoveElement(targets ...string) (err error) {
	return je.updateElement(updateOptionRemove, nil, targets...)
}

func (je *jsonElement) EscapeElement(targets ...string) error {
	element, err := je.getElement(targets...)
	if err != nil {
		return err
	}

	quotedValue := strconv.Quote(element.String())
	if err != nil {
		return fmt.Errorf("element value is not quoted. value: %v", element)
	}

	var nVal interface{}
	if err := json.Unmarshal([]byte(quotedValue), &nVal); err != nil {
		return err
	}

	if err := je.SetElement(nVal, targets...); err != nil {
		return err
	}

	return nil
}

func (je *jsonElement) UnescapeElement(targets ...string) error {
	element, err := je.getElement(targets...)
	if err != nil {
		return err
	}

	unquotedValue, err := strconv.Unquote(element.String())
	if err != nil {
		return fmt.Errorf("element value is not quoted. value: %v", element)
	}

	var nVal interface{}
	if err := json.Unmarshal([]byte(unquotedValue), &nVal); err != nil {
		return err
	}

	if err := je.SetElement(nVal, targets...); err != nil {
		return err
	}

	return nil
}

func (je *jsonElement) Copy() JSONElement {
	nVal, _ := deepCopy(je.value)
	return &jsonElement{value: nVal}
}

func (je *jsonElement) String() string {
	return string(je.Value())
}

func (je *jsonElement) Value() []byte {
	data, _ := json.Marshal(je.value)
	return data
}

func (je *jsonElement) Len() int {
	switch valObj := je.value.(type) {
	case map[string]interface{}:
		return len(valObj)
	case []interface{}:
		return len(valObj)
	case string:
		return len(valObj)
	}

	return 0
}

func (je *jsonElement) Marshal(isPretty bool, targets ...string) ([]byte, error) {
	sel, err := je.getElement(targets...)
	if err != nil {
		return nil, err
	}

	if isPretty {
		return json.MarshalIndent(sel.value, "", "\t")
	}

	return json.Marshal(sel.value)
}

func (je *jsonElement) MarshalWrite(path string, isPretty bool, targets ...string) error {
	data, err := je.Marshal(isPretty)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, os.ModePerm)
}

type JSONPath []string

func (path JSONPath) String() string {
	return "[" + strconv.QuoteToASCII(path[0]) + "]"
}

func (path JSONPath) Append(target string) JSONPath {
	return append(path, target)
}

func (path JSONPath) AppendIndex(index int) JSONPath {
	return append(path, strconv.Itoa(index))
}

func (je *jsonElement) getElement(targets ...string) (*jsonElement, error) {
	var (
		selector = je.value
		path     = JSONPath{}
		err      error
	)

	for _, target := range targets {
		path = path.Append(target)

		switch obj := selector.(type) {
		case map[string]interface{}:
			selector, err = je.getElementFromMap(obj, target, path.String())
		case []interface{}:
			selector, err = je.getElementFromArray(obj, target, path.String())
		default:
			err = fmt.Errorf("element %v is not found", path)
		}

		if err != nil {
			return nil, err
		}
	}

	return &jsonElement{value: selector}, nil
}

func (je *jsonElement) getElementFromMap(obj map[string]interface{}, target, location string) (interface{}, error) {
	selector, ok := obj[target]
	if !ok {
		return nil, fmt.Errorf("element %v is not found", location)
	}
	return selector, nil
}

func (je *jsonElement) getElementFromArray(arr []interface{}, target, location string) (interface{}, error) {
	idx, err := strconv.Atoi(target)
	if err != nil || idx < 0 || idx >= len(arr) {
		return nil, fmt.Errorf("element %v is not valid index for JSON array", location)
	}
	return arr[idx], nil
}

type updateOption int

const (
	updateOptionAdd    updateOption = iota
	updateOptionSet    updateOption = iota
	updateOptionRemove updateOption = iota
)

func (je *jsonElement) updateElement(option updateOption, value interface{}, targets ...string) error {
	if value != nil {
		nValue, err := deepCopy(value)
		if err != nil {
			return err
		}
		value = nValue
	}

	if len(targets) == 0 {
		return je.updateTopLevelElement(option, value)
	}

	// Append a dummy element
	targets = append(targets, "")

	nValue, err := je.recursiveUpdateElement(option, je.value, value, targets[0], "JSON", targets[1:]...)
	if err != nil {
		return err
	}

	je.value = nValue
	return nil
}

func (je *jsonElement) updateTopLevelElement(option updateOption, value interface{}) error {
	if parentObj, ok := je.value.([]interface{}); ok && option == updateOptionAdd {
		je.value = append(parentObj, value)
		return nil
	}

	return fmt.Errorf("invalid operation for %T", je.value)
}

func (je *jsonElement) recursiveUpdateElement(option updateOption, parent interface{}, value interface{}, currentTarget string, location string, remainingTargets ...string) (interface{}, error) {
	isTail := len(remainingTargets) == 1
	switch parentObj := parent.(type) {
	case map[string]interface{}:
		location += fmt.Sprintf(`["%v"]`, currentTarget)

		child, isExist := parentObj[currentTarget]
		if !isExist && (option == updateOptionSet || option == updateOptionRemove) {
			return nil, fmt.Errorf("key at %v is not found", location)
		}

		if isTail {
			return je.updateTailMapElement(option, parentObj, value, currentTarget, location, child, isExist)
		}

		updatedChild, err := je.recursiveUpdateElement(option, child, value, remainingTargets[0], location, remainingTargets[1:]...)
		if err != nil {
			return nil, err
		}

		parentObj[currentTarget] = updatedChild
		return parent, nil

	case []interface{}:
		location += fmt.Sprintf("[%v]", currentTarget)

		idx, err := strconv.Atoi(currentTarget)
		if err != nil {
			return nil, fmt.Errorf("element %v is not valid index (int) for JSON array. %v", location, err)
		}

		if len(parentObj) == 0 || idx >= len(parentObj) {
			return nil, fmt.Errorf("invalid index for json array at %v", location)
		}

		if isTail {
			return je.updateTailArrayElement(option, parentObj, value, idx, location)
		}

		parentObj[idx], err = je.recursiveUpdateElement(option, parentObj[idx], value, remainingTargets[0], location, remainingTargets[1:]...)
		if err != nil {
			return nil, err
		}

		return parentObj, nil

	default:
		return nil, fmt.Errorf("element '%v' is not found at '%v'", currentTarget, location)
	}
}

func (je *jsonElement) updateTailMapElement(option updateOption, parentObj map[string]interface{}, value interface{}, currentTarget string, location string, child interface{}, isExist bool) (interface{}, error) {
	if arr, ok := child.([]interface{}); (option == updateOptionAdd || option == updateOptionSet) && ok {
		if _, ok := parentObj[currentTarget]; ok && option == updateOptionAdd {
			parentObj[currentTarget] = append(arr, value)
			return parentObj, nil
		}
	}

	if isExist && option == updateOptionAdd {
		return nil, fmt.Errorf("key at %v is already exists", location)
	}

	if option == updateOptionRemove {
		delete(parentObj, currentTarget)
		return parentObj, nil
	}

	parentObj[currentTarget] = value
	return parentObj, nil
}

func (je *jsonElement) updateTailArrayElement(option updateOption, parentObj []interface{}, value interface{}, idx int, location string) (interface{}, error) {
	child := parentObj[idx]
	if arr, ok := child.([]interface{}); option == updateOptionAdd && ok {
		switch option {
		case updateOptionAdd:
			parentObj[idx] = append(arr, value)
			return parentObj, nil
		}
	}

	if option == updateOptionSet {
		parentObj[idx] = value
		return parentObj, nil
	}

	if option == updateOptionRemove {
		parentObj = append(parentObj[:idx], parentObj[idx+1:]...)
		return parentObj, nil
	}

	return nil, fmt.Errorf("cannot update element at: %v", location)
}

func deepCopy(value interface{}) (interface{}, error) {
	var ret interface{}
	switch v := value.(type) {
	case []byte:
		if err := json.Unmarshal(v, &ret); err != nil {
			return nil, err
		}
		return ret, nil

	case *jsonElement:
		value = v.value
	}

	data, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &ret); err != nil {
		return nil, err
	}

	return ret, nil
}
