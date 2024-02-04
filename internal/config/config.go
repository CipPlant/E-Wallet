package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"log/slog"
	"os"
	"time"
)

type Config struct {
	PGSQL      `yaml:"pgsql"`
	HTTPServer `yaml:"httpServer"`
	Logger     `yaml:"logger"`
}

type Logger struct {
	LogLevel string `yaml:"LOG_LEVEL" env-default:"dev"`
}

type PGSQL struct {
	UserName string        `yaml:"PG_USER" env-default:"postgres"`
	Password string        `yaml:"PG_PASSWORD" env-required:"true"`
	DbName   string        `yaml:"PG_DATABASE" env-required:"true"`
	Port     string        `yaml:"PG_PORT" env-default:"8080"`
	Host     string        `yaml:"PG_HOST" env-default:"localhost:8080"`
	PingTime time.Duration `yaml:"PG_PING_TIME" env-default:"3s"`
}

type HTTPServer struct {
	Host string `yaml:"HTTP_HOST" env-default:"localhost"`
	Port string `yaml:"HTTP_PORT" env-required:"80"`
}

func MustLoad() (*Config, error) {

	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	var cfg Config

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		slog.Error("error with file: %v", err)
		return nil, err
	}
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		slog.Error("error with parse: %v", err)
		return nil, err
	}
	return &cfg, nil
}
