package main

import (
	"encoding/json"
	"os"
)

// Config maintains the configuration for the application and manages its
// storage and retrieval from disk (in JSON format).
type Config struct {
	Addr          string `json:"addr"`           // Address to listen on
	RootPath      string `json:"root_path"`      // Path to www directory
	DataPath      string `json:"data_path"`      // Path to data files
	AdminPassword string `json:"admin_password"` // Password for admin user
}

// LoadConfig loads configuration from disk.
func LoadConfig(name string) (*Config, error) {
	r, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	c := &Config{}
	if err := json.NewDecoder(r).Decode(c); err != nil {
		return nil, err
	}
	return c, nil
}
