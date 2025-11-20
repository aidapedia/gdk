package consul

import (
	"context"

	"github.com/aidapedia/gdk/featureflag/module"
	"github.com/aidapedia/gdk/util"
	"github.com/bytedance/sonic"
	consulCli "github.com/hashicorp/consul/api"
)

type FeatureFlag struct {
	cli    *consulCli.Client
	prefix string
}

// New creates a new Consul client.
func New(address, prefix string) module.Interface {
	cfg := consulCli.DefaultConfig()
	cfg.Address = address
	client, err := consulCli.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	return &FeatureFlag{
		cli:    client,
		prefix: prefix,
	}
}

// GetValue gets the value of the key from Consul.
func (i *FeatureFlag) GetValue(ctx context.Context, key string) (interface{}, error) {
	if i.prefix != "" {
		key = i.prefix + "/" + key
	}
	pair, _, err := i.cli.KV().Get(key, nil)
	if err != nil {
		return nil, err
	}
	return pair.Value, nil
}

func (i *FeatureFlag) GetString(ctx context.Context, key string) (string, error) {
	value, err := i.GetValue(ctx, key)
	if err != nil {
		return "", err
	}
	return util.ToStr(value), nil
}

func (i *FeatureFlag) GetInt(ctx context.Context, key string) (int, error) {
	value, err := i.GetValue(ctx, key)
	if err != nil {
		return 0, err
	}
	return util.ToInt(value), nil
}

func (i *FeatureFlag) GetBool(ctx context.Context, key string) (bool, error) {
	value, err := i.GetValue(ctx, key)
	if err != nil {
		return false, err
	}
	return util.ToBool(value), nil
}

func (i *FeatureFlag) GetStruct(ctx context.Context, key string, v interface{}) error {
	value, err := i.GetValue(ctx, key)
	if err != nil {
		return err
	}
	return sonic.Unmarshal(value.([]byte), v)
}

// Watch watches for feature flag changes.
func (i *FeatureFlag) Watch(ctx context.Context) (chan bool, error) {
	var opts = &consulCli.QueryOptions{}
	ch := make(chan bool)
	for {
		select {
		case <-ctx.Done():
			return ch, ctx.Err()
		default:
		}

		pairs, meta, err := i.cli.KV().List(i.prefix, opts)
		if err != nil {
			return ch, err
		}

		for _, pair := range pairs {
			if pair.Key == i.prefix {
				continue
			}
		}
		ch <- true

		opts.WaitIndex = meta.LastIndex
	}
}
