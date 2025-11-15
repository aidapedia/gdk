package secret

import (
	"context"

	"github.com/spf13/viper"
)

type File struct {
	viper    *viper.Viper
	fileName string
}

func NewSecretFile(fileName string) Interface {
	return &File{
		viper:    viper.New(),
		fileName: fileName,
	}
}

func (f *File) GetSecret(ctx context.Context, target interface{}) error {
	f.viper.SetConfigFile(f.fileName)
	if err := f.viper.MergeInConfig(); err != nil {
		return err
	}
	if err := f.viper.Unmarshal(target); err != nil {
		return err
	}
	return nil
}
