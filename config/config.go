package config

import (
	"encoding/json"
	"fmt"

	"github.com/awlsring/action-runner/api"
	"github.com/awlsring/action-runner/runner"
	"github.com/awlsring/surreal-db-client/surreal"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

type DatabaseConfig struct {
	Surreal surreal.SurrealConfig `mapstructure:"surreal"`
}

type Config struct {
	Runner runner.Config `mapstructure:"runner"`
	Database DatabaseConfig `mapstructure:"db"`
	Api api.Config `mapstructure:"api"`
}

func LoadConfig() (Config, error) {
	vp := viper.New()

	var config Config

	vp.SetConfigName("config")
	vp.SetConfigType("yaml")
	vp.AddConfigPath(".")
	err := vp.ReadInConfig()
	if err != nil {
		return Config{}, err
	}

	err = vp.Unmarshal(&config)
	if err != nil {
		return Config{}, err
	}

	j, _ := json.Marshal(config)
	log.Info("Loaded config: %v", string(j))

	validateDbConfig(config.Database)

	return config, nil
}

func validateDbConfig(cfg DatabaseConfig) {
	m, _ := json.Marshal(cfg)
	var dbMap map[string]interface{}
	json.Unmarshal(m, &dbMap)

	fmt.Println(len(dbMap))

	if len(dbMap) > 1	{
		log.Fatal("Multiple databases are set")
	} else if len(dbMap) == 0 {
		log.Fatal("No database is set")
	}

	if (surreal.SurrealConfig{}) != cfg.Surreal {
		validateSurrealDbConfig(cfg.Surreal)
	}
}

func validateSurrealDbConfig(cfg surreal.SurrealConfig) {
	if cfg.Address == "" {
		log.Fatal("Surreal host is required")
	}
	if cfg.User == "" {
		log.Fatal("Surreal user is required")
	}
	if cfg.Password == "" {
		log.Fatal("Surreal password is required")
	}
}