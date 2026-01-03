package config

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/aidapedia/gdk/config/secret"
	"github.com/aidapedia/gdk/environment"
	"github.com/spf13/viper"
)

type Manager struct {
	// Store is the config value store to use.
	store       interface{}
	secretStore interface{}

	// secretType is the type of secret store to use.
	secretType SecretType

	// fileName is the list of config file names to read.
	// The file path will read from CONFIG_FILE_PATH environment variable.
	fileName []string

	// key is the key to unmarshal the config value.
	key string
}

type Option struct {
	TargetStore interface{}
	ConfigKey   string
	FileName    []string

	WithSecret   SecretType
	TargetSecret interface{}
}

func (o *Option) Validate() error {
	// Config Validation
	if o.TargetStore == nil {
		return errors.New("target store cannot be nil")
	}
	v := reflect.ValueOf(o.TargetStore)
	if v.Kind() != reflect.Ptr {
		return errors.New("target store should be pointer")
	}
	if o.WithSecret == "" {
		return nil
	}
	// Secret Validation
	if o.TargetSecret == nil {
		return errors.New("target secret cannot be nil")
	}
	v = reflect.ValueOf(o.TargetSecret)
	if v.Kind() != reflect.Ptr {
		return errors.New("target store should be pointer")
	}
	return nil
}

func New(opt Option) *Manager {
	if err := opt.Validate(); err != nil {
		panic(err)
	}
	return &Manager{
		store:       opt.TargetStore,
		secretStore: opt.TargetSecret,
		secretType:  opt.WithSecret,
		fileName:    opt.FileName,
		key:         opt.ConfigKey,
	}
}

// SetConfig sets the config value store.
func (m *Manager) SetConfig(ctx context.Context) error {
	s := viper.New()
	path := environment.GetConfigPath()
	if path == "" {
		return fmt.Errorf("CONFIG_FILE_PATH environment variable is not set")
	}
	s.AddConfigPath(path)
	for _, fileName := range m.fileName {
		s.SetConfigName(fileName)
		if err := s.MergeInConfig(); err != nil {
			return err
		}
		if err := s.UnmarshalKey(m.key, &m.store); err != nil {
			return err
		}
	}
	return nil
}

// SetSecretStore sets the secret value store.
func (m *Manager) SetSecretStore(ctx context.Context) error {
	switch m.secretType {
	case SecretTypeFile:
		filePath := environment.GetSecretFilePath()
		if filePath == "" {
			return fmt.Errorf("SECRET_FILE_PATH environment variable is not set")
		}
		s := secret.NewSecretFile(filePath)
		return s.GetSecret(ctx, m.secretStore)
	case SecretTypeGSM:
		projectID := environment.GetSecretGSMProjectID()
		if projectID == "" {
			return fmt.Errorf("SECRET_GSM_PROJECT_ID environment variable is not set")
		}
		s := secret.NewSecretGSM(projectID)
		return s.GetSecret(ctx, m.secretStore)
	case SecretTypeVault:
		address := environment.GetSecretVaultAddress()
		if address == "" {
			return fmt.Errorf("SECRET_VAULT_ADDRESS environment variable is not set")
		}
		token := environment.GetSecretVaultToken()
		if token == "" {
			return fmt.Errorf("SECRET_VAULT_TOKEN environment variable is not set")
		}
		path := environment.GetSecretVaultPath()
		engine := environment.GetSecretVaultEngine()
		s := secret.NewSecretVault(address, engine, token, path)
		return s.GetSecret(ctx, m.secretStore)
	default:
		return fmt.Errorf("unknown secret type: %s", m.secretType)
	}
}
