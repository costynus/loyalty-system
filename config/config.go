package config

import (
	"flag"
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		App `yaml:"app"`
        HTTP `yaml:"http"`
        Log `yaml:"logger"`
        PG `yaml:"postgres"`
	}

	App struct {
		Name    string `env-required:"true" yaml:"name" env:"NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
        AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
        WorkersCount int `env-required:"true" yaml:"workers_count"`
        MarshalJSONWithoutQuotes bool `yaml:"marshal_json_without_quotes"`
	}

    HTTP struct {
        Address string `env-required:"true" yaml:"address" env:"RUN_ADDRESS"`
    }

    Log struct {
        Level string `env-required:"true" yaml:"log_level" env:"LOG_LEVEL"`
    }

    PG struct {
        PollMax int `env-required:"true" yaml:"pool_max" env:"PG_POOL_MAX"`
        URL string `env-required:"true" env:"DATABASE_URI"`
        MigDir string `env-required:"true" env:"MIG_DIR" yaml:"migration_dir"`
    }
)

func NewConfig() (*Config, error) {
	cfg := &Config{}

    flag.StringVar(&cfg.HTTP.Address, "a", cfg.HTTP.Address, "address to listen on")
    flag.StringVar(&cfg.PG.URL, "d", cfg.PG.URL, "db url")
    flag.StringVar(&cfg.App.AccrualSystemAddress, "r", cfg.App.AccrualSystemAddress, "accreal system address")

	err := cleanenv.ReadConfig("./config/config.yaml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config: error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
