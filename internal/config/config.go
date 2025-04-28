package config

import (
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env          string `yaml:"env" env-default:"local"`
	DatabasePath string `yaml:"database-path" env-required:"true"`
	HTTPServer   struct {
		Host        string        `yaml:"host" env-default:"localhost"`
		Port        int           `yaml:"port" env-default:"8080"`
		Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
		IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
	} `yaml:"http-server"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		panic("config path is not specified")
	}

	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		panic(fmt.Errorf("failed to read config: %w", err))
	}

	return &cfg
}
