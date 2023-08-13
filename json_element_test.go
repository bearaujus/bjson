package bjson

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
)

func Test_jsonElement_AddElement(t *testing.T) {
	type fields struct {
		value interface{}
	}
	type args struct {
		value   interface{}
		targets []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// STRING
		{
			name:   "fail - add to string",
			fields: fields{value: `"test"`},
			args: args{
				value:   "test",
				targets: []string{},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name:   "fail - add to string at json object",
			fields: fields{value: `{"v1":"test"}`},
			args: args{
				value:   "test",
				targets: []string{"v1"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name:   "fail - add to string at json array",
			fields: fields{value: `["test"]`},
			args: args{
				value:   "asd",
				targets: []string{"0"},
			},
			want:    ``,
			wantErr: true,
		},

		// NUMBER
		{
			name:   "fail - add to number",
			fields: fields{value: `10`},
			args: args{
				value:   "test",
				targets: []string{},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name:   "fail - add to number at json object",
			fields: fields{value: `{"v1":10}`},
			args: args{
				value:   "test",
				targets: []string{"v1"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name:   "fail - add to number at json array",
			fields: fields{value: `[10]`},
			args: args{
				value:   "asd",
				targets: []string{"0"},
			},
			want:    ``,
			wantErr: true,
		},

		// JSON OBJECT
		// - JSON OBJECT ROOT
		{
			name:   "success - add string to json object at root",
			fields: fields{value: `{"v1":"str","v2":0,"v3":[],"v4":{}}`},
			args: args{
				value:   "test",
				targets: []string{"z"},
			},
			want:    `{"v1":"str","v2":0,"v3":[],"v4":{},"z":"test"}`,
			wantErr: false,
		},
		{
			name:   "success - add number to json object at root",
			fields: fields{value: `{"v1":"str","v2":0,"v3":[],"v4":{}}`},
			args: args{
				value:   10,
				targets: []string{"z"},
			},
			want:    `{"v1":"str","v2":0,"v3":[],"v4":{},"z":10}`,
			wantErr: false,
		},
		{
			name:   "success - add json object to json object at root",
			fields: fields{value: `{"v1":"str","v2":0,"v3":[],"v4":{}}`},
			args: args{
				value:   map[string]interface{}{"z": "test"},
				targets: []string{"z"},
			},
			want:    `{"v1":"str","v2":0,"v3":[],"v4":{},"z":{"z":"test"}}`,
			wantErr: false,
		},
		{
			name:   "success - add json arr to json object at root",
			fields: fields{value: `{"v1":"str","v2":0,"v3":[],"v4":{}}`},
			args: args{
				value:   []interface{}{"test"},
				targets: []string{"z"},
			},
			want:    `{"v1":"str","v2":0,"v3":[],"v4":{},"z":["test"]}`,
			wantErr: false,
		},
		{
			name:   "fail - add to existing json object at root",
			fields: fields{value: `{"v1":"str","v2":0,"v3":[],"v4":{}}`},
			args: args{
				value:   "test",
				targets: []string{"v1"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name:   "fail - add to invalid json object at root",
			fields: fields{value: `{"v1":"str","v2":0,"v3":[],"v4":{}}`},
			args: args{
				value:   "test",
				targets: []string{},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name:   "fail - add to invalid json object at root child",
			fields: fields{value: `{"v1":"str","v2":0,"v3":[],"v4":{}}`},
			args: args{
				value:   "test",
				targets: []string{"z", "z"},
			},
			want:    ``,
			wantErr: true,
		},

		// - JSON OBJECT CHILD
		{
			name:   "success - add string to json object at child json object",
			fields: fields{value: `{"v1":"str","v2":0,"v3":[],"v4":{}}`},
			args: args{
				value:   "test",
				targets: []string{"v4", "z"},
			},
			want:    `{"v1":"str","v2":0,"v3":[],"v4":{"z":"test"}}`,
			wantErr: false,
		},
		{
			name:   "success - add number to json object at child json object",
			fields: fields{value: `{"v1":"str","v2":0,"v3":[],"v4":{}}`},
			args: args{
				value:   10,
				targets: []string{"v4", "z"},
			},
			want:    `{"v1":"str","v2":0,"v3":[],"v4":{"z":10}}`,
			wantErr: false,
		},
		{
			name:   "success - add json object to json object at child json object",
			fields: fields{value: `{"v1":"str","v2":0,"v3":[],"v4":{}}`},
			args: args{
				value:   map[string]interface{}{"z": "test"},
				targets: []string{"v4", "z"},
			},
			want:    `{"v1":"str","v2":0,"v3":[],"v4":{"z":{"z":"test"}}}`,
			wantErr: false,
		},
		{
			name:   "success - add json arr to json object at child json object",
			fields: fields{value: `{"v1":"str","v2":0,"v3":[],"v4":{}}`},
			args: args{
				value:   []interface{}{"test"},
				targets: []string{"v4", "z"},
			},
			want:    `{"v1":"str","v2":0,"v3":[],"v4":{"z":["test"]}}`,
			wantErr: false,
		},
		{
			name:   "fail - add to existing json object at child json object",
			fields: fields{value: `{"v1":"str","v2":0,"v3":[],"v4":{"z":"test"}}`},
			args: args{
				value:   "test",
				targets: []string{"v4", "z"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name:   "fail - add to invalid json object at child json object",
			fields: fields{value: `{"v1":"str","v2":0,"v3":[],"v4":{"z":"test"}}`},
			args: args{
				value:   "test",
				targets: []string{"v4"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name:   "fail - add to invalid json object at child json object child",
			fields: fields{value: `{"v1":"str","v2":0,"v3":[],"v4":{"z":"test"}}`},
			args: args{
				value:   "test",
				targets: []string{"v4", "z", "z"},
			},
			want:    ``,
			wantErr: true,
		},

		{
			name:   "success - add string to json object at child json array",
			fields: fields{value: `{"v1":"str","v2":0,"v3":[],"v4":{}}`},
			args: args{
				value:   "test",
				targets: []string{"v3"},
			},
			want:    `{"v1":"str","v2":0,"v3":["test"],"v4":{}}`,
			wantErr: false,
		},
		{
			name:   "success - add number to json object at child json array",
			fields: fields{value: `{"v1":"str","v2":0,"v3":[],"v4":{}}`},
			args: args{
				value:   10,
				targets: []string{"v3"},
			},
			want:    `{"v1":"str","v2":0,"v3":[10],"v4":{}}`,
			wantErr: false,
		},
		{
			name:   "success - add json object to json object at child json array",
			fields: fields{value: `{"v1":"str","v2":0,"v3":[],"v4":{}}`},
			args: args{
				value:   map[string]interface{}{"z": "test"},
				targets: []string{"v3"},
			},
			want:    `{"v1":"str","v2":0,"v3":[{"z":"test"}],"v4":{}}`,
			wantErr: false,
		},
		{
			name:   "success - add json arr to json object at child json array",
			fields: fields{value: `{"v1":"str","v2":0,"v3":[],"v4":{}}`},
			args: args{
				value:   []interface{}{"test"},
				targets: []string{"v3"},
			},
			want:    `{"v1":"str","v2":0,"v3":[["test"]],"v4":{}}`,
			wantErr: false,
		},
		{
			name:   "fail - add to existing json object at child json array",
			fields: fields{value: `{"v1":"str","v2":0,"v3":[{"z":"test"}],"v4":{}}`},
			args: args{
				value:   "test",
				targets: []string{"v3", "0", "z"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name:   "fail - add to invalid json object at child json array",
			fields: fields{value: `{"v1":"str","v2":0,"v3":[{"z":"test"}],"v4":{}}`},
			args: args{
				value:   "test",
				targets: []string{"v3", "0", "z", "z"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name:   "fail - add to invalid json object at child json array with index outbound",
			fields: fields{value: `{"v1":"str","v2":0,"v3":[{"z":"test"}],"v4":{}}`},
			args: args{
				value:   "test",
				targets: []string{"v3", "99", "z", "z"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name:   "fail - add to invalid json object at child json array with invalid index",
			fields: fields{value: `{"v1":"str","v2":0,"v3":[{"z":"test"}],"v4":{}}`},
			args: args{
				value:   "test",
				targets: []string{"v3", "invalid", "z", "z"},
			},
			want:    ``,
			wantErr: true,
		},

		// JSON ARRAY
		// - JSON ARRAY ROOT
		{
			name:   "success - add string to json array at root",
			fields: fields{value: `["str",0,[],{}]`},
			args: args{
				value:   "test",
				targets: []string{},
			},
			want:    `["str",0,[],{},"test"]`,
			wantErr: false,
		},
		{
			name:   "success - add number to json array at root",
			fields: fields{value: `["str",0,[],{}]`},
			args: args{
				value:   10,
				targets: []string{},
			},
			want:    `["str",0,[],{},10]`,
			wantErr: false,
		},
		{
			name:   "success - add json object to json array at root",
			fields: fields{value: `["str",0,[],{}]`},
			args: args{
				value:   map[string]interface{}{"z": "test"},
				targets: []string{},
			},
			want:    `["str",0,[],{},{"z":"test"}]`,
			wantErr: false,
		},
		{
			name:   "success - add json arr to json array at root",
			fields: fields{value: `["str",0,[],{}]`},
			args: args{
				value:   []interface{}{"test"},
				targets: []string{},
			},
			want:    `["str",0,[],{},["test"]]`,
			wantErr: false,
		},

		// - JSON ARRAY CHILD
		{
			name:   "success - add string to json array at child json object",
			fields: fields{value: `["str",0,[],{}]`},
			args: args{
				value:   "test",
				targets: []string{"3", "z"},
			},
			want:    `["str",0,[],{"z":"test"}]`,
			wantErr: false,
		},
		{
			name:   "success - add number to json array at child json object",
			fields: fields{value: `["str",0,[],{}]`},
			args: args{
				value:   10,
				targets: []string{"3", "z"},
			},
			want:    `["str",0,[],{"z":10}]`,
			wantErr: false,
		},
		{
			name:   "success - add json object to array object at child json object",
			fields: fields{value: `["str",0,[],{}]`},
			args: args{
				value:   map[string]interface{}{"z": "test"},
				targets: []string{"3", "z"},
			},
			want:    `["str",0,[],{"z":{"z":"test"}}]`,
			wantErr: false,
		},
		{
			name:   "success - add json arr to json array at child json object",
			fields: fields{value: `["str",0,[],{}]`},
			args: args{
				value:   []interface{}{"test"},
				targets: []string{"3", "z"},
			},
			want:    `["str",0,[],{"z":["test"]}]`,
			wantErr: false,
		},
		{
			name:   "fail - add to existing json array at child json object",
			fields: fields{value: `["str",0,[],{"z":"test"}]`},
			args: args{
				value:   "test",
				targets: []string{"3", "z"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name:   "fail - add to invalid json array at child json object",
			fields: fields{value: `["str",0,[],{"z":"test"}]`},
			args: args{
				value:   "test",
				targets: []string{"3", "z", "z"},
			},
			want:    ``,
			wantErr: true,
		},

		{
			name:   "success - add string to json array at child json object",
			fields: fields{value: `["str",0,[],{}]`},
			args: args{
				value:   "test",
				targets: []string{"2"},
			},
			want:    `["str",0,["test"],{}]`,
			wantErr: false,
		},
		{
			name:   "success - add number to json array at child json array",
			fields: fields{value: `["str",0,[],{}]`},
			args: args{
				value:   10,
				targets: []string{"2"},
			},
			want:    `["str",0,[10],{}]`,
			wantErr: false,
		},
		{
			name:   "success - add json object to json array at child json array",
			fields: fields{value: `["str",0,[],{}]`},
			args: args{
				value:   map[string]interface{}{"z": "test"},
				targets: []string{"2"},
			},
			want:    `["str",0,[{"z":"test"}],{}]`,
			wantErr: false,
		},
		{
			name:   "success - add json arr to json array at child json array",
			fields: fields{value: `["str",0,[],{}]`},
			args: args{
				value:   []interface{}{"test"},
				targets: []string{"2"},
			},
			want:    `["str",0,[["test"]],{}]`,
			wantErr: false,
		},
		{
			name:   "fail - add to existing json array at child json array",
			fields: fields{value: `["str",0,[],{}]`},
			args: args{
				value:   "test",
				targets: []string{"2", "0", "z"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name:   "fail - add to invalid json array at child json array with index outbound",
			fields: fields{value: `["str",0,[],{}]`},
			args: args{
				value:   "test",
				targets: []string{"2", "99", "z", "z"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name:   "fail - add to invalid json array at child json array with invalid index",
			fields: fields{value: `["str",0,[],{}]`},
			args: args{
				value:   "test",
				targets: []string{"2", "invalid", "z", "z"},
			},
			want:    ``,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			je, err := NewJSONElement(tt.fields.value)
			if err != nil {
				t.Fatal(err)
			}

			err = je.AddElement(tt.args.value, tt.args.targets...)
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

func Test_jsonElement_GetElement(t *testing.T) {
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			je, err := NewJSONElement(tt.fields.value)
			if err != nil {
				t.Fatal(err)
			}

			got, err := je.GetElement(tt.args.targets...)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}

			var strGot string
			if err == nil {
				strGot = got.String()
			}

			assert.Equal(t, tt.want, strGot)
		})
	}
}

func Test_jsonElement_SetElement(t *testing.T) {
	type fields struct {
		value interface{}
	}
	type args struct {
		value   interface{}
		targets []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// JSON ARRAY SET ELEMENT
		// - JSON ARRAY ROOT
		{
			name:   "success - set string in json array at root",
			fields: fields{value: `["str",0,[],{}]`},
			args: args{
				value:   "new",
				targets: []string{"0"},
			},
			want:    `["new",0,[],{}]`,
			wantErr: false,
		},
		{
			name:   "success - set number in json array at root",
			fields: fields{value: `["str",0,[],{}]`},
			args: args{
				value:   100,
				targets: []string{"1"},
			},
			want:    `["str",100,[],{}]`,
			wantErr: false,
		},
		{
			name:   "success - set json object in json array at root",
			fields: fields{value: `["str",0,[],{}]`},
			args: args{
				value:   map[string]interface{}{"newKey": "newValue"},
				targets: []string{"2"},
			},
			want:    `["str",0,{"newKey":"newValue"},{}]`,
			wantErr: false,
		},
		{
			name:   "success - set json array in json array at root",
			fields: fields{value: `["str",0,[],{}]`},
			args: args{
				value:   []interface{}{"newElement"},
				targets: []string{"3"},
			},
			want:    `["str",0,[],["newElement"]]`,
			wantErr: false,
		},
		{
			name:   "fail - set element in non-existing json array at root",
			fields: fields{value: `["str",0,[],{}]`},
			args: args{
				value:   "new",
				targets: []string{"99"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name:   "fail - set element in invalid json array at root",
			fields: fields{value: `["str",0,[],{}]`},
			args: args{
				value:   "new",
				targets: []string{},
			},
			want:    ``,
			wantErr: true,
		},

		// - JSON ARRAY CHILD
		{
			name:   "success - set string in json array at child json object",
			fields: fields{value: `{"parent":["str",0,[],{}]}`},
			args: args{
				value:   "new",
				targets: []string{"parent", "0"},
			},
			want:    `{"parent":["new",0,[],{}]}`,
			wantErr: false,
		},
		{
			name:   "success - set number in json array at child json object",
			fields: fields{value: `{"parent":["str",0,[],{}]}`},
			args: args{
				value:   100,
				targets: []string{"parent", "1"},
			},
			want:    `{"parent":["str",100,[],{}]}`,
			wantErr: false,
		},
		{
			name:   "success - set json object in json array at child json object",
			fields: fields{value: `{"parent":["str",0,[],{}]}`},
			args: args{
				value:   map[string]interface{}{"newKey": "newValue"},
				targets: []string{"parent", "2"},
			},
			want:    `{"parent":["str",0,{"newKey":"newValue"},{}]}`,
			wantErr: false,
		},
		{
			name:   "success - set json array in json array at child json object",
			fields: fields{value: `{"parent":["str",0,[],{}]}`},
			args: args{
				value:   []interface{}{"newElement"},
				targets: []string{"parent", "3"},
			},
			want:    `{"parent":["str",0,[],["newElement"]]}`,
			wantErr: false,
		},
		{
			name:   "fail - set element in non-existing json array at child json object",
			fields: fields{value: `{"parent":["str",0,[],{}]}`},
			args: args{
				value:   "new",
				targets: []string{"parent", "99"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name:   "success - set element in json object (not array) at child json object",
			fields: fields{value: `{"parent":["str",0,[],{}]}`},
			args: args{
				value:   "new",
				targets: []string{"parent"},
			},
			want:    `{"parent":"new"}`,
			wantErr: false,
		},
		{
			name:   "fail - set element in invalid json array at child json object child",
			fields: fields{value: `{"parent":["str",0,[],{}]}`},
			args: args{
				value:   "new",
				targets: []string{"parent", "2", "invalid"},
			},
			want:    ``,
			wantErr: true,
		},

		// Edge Cases
		{
			name:   "fail - set to invalid target",
			fields: fields{value: `{"v1":"str","v2":0,"v3":[],"v4":{}}`},
			args: args{
				value:   "new",
				targets: []string{"v4", "z", "nonExistingKey"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name:   "success - set element to invalid type",
			fields: fields{value: `{"v1":"str","v2":0,"v3":[],"v4":{}}`},
			args: args{
				value:   map[string]interface{}{"newKey": "newValue"},
				targets: []string{"v2"},
			},
			want:    `{"v1":"str","v2":{"newKey":"newValue"},"v3":[],"v4":{}}`,
			wantErr: false,
		},
		{
			name:   "fail - set element to invalid index",
			fields: fields{value: `["str",0,[],{}]`},
			args: args{
				value:   "new",
				targets: []string{"5"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name:   "success - set string in nested json array at root",
			fields: fields{value: `{"parent":{"array":["value1","value2"]}}`},
			args: args{
				value:   "newStringValue",
				targets: []string{"parent", "array", "0"},
			},
			want:    `{"parent":{"array":["newStringValue","value2"]}}`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			je, err := NewJSONElement(tt.fields.value)
			if err != nil {
				t.Fatal(err)
			}

			err = je.SetElement(tt.args.value, tt.args.targets...)
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

func Test_jsonElement_RemoveElement(t *testing.T) {
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
			je, err := NewJSONElement(tt.fields.value)
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

func Test_jsonElement_EscapeElement(t *testing.T) {
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
			name:   "success - escape JSON object child",
			fields: fields{value: `{"a":{"b":"value"}}`},
			args:   args{targets: []string{"a"}},
			want:   `{"a":"{\"b\":\"value\"}"}`,
		},
		{
			name:   "success - escape JSON array child",
			fields: fields{value: `{"arr":[1,2,3]}`},
			args:   args{targets: []string{"arr"}},
			want:   `{"arr":"[1,2,3]"}`,
		},
		{
			name:   "success - escape nested JSON object child",
			fields: fields{value: `{"a":{"b":{"c":"value"}}}`},
			args:   args{targets: []string{"a", "b"}},
			want:   `{"a":{"b":"{\"c\":\"value\"}"}}`,
		},
		{
			name:   "success - escape nested JSON array child",
			fields: fields{value: `{"arr":[[1,2,3],[4,5,6]]}`},
			args:   args{targets: []string{"arr", "0"}},
			want:   `{"arr":["[1,2,3]",[4,5,6]]}`,
		},
		{
			name:   "success - escape JSON object with nested JSON object, array, string, and boolean at depth 2",
			fields: fields{value: `{"nested":{"obj":{"key":true},"arr":[1,2,3],"str":"value"}}`},
			args:   args{targets: []string{"nested"}},
			want:   `{"nested":"{\"arr\":[1,2,3],\"obj\":{\"key\":true},\"str\":\"value\"}"}`,
		},
		{
			name:    "fail - escape non-existent JSON element",
			fields:  fields{value: `{"a":"value"}`},
			args:    args{targets: []string{"b"}},
			want:    `{"a":"value"}`,
			wantErr: true,
		},
		{
			name:   "success - escape JSON boolean",
			fields: fields{value: `{"bool":true}`},
			args:   args{targets: []string{"bool"}},
			want:   `{"bool":"true"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			je, err := NewJSONElement(tt.fields.value)
			if err != nil {
				t.Fatal(err)
			}

			err = je.EscapeElement(tt.args.targets...)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.want, je.String())
			}
		})
	}
}

func Test_jsonElement_UnescapeElement(t *testing.T) {
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
			name:   "success - unescape JSON object child",
			fields: fields{value: `{"a":"{\"b\":\"value\"}"}`},
			args:   args{targets: []string{"a"}},
			want:   `{"a":{"b":"value"}}`,
		},
		{
			name:   "success - unescape JSON array child",
			fields: fields{value: `{"arr":"[1,2,3]"}`},
			args:   args{targets: []string{"arr"}},
			want:   `{"arr":[1,2,3]}`,
		},
		{
			name:   "success - unescape nested JSON object child",
			fields: fields{value: `{"a":{"b":"{\"c\":\"value\"}"}}`},
			args:   args{targets: []string{"a", "b"}},
			want:   `{"a":{"b":{"c":"value"}}}`,
		},
		{
			name:   "success - unescape nested JSON array child",
			fields: fields{value: `{"arr":"[[1,2,3],[4,5,6]]"}`},
			args:   args{targets: []string{"arr"}},
			want:   `{"arr":[[1,2,3],[4,5,6]]}`,
		},
		{
			name:   "success - unescape JSON object with nested JSON object, array, string, and boolean at depth 2",
			fields: fields{value: `{"nested":"{\"obj\":{\"key\":true},\"arr\":[1,2,3],\"str\":\"value\"}"}`},
			args:   args{targets: []string{"nested"}},
			want:   `{"nested":{"arr":[1,2,3],"obj":{"key":true},"str":"value"}}`,
		},
		{
			name:    "fail - unescape non-existent JSON element",
			fields:  fields{value: `{"a":"value"}`},
			args:    args{targets: []string{"b"}},
			want:    `{"a":"value"}`,
			wantErr: true,
		},
		{
			name:   "success - unescape JSON boolean",
			fields: fields{value: `{"bool":"true"}`},
			args:   args{targets: []string{"bool"}},
			want:   `{"bool":true}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			je, err := NewJSONElement(tt.fields.value)
			if err != nil {
				t.Fatal(err)
			}

			err = je.UnescapeElement(tt.args.targets...)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.want, je.String())
			}
		})
	}
}

func Test_jsonElement_Marshal(t *testing.T) {
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
			je, err := NewJSONElement(tt.fields.value)
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

func Test_jsonElement_MarshalWrite(t *testing.T) {
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
			want:   `{"a":1,"b":2}`,
		},
		{
			name:   "success - marshal and write JSON array",
			fields: fields{value: `{"arr":[1, 2, 3]}`},
			args:   args{path: path.Join(os.TempDir(), "test.json"), isPretty: true, targets: []string{"arr"}},
			want:   "{\n\t\"arr\": [\n\t\t1,\n\t\t2,\n\t\t3\n\t]\n}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			je, err := NewJSONElement(tt.fields.value)
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

func Test_jsonElement_Copy(t *testing.T) {
	type fields struct {
		value interface{}
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name:   "success - copy JSON object",
			fields: fields{value: `{"a":1,"b":2}`},
		},
		{
			name:   "success - copy JSON array",
			fields: fields{value: `{"a":1,"arr":[1, 2, 3]}`},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			je, err := NewJSONElement(tt.fields.value)
			if err != nil {
				t.Fatal(err)
			}

			copyJe := je.Copy()

			// Modify the original and verify that the copy remains unchanged
			je.SetElement(42, "a")
			assert.NotEqual(t, je.String(), copyJe.String())
		})
	}
}

func Test_jsonElement_String(t *testing.T) {
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
			je, err := NewJSONElement(tt.fields.value)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.want, je.String())
		})
	}
}

func Test_jsonElement_Len(t *testing.T) {
	type fields struct {
		value interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name:   "success - JSON object",
			fields: fields{value: `{"a": 1, "b": 2}`},
			want:   2,
		},
		{
			name:   "success - JSON array",
			fields: fields{value: `[1, 2, 3]`},
			want:   3,
		},
		// ... (other test cases)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			je, err := NewJSONElement(tt.fields.value)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.want, je.Len())
		})
	}
}
