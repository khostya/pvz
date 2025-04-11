package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type (
	Config struct {
		Env     string  `env-required:"true" yaml:"env"`
		API     API     `yaml:"api" env-prefix:"api"`
		Auth    Auth    `yaml:"auth" env-prefix:"auth"`
		PG      PG      `yaml:"postgres"`
		Swagger Swagger `yaml:"swagger"`
	}

	Swagger struct {
		Port uint16 `yaml:"port" env-required:"true"`
	}

	GRPC struct {
		Port uint16 `yaml:"port" env-required:"true"`
	}

	API struct {
		HTTP HTTP `yaml:"http" env-required:"true"`
		GRPC GRPC `yaml:"grpc" env-required:"true"`
	}

	HTTP struct {
		Port           uint16        `env-required:"true" yaml:"port"`
		ReadTimeout    time.Duration `env-required:"true" yaml:"read_timeout"`
		WriteTimeout   time.Duration `env-required:"true" yaml:"write_timeout"`
		IdleTimeout    time.Duration `env-required:"true" yaml:"idle_timeout"`
		MaxHeaderBytes int           `env-required:"true" yaml:"max_header_bytes"`
	}

	PG struct {
		URL             string        `env:"DATABASE_URL" env-default:"postgres://user:password@postgres:5431/postgres?sslmode=disable"`
		MaxOpenConns    int           `env-required:"true" yaml:"max_open_conns"`
		MaxIdleConns    int           `env-required:"true" yaml:"max_idle_conns"`
		ConnMaxIdleTime time.Duration `env-required:"true" yaml:"conn_max_idle_time"`
		ConnMaxLifetime time.Duration `env-required:"true" yaml:"conn_max_lifetime"`
	}

	Auth struct {
		PasswordCostBcrypt int           `env-required:"true" yaml:"password_cost_bcrypt"`
		AccessTokenTTL     time.Duration `env-required:"true" yaml:"access_token_ttl"`
		SigningKey         string        `env:"JWT_SIGNING_KEY" env-default:"fdsfsfsawf"`
	}
)

const (
	configPath = "./config/config.yaml"
)

func NewConfig() (Config, error) {
	cfg := Config{}

	path, exists := os.LookupEnv("CONFIG_PATH")
	if !exists {
		path = configPath
	}

	err := cleanenv.ReadConfig(path, &cfg)
	if err != nil {
		return Config{}, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}

func MustNewConfig() Config {
	cfg, err := NewConfig()
	if err != nil {
		panic(err)
	}
	return cfg
}
