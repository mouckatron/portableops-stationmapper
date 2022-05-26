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
  err := viper.ReadInConfig()
  if err != nil {
    panic(fmt.Errorf("Fatal error config file: %w \n", err))
  }
}
 
