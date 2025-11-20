package file

import (
	"context"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/aidapedia/gdk/featureflag/module"
	"github.com/aidapedia/gdk/util"
	"github.com/bytedance/sonic"
)

type Root struct {
	Dir `json:"root"`
}

type Dir struct {
	KVs   map[string]interface{} `json:"keys,omitempty"`
	Child map[string]Dir         `json:"child,omitempty"`
}

// File is the file module.
type FeatureFlag struct {
	// root is the root of the feature flag.
	// Contain local cache data of featrue flag.
	root Dir
	// filepath is the path to the feature flag file.
	address string

	prefix string
}

// New creates a new file module.
func New(address, prefix string) module.Interface {
	configKeys, err := readConfigRoot(address, prefix)
	if err != nil {
		log.Fatalf("failed to read config root: %v", err)
	}

	return &FeatureFlag{
		root:    configKeys,
		address: address,
		prefix:  prefix,
	}
}

// GetValue returns the value of the key.
func (i *FeatureFlag) GetValue(ctx context.Context, key string) (interface{}, error) {
	dir := strings.Split(key, "/")
	if len(dir) == 0 {
		return nil, module.ErrKeyNotFound
	}

	var KVs map[string]interface{}
	if len(dir) > 1 {
		dirs, err := findDir(i.root, dir[:len(dir)-1])
		if err != nil {
			return nil, err
		}
		KVs = dirs.Child[dir[len(dir)-2]].KVs
	} else {
		KVs = i.root.KVs
	}

	k := dir[len(dir)-1]
	if value, ok := KVs[k]; ok {
		return value, nil
	}

	return nil, module.ErrKeyNotFound
}

func (i *FeatureFlag) GetBool(ctx context.Context, key string) (bool, error) {
	value, err := i.GetValue(ctx, key)
	if err != nil {
		return false, err
	}
	return util.ToBool(value), nil
}

func (i *FeatureFlag) GetInt(ctx context.Context, key string) (int, error) {
	value, err := i.GetValue(ctx, key)
	if err != nil {
		return 0, err
	}
	return util.ToInt(value), nil
}

func (i *FeatureFlag) GetString(ctx context.Context, key string) (string, error) {
	value, err := i.GetValue(ctx, key)
	if err != nil {
		return "", err
	}
	return util.ToStr(value), nil
}

func (i *FeatureFlag) GetStruct(ctx context.Context, key string, v interface{}) error {
	value, err := i.GetValue(ctx, key)
	if err != nil {
		return err
	}
	return sonic.UnmarshalString(util.ToStr(value), v)
}

// Watch watches for feature flag changes.
// Each time the feature flag changes, it will send true to the channel.
func (i *FeatureFlag) Watch(ctx context.Context) (chan bool, error) {
	var (
		err error
		do  = make(chan bool)
	)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				var root Dir
				root, err = readConfigRoot(i.address, i.prefix)
				if err != nil {
					return
				}
				if reflect.DeepEqual(i.root, root) {
					time.Sleep(time.Second * 5)
					continue
				}
				i.root = root
				do <- true
				time.Sleep(time.Second * 5)
			}
		}
	}()
	return do, err
}

// findDir finds the directory in the config.
func findDir(config Dir, filepath []string) (Dir, error) {
	for index := 0; index < len(filepath); index++ {
		_, found := config.Child[filepath[index]]
		if !found {
			return Dir{}, fmt.Errorf("failed to find dir %s in file %s", filepath[index], filepath)
		}
		if index == len(filepath)-1 {
			break
		}
		config = config.Child[filepath[index]]
	}
	return config, nil
}

func readConfigRoot(filepath, prefix string) (Dir, error) {
	openFile, err := os.Open(filepath)
	if err != nil {
		log.Fatalf("failed to open file %s: %v", filepath, err)
	}
	defer openFile.Close()

	stat, err := openFile.Stat()
	if err != nil {
		log.Fatalf("failed to stat file %s: %v", filepath, err)
	}

	byteValue := make([]byte, stat.Size())
	_, err = openFile.Read(byteValue)
	if err != nil {
		log.Fatalf("failed to read file %s: %v", filepath, err)
	}

	// parse byteValue to KV
	var root Root
	err = sonic.Unmarshal(byteValue, &root)
	if err != nil {
		log.Fatalf("failed to unmarshal file %s: %v", filepath, err)
	}

	var configKeys Dir
	prefix = strings.TrimPrefix(prefix, "/")
	if prefix != "" {
		configKeys, err = findDir(root.Dir, strings.Split(prefix, "/"))
		if err != nil {
			return Dir{}, err
		}
	} else {
		configKeys = root.Dir
	}
	return configKeys, nil
}
