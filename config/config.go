package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	LogLevel     string `json:"log_level"`
	DiscordID    string `json:"discord_id"`
	DiscordToken string `json:"discord_token"`
	SleepTime    int    `json:"sleep_time"`
	Tokens       []struct {
		Token   string  `json:"token"`
		Price   float64 `json:"price"`
		IsBelow bool    `json:"is_below"`
		IsAbove bool    `json:"is_above"`
	} `json:"tokens"`
}

func GetConfig() (Config, error) {
	var config Config
	err := ReadJSON("config.json", &config)
	if err != nil {
		return config, err
	}
	return config, nil
}

func ReadJSON(path string, v interface{}) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&v)
	if err != nil {
		return err
	}
	return nil
}
