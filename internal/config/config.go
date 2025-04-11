package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type (
	App struct {
		Name string `yaml:"name" env-default:"pvz"`
		Env  string `yaml:"env" env-default:"dev"`
	}

	Cache struct {
		PvzListTTL time.Duration `yaml:"pvz_list_ttl" env-default:"5m"`
	}

	Config struct {
		App     App     `yaml:"app"`
		API     API     `yaml:"api"`
		Cache   Cache   `yaml:"cache"`
		Auth    Auth    `yaml:"auth" env-prefix:"auth"`
		PG      PG      `yaml:"postgres"`
		Swagger Swagger `yaml:"swagger"`
	}

	Prometheus struct {
		Port uint16 `yaml:"port" env-default:"3000"`
	}

	Swagger struct {
		Port uint16 `yaml:"port" env-required:"true"`
	}

	GRPC struct {
		Port uint16 `yaml:"port" env-required:"true"`
	}

	API struct {
		HTTP       HTTP       `yaml:"http" env-required:"true"`
		GRPC       GRPC       `yaml:"grpc" env-required:"true"`
		Prometheus Prometheus `yaml:"prometheus" env-required:"true"`
	}

	HTTP struct {
		Port         uint16        `env-required:"true" yaml:"port"`
		ReadTimeout  time.Duration `env-required:"true" yaml:"read_timeout"`
		WriteTimeout time.Duration `env-required:"true" yaml:"write_timeout"`
		IdleTimeout  time.Duration `env-required:"true" yaml:"idle_timeout"`
	}

	PG struct {
		URL             string        `env:"DATABASE_URL" env-default:"postgres://user:password@localhost:5432/pvz?sslmode=disable"`
		MaxOpenConns    int32         `env-required:"true" yaml:"max_open_conns"`
		MinIdleConns    int32         `env-required:"true" yaml:"min_idle_conns"`
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
