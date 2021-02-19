package config

import (
	"fmt"
	"regexp"

	"github.com/spf13/viper"
)

var (
	Loaded bool = false
)

func Init() error {
	viper.GetViper().SetConfigFile("/etc/dnsfsd/config.yml")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err == nil {
		Loaded = true
	} else {
		return fmt.Errorf("configuration could not be read")
	}

	return nil
}

func LoadPatterns() ([]*regexp.Regexp, error) {
	strs := viper.GetStringSlice("patterns")
	compileds := make([]*regexp.Regexp, 0)

	for _, v := range strs {
		if compiled, err := regexp.Compile(v); err != nil {
			return nil, fmt.Errorf("could not read pattern '%v'", v)
		} else {
			compileds = append(compileds, compiled)
		}
	}

	return compileds, nil
}
