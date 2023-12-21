package config

import (
	"fmt"
	"log"
	"os"
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
	// Загружаем переменные окружения из файла .env
	err := godotenv.Load("C:/Users/User/GoProjects/url-shortener/local.env")
	if err != nil {
		fmt.Println("Error loading local.env file")
		// обработка ошибки
	}

	// // Теперь вы можете использовать переменные окружения
	// configPath := os.Getenv("CONFIG_PATH")
	// fmt.Println("CONFIG_PATH:", configPath)

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}
	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal("config file does not exist: %s", configPath)
	}

	var cfg Config

	// fmt.Println(configPath)

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatal("cannot read config: %s", err)
	}
	return &cfg
}
