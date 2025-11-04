package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Config represents the application configuration
type Config struct {
	Server    ServerConfig  `json:"server"`
	Services  []Service     `json:"services"`
	RateLimit int           `json:"rate_limit"`
}

// ServerConfig contains server-specific configuration
type ServerConfig struct {
	Address      string        `json:"address"`
	ReadTimeout  time.Duration `json:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout"`
	IdleTimeout  time.Duration `json:"idle_timeout"`
}

// Service represents a backend service configuration
type Service struct {
	Name   string `json:"name"`
	Path   string `json:"path"`
	Target string `json:"target"`
}

// Load loads the configuration from a file or returns default config
func Load() (*Config, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.json"
	}

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Return default configuration
		return defaultConfig(), nil
	}

	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &cfg, nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.Server.Address == "" {
		return fmt.Errorf("server address is required")
	}

	if len(c.Services) == 0 {
		return fmt.Errorf("at least one service must be configured")
	}

	for _, svc := range c.Services {
		if svc.Name == "" {
			return fmt.Errorf("service name is required")
		}
		if svc.Path == "" {
			return fmt.Errorf("service path is required for %s", svc.Name)
		}
		if svc.Target == "" {
			return fmt.Errorf("service target is required for %s", svc.Name)
		}
	}

	return nil
}

// defaultConfig returns a default configuration
func defaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Address:      ":8080",
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
		Services: []Service{
			{
				Name:   "example-service",
				Path:   "/api/v1/*",
				Target: "http://localhost:9000",
			},
		},
		RateLimit: 100,
	}
}
