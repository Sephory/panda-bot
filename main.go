package main

import (
	"io/ioutil"

	"github.com/apex/log"
	"github.com/sephory/panda-bot/internal/panda"
	"github.com/sephory/panda-bot/pkg/chat"
	"github.com/sephory/panda-bot/pkg/chat/twitch"
	"github.com/sephory/panda-bot/pkg/chat/youtube"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Bot     *panda.BotConfig
	Twitch  *twitch.TwitchClientConfiguration
	YouTube *youtube.YouTubeClientConfiguration
}

func main() {
	log.SetLevel(log.DebugLevel)
	panda, err := InitializePanda()
	if err != nil {
		log.Fatal(err.Error())
	}
	err = panda.Start()
	if err != nil {
		log.Fatal(err.Error())
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

func GetBotConfig(config *Config) *panda.BotConfig {
	return config.Bot
}

func GetChatClients(config *Config) []chat.ChatClient {
	return []chat.ChatClient{
		twitch.New(config.Twitch),
		youtube.New(config.YouTube),
	}
}
