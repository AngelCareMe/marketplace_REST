package config

import (
	"fmt"
	"marketplace/pkg/migrate"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Database struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		DBName   string `yaml:"dbname"`
		SSLMode  string `yaml:"sslmode"`
	} `yaml:"database"`
	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`
	Logger struct {
		Level  string `yaml:"level"`
		Format string `yaml:"format"`
	} `yaml:"logger"`
	Migrations struct {
		Dir     string `yaml:"dir"`
		Enabled bool   `yaml:"enabled"`
	} `yaml:"migrations"`
	JWT struct {
		SecretKey string `yaml:"secret_key"`
	} `yaml:"jwt"`
	DatabaseDSN string
}

func LoadConfig() (*Config, error) {
	data, err := os.ReadFile("pkg/config/config.yaml")
	if err != nil {
		logrus.WithError(err).Error("Failed to read config.yaml")
		return nil, fmt.Errorf("failed to read config.yaml: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		logrus.WithError(err).Error("Failed to parse config.yaml")
		return nil, fmt.Errorf("failed to parse config.yaml: %w", err)
	}

	cfg.DatabaseDSN = fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	if cfg.JWT.SecretKey == "" {
		logrus.Error("JWT secret key is required in config.yaml")
		return nil, fmt.Errorf("jwt.secret_key cannot be empty")
	}

	if cfg.Migrations.Enabled {
		if err := migrate.RunMigrations(cfg.DatabaseDSN, cfg.Migrations.Dir); err != nil {
			logrus.WithError(err).Error("Failed to run migrations")
			return nil, fmt.Errorf("failed to run migrations: %w", err)
		}
		logrus.Info("Database migrations applied successfully")
	}

	return &cfg, nil
}
