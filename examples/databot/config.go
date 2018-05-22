package main

import (
	"io/ioutil"
	"log"

	"github.com/BurntSushi/toml"
)

// Set up a developer account at twitter
// https://apps.twitter.com/
// And don't check your keys into prod :)
type TwitterAuth struct {
	AccessToken    string
	AccessSecret   string
	ConsumerKey    string
	ConsumerSecret string
}

// To avoid burning my AWS budget on what will
// certainly be a popular bot, restrict who can
// message our bot. CAT PICS AREN'T FREE
type Whitelist struct {
	AllowedUsers []string `toml:"allowed_users"`
}

type TomlConfig struct {
	twitter   TwitterAuth `toml:"twitter"`
	whitelist Whitelist
}

func ReadConfig(filename string) *TomlConfig {
	bytes, err := ioutil.ReadFile(filename)
	tomlData := string(bytes)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	var conf TomlConfig
	if _, err := toml.Decode(tomlData, &conf); err != nil {
		log.Fatal(err)
		return nil
	}
	return &conf
}
