package utils

import "github.com/caarlos0/env/v10"

// ReadConfig reads the configuration from environment variables and parses it into the provided cfg structure.
func ReadConfig(cfg interface{}) error {
	return env.Parse(cfg)
}
