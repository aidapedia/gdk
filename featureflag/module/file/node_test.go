package file

import (
	"reflect"
	"testing"
)

func Test_getKeyValue(t *testing.T) {
	rootFolder := NewFolder("root")
	rootFolder.Add(NewKey("key", "valkey"))
	child1 := NewFolder("child1")
	child1.Add(NewKey("keyChild1", "valkeyChild1"))
	rootFolder.Add(child1)

	t.Run("nil node", func(t *testing.T) {
		value, err := getKeyValue(nil, "key")
		if err != nil && err.Error() != "node is nil" {
			t.Errorf("GetKeyValue() error = %v", err)
		}
		if value != nil {
			t.Errorf("GetKeyValue() = %v, want %v", value, "valkey")
		}
	})

	t.Run("root key value found", func(t *testing.T) {
		value, err := getKeyValue(rootFolder, "key")
		if err != nil && err.Error() != "key not found" {
			t.Errorf("GetKeyValue() error = %v", err)
		}
		if value != "valkey" {
			t.Errorf("GetKeyValue() = %v, want %v", value, "valkey")
		}
	})

	t.Run("child key value found", func(t *testing.T) {
		value, err := getKeyValue(rootFolder, "child1/keyChild1")
		if err != nil && err.Error() != "key not found" {
			t.Errorf("GetKeyValue() error = %v", err)
		}
		if value != "valkeyChild1" {
			t.Errorf("GetKeyValue() = %v, want %v", value, "valkeyChild1")
		}
	})

	t.Run("key value not found", func(t *testing.T) {
		value, err := getKeyValue(rootFolder, "child1/keyChild2")
		if err != nil && err.Error() != "key not found" {
			t.Errorf("GetKeyValue() error = %v", err)
		}
		if value != nil {
			t.Errorf("GetKeyValue() = %v, want %v", value, nil)
		}
	})

	t.Run("last node not key", func(t *testing.T) {
		value, err := getKeyValue(rootFolder, "child1")
		if err != nil && err.Error() != "key not found" {
			t.Errorf("GetKeyValue() error = %v", err)
		}
		if value != nil {
			t.Errorf("GetKeyValue() = %v, want %v", value, nil)
		}
	})
}

func Test_parseFromMap(t *testing.T) {
	type args struct {
		root FolderItf
		data map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want FolderItf
	}{
		{
			name: "with child",
			args: args{
				root: NewFolder("root"),
				data: map[string]interface{}{
					"key": "value",
				},
			},
			want: NewFolder("root").Add(NewKey("key", "value")),
		},
		{
			name: "with child",
			args: args{
				root: NewFolder("root"),
				data: map[string]interface{}{
					"folder": map[string]interface{}{
						"key": "value",
					},
				},
			},
			want: NewFolder("root").Add(NewFolder("folder").Add(NewKey("key", "value"))),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseFromMap(tt.args.root, tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseFromMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
