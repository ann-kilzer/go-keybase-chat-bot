package toml

import (
	"log"

	"github.com/BurntSushi/toml"
)

type TomlConfig struct {
	Title     string
	Twitter   TwitterAuth     `toml:"twitter"`
	Whitelist WhitelistConfig `toml:"whitelist"`
}

// Set up a developer account at twitter
// https://apps.twitter.com/
// And don't check your keys into prod :)
type TwitterAuth struct {
	AccessToken    string `toml:"access_token"`
	AccessSecret   string `toml:"access_secret"`
	ConsumerKey    string `toml:"consumer_key"`
	ConsumerSecret string `toml:"consumer_secret"`
}

// To avoid burning my AWS budget on what will
// certainly be a popular bot, restrict who can
// message our bot. CAT PICS AREN'T FREE
type WhitelistConfig struct {
	AllowedUsers []string `toml:"allowed_users"`
}

func ReadConfig(filename string) *TomlConfig {
	var conf TomlConfig
	if _, err := toml.DecodeFile(filename, &conf); err != nil {
		log.Fatal(err)
		return nil
	}
	return &conf
}
