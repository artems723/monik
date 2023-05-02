// Package config provides the configuration for the server.
package config

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

type Config struct {
	Address       string        `env:"ADDRESS" json:"address"`
	CertFile      string        `env:"CERT_FILE" json:"cert_file"`
	ConfigFile    string        `env:"CONFIG" json:"-"`
	CryptoKey     string        `env:"CRYPTO_KEY" json:"crypto_key"`
	DatabaseDSN   string        `env:"DATABASE_DSN" json:"database_dsn"`
	EnableHTTPS   bool          `env:"ENABLE_HTTPS" json:"enable_https"`
	Key           string        `env:"KEY" json:"hash_key"`
	Restore       bool          `env:"RESTORE" json:"restore"`
	StoreInterval time.Duration `env:"STORE_INTERVAL" json:"store_interval"`
	StoreFile     string        `env:"STORE_FILE" json:"store_file"`
}

func LoadJSONConfig(configFile string, config *Config) error {
	raw, err := os.ReadFile(configFile)
	if err != nil {
		log.Println("Error occurred while reading config")
		return err
	}
	err = json.Unmarshal(raw, &config)
	if err != nil {
		return err
	}

	return nil
}
