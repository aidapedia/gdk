package secret

import (
	"context"
	"errors"

	vault "github.com/hashicorp/vault/api"
	"github.com/mitchellh/mapstructure"
)

type VaultEngine string

const (
	VaultEngineCubbyHole VaultEngine = "cubbyhole"
)

type Vault struct {
	config *vault.Config
	engine VaultEngine
	token  string
	path   string
}

func NewSecretVault(address string, engine, token, path string) Interface {
	vaultEngines := map[string]VaultEngine{
		"cubbyhole": VaultEngineCubbyHole,
	}

	config := vault.DefaultConfig()
	config.Address = address
	engineVal, ok := vaultEngines[engine]
	if !ok {
		return nil
	}
	return &Vault{
		config: config,
		token:  token,
		path:   path,
		engine: engineVal,
	}
}

// Be careful if the
func (v *Vault) GetSecret(ctx context.Context, target interface{}) error {
	if v == nil {
		return errors.New("vault is nil")
	}

	client, err := vault.NewClient(v.config)
	if err != nil {
		return err
	}
	client.SetToken(v.token)

	secret, err := client.KVv1(string(v.engine)).Get(ctx, v.path)
	if err != nil {
		return err
	}

	config := &mapstructure.DecoderConfig{
		Result:               target,
		ErrorUnused:          true, // This is the crucial setting
		TagName:              "mapstructure",
		IgnoreUntaggedFields: true,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	err = decoder.Decode(secret.Data)
	if err != nil {
		return err
	}

	return nil
}
