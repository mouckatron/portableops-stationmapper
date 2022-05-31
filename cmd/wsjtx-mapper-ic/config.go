package main

import (
	"flag"
	"strings"
)

type Config struct {
	stationmapperURL []string
}

var config Config

func loadConfig() {
	var stationmapperURLs ArrayOfHostURLs
	const stationmapperURLHelp = "A stationmapper URL to send data to, can be used multiple times to send to multiple stationmapper instances"

	flag.Var(&stationmapperURLs, "url", stationmapperURLHelp)
	flag.Parse()

	if len(stationmapperURLs) == 0 {
		stationmapperURLs.Set("http://localhost:8080")
	}

	config = Config{
		stationmapperURLs,
	}
}

type ArrayOfHostURLs []string

func (a *ArrayOfHostURLs) String() string {
	return strings.Join(*a, ",")
}

func (a *ArrayOfHostURLs) Set(value string) error {
	*a = append(*a, value)
	return nil
}
