package main

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

type ServerConfig struct {
	Port              int           `yaml:"port"`
	ReadTimeout       time.Duration `yaml:"readTimeout"`
	WriteTimeout      time.Duration `yaml:"writeTimeout"`
	IdleTimeout       time.Duration `yaml:"idleTimeout"`
	ReadHeaderTimeout time.Duration `yaml:"readHeaderTimeout"`
}

type Config struct {
	App    AppConfig    `yaml:"app"`
	Server ServerConfig `yaml:"server"`
}

func LoadConfig() (*Config, error) {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
