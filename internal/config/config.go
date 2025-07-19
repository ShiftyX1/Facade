package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)


type Config struct {
	Server ServerConfig    `yaml:"server"`
	Routes []Route         `yaml:"routes"`
	Global GlobalSettings  `yaml:"global"`
}


type ServerConfig struct {
	Port int    `yaml:"port"`
	Host string `yaml:"host"`
}


type Route struct {
	Path     string    `yaml:"path"`
	Method   string    `yaml:"method"`
	Response Response  `yaml:"response"`
}


type Response struct {
	Status     int                    `yaml:"status"`
	Body       string                 `yaml:"body,omitempty"`
	Schema     *SchemaDefinition      `yaml:"schema,omitempty"`
	Headers    map[string]string      `yaml:"headers,omitempty"`
	Conditions []RouteCondition       `yaml:"conditions,omitempty"`
}


type RouteCondition struct {
	If      string            `yaml:"if,omitempty"`
	Default bool              `yaml:"default,omitempty"`
	Body    string            `yaml:"body,omitempty"`
	Schema  *SchemaDefinition `yaml:"schema,omitempty"`
	Status  int               `yaml:"status,omitempty"`
}


type SchemaDefinition struct {
	Type       string                         `yaml:"type"`
	Items      *SchemaDefinition              `yaml:"items,omitempty"`
	Properties map[string]SchemaProperty      `yaml:"properties,omitempty"`
	Count      int                            `yaml:"count,omitempty"`
	Min        *int                           `yaml:"min,omitempty"`
	Max        *int                           `yaml:"max,omitempty"`
	Faker      string                         `yaml:"faker,omitempty"`
}


type SchemaProperty struct {
	Type  string `yaml:"type"`
	Min   *int   `yaml:"min,omitempty"`
	Max   *int   `yaml:"max,omitempty"`
	Faker string `yaml:"faker,omitempty"`
}


type GlobalSettings struct {
	CORS        bool          `yaml:"cors"`
	Delay       time.Duration `yaml:"delay"`
	LogRequests bool          `yaml:"log_requests"`
}


func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	
	if config.Server.Port == 0 {
		config.Server.Port = 8080
	}
	if config.Server.Host == "" {
		config.Server.Host = "localhost"
	}

	
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &config, nil
}


func validateConfig(config *Config) error {
	if config.Server.Port < 1 || config.Server.Port > 65535 {
		return fmt.Errorf("invalid port: %d", config.Server.Port)
	}

	for i, route := range config.Routes {
		if route.Path == "" {
			return fmt.Errorf("route %d: path is required", i)
		}
		if route.Method == "" {
			return fmt.Errorf("route %d: method is required", i)
		}
		if route.Response.Status == 0 {
			return fmt.Errorf("route %d: response status is required", i)
		}
	}

	return nil
}