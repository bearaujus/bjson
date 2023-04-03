package bjson

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestNewBJSONFromByte(t *testing.T) {
	// Test valid JSON
	validJSON := []byte(`{
		"key1": "value1",
		"key2": {
			"subkey1": "subvalue1"
		}
	}`)
	_, err := NewBJSONFromByte(validJSON)
	if err != nil {
		t.Errorf("NewBJSONFromByte failed for valid JSON: %v", err)
	}

	// Test invalid JSON
	invalidJSON := []byte(`{
		"key1": "value1",
		"key2": {
			"subkey1": "subvalue1",
		}
	}`)
	_, err = NewBJSONFromByte(invalidJSON)
	if err == nil {
		t.Error("NewBJSONFromByte should have failed for invalid JSON")
	}
}

func TestNewBJSONFromString(t *testing.T) {
	// Test valid JSON
	validJSON := `{
		"key1": "value1",
		"key2": {
			"subkey1": "subvalue1"
		}
	}`
	_, err := NewBJSONFromString(validJSON)
	if err != nil {
		t.Errorf("NewBJSONFromString failed for valid JSON: %v", err)
	}

	// Test invalid JSON
	invalidJSON := `{
		"key1": "value1",
		"key2": {
			"subkey1": "subvalue1",
		}
	}`
	_, err = NewBJSONFromString(invalidJSON)
	if err == nil {
		t.Error("NewBJSONFromString should have failed for invalid JSON")
	}
}

func TestNewBJSONFromFile(t *testing.T) {
	// Create a temporary file with valid JSON content
	tmpFile, err := ioutil.TempFile("", "valid_json_*.json")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	validJSON := `{
		"key1": "value1",
		"key2": {
			"subkey1": "subvalue1"
		}
	}`

	if _, err := tmpFile.WriteString(validJSON); err != nil {
		t.Fatalf("Failed to write to temporary file: %v", err)
	}

	_, err = NewBJSONFromFile(tmpFile.Name())
	if err != nil {
		t.Errorf("NewBJSONFromFile failed for valid JSON file: %v", err)
	}

	// Create a temporary file with invalid JSON content
	tmpFileInvalid, err := ioutil.TempFile("", "invalid_json_*.json")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpFileInvalid.Name())

	invalidJSON := `{
		"key1": "value1",
		"key2": {
			"subkey1": "subvalue1",
		}
	}`

	if _, err := tmpFileInvalid.WriteString(invalidJSON); err != nil {
		t.Fatalf("Failed to write to temporary file: %v", err)
	}

	_, err = NewBJSONFromFile(tmpFileInvalid.Name())
	if err == nil {
		t.Error("NewBJSONFromFile should have failed for invalid JSON file")
	}
}

func TestUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name: "valid JSON",
			data: `{"key1": "value1", "key2": 42, "key3": [1, 2, 3], "key4": {"key5": "value5"}}`,
			want: map[string]interface{}{
				"key1": "value1",
				"key2": float64(42),
				"key3": []interface{}{float64(1), float64(2), float64(3)},
				"key4": map[string]interface{}{"key5": "value5"},
			},
			wantErr: false,
		},
		{
			name:    "invalid JSON",
			data:    `{"key1": "value1",}`,
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bj := newBJSON()
			err := bj.UnmarshalJSON([]byte(tt.data))
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(bj.value, tt.want) {
				t.Errorf("UnmarshalJSON() = %v, want %v", bj.value, tt.want)
			}
		})
	}
}

func TestMarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		data    map[string]interface{}
		want    string
		wantErr bool
	}{
		{
			name: "valid JSON object",
			data: map[string]interface{}{
				"key1": "value1",
				"key2": float64(42),
				"key3": []interface{}{float64(1), float64(2), float64(3)},
				"key4": map[string]interface{}{"key5": "value5"},
			},
			want:    `{"key1":"value1","key2":42,"key3":[1,2,3],"key4":{"key5":"value5"}}`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bj := newBJSON()
			bj.value = tt.data
			got, err := bj.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				var gotJSON, wantJSON map[string]interface{}
				json.Unmarshal(got, &gotJSON)
				json.Unmarshal([]byte(tt.want), &wantJSON)
				if !reflect.DeepEqual(gotJSON, wantJSON) {
					t.Errorf("MarshalJSON() = %v, want %v", string(got), tt.want)
				}
			}
		})
	}
}

