package bjson

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

func (bj *bjson) AddElement(value interface{}, targets ...string) (err error) {
	return bj.updateElement(uoAdd, value, newTracer(targets))
}

func (bj *bjson) GetElement(targets ...string) (BJSON, error) {
	return bj.getElement(newTracer(targets))
}

func (bj *bjson) SetElement(value interface{}, targets ...string) (err error) {
	return bj.updateElement(uoSet, value, newTracer(targets))
}

func (bj *bjson) RemoveElement(targets ...string) (err error) {
	return bj.updateElement(uoRemove, nil, newTracer(targets))
}

func (bj *bjson) EscapeElement(targets ...string) error {
	element, err := bj.getElement(newTracer(targets))
	if err != nil {
		return err
	}

	elementStr := element.String()
	if elementStr == `""` {
		return nil
	}

	quotedValue := strconv.Quote(elementStr)
	if err != nil {
		return fmt.Errorf("element value is not quoted. value: %v", element)
	}

	var nVal interface{}
	if err = json.Unmarshal([]byte(quotedValue), &nVal); err != nil {
		return err
	}

	if err = bj.SetElement(nVal, targets...); err != nil {
		return err
	}

	return nil
}

func (bj *bjson) UnescapeElement(targets ...string) error {
	element, err := bj.getElement(newTracer(targets))
	if err != nil {
		return err
	}

	elementStr := element.String()
	if elementStr == `""` {
		return nil
	}

	unquotedValue, err := strconv.Unquote(elementStr)
	if err != nil {
		return fmt.Errorf("element value is not quoted. value: %v", element)
	}

	var nVal interface{}
	if err = json.Unmarshal([]byte(unquotedValue), &nVal); err != nil {
		return err
	}

	if err = bj.SetElement(nVal, targets...); err != nil {
		return err
	}

	return nil
}

func (bj *bjson) Len() int {
	switch valObj := bj.value.(type) {
	case map[string]interface{}:
		return len(valObj)
	case []interface{}:
		return len(valObj)
	}

	return 0
}

func (bj *bjson) Copy() (BJSON, error) {
	nVal, err := deepCopy(bj.value)
	if err != nil {
		return nil, err
	}

	return &bjson{value: nVal}, nil
}

func (bj *bjson) String() string {
	ret, _ := bj.Marshal(false)
	return string(ret)
}

func (bj *bjson) Marshal(isPretty bool, targets ...string) ([]byte, error) {
	sel, err := bj.getElement(newTracer(targets))
	if err != nil {
		return nil, err
	}

	if isPretty {
		return json.MarshalIndent(sel.value, "", "\t")
	}

	return json.Marshal(sel.value)
}

func (bj *bjson) MarshalWrite(path string, isPretty bool, targets ...string) error {
	data, err := bj.Marshal(isPretty, targets...)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, os.ModePerm)
}

func (bj *bjson) Unmarshal(v any, targets ...string) error {
	d, err := bj.Marshal(false, targets...)
	if err != nil {
		return err
	}

	return json.Unmarshal(d, v)
}

func (bj *bjson) getElement(tc *tracer) (*bjson, error) {
	sel := bj.value
	for tc.next() {
		switch obj := sel.(type) {
		case map[string]interface{}:
			var ok bool
			sel, ok = obj[tc.currTarget()]
			if !ok {
				return nil, fmt.Errorf("element %v is not found at %v", tc.currTarget(), tc.passedPath())
			}

		case []interface{}:
			idx, err := strconv.Atoi(tc.currTarget())
			if err != nil {
				return nil, fmt.Errorf("element %v is not valid index (int) for JSON array. %v", tc.passedPath(), err)
			}

			if idx < 0 || idx > len(obj)-1 {
				return nil, fmt.Errorf("invalid index for json array at %v", tc.passedPath())
			}

			sel = obj[idx]

		default:
			return nil, fmt.Errorf("element %v is not found. target: %v", tc.passedPath(), tc.originPath())
		}
	}

	return &bjson{value: sel}, nil
}

