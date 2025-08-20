package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	ClientID     string
	ClientSecret string
	ChannelID    string
}

func LoadEnv(envFile string) Config {
	if err := godotenv.Load(envFile); err != nil {
		log.Print(err)
	}

	return Config{
		Port:         getEnv("PORT", "8080"),
		ClientID:     getEnv("TWITCH_CLIENT_ID", ""),
		ClientSecret: getEnv("TWITCH_CLIENT_SECRET", ""),
		ChannelID:    getEnv("TWITCH_CHANNEL_ID", ""),
	}
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
