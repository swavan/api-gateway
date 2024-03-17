package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type ConfigProp struct {
	ConfigFilePath  string
	ConfigFileName  string
	ConfigExtension string
}

func Configuration() *ConfigProp {
	return &ConfigProp{
		ConfigFilePath:  os.Getenv("CONFIG_FILE_PATH"),
		ConfigFileName:  os.Getenv("CONFIG_FILE_NAME"),
		ConfigExtension: os.Getenv("CONFIG_FILE_EXTENSION"),
	}
}

func (prop *ConfigProp) SetFilePath(path string) *ConfigProp {
	prop.ConfigFilePath = path
	return prop
}

func (prop *ConfigProp) SetFileName(data string) *ConfigProp {
	prop.ConfigFileName = data
	return prop
}

func (prop *ConfigProp) SetFileExtension(data string) *ConfigProp {
	prop.ConfigExtension = data
	return prop
}

func New(prop *ConfigProp, data any) error {
	viper.AddConfigPath(prop.ConfigFilePath)
	viper.SetConfigName(prop.ConfigFileName)
	viper.SetConfigType(prop.ConfigExtension)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(`.`, `_`))
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("error loading config file: %s", err)
	}
	if err := viper.Unmarshal(data); err != nil {
		return fmt.Errorf("error reading config file: %s", err)
	}
	return nil
}
