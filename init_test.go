package bjson

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
)

func TestNewBJSON(t *testing.T) {
	type args struct {
		data interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "success - from string",
			args:    args{data: `{"a":"str","b":123,"c":true,"d":[],"e":{}}`},
			want:    `{"a":"str","b":123,"c":true,"d":[],"e":{}}`,
			wantErr: false,
		},
		{
			name:    "success - from byte",
			args:    args{data: []byte(`{"a":"str","b":123,"c":true,"d":[],"e":{}}`)},
			want:    `{"a":"str","b":123,"c":true,"d":[],"e":{}}`,
			wantErr: false,
		},
		{
			name: "success - from struct",
			args: args{data: struct {
				Name  string  `json:"name"`
				Score float64 `json:"score"`
			}{
				Name:  "t1",
				Score: 0.95,
			}},
			want:    `{"name":"t1","score":0.95}`,
			wantErr: false,
		},
		{
			name: "success - from list",
			args: args{data: []struct {
				Name  string  `json:"name"`
				Score float64 `json:"score"`
			}{{
				Name:  "t1",
				Score: 0.95,
			}, {
				Name:  "t2",
				Score: 0.52,
			}, {
				Name:  "t3",
				Score: 0.22,
			}}},
			want:    `[{"name":"t1","score":0.95},{"name":"t2","score":0.52},{"name":"t3","score":0.22}]`,
			wantErr: false,
		},
		{
			name: "success - from map",
			args: args{data: map[string]struct {
				Name  string  `json:"name"`
				Score float64 `json:"score"`
			}{"a": {
				Name:  "t1",
				Score: 0.95,
			}, "b": {
				Name:  "t2",
				Score: 0.52,
			}, "c": {
				Name:  "t3",
				Score: 0.22,
			}}},
			want:    `{"a":{"name":"t1","score":0.95},"b":{"name":"t2","score":0.52},"c":{"name":"t3","score":0.22}}`,
			wantErr: false,
		},
		{
			name: "success - from bjson obj itself",
			args: args{data: func() BJSON {
				bj, err := NewBJSON(`{"a":"str","b":123,"c":true,"d":[],"e":{}}`)
				if err != nil {
					t.Fatal(err)
				}
				return bj
			}()},
			want:    `{"a":"str","b":123,"c":true,"d":[],"e":{}}`,
			wantErr: false,
		},
		{
			name:    "success - from empty string",
			args:    args{data: `""`},
			want:    `""`,
			wantErr: false,
		},
		{
			name:    "success - from empty json object",
			args:    args{data: `{}`},
			want:    `{}`,
			wantErr: false,
		},
		{
			name:    "success - from empty json array",
			args:    args{data: `[]`},
			want:    `[]`,
			wantErr: false,
		},
		{
			name:    "success - from boolean",
			args:    args{data: `true`},
			want:    `true`,
			wantErr: false,
		},
		{
			name:    "success - from number",
			args:    args{data: `13.5`},
			want:    `13.5`,
			wantErr: false,
		},
		{
			name:    "fail - invalid json",
			args:    args{data: "asd"},
			want:    "",
			wantErr: true,
		},
		{
			name:    "success - from escaped json",
			args:    args{data: `"{\"arr\":[1,2,3]}"`},
			want:    `"{\"arr\":[1,2,3]}"`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewBJSON(tt.args.data)
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

func TestNewJSONElementFromFile(t *testing.T) {
	// add valid json
	validPath := path.Join(os.TempDir(), "bjson_test_valid.json")
	if err := os.WriteFile(validPath, []byte(`{"a":"str","b":123,"c":true,"d":[],"e":{}}`), os.ModePerm); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(validPath)

	// add invalid json
	invalidPath := path.Join(os.TempDir(), "bjson_test_invalid.json")
	if err := os.WriteFile(invalidPath, []byte("asd"), os.ModePerm); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(invalidPath)

	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "success",
			args:    args{path: validPath},
			want:    `{"a":"str","b":123,"c":true,"d":[],"e":{}}`,
			wantErr: false,
		},
		{
			name:    "fail - invalid json",
			args:    args{path: invalidPath},
			want:    "",
			wantErr: true,
		},
		{
			name:    "fail - file is not found",
			args:    args{path: path.Join(os.TempDir(), "bjson_invalid_path_test")},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewBJSONFromFile(tt.args.path)
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
