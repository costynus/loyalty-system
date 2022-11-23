package config
import (
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
		Name    string `env-required:"true" yaml:"name" env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

    HTTP struct {
        Address string `env-required:"true" yaml:"address" env:"ADDRESS"`
    }

    Log struct {
        Level string `env-required:"true" yaml:"log_level" env:"LOG_LEVEL"`
    }

    PG struct {
        PollMax int `env-required:"true" yaml:"pool_max" env:"PG_POOL_MAX"`
        URL string `env-required:"true" env:"DATABASE_DSN"`
        MigDir string `env-required:"true" env:"MIG_DIR" yaml:"migration_dir"`
    }
)

func NewConfig() (*Config, error) {
	cfg := &Config{}

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
