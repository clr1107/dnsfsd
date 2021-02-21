package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

var (
	// ConfigLoaded is a flag for whether the configuration has been loaded.
	ConfigLoaded bool = false
)

// InitConfig initialises the Viper configuration with default values and the
// path.
func InitConfig() error {
	viper.GetViper().SetConfigFile("/etc/dnsfsd/config.yml")
	viper.SetConfigType("yaml")

	viper.Sub("server").SetDefault("port", 53)
	viper.Sub("dns").SetDefault("forwards", []string{"1.0.0.1:53", "1.1.1.1:53"})
	viper.Sub("log").SetDefault("path", "/var/log/dnsfsd/log.txt")
	viper.Sub("log").SetDefault("verbose", false)
	viper.Sub("dns").SetDefault("cache", 86400)

	if err := viper.ReadInConfig(); err == nil {
		ConfigLoaded = true
	} else {
		return fmt.Errorf("configuration could not be read")
	}

	return nil
}

func GetCacheTime() time.Duration {
	x := viper.GetInt("cache")
	return time.Duration(x) * time.Second
}
