package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env string `env:"ENV" env-default:"local" env:"ENV"`
	HTTPServer
	Database
}

type HTTPServer struct {
	Address     string        `env-default:"localhost:8082" env:"ADDRESS"`
	Timeout     time.Duration `env-default:"5s" env:"TIMEOUT"`
	IdleTimeout time.Duration `env-default:"60s" env:"IDLE_TIMEOUT"`
	AliasLength int           `env-default:"6" env:"ALIAS_LEN"`
}

type Database struct {
	Type string `env-default:"sqlite" env:"DB_TYPE"`
	DSN  string `env:"DB_DSN"`
}

func GetConfig() *Config {
	var cfg Config
	errEnv := cleanenv.ReadEnv(&cfg)
	if errEnv != nil {
		return nil
	}
	return &cfg
}
