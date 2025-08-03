package config

import (
	"errors"
	"io"
	"log"

	"github.com/goccy/go-yaml"
)

type Config struct {
	Generate []GenerateConfig `yaml:"generate"`
	Process  ProcessConfig    `yaml:"process"`
}

// GenerateConfig represents a configuration for generating files
type GenerateConfig struct {
	Name     string      `yaml:"name"`
	Path     string      `yaml:"path"`
	Strategy string      `yaml:"strategy"` // "append" or "template"
	Template string      `yaml:"template"` // Used when strategy=template
	Inputs   []InputFile `yaml:"inputs"`
}

// InputFile represents an input file to be watched
type InputFile struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
}

// ProcessConfig represents configuration for the managed process
type ProcessConfig struct {
	Path   string       `yaml:"path"`
	Reload ReloadConfig `yaml:"reload"`
	Args   []string     `yaml:"args"`
}

// ReloadConfig represents reload configuration for the managed process
type ReloadConfig struct {
	Enabled bool   `yaml:"enabled"`
	Method  string `yaml:"method"` // "restart" or "signal"
	Signal  string `yaml:"signal"` // E.g., "SIGHUP"
}

func LoadConfig(r io.Reader) (*Config, error) {
	appConfig := &Config{}
	configData, err := io.ReadAll(r)
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	err = yaml.Unmarshal(configData, appConfig)
	if err != nil {
		return nil, errors.Join(ErrorParseConfig, err)
	}

	// Validate configuration
	for _, gen := range appConfig.Generate {
		if gen.Strategy != "append" && gen.Strategy != "template" {
			return nil, &ErrorInvalidStrategy{Strategy: gen.Strategy, Name: gen.Name}
		}
		if gen.Strategy == "template" && gen.Template == "" {
			return nil, &ErrorMissingTemplate{Name: gen.Name}
		}
	}

	if appConfig.Process.Reload.Enabled {
		if appConfig.Process.Reload.Method != "restart" && appConfig.Process.Reload.Method != "signal" {
			return nil, &ErrorInvalidReloadMethod{Method: appConfig.Process.Reload.Method}
		}
		if appConfig.Process.Reload.Method == "signal" && appConfig.Process.Reload.Signal == "" {
			return nil, &ErrorMissingSignal{}
		}
	}
	return appConfig, nil
}
