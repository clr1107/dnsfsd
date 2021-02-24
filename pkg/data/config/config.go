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

func setNestedDefault(key string, val interface{}) {
	if viper.Get(key) == nil {
		viper.SetDefault(key, val)
	}
}

// InitConfig initialises the Viper configuration with default values and the
// path.
func InitConfig() error {
	viper.GetViper().SetConfigFile("/etc/dnsfsd/config.yml")
	viper.SetConfigType("yaml")

	setNestedDefault("server.port", 53)
	setNestedDefault("dns.forwards", []string{"1.0.0.1:53", "1.1.1.1:53"})
	setNestedDefault("log.path", "/var/log/dnsfsd/log.txt")
	setNestedDefault("log.verbose", false)
	setNestedDefault("dns.cache", 86400)

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
