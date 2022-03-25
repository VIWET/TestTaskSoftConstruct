package repository

import "fmt"

type Config struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Database string `yaml:"database"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

func NewConfig() *Config {
	return &Config{
		Host:     "localhost",
		Port:     "3306",
		Database: "game_chat_db",
		User:     "game_chat",
		Password: "",
	}
}

func (c *Config) GetConnectionString() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Database)
}
