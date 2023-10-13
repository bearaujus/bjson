package bjson

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
)

func Test_bjson_AddElement(t *testing.T) {
	type fields struct {
		value interface{}
	}
	type args struct {
		value   interface{}
		targets []string
	}
	types := []struct {
		name  string
		value interface{}
		want  string
	}{
		{
			name:  "string",
			value: "test",
			want:  `"test"`,
		},
		{
			name:  "number",
			value: 23,
			want:  `23`,
		},
		{
			name:  "boolean",
			value: true,
			want:  `true`,
		},
		{
			name:  "json array",
			value: []interface{}{"json_array"},
			want:  `["json_array"]`,
		},
		{
			name:  "json object",
			value: map[string]interface{}{"json_object": "value"},
			want:  `{"json_object":"value"}`,
		},
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name:   "fail - add %v to string",
			fields: fields{value: `"test"`},
			args: args{
				targets: []string{},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name:   "fail - add %v to number",
			fields: fields{value: `123.3`},
			args: args{
				targets: []string{},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name:   "fail - add %v to boolean",
			fields: fields{value: `true`},
			args: args{
				targets: []string{},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name:   "fail - add %v to null",
			fields: fields{value: `null`},
			args: args{
				targets: []string{},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name:   "fail - add %v to root json object",
			fields: fields{value: `{}`},
			args: args{
				targets: []string{},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name:   "fail - add %v to root json object inside json array",
			fields: fields{value: `[{}]`},
			args: args{
				targets: []string{"0"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name:   "fail - add %v to root json object inside json array without an existing index",
			fields: fields{value: `[{}]`},
			args: args{
				targets: []string{"1"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name:   "success - add %v to json object",
			fields: fields{value: `{}`},
			args: args{
				targets: []string{"v1"},
			},
			want:    `{"v1":%v}`,
			wantErr: false,
		},
		{
			name:   "success - add %v to json object inside json array",
			fields: fields{value: `[{}]`},
			args: args{
				targets: []string{"0", "v1"},
			},
			want:    `[{"v1":%v}]`,
			wantErr: false,
		},
		{
			name:   "fail - add %v to json object inside json array without an existing index",
			fields: fields{value: `[{}]`},
			args: args{
				targets: []string{"1", "v1"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name:   "fail - add %v to existing json object",
			fields: fields{value: `{"v1":"val"}`},
			args: args{
				targets: []string{"v1"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name:   "fail - add %v to existing json object inside json array",
			fields: fields{value: `[{"v1":"val"}]`},
			args: args{
				targets: []string{"0", "v1"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name:   "success - add %v to root json array",
			fields: fields{value: `["val"]`},
			args: args{
				targets: []string{},
			},
			want:    `["val",%v]`,
			wantErr: false,
		},
		{
			name:   "success - add %v to root json array inside json object",
			fields: fields{value: `{"test":[]}`},
			args: args{
				targets: []string{"test"},
			},
			want:    `{"test":[%v]}`,
			wantErr: false,
		},
		{
			name:   "fail - add %v to json array",
			fields: fields{value: `[]`},
			args: args{
				targets: []string{"0"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name:   "fail - add %v to json array with invalid index",
			fields: fields{value: `["val"]`},
			args: args{
				targets: []string{"invalid"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name:   "fail - add %v to json array inside json object",
			fields: fields{value: `{"test":[]}`},
			args: args{
				targets: []string{"test", "0"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name:   "fail - add %v to json array inside json object without an existing index",
			fields: fields{value: `{"test":[]}`},
			args: args{
				targets: []string{"test_2", "0"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name:   "fail - add %v to existing json array",
			fields: fields{value: `["val"]`},
			args: args{
				targets: []string{"0"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name:   "fail - add %v to existing json array inside json object",
			fields: fields{value: `{"test":["val"]}`},
			args: args{
				targets: []string{"test", "0"},
			},
			want:    ``,
			wantErr: true,
		},
	}
	for _, ty := range types {
		for _, tt := range tests {
			tt.name = fmt.Sprintf(tt.name, ty.name)
			tt.args.value = ty.value
			tt.want = fmt.Sprintf(tt.want, ty.want)
			t.Run(tt.name, func(t *testing.T) {
				je, err := NewBJSON(tt.fields.value)
				if err != nil {
					assert.FailNow(t, err.Error())
				}

				err = je.AddElement(tt.args.value, tt.args.targets...)
				if tt.wantErr {
					assert.Error(t, err)
					return
				}

				assert.NoError(t, err)
				assert.Equal(t, tt.want, je.String())
			})
		}
	}
}

func Test_bjson_GetElement(t *testing.T) {
	type fields struct {
		value interface{}
	}
	type args struct {
		targets []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// Basic JSON element retrieval cases
		{
			name:   "success - get root string element",
			fields: fields{value: `"test"`},
			args: args{
				targets: []string{},
			},
			want:    `"test"`,
			wantErr: false,
		},
		{
			name:   "success - get root number element",
			fields: fields{value: `10`},
			args: args{
				targets: []string{},
			},
			want:    `10`,
			wantErr: false,
		},
		{
			name:   "success - get root array element",
			fields: fields{value: `[1, 2, 3]`},
			args: args{
				targets: []string{},
			},
			want:    `[1,2,3]`,
			wantErr: false,
		},
		{
			name:   "success - get root object element",
			fields: fields{value: `{"key": "value"}`},
			args: args{
				targets: []string{},
			},
			want:    `{"key":"value"}`,
			wantErr: false,
		},
		{
			name:   "success - get array element by index",
			fields: fields{value: `[1, 2, 3]`},
			args: args{
				targets: []string{"1"},
			},
			want:    `2`,
			wantErr: false,
		},
		{
			name:   "success - get nested object element by path",
			fields: fields{value: `{"outer": {"inner": "value"}}`},
			args: args{
				targets: []string{"outer", "inner"},
			},
			want:    `"value"`,
			wantErr: false,
		},

		// Error cases
		{
			name:   "fail - get non-existent root element",
			fields: fields{value: `{"key": "value"}`},
			args: args{
				targets: []string{"nonexistent"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name:   "fail - get non-existent nested element",
			fields: fields{value: `{"outer": {"inner": "value"}}`},
			args: args{
				targets: []string{"outer", "nonexistent"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name:   "fail - get element from non-object root",
			fields: fields{value: `10`},
			args: args{
				targets: []string{"key"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name:   "fail - get element from non-array root",
			fields: fields{value: `{"key": "value"}`},
			args: args{
				targets: []string{"0"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name:   "fail - get element from non-array non-object nested",
			fields: fields{value: `{"key": "value"}`},
			args: args{
				targets: []string{"key", "inner"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name:   "fail - get element from non-array object nested",
			fields: fields{value: `{"key": {"inner": "value"}}`},
			args: args{
				targets: []string{"key", "inner", "element"},
			},
			want:    ``,
			wantErr: true,
		},

		// Nested JSON element retrieval cases
		{
			name:   "success - get nested object within root object",
			fields: fields{value: `{"parent": {"child": "value"}}`},
			args: args{
				targets: []string{"parent"},
			},
			want:    `{"child":"value"}`,
			wantErr: false,
		},
		{
			name:   "success - get nested string element within root object",
			fields: fields{value: `{"parent": {"child": "value"}}`},
			args: args{
				targets: []string{"parent", "child"},
			},
			want:    `"value"`,
			wantErr: false,
		},
		{
			name:   "success - get nested number element within root object",
			fields: fields{value: `{"parent": {"child": 42}}`},
			args: args{
				targets: []string{"parent", "child"},
			},
			want:    `42`,
			wantErr: false,
		},
		{
			name:   "success - get nested boolean element within root object",
			fields: fields{value: `{"parent": {"child": true}}`},
			args: args{
				targets: []string{"parent", "child"},
			},
			want:    `true`,
			wantErr: false,
		},
		{
			name:   "success - get nested array within root object",
			fields: fields{value: `{"parent": {"child": [1, 2, 3]}}`},
			args: args{
				targets: []string{"parent", "child"},
			},
			want:    `[1,2,3]`,
			wantErr: false,
		},
		{
			name:   "success - get element within nested array",
			fields: fields{value: `{"parent": {"child": [1, 2, 3]}}`},
			args: args{
				targets: []string{"parent", "child", "1"},
			},
			want:    `2`,
			wantErr: false,
		},
		{
			name:    "fail - invalid json array index",
			fields:  fields{value: "[2, 3, 4]"},
			args:    args{targets: []string{"5"}},
			want:    "",
			wantErr: true,
		},
		{
			name:    "fail - invalid json array index 2",
			fields:  fields{value: "[2, 3, 4]"},
			args:    args{targets: []string{"test"}},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			je, err := NewBJSON(tt.fields.value)
			if err != nil {
				t.Fatal(err)
			}

			got, err := je.GetElement(tt.args.targets...)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got.String())
		})
	}
}

func Test_bjson_SetElement(t *testing.T) {
	type fields struct {
		value interface{}
	}
	type args struct {
		value   interface{}
		targets []string
	}
	types := []struct {
		name  string
		value interface{}
		want  string
	}{
		{
			name:  "string",
			value: "test",
			want:  `"test"`,
		},
		{
			name:  "number",
			value: 23,
			want:  `23`,
		},
		{
			name:  "boolean",
			value: true,
			want:  `true`,
		},
		{
			name:  "json array",
			value: []interface{}{"json_array"},
			want:  `["json_array"]`,
		},
		{
			name:  "json object",
			value: map[string]interface{}{"json_object": "value"},
			want:  `{"json_object":"value"}`,
		},
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name:   "success - set %v to string",
			fields: fields{value: `"test"`},
			args: args{
				targets: []string{},
			},
			want:    `%v`,
			wantErr: false,
		},
		{
			name:   "success - set %v to number",
			fields: fields{value: `123.3`},
			args: args{
				targets: []string{},
			},
			want:    `%v`,
			wantErr: false,
		},
		{
			name:   "success - set %v to boolean",
			fields: fields{value: `true`},
			args: args{
				targets: []string{},
			},
			want:    `%v`,
			wantErr: false,
		},
		{
			name:   "success - set %v to null",
			fields: fields{value: `null`},
			args: args{
				targets: []string{},
			},
			want:    `%v`,
			wantErr: false,
		},
		{
			name:   "success - set %v to root json object",
			fields: fields{value: `{}`},
			args: args{
				targets: []string{},
			},
			want:    `%v`,
			wantErr: false,
		},
		{
			name:   "success - set %v to root json object inside json array",
			fields: fields{value: `[{}]`},
			args: args{
				targets: []string{"0"},
			},
			want:    `[%v]`,
			wantErr: false,
		},
		{
			name:   "fail - set %v to root json object inside json array without an existing index",
			fields: fields{value: `[{}]`},
			args: args{
				targets: []string{"1"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name:   "fail - set %v to json object",
			fields: fields{value: `{}`},
			args: args{
				targets: []string{"v1"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name:   "fail - set %v to json object inside json array",
			fields: fields{value: `[{}]`},
			args: args{
				targets: []string{"0", "v1"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name:   "fail - set %v to json object inside json array without an existing index",
			fields: fields{value: `[{}]`},
			args: args{
				targets: []string{"1", "v1"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name:   "success - set %v to existing json object",
			fields: fields{value: `{"v1":"val"}`},
			args: args{
				targets: []string{"v1"},
			},
			want:    `{"v1":%v}`,
			wantErr: false,
		},
		{
			name:   "success - set %v to existing json object inside json array",
			fields: fields{value: `[{"v1":"val"}]`},
			args: args{
				targets: []string{"0", "v1"},
			},
			want:    `[{"v1":%v}]`,
			wantErr: false,
		},
		{
			name:   "success - set %v to root json array",
			fields: fields{value: `["val"]`},
			args: args{
				targets: []string{},
			},
			want:    `%v`,
			wantErr: false,
		},
		{
			name:   "success - set %v to root json array inside json object",
			fields: fields{value: `{"test":[]}`},
			args: args{
				targets: []string{"test"},
			},
			want:    `{"test":%v}`,
			wantErr: false,
		},
		{
			name:   "fail - set %v to json array",
			fields: fields{value: `[]`},
			args: args{
				targets: []string{"0"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name:   "fail - set %v to json array with invalid index",
			fields: fields{value: `["val"]`},
			args: args{
				targets: []string{"invalid"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name:   "fail - set %v to json array inside json object",
			fields: fields{value: `{"test":[]}`},
			args: args{
				targets: []string{"test", "0"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name:   "fail - set %v to json array inside json object without an existing index",
			fields: fields{value: `{"test":[]}`},
			args: args{
				targets: []string{"test_2", "0"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name:   "success - set %v to existing json array",
			fields: fields{value: `["val"]`},
			args: args{
				targets: []string{"0"},
			},
			want:    `[%v]`,
			wantErr: false,
		},
		{
			name:   "success - set %v to existing json array inside json object",
			fields: fields{value: `{"test":["val"]}`},
			args: args{
				targets: []string{"test", "0"},
			},
			want:    `{"test":[%v]}`,
			wantErr: false,
		},
	}

	for _, ty := range types {
		for _, tt := range tests {
			tt.name = fmt.Sprintf(tt.name, ty.name)
			tt.args.value = ty.value
			tt.want = fmt.Sprintf(tt.want, ty.want)
			t.Run(tt.name, func(t *testing.T) {
				je, err := NewBJSON(tt.fields.value)
				if err != nil {
					assert.FailNow(t, err.Error())
				}

				err = je.SetElement(tt.args.value, tt.args.targets...)
				if tt.wantErr {
					assert.Error(t, err)
					return
				}

				assert.NoError(t, err)
				assert.Equal(t, tt.want, je.String())
			})
		}
	}
}

func Test_bjson_RemoveElement(t *testing.T) {
	type fields struct {
		value interface{}
	}
	type args struct {
		targets []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "success - remove from root",
			fields:  fields{value: `{"a":"str","b":123,"c":true,"d":["f",123,456],"e":{"g":"test","h":456,"i":777}}`},
			args:    args{targets: []string{"d"}},
			want:    `{"a":"str","b":123,"c":true,"e":{"g":"test","h":456,"i":777}}`,
			wantErr: false,
		},
		{
			name:    "success - remove json object child",
			fields:  fields{value: `{"a":"str","b":123,"c":true,"d":["f",123,456],"e":{"g":"test","h":456,"i":777}}`},
			args:    args{targets: []string{"e", "g"}},
			want:    `{"a":"str","b":123,"c":true,"d":["f",123,456],"e":{"h":456,"i":777}}`,
			wantErr: false,
		},
		{
			name:    "success - remove json array child",
			fields:  fields{value: `{"a":"str","b":123,"c":true,"d":["f",123,456],"e":{"g":"test","h":456,"i":777}}`},
			args:    args{targets: []string{"d", "1"}},
			want:    `{"a":"str","b":123,"c":true,"d":["f",456],"e":{"g":"test","h":456,"i":777}}`,
			wantErr: false,
		},
		{
			name:    "fail - remove not exist element from root",
			fields:  fields{value: `{"a":"str","b":123,"c":true,"d":["f",123,456],"e":{"g":"test","h":456,"i":777}}`},
			args:    args{targets: []string{"z"}},
			want:    ``,
			wantErr: true,
		},
		{
			name:    "fail - remove not exist element from json object child",
			fields:  fields{value: `{"a":"str","b":123,"c":true,"d":["f",123,456],"e":{"g":"test","h":456,"i":777}}`},
			args:    args{targets: []string{"e", "z"}},
			want:    ``,
			wantErr: true,
		},
		{
			name:    "fail - remove not exist element from json array child",
			fields:  fields{value: `{"a":"str","b":123,"c":true,"d":["f",123,456],"e":{"g":"test","h":456,"i":777}}`},
			args:    args{targets: []string{"d", "99"}},
			want:    ``,
			wantErr: true,
		},
		{
			name:    "fail - remove not exist element from json array child 2",
			fields:  fields{value: `{"a":"str","b":123,"c":true,"d":["f",123,456],"e":{"g":"test","h":456,"i":777}}`},
			args:    args{targets: []string{"d", "z"}},
			want:    ``,
			wantErr: true,
		},

		// Nested JSON element removal cases
		{
			name:    "success - remove nested element from JSON object child",
			fields:  fields{value: `{"a":{"b":{"c":"value"}}}`},
			args:    args{targets: []string{"a", "b", "c"}},
			want:    `{"a":{"b":{}}}`,
			wantErr: false,
		},
		{
			name:    "success - remove element from nested JSON array",
			fields:  fields{value: `{"arr":[1,{"nested":"value"},3]}`},
			args:    args{targets: []string{"arr", "1"}},
			want:    `{"arr":[1,3]}`,
			wantErr: false,
		},
		{
			name:    "success - remove multiple levels of nested elements",
			fields:  fields{value: `{"a":{"b":{"c":{"d":1}}}}`},
			args:    args{targets: []string{"a", "b", "c", "d"}},
			want:    `{"a":{"b":{"c":{}}}}`,
			wantErr: false,
		},
		{
			name:    "fail - remove from non-existent root element",
			fields:  fields{value: `{"a":{"b":{"c":"value"}}}`},
			args:    args{targets: []string{"x"}},
			want:    ``,
			wantErr: true,
		},
		{
			name:    "fail - remove from non-object root element",
			fields:  fields{value: `10`},
			args:    args{targets: []string{"key"}},
			want:    ``,
			wantErr: true,
		},
		{
			name:    "fail - remove from non-object nested element",
			fields:  fields{value: `{"a":{"b":1}}`},
			args:    args{targets: []string{"a", "b", "c"}},
			want:    ``,
			wantErr: true,
		},
		{
			name:    "fail - remove from non-array root element",
			fields:  fields{value: `{"a":{"b":1}}`},
			args:    args{targets: []string{"a", "1"}},
			want:    ``,
			wantErr: true,
		},
		{
			name:    "fail - remove from non-array nested element",
			fields:  fields{value: `{"a":{"b":[1,2,3]}}`},
			args:    args{targets: []string{"a", "b", "3"}},
			want:    ``,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			je, err := NewBJSON(tt.fields.value)
			if err != nil {
				t.Fatal(err)
			}

			err = je.RemoveElement(tt.args.targets...)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}

			var strGot string
			if err == nil {
				strGot = je.String()
			}

			assert.Equal(t, tt.want, strGot)
		})
	}
}

func Test_bjson_Marshal(t *testing.T) {
	type fields struct {
		value interface{}
	}
	type args struct {
		isPretty bool
		targets  []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name:   "success - marshal JSON object",
			fields: fields{value: `{"a": 1, "b": 2}`},
			args:   args{isPretty: false, targets: []string{"a"}},
			want:   `1`,
		},
		{
			name:   "success - marshal JSON array",
			fields: fields{value: `{"arr": [1, 2, 3]}`},
			args:   args{isPretty: true, targets: []string{"arr"}},
			want:   "[\n\t1,\n\t2,\n\t3\n]",
		},
		// ... (other test cases)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			je, err := NewBJSON(tt.fields.value)
			if err != nil {
				t.Fatal(err)
			}

			got, err := je.Marshal(tt.args.isPretty, tt.args.targets...)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.want, string(got))
			}
		})
	}
}

func Test_bjson_MarshalWrite(t *testing.T) {
	type fields struct {
		value interface{}
	}
	type args struct {
		path     string
		isPretty bool
		targets  []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name:   "success - marshal and write JSON object",
			fields: fields{value: `{"a":1,"b":2}`},
			args:   args{path: path.Join(os.TempDir(), "test.json"), isPretty: false, targets: []string{"a"}},
			want:   `1`,
		},
		{
			name:   "success - marshal and write JSON array",
			fields: fields{value: `{"arr":[1, 2, 3]}`},
			args:   args{path: path.Join(os.TempDir(), "test.json"), isPretty: true, targets: nil},
			want:   "{\n\t\"arr\": [\n\t\t1,\n\t\t2,\n\t\t3\n\t]\n}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			je, err := NewBJSON(tt.fields.value)
			if err != nil {
				t.Fatal(err)
			}

			err = je.MarshalWrite(tt.args.path, tt.args.isPretty, tt.args.targets...)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)

				// Read the written file and compare its content
				data, err := os.ReadFile(tt.args.path)
				if err != nil {
					t.Fatal(err)
				}

				assert.Equal(t, tt.want, string(data))
				os.Remove(tt.args.path) // Clean up the temporary file
			}
		})
	}
}

func Test_bjson_Unmarshal(t *testing.T) {
	type fields struct {
		value interface{}
	}
	type args struct {
		v       any
		targets []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "success - string",
			fields: fields{
				value: `"test"`,
			},
			args: args{
				v:       "",
				targets: nil,
			},
			wantErr: false,
		},
		{
			name: "fail - string type not match",
			fields: fields{
				value: `"test"`,
			},
			args: args{
				v:       true,
				targets: nil,
			},
			wantErr: true,
		},
		{
			name: "fail - %v not found",
			fields: fields{
				value: `"test"`,
			},
			args: args{
				v:       "",
				targets: []string{"a", "b", "c", "d"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bj, err := NewBJSON(tt.fields.value)
			if err != nil {
				assert.FailNow(t, err.Error())
			}

			var got interface{}
			switch obj := tt.args.v.(type) {
			case map[string]interface{}:
				err = bj.Unmarshal(&obj, tt.args.targets...)
				got = obj
			case []interface{}:
				err = bj.Unmarshal(&obj, tt.args.targets...)
				got = obj
			case string:
				err = bj.Unmarshal(&obj, tt.args.targets...)
				got = obj
			case bool:
				err = bj.Unmarshal(&obj, tt.args.targets...)
				got = obj
			case float64:
				err = bj.Unmarshal(&obj, tt.args.targets...)
				got = obj
			case nil:
				err = bj.Unmarshal(&obj, tt.args.targets...)
				got = obj
			default:
				err = bj.Unmarshal(&obj, tt.args.targets...)
				got = obj
			}
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			d, err := json.Marshal(got)
			if err != nil {
				assert.FailNow(t, err.Error())
			}

			assert.Equal(t, bj.String(), string(d))
		})
	}
}

func Test_bjson_EscapeElement(t *testing.T) {
	type fields struct {
		value interface{}
	}
	type args struct {
		targets []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "success - escape string",
			fields:  fields{value: `"test"`},
			args:    args{targets: []string{}},
			want:    `"\"test\""`,
			wantErr: false,
		},
		{
			name:    "fail - escape non-existent JSON element",
			fields:  fields{value: `{"a":"value"}`},
			args:    args{targets: []string{"b"}},
			want:    `{"a":"value"}`,
			wantErr: true,
		},
		{
			name:    "success - escape JSON object child",
			fields:  fields{value: `{"a":{"b":"value"}}`},
			args:    args{targets: []string{"a"}},
			want:    `{"a":"{\"b\":\"value\"}"}`,
			wantErr: false,
		},
		{
			name:    "success - escape JSON array child",
			fields:  fields{value: `{"arr":[1,2,3]}`},
			args:    args{targets: []string{"arr"}},
			want:    `{"arr":"[1,2,3]"}`,
			wantErr: false,
		},
		{
			name:    "success - escape nested JSON object child",
			fields:  fields{value: `{"a":{"b":{"c":"value"}}}`},
			args:    args{targets: []string{"a", "b"}},
			want:    `{"a":{"b":"{\"c\":\"value\"}"}}`,
			wantErr: false,
		},
		{
			name:    "success - escape nested JSON array child",
			fields:  fields{value: `{"arr":[[1,2,3],[4,5,6]]}`},
			args:    args{targets: []string{"arr", "0"}},
			want:    `{"arr":["[1,2,3]",[4,5,6]]}`,
			wantErr: false,
		},
		{
			name:    "success - escape JSON object with nested JSON object, array, string, and boolean at depth 2",
			fields:  fields{value: `{"nested":{"obj":{"key":true},"arr":[1,2,3],"str":"value"}}`},
			args:    args{targets: []string{"nested"}},
			want:    `{"nested":"{\"arr\":[1,2,3],\"obj\":{\"key\":true},\"str\":\"value\"}"}`,
			wantErr: false,
		},
		{
			name:    "success - escape JSON boolean",
			fields:  fields{value: `{"bool":true}`},
			args:    args{targets: []string{"bool"}},
			want:    `{"bool":"true"}`,
			wantErr: false,
		},
		{
			name:    "success - escape empty string",
			fields:  fields{value: `""`},
			args:    args{targets: []string{}},
			want:    `""`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			je, err := NewBJSON(tt.fields.value)
			if err != nil {
				t.Fatal(err)
			}

			err = je.EscapeElement(tt.args.targets...)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, je.String())
		})
	}
}

func Test_bjson_UnescapeElement(t *testing.T) {
	type fields struct {
		value interface{}
	}
	type args struct {
		targets []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "success - unescape JSON object child",
			fields:  fields{value: `{"a":"{\"b\":\"value\"}"}`},
			args:    args{targets: []string{"a"}},
			want:    `{"a":{"b":"value"}}`,
			wantErr: false,
		},
		{
			name:    "success - unescape JSON array child",
			fields:  fields{value: `{"arr":"[1,2,3]"}`},
			args:    args{targets: []string{"arr"}},
			want:    `{"arr":[1,2,3]}`,
			wantErr: false,
		},
		{
			name:    "success - unescape nested JSON object child",
			fields:  fields{value: `{"a":{"b":"{\"c\":\"value\"}"}}`},
			args:    args{targets: []string{"a", "b"}},
			want:    `{"a":{"b":{"c":"value"}}}`,
			wantErr: false,
		},
		{
			name:    "success - unescape nested JSON array child",
			fields:  fields{value: `{"arr":"[[1,2,3],[4,5,6]]"}`},
			args:    args{targets: []string{"arr"}},
			want:    `{"arr":[[1,2,3],[4,5,6]]}`,
			wantErr: false,
		},
		{
			name:    "success - unescape JSON object with nested JSON object, array, string, and boolean at depth 2",
			fields:  fields{value: `{"nested":"{\"obj\":{\"key\":true},\"arr\":[1,2,3],\"str\":\"value\"}"}`},
			args:    args{targets: []string{"nested"}},
			want:    `{"nested":{"arr":[1,2,3],"obj":{"key":true},"str":"value"}}`,
			wantErr: false,
		},
		{
			name:    "fail - unescape non-existent JSON element",
			fields:  fields{value: `{"a":"value"}`},
			args:    args{targets: []string{"b"}},
			want:    `{"a":"value"}`,
			wantErr: true,
		},
		{
			name:    "success - unescape JSON boolean",
			fields:  fields{value: `{"bool":"true"}`},
			args:    args{targets: []string{"bool"}},
			want:    `{"bool":true}`,
			wantErr: false,
		},
		{
			name:    "success - unescape empty string",
			fields:  fields{value: `""`},
			args:    args{targets: []string{}},
			want:    `""`,
			wantErr: false,
		},
		{
			name:    "fail - unescape string",
			fields:  fields{value: `"test"`},
			args:    args{targets: []string{}},
			want:    ``,
			wantErr: true,
		},
		{
			name:    "success - from escaped root",
			fields:  fields{value: `"{\"arr\":[1,2,3]}"`},
			args:    args{},
			want:    `{"arr":[1,2,3]}`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			je, err := NewBJSON(tt.fields.value)
			if err != nil {
				t.Fatal(err)
			}

			err = je.UnescapeElement(tt.args.targets...)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, je.String())
		})
	}
}

func Test_bjson_Len(t *testing.T) {
	type fields struct {
		value interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name:   "success - from json object",
			fields: fields{value: `{"a": 1, "b": 2}`},
			want:   2,
		},
		{
			name:   "success - from json array",
			fields: fields{value: `[1, 2, 3]`},
			want:   3,
		},
		{
			name:   "fail - from string",
			fields: fields{value: `"hello"`},
			want:   0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			je, err := NewBJSON(tt.fields.value)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.want, je.Len())
		})
	}
}

func Test_bjson_Copy(t *testing.T) {
	type fields struct {
		value interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "success - copy JSON object",
			fields:  fields{value: `{"a":1,"b":2}`},
			wantErr: false,
		},
		{
			name:    "success - copy JSON array",
			fields:  fields{value: `{"a":1,"arr":[1, 2, 3]}`},
			wantErr: false,
		},
		{
			name:    "fail - corner case: copy invalid data",
			fields:  fields{value: func() {}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var je BJSON
			if tt.wantErr {
				je = &bjson{value: tt.fields.value}
			} else {
				var err error
				je, err = NewBJSON(tt.fields.value)
				if err != nil {
					t.Fatal(err)
				}
			}

			got, err := je.Copy()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, je.String(), got.String())

			// modify the original and verify that the copy remains unchanged
			if err = je.SetElement(42, "a"); err != nil {
				assert.FailNow(t, err.Error())
			}
			assert.NotEqual(t, je.String(), got.String())
		})
	}
}

func Test_bjson_String(t *testing.T) {
	type fields struct {
		value interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "success - JSON object",
			fields: fields{value: `{"a": 1, "b": 2}`},
			want:   `{"a":1,"b":2}`,
		},
		{
			name:   "success - JSON array",
			fields: fields{value: `{"arr": [1, 2, 3]}`},
			want:   `{"arr":[1,2,3]}`,
		},
		// ... (other test cases)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			je, err := NewBJSON(tt.fields.value)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.want, je.String())
		})
	}
}
