package config

import (
	"github.com/spf13/viper"
	"testing"
)

func TestInit(t *testing.T) {
	if err := InitConfig(); err != nil {
		t.Fatalf("error on #InitConfig: %v", err)
	}

	if !ConfigLoaded {
		t.Fatalf("ConfigLoaded is false but should be true")
	}
}

func TestDefault(t *testing.T) {
	const key string = "_testing.non.nested.str"
	setNestedDefault(key, "hello")

	if viper.GetString(key) != "hello" {
		t.Fatalf("getting key '%v' did not return \"hello\", as expected", key)
	}
}

func TestDefaultNested(t *testing.T) {
	const key string = "_testing.nested.int"
	setNestedDefault(key, 4)

	if viper.GetInt(key) != 4 {
		t.Fatalf("getting key '%v' did not return 4, as expected", key)
	}
}

