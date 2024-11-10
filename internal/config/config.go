package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env      string   `yaml:"env" env-default:"local"`
	Server   server   `yaml:"http_server" env-required:"true"`
	Database Database `yaml:"database" env-required:"true"`
}

type server struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type Database struct {
	Username string `yaml:"db_username" env-required:"true"`
	Password string `yaml:"db_password" env-required:"true"`
	Name     string `yaml:"db_name" env-required:"true"`
	Host     string `yaml:"db_host" env-default:"localhost"`
	Port     int    `yaml:"db_port" env-default:"5432"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist: " + path)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string
	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		if err := godotenv.Load(); err != nil {
			log.Println("Error loading .env file")
		}

		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
