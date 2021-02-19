package persistence

import (
	"fmt"

	"github.com/spf13/viper"
)

var (
	ConfigLoaded bool = false
)

func InitConfig() error {
	viper.GetViper().SetConfigFile("/etc/dnsfsd/config.yml")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err == nil {
		ConfigLoaded = true
	} else {
		return fmt.Errorf("configuration could not be read")
	}

	return nil
}
