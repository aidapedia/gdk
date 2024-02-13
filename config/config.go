package config

import (
	"github.com/spf13/viper"
)

type databaseInstances struct {
	files  []FileConfig
	key    string
	target interface{}
}

type FileConfig struct {
	FilePath string
	Files    []string
}

type Config struct {
	viper  *viper.Viper
	config databaseInstances
}

func NewConfig(file []FileConfig, key string, target interface{}) *Config {
	return &Config{
		viper: viper.New(),
		config: databaseInstances{
			files:  file,
			key:    key,
			target: target,
		},
	}
}

func (c *Config) ReadConfigFiles(names []string) error {
	for _, name := range names {
		err := c.ReadConfigFile(name)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Config) ReadConfigFile(name string) error {
	c.viper.SetConfigName(name)
	return c.viper.MergeInConfig()
}

func (c *Config) AddConfigPaths(paths []string) {
	for i := range paths {
		c.viper.AddConfigPath(paths[i])
	}
}

func (c *Config) SetConfig() error {
	var filePaths, files []string
	for i := range c.config.files {
		filePaths = append(filePaths, c.config.files[i].FilePath)
		files = append(files, c.config.files[i].Files...)
	}

	c.AddConfigPaths(filePaths)
	err := c.ReadConfigFiles(files)
	if err != nil {
		return err
	}

	err = c.viper.UnmarshalKey(c.config.key, c.config.target)
	return err
}

func (c *Config) GetConfig() interface{} {
	return c.config.target
}
