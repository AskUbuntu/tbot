package config

import (
	"encoding/json"
	"os"
)

// Config maintains the configuration for the application and manages its
// storage and retrieval from disk (in JSON format).
type Config struct {
	Addr                  string `json:"addr"`           // Address to listen on
	RootPath              string `json:"root_path"`      // Path to www directory
	DataPath              string `json:"data_path"`      // Path to data files
	AdminPassword         string `json:"admin_password"` // Password for admin user
	TwitterConsumerKey    string `json:"twitter_consumer_key"`
	TwitterConsumerSecret string `json:"twitter_consumer_secret"`
	TwitterAccessToken    string `json:"twitter_access_token"`
	TwitterAccessSecret   string `json:"twitter_access_secret"`
}

// Load loads configuration from disk.
func Load(name string) (*Config, error) {
	r, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	c := &Config{}
	if err := json.NewDecoder(r).Decode(c); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(c.DataPath, 0700); err != nil {
		return nil, err
	}
	return c, nil
}
