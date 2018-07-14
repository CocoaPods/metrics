package config

import (
	"github.com/spf13/viper"
	"strings"
)

// DBConfig is a struct that holds the values required to connect to a
// postgresish database.
type DBConfig struct {
	Username     string
	Password     string
	Host         string `mapstructure:"host"`
	Port         int
	DatabaseName string `mapstructure:"database_name"`
	Schema       string
}

// Config is a struct that holds the overall config for the metrics application.
type Config struct {
	WarehouseDB *DBConfig
	MetricsDB   *DBConfig
}

var defaults = map[string]string{
	"warehouse.name": "events",

	"metrics.schema": "metrics_v2",
}

// Parse will resolve the configuration from any provided config paths and the
// process environment, and return an error or populated configuration.
func Parse(name string, paths []string) (*Config, error) {
	v := viper.New()
	v.SetConfigType("toml")

	for key, val := range defaults {
		v.SetDefault(key, val)
	}

	// len(paths) == 0 is ok when reading only from env.
	if len(paths) == 1 {
		v.SetConfigFile(paths[0])
	} else if len(paths) > 1 {
		v.SetConfigName(name)
		for _, p := range paths {
			v.AddConfigPath(p)
		}
	}

	v.SetEnvPrefix("COCOAPODS_")
	v.AutomaticEnv()
	replacer := strings.NewReplacer(".", "__")
	v.SetEnvKeyReplacer(replacer)

	if err := v.ReadInConfig(); err != nil {
		switch err.(type) {
		// Env only config is ok, so if we don't find a config file, it's fine.
		case viper.ConfigFileNotFoundError:
		default:
			return nil, err
		}
	}

	c := &Config{
		WarehouseDB: &DBConfig{
			Username:     v.GetString("warehouse.username"),
			Password:     v.GetString("warehouse.password"),
			Host:         v.GetString("warehouse.host"),
			Port:         v.GetInt("warehouse.port"),
			DatabaseName: v.GetString("warehouse.name"),
			Schema:       v.GetString("warehouse.schema"),
		},
		MetricsDB: &DBConfig{
			Username:     v.GetString("metrics.username"),
			Password:     v.GetString("metrics.password"),
			Host:         v.GetString("metrics.host"),
			Port:         v.GetInt("metrics.port"),
			DatabaseName: v.GetString("metrics.name"),
			Schema:       v.GetString("metrics.schema"),
		},
	}

	return c, nil
}
