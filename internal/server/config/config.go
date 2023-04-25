// Package config provides the configuration for the server.
package config

import "time"

type Config struct {
	Address       string        `env:"ADDRESS"`
	CertFile      string        `env:"CERT_FILE"`
	CryptoKey     string        `env:"CRYPTO_KEY"`
	DatabaseDSN   string        `env:"DATABASE_DSN"`
	EnableHTTPS   bool          `env:"ENABLE_HTTPS"`
	Key           string        `env:"KEY"`
	Restore       bool          `env:"RESTORE"`
	StoreInterval time.Duration `env:"STORE_INTERVAL"`
	StoreFile     string        `env:"STORE_FILE"`
}
