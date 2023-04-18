// Package config provides the configuration for the server.
package config

import "time"

type Config struct {
	Address       string        `env:"ADDRESS"`
	DatabaseDSN   string        `env:"DATABASE_DSN"`
	Key           string        `env:"KEY"`
	Restore       bool          `env:"RESTORE"`
	StoreInterval time.Duration `env:"STORE_INTERVAL"`
	StoreFile     string        `env:"STORE_FILE"`
}
