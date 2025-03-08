package config

import (
	"os"
)

type Config struct {
	BotToken    string
	ChannelID   string
	DatabaseURL string
	Debug       bool
}

func LoadConfig() (*Config, error) {
	return &Config{
		BotToken:    os.Getenv("BOT_TOKEN"),
		ChannelID:   os.Getenv("CHANNEL_ID"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		Debug:       os.Getenv("DEBUG") == "true",
	}, nil
}
