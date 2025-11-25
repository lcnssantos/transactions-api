package postgres

import (
	"fmt"
	"time"
)

type poolConfig struct {
	maxIdle     int
	maxOpen     int
	maxLifeTime time.Duration
}

func NewPoolConfig(maxIdle, maxOpen int, maxLifeTime time.Duration) poolConfig {
	return poolConfig{
		maxIdle,
		maxOpen,
		maxLifeTime,
	}
}

func NewConfig(Host, Port, User, Password, Database, SSLMode string) config {
	return config{
		Host,
		Port,
		User,
		Password,
		Database,
		SSLMode,
	}
}

type config struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

func (c config) string() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=America/Sao_Paulo",
		c.Host,
		c.User,
		c.Password,
		c.Database,
		c.Port,
		c.SSLMode,
	)
}
