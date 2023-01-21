package main

import (
	"io/ioutil"
	"log"

	"github.com/sephory/panda-bot/pkg/chat"
	"github.com/sephory/panda-bot/pkg/chat/twitch"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Twitch *twitch.TwitchClientConfiguration
}

func main() {
	panda, err := InitializePanda()
	if err != nil {
		log.Fatal(err)
	}
	err = panda.Start()
	if err != nil {
		log.Fatal(err)
	}
}

func ReadConfig() *Config {
	configFile, err := ioutil.ReadFile("./config.yml")
	if err != nil {
		log.Fatal("Could not read configuration")
	}
	config := &Config{}
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		log.Fatal("Could not read configuration")
	}
	return config
}

func GetChatClients(config *Config) []chat.ChatClient {
	return []chat.ChatClient{
		twitch.New(config.Twitch),
	}
}
