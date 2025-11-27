package file

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
)

type nodeType string

const (
	nodeTypeFolder nodeType = "folder"
	nodeTypeKey    nodeType = "key"
)

type Node interface {
	GetName() string
	GetType() nodeType
}

func getKeyValue(node Node, pathKey string) (interface{}, error) {
	if node == nil {
		return nil, errors.New("node is nil")
	}
	keys := strings.Split(pathKey, "/")
	for index, key := range keys {
		if node.GetType() != nodeTypeFolder {
			return nil, errors.New("key not found")
		}
		node = node.(FolderItf).GetChild(key)
		if node == nil {
			return nil, errors.New("key not found")
		}
		if index == len(keys)-1 {
			if node.GetType() != nodeTypeKey {
				return nil, errors.New("key not found")
			}
			return node.(KeyItf).GetValue(), nil
		}
	}
	return nil, errors.New("key not found")
}

func parseFromFileJSON(root FolderItf, path string) (FolderItf, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	err = json.Unmarshal(file, &data)
	if err != nil {
		return nil, err
	}
	if root == nil {
		root = NewFolder("root")
	}
	return parseFromMap(root, data), nil
}

func parseFromMap(root FolderItf, data map[string]interface{}) FolderItf {
	for key, value := range data {
		if _, ok := value.(map[string]interface{}); ok {
			root.Add(parseFromMap(NewFolder(key), value.(map[string]interface{})))
		} else {
			root.Add(NewKey(key, value))
		}
	}
	return root
}