func (bj *bjson) updateElement(opt updateOption, value interface{}, tc *tracer) error {
	if value != nil {
		var err error
		value, err = deepCopy(value)
		if err != nil {
			return err
		}
	}

	if tc.isTail() {
		return bj.updateTopLevelElement(opt, value)
	}

	nValue, err := bj.recursiveUpdateElement(opt, bj.value, value, tc)
	if err != nil {
		return err
	}

	bj.value = nValue
	return nil
}

func (bj *bjson) updateTopLevelElement(opt updateOption, value interface{}) error {
	switch opt {
	case uoAdd:
		if parentObj, ok := bj.value.([]interface{}); ok {
			bj.value = append(parentObj, value)
			return nil
		}

	case uoSet:
		bj.value = value
		return nil
	}

	return fmt.Errorf("cannot %v top level element with type %T", opt, bj.value)
}

func (bj *bjson) recursiveUpdateElement(opt updateOption, parent interface{}, value interface{}, tc *tracer) (interface{}, error) {
	for tc.next() {
		target := tc.currTarget()
		switch obj := parent.(type) {
		case map[string]interface{}:
			child, isExist := obj[target]
			if !isExist && (opt == uoSet || opt == uoRemove) {
				return nil, fmt.Errorf("element %v is not found. target: %v", tc.passedPath(), tc.originPath())
			}

			if tc.isTail() {
				return bj.updateTailMapElement(opt, obj, value, child, isExist, tc)
			}

			updatedChild, err := bj.recursiveUpdateElement(opt, child, value, tc)
			if err != nil {
				return nil, err
			}

			obj[target] = updatedChild

		case []interface{}:
			idx, err := strconv.Atoi(target)
			if err != nil {
				return nil, fmt.Errorf("element %v is not valid index (int) for JSON array. %v", tc.passedPath(), err)
			}

			if idx < 0 || idx > len(obj)-1 {
				return nil, fmt.Errorf("invalid index for json array at %v", tc.passedPath())
			}

			if tc.isTail() {
				return bj.updateTailArrayElement(opt, obj, value, idx, tc)
			}

			obj[idx], err = bj.recursiveUpdateElement(opt, obj[idx], value, tc)
			if err != nil {
				return nil, err
			}

		case nil:
			return nil, fmt.Errorf("element %v is not found. target: %v", tc.passedPath(), tc.originPath())

		default:
			return nil, fmt.Errorf("cannot %v element at %v. operation is not allowed for element %T. target: %v", opt, tc.passedPath(), obj, tc.originPath())
		}
	}

	return parent, nil
}

func (bj *bjson) updateTailMapElement(opt updateOption, obj map[string]interface{}, value interface{}, child interface{}, isExist bool, tc *tracer) (interface{}, error) {
	arr, isArr := child.([]interface{})
	switch opt {
	case uoAdd:
		if isArr {
			obj[tc.currTarget()] = append(arr, value)
			break
		}

		if isExist {
			return nil, fmt.Errorf("key %v is already exist", tc.passedPath())
		}

		fallthrough

	case uoSet:
		obj[tc.currTarget()] = value

	case uoRemove:
		delete(obj, tc.currTarget())
	}

	return obj, nil
}

func (bj *bjson) updateTailArrayElement(opt updateOption, parentObj []interface{}, value interface{}, idx int, tc *tracer) (interface{}, error) {
	child := parentObj[idx]
	switch opt {
	case uoAdd:
		arr, ok := child.([]interface{})
		if !ok {
			return nil, fmt.Errorf("cannot add element at: %v", tc.passedPath())
		}

		parentObj[idx] = append(arr, value)

	case uoSet:
		parentObj[idx] = value

	case uoRemove:
		parentObj = append(parentObj[:idx], parentObj[idx+1:]...)
	}

	return parentObj, nil
}

func deepCopy(data interface{}) (interface{}, error) {
	var (
		ret       interface{}
		dataBytes []byte
		typeBytes = false
	)

	switch obj := data.(type) {
	case *bjson:
		data = obj.value
		return deepCopy(data)

	case []byte:
		typeBytes = true
		dataBytes = obj
	}

	if !typeBytes {
		var err error
		dataBytes, err = json.Marshal(data)
		if err != nil {
			return nil, err
		}
	}

	if err := json.Unmarshal(dataBytes, &ret); err != nil {
		return nil, err
	}

	return ret, nil
}
