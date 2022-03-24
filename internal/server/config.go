package server

import "github.com/VIWET/TestTaskSoftConstruct/internal/repository"

type Config struct {
	Addr           string             `yaml:"addr"`
	LogLevel       string             `yaml:"logLevel"`
	DatabaseConfig *repository.Config `yaml:"db"`
}

func NewConfig() *Config {
	return &Config{
		Addr:           ":8080",
		LogLevel:       "debug",
		DatabaseConfig: repository.NewConfig(),
	}
}
