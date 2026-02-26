package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
}

type ServerConfig struct {
	Port string `yaml:"port" env:"SERVER_PORT" envDefault:"8080"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host" env:"DB_HOST" envDefault:"localhost"`
	Port     string `yaml:"port" env:"DB_PORT" envDefault:"5432"`
	User     string `yaml:"user" env:"DB_USER" envDefault:"postgres"`
	Password string `yaml:"password" env:"DB_PASSWORD" envDefault:"postgres"`
	DBName   string `yaml:"dbname" env:"DB_NAME" envDefault:"subscriptions"`
	SSLMode  string `yaml:"sslmode" env:"DB_SSLMODE" envDefault:"disable"`
}

func Load() (*Config, error) {
	// Попытка загрузить .env файл
	_ = godotenv.Load()

	// Попытка загрузить config.yaml
	cfg := &Config{}
	if yamlFile, err := os.ReadFile("config.yaml"); err == nil {
		if err := yaml.Unmarshal(yamlFile, cfg); err != nil {
			return nil, fmt.Errorf("failed to parse config.yaml: %w", err)
		}
	}

	// Переопределение из переменных окружения
	if port := os.Getenv("SERVER_PORT"); port != "" {
		cfg.Server.Port = port
	}
	if cfg.Server.Port == "" {
		cfg.Server.Port = "8080"
	}

	if host := os.Getenv("DB_HOST"); host != "" {
		cfg.Database.Host = host
	}
	if cfg.Database.Host == "" {
		cfg.Database.Host = "localhost"
	}

	if port := os.Getenv("DB_PORT"); port != "" {
		cfg.Database.Port = port
	}
	if cfg.Database.Port == "" {
		cfg.Database.Port = "5432"
	}

	if user := os.Getenv("DB_USER"); user != "" {
		cfg.Database.User = user
	}
	if cfg.Database.User == "" {
		cfg.Database.User = "postgres"
	}

	if password := os.Getenv("DB_PASSWORD"); password != "" {
		cfg.Database.Password = password
	}
	if cfg.Database.Password == "" {
		cfg.Database.Password = "postgres"
	}

	if dbname := os.Getenv("DB_NAME"); dbname != "" {
		cfg.Database.DBName = dbname
	}
	if cfg.Database.DBName == "" {
		cfg.Database.DBName = "subscriptions"
	}

	if sslmode := os.Getenv("DB_SSLMODE"); sslmode != "" {
		cfg.Database.SSLMode = sslmode
	}
	if cfg.Database.SSLMode == "" {
		cfg.Database.SSLMode = "disable"
	}

	return cfg, nil
}
