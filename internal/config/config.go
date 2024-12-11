package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"time"
)

type Config struct {
	Env      string   `yaml:"env" env:"ENV" env-default:"local"`
	Server   Server   `yaml:"http_server" env-required:"true"`
	Database Database `yaml:"database" env-required:"true"`
	Cache    Cache    `yaml:"cache" env-required:"true"`
}

type Server struct {
	Address     string        `yaml:"address" env:"HTTP_SERVER_ADDRESS" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env:"HTTP_SERVER_TIMEOUT" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env:"HTTP_SERVER_IDLE_TIMEOUT" env-default:"60s"`
}

type Database struct {
	Username string `yaml:"db_username" env:"DB_USERNAME" env-required:"true"`
	Password string `yaml:"db_password" env:"DB_PASSWORD" env-required:"true"`
	Name     string `yaml:"db_name" env:"DB_NAME" env-required:"true"`
	Host     string `yaml:"db_host" env:"DB_HOST" env-default:"localhost"`
	Port     int    `yaml:"db_port" env:"DB_PORT" env-default:"5432"`
}

type Cache struct {
	Address  string        `yaml:"cache_address" env:"CACHE_ADDRESS" env-default:"localhost:6379"`
	Username string        `yaml:"cache_user" env:"CACHE_USER" env-required:"true"`
	Password string        `yaml:"cache_password" env:"CACHE_PASSWORD" env-default:""`
	DB       int           `yaml:"cache_db" env:"CACHE_DB" env-default:"0"`
	Timeout  time.Duration `yaml:"cache_timeout" env:"CACHE_TIMEOUT" env-default:"5s"`
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
