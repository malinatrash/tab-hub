package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env      string `env:"ENV" env-default:"local"`
	Server   server
	Database Database
	Cache    Cache
}

type server struct {
	Address     string        `env:"HTTP_SERVER_ADDRESS" env-default:"localhost:8000"`
	Timeout     time.Duration `env:"HTTP_SERVER_TIMEOUT" env-default:"4s"`
	IdleTimeout time.Duration `env:"HTTP_SERVER_IDLE_TIMEOUT" env-default:"60s"`
}

type Database struct {
	Username string `env:"DB_USERNAME" env-required:"true"`
	Password string `env:"DB_PASSWORD" env-required:"true"`
	Name     string `env:"DB_NAME" env-required:"true"`
	Host     string `env:"DB_HOST" env-default:"localhost"`
	Port     int    `env:"DB_PORT" env-default:"5432"`
}

type Cache struct {
	Address  string        `env:"CACHE_ADDRESS" env-default:"localhost:6379"`
	Username string        `env:"CACHE_USER" env-required:"true"`
	Password string        `env:"CACHE_PASSWORD" env-default:""`
	DB       int           `env:"CACHE_DB" env-default:"0"`
	Timeout  time.Duration `env:"CACHE_TIMEOUT" env-default:"5s"`
}

func MustLoad() *Config {
	if err := godotenv.Load(); err != nil {
		panic("Warning: .env file not found, using environment variables")
	}

	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic("failed to read environment variables: " + err.Error())
	}

	return &cfg
}
