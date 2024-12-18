package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ServerConfiguration struct {
	HTTPPort     string `mapstructure:"HTTP_PORT"`
	AllowedHosts string `mapstructure:"ALLOWED_HOSTS"`
}

type DatabaseConfiguration struct {
	Name         string `mapstructure:"DATABASE_NAME"`
	DoMigrations bool   `mapstructure:"DO_MIGRATIONS"`
}

type SyncConfiguration struct {
	CheckpointInSeconds   int `mapstructure:"CHECKPOINT_IN_SECONDS"`
	WorkersAllowedForSync int `mapstructure:"WORKERS_ALLOWED_FOR_SYNC"`
}

type WorkerConfiguration struct {
	MaxWorkersAllowedConcurrentlyForRealtimeUpdates int `mapstructure:"MAX_WORKERS_ALLOWED_CONCURRENTLY_FOR_REAL_TIME_UPDATES"`
}

type Config struct {
	Server            ServerConfiguration   `mapstructure:",squash"`
	Database          DatabaseConfiguration `mapstructure:",squash"`
	SyncConfiguration SyncConfiguration     `mapstructure:",squash"`
	WorkerConfig      WorkerConfiguration   `mapstructure:",squash"`
}

func Load() (*Config, error) {
	var config Config

	// Set config file and read it
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		logrus.Warn("Could not read .env file, falling back to environment variables only")
	}

	// Bind environment variables
	viper.AutomaticEnv()

	// Optionally, set defaults
	viper.SetDefault("DO_MIGRATIONS", false)

	// Decode into the struct
	if err := viper.Unmarshal(&config); err != nil {
		logrus.Errorf("Unable to decode config: %v", err)
		return nil, err
	}

	return &config, nil
}
