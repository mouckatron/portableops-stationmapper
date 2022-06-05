package main

import (
	"fmt"

	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("stationmapper")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/stationmapper")
	viper.AddConfigPath("$HOME/.stationmapper")
	viper.AddConfigPath(".")
}

func readConfig() {

	viper.SetDefault("db.type", "sqlite")
	viper.SetDefault("db.path", "stationmapper.sqlite")

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		} else {
			panic(fmt.Errorf("Fatal error config file: %w \n", err))
		}
	}

	fmt.Println(viper.GetString("db.type"))
}
