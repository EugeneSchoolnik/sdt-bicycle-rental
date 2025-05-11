package config

import (
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env        string     `yaml:"env" env-default:"local"`
	HTTPServer HTTPServer `yaml:"http-server"`
	Postgres   Postgres   `yaml:"postgres"`
}

type HTTPServer struct {
	Host        string        `yaml:"host" env-default:"localhost"`
	Port        int           `yaml:"port" env-default:"8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type Postgres struct {
	Host         string `yaml:"host" env-default:"localhost"`
	Port         string `yaml:"port" env-default:"5432"`
	User         string `yaml:"user" env-required:"true"`
	Password     string `yaml:"password" env:"POSTGRES_PASSWORD" env-required:"true"`
	DBName       string `yaml:"db-name" env-required:"true"`
	SSLMode      string `yaml:"ssl-mode" env-default:"disable"` // disable / require / verify-full
	TimeZone     string `yaml:"time-zone" env-default:"UTC"`    // i.g. "UTC" or "Europe/Kyiv"
	MaxOpenConns int    `yaml:"max-open-conns" env-default:"10"`
	MaxIdleConns int    `yaml:"max-idle-conns" env-default:"10"`
}

func MustLoad() *Config {
	err := godotenv.Load()
	if err != nil {
		panic(fmt.Errorf("failed to load .env file: %w", err))
	}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		panic("config path is not specified")
	}

	var cfg Config

	err = cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		panic(fmt.Errorf("failed to read config: %w", err))
	}

	return &cfg
}
