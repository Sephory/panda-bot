package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/signal"

	"github.com/sephory/panda-bot/internal/panda"
	"gopkg.in/yaml.v3"
)


func main() {
	config := getConfig()
	panda := panda.New(config)
	if err := panda.Start(); err != nil {
		log.Fatal(err)
	}
}

func getConfig() panda.Config {
	configFile, err := ioutil.ReadFile("./config.yml")
	if err != nil {
		log.Fatal("Could not read configuration")
	}
	config := &panda.Config{}
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		log.Fatal("Could not read configuration")
	}
	return *config
}

func waitForInterrupt() chan os.Signal {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	return interrupt
}