func TestMarshalJSONPretty(t *testing.T) {
	input := `{"name": "John Doe","age": 30,"address": {"city": "New York","country": "USA"}}`

	expectedPrettyJSON := `{
	"name": "John Doe",
	"age": 30,
	"address": {
		"city": "New York",
		"country": "USA"
	}
}`

	testCases := []struct {
		name           string
		input          string
		expectedOutput string
	}{
		{
			name:           "Test Pretty JSON",
			input:          input,
			expectedOutput: expectedPrettyJSON,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bj, err := NewBJSONFromString(tc.input)
			assert.NoError(t, err)

			prettyJSON, err := bj.MarshalJSONPretty()
			assert.NoError(t, err)

			var expectedMap, actualMap map[string]interface{}
			err = json.Unmarshal([]byte(tc.expectedOutput), &expectedMap)
			assert.NoError(t, err)
			err = json.Unmarshal(prettyJSON, &actualMap)
			assert.NoError(t, err)

			assert.Equal(t, expectedMap, actualMap)
		})
	}
}

func TestWriteMarshalJSON(t *testing.T) {
	input := `{
		"name": "John Doe",
		"age": 30,
		"address": {
			"city": "New York",
			"country": "USA"
		}
	}`

	expectedNonPrettyJSON := `{"name":"John Doe","age":30,"address":{"city":"New York","country":"USA"}}`
	expectedPrettyJSON := `{
	"name": "John Doe",
	"age": 30,
	"address": {
		"city": "New York",
		"country": "USA"
	}
}`

	testCases := []struct {
		name           string
		input          string
		expectedOutput string
		isPretty       bool
	}{
		{
			name:           "Test Write Non-Pretty JSON",
			input:          input,
			expectedOutput: expectedNonPrettyJSON,
			isPretty:       false,
		},
		{
			name:           "Test Write Pretty JSON",
			input:          input,
			expectedOutput: expectedPrettyJSON,
			isPretty:       true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bj, err := NewBJSONFromString(tc.input)
			assert.NoError(t, err)

			tmpfile, err := ioutil.TempFile("", "bjson-test")
			assert.NoError(t, err)
			defer os.Remove(tmpfile.Name())

			err = bj.WriteMarshalJSON(tmpfile.Name(), tc.isPretty)
			assert.NoError(t, err)

			actualOutput, err := ioutil.ReadFile(tmpfile.Name())
			assert.NoError(t, err)

			var expectedMap, actualMap map[string]interface{}
			err = json.Unmarshal([]byte(tc.expectedOutput), &expectedMap)
			assert.NoError(t, err)
			err = json.Unmarshal(actualOutput, &actualMap)
			assert.NoError(t, err)

			assert.Equal(t, expectedMap, actualMap)
		})
	}
}

func TestSetMarshalRootJSONElementAndReset(t *testing.T) {
	input := `{
		"name": "John Doe",
		"age": 30,
		"address": {
			"city": "New York",
			"country": "USA"
		}
	}`

	testCases := []struct {
		name           string
		input          string
		expectedOutput string
		targetElement  []string
	}{
		{
			name:           "Test SetMarshalRootJSONElement",
			input:          input,
			expectedOutput: `{"city": "New York","country": "USA"}`,
			targetElement:  []string{"address"},
		},
		{
			name:           "Test ResetMarshalRootJSONElement",
			input:          input,
			expectedOutput: `{"name":"John Doe","age":30,"address":{"city":"New York","country":"USA"}}`,
			targetElement:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bj, err := NewBJSONFromString(tc.input)
			assert.NoError(t, err)

			if tc.targetElement != nil {
				err = bj.SetMarshalRootJSONElement(tc.targetElement)
				assert.NoError(t, err)
			} else {
				bj.ResetMarshalRootJSONElement()
			}

			actualOutput, err := bj.MarshalJSON()
			assert.NoError(t, err)

			var expectedMap, actualMap map[string]interface{}
			err = json.Unmarshal([]byte(tc.expectedOutput), &expectedMap)
			assert.NoError(t, err)
			err = json.Unmarshal(actualOutput, &actualMap)
			assert.NoError(t, err)

			assert.Equal(t, expectedMap, actualMap)
		})
	}
}

