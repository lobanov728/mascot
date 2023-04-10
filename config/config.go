package config

import (
	"github.com/Lobanov728/mascot/internal/billing/adapters"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Config struct {
	DB DbConfig
}

type DbConfig struct {
	Postgres adapters.Postgres
}

func Init(configPath string, cfg *Config) error {
	vp := viper.New()

	vp.SetConfigFile(configPath)

	if err := vp.MergeInConfig(); err != nil {
		return errors.Wrap(err, "read config")
	}

	if err := vp.Unmarshal(&cfg); err != nil {
		return errors.Wrap(err, "unmarshal config to obj")
	}

	return nil
}
