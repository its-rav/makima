package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/caarlos0/env/v8"
)

const (
	DefaultConfigFile = "config.json"
)

func processError(err error) {
	fmt.Println(err)
	os.Exit(2)
}

func (cfg *ConsumerConfig) FromFile(fileName string) {
	f, err := os.Open(fileName)
	if err != nil {
		processError(err)
	}
	defer f.Close()

	// Parse json file
	err = json.NewDecoder(f).Decode(cfg)
	if err != nil {
		processError(err)
	}
}

func (cfg *ConsumerConfig) FromEnv() {
	if err := env.Parse(cfg); err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func (cfg *ConsumerConfig) Load() {
	// if file exists, load from file
	// else load from env
	if _, err := os.Stat(DefaultConfigFile); err == nil {
		cfg.FromFile(DefaultConfigFile)
	} else {
		cfg.FromEnv()
	}
}

func (cfg *CollectorConfig) FromFile(fileName string) {
	f, err := os.Open(fileName)
	if err != nil {
		processError(err)
	}
	defer f.Close()

	// Parse json file
	err = json.NewDecoder(f).Decode(cfg)
	if err != nil {
		processError(err)
	}
}

func (cfg *CollectorConfig) FromEnv() {
	if err := env.Parse(cfg); err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func (cfg *CollectorConfig) Load() {
	// if file exists, load from file ( from running file folder )
	// else load from env
	if _, err := os.Stat(DefaultConfigFile); err == nil {
		cfg.FromFile(DefaultConfigFile)
	} else {
		cfg.FromEnv()
	}
}

func (cfg *BaseLoggerConfig) FromFile(fileName string) {
	f, err := os.Open(fileName)
	if err != nil {
		processError(err)
	}
	defer f.Close()

	// Parse json file
	err = json.NewDecoder(f).Decode(cfg)
	if err != nil {
		processError(err)
	}
}

func (cfg *BaseLoggerConfig) FromEnv() {
	if err := env.Parse(cfg); err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func (cfg *BaseLoggerConfig) Load() {
	// if file exists, load from file ( from running file folder )
	// else load from env
	if _, err := os.Stat(DefaultConfigFile); err == nil {
		cfg.FromFile(DefaultConfigFile)
	} else {
		cfg.FromEnv()
	}
}