func TestRemoveElement(t *testing.T) {
	input := `{
		"name": "John Doe",
		"age": 30,
		"address": {
			"city": "New York",
			"country": "USA"
		}
	}`

	testCases := []struct {
		name           string
		input          string
		expectedOutput string
		targetElement  []string
	}{
		{
			name:           "Test RemoveElement: Remove address",
			input:          input,
			expectedOutput: `{"name":"John Doe","age":30}`,
			targetElement:  []string{"address"},
		},
		{
			name:           "Test RemoveElement: Remove city from address",
			input:          input,
			expectedOutput: `{"name":"John Doe","age":30,"address":{"country":"USA"}}`,
			targetElement:  []string{"address", "city"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bj, err := NewBJSONFromString(tc.input)
			assert.NoError(t, err)

			err = bj.RemoveElement(tc.targetElement)
			assert.NoError(t, err)

			actualOutput, err := bj.MarshalJSON()
			assert.NoError(t, err)

			var expectedMap, actualMap map[string]interface{}
			err = json.Unmarshal([]byte(tc.expectedOutput), &expectedMap)
			assert.NoError(t, err)
			err = json.Unmarshal(actualOutput, &actualMap)
			assert.NoError(t, err)

			assert.Equal(t, expectedMap, actualMap)
		})
	}
}

func TestEscapeUnescapeJSONElement(t *testing.T) {
	input := `{
		"name": "John Doe",
		"age": 30,
		"address": {
			"city": "New York",
			"country": "USA"
		}
	}`

	testCases := []struct {
		name                string
		input               string
		expectedOutput      string
		targetElement       []string
		expectEscapeError   bool
		expectUnescapeError bool
	}{
		{
			name:                "Test Escape and Unescape: Escape address, Unescape address",
			input:               input,
			expectedOutput:      `{"name":"John Doe","age":30,"address":"{\"city\":\"New York\",\"country\":\"USA\"}"}`,
			targetElement:       []string{"address"},
			expectEscapeError:   false,
			expectUnescapeError: false,
		},
		{
			name:                "Test Escape and Unescape: Error when escaping non-object/non-array",
			input:               input,
			targetElement:       []string{"name"},
			expectEscapeError:   true,
			expectUnescapeError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bj, err := NewBJSONFromString(tc.input)
			assert.NoError(t, err)

			err = bj.EscapeJSONElement(tc.targetElement)
			if tc.expectEscapeError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				actualOutput, err := bj.MarshalJSON()
				assert.NoError(t, err)

				var expectedMap, actualMap map[string]interface{}
				err = json.Unmarshal([]byte(tc.expectedOutput), &expectedMap)
				assert.NoError(t, err)
				err = json.Unmarshal(actualOutput, &actualMap)
				assert.NoError(t, err)

				assert.Equal(t, expectedMap, actualMap)

				err = bj.UnescapeJSONElement(tc.targetElement)
				if tc.expectUnescapeError {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)

					actualOutput, err = bj.MarshalJSON()
					assert.NoError(t, err)

					err = json.Unmarshal(actualOutput, &actualMap)
					assert.NoError(t, err)

					// Use the original input to compare unescaped JSON
					err = json.Unmarshal([]byte(input), &expectedMap)
					assert.NoError(t, err)

					assert.Equal(t, expectedMap, actualMap)
				}
			}
		})
	}
}
