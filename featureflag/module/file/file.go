package file

import (
	"context"
	"log"
	"reflect"
	"time"

	"github.com/aidapedia/gdk/featureflag/module"
	"github.com/aidapedia/gdk/util"
	"github.com/bytedance/sonic"
)

// File is the file module.
type FeatureFlag struct {
	// root is the root of the feature flag.
	// Contain local cache data of featrue flag.
	root FolderItf

	address string
	prefix  string
}

// New creates a new file module.
func New(address, prefix string) module.Interface {
	root, err := parseFromFileJSON(nil, address)
	if err != nil {
		log.Fatalf("failed to read config root: %v", err)
	}
	if prefix != "" {
		node := root.GetChild(prefix)
		if node.GetType() != nodeTypeFolder {
			log.Fatalf("failed to get config root: %v", err)
		}
		root = node.(FolderItf)
	}

	return &FeatureFlag{
		root:    root,
		address: address,
		prefix:  prefix,
	}
}

func (i *FeatureFlag) GetValue(ctx context.Context, key string) (interface{}, error) {
	return getKeyValue(i.root, key)
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
				root, err := parseFromFileJSON(nil, i.address)
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
