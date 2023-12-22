package config

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env         string `yaml:"env" env-default:"local"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer  `yaml:"http_server"`
}

type HTTPServer struct {
	Addres      string        `yaml:"addreas" env-default:"localhost:8081`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default"60s"`
}

func MustLoad() *Config {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal("Error getting current directory:", err)
	}
	//Загружаем переменные окружения из файла .env
	err = godotenv.Load(filepath.Join(currentDir[0:len(currentDir)-18], "local.env"))
	if err != nil {
		log.Fatal("Error loading local.env file whith error:", err)
	}
	// Теперь вы можете использовать переменные окружения
	configPath := filepath.Join(currentDir[0:len(currentDir)-18], os.Getenv("CONFIG_PATH"))
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatal("cannot read config: %s", err)
	}
	return &cfg
}
