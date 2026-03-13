package config

import (
	"encoding/json"
	"os"
)

// Config represents the structure of setting file
type Config struct {
	ListenPort string 		`json:"listen_port"`
	Backends   []string		`json:"backends"`
}

// Reads a JSON file and returns a Config struct
func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var cfg Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}