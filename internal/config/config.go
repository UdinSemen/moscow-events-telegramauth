package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string      `yaml:"env"`
	Postgres    postgres    `yaml:"postgres"`
	TelegramBot telegramBot `yaml:"telegram-bot"`
	Redis       redis       `yaml:"redis"`
}

type telegramBot struct {
	Token   string `yaml:"token"`
	TimeOut int    `yaml:"time-out"`
}

type postgres struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DbName   string `yaml:"db-name"`
	SslMode  string `yaml:"ssl-mode"`
}

type redis struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Password string `yaml:"password"`
	DbName   int    `yaml:"db-name"`
}

func MustLoad() *Config {
	var cfg Config

	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		panic("CONFIG_PATH isn't set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Panicf("config file doesn't exist: %s", configPath)
	}

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Panicf("can't read config: %s", err)
	}

	return &cfg
}
