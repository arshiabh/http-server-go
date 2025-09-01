package config

import (
	"fmt"
	"time"
)

type Config struct {
	Host         string
	Port         int
	WriteTimeout time.Duration
	ReadTimeout  time.Duration
	LogLevel     string
}

func DefaultConfig() *Config {
	return &Config{
		Host:         "localhost",
		Port:         8000,
		WriteTimeout: time.Second * 10,
		ReadTimeout:  time.Second * 30,
		LogLevel:     "INFO",
	}
}

func (c *Config) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
