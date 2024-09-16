package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Auth struct {
		WebAppInit string `yaml:"WebAppInit"`
		TokenFile  string `yaml:"tokenFile"`
		Folder     string `yaml:"folder"`
	} `yaml:"auth"`

	Settings struct {
		Daemon  bool `yaml:"daemon"`
		Delay   int  `yaml:"delay"`
		Tasks   bool `yaml:"tasks"`
		Farming bool `yaml:"farming"`
		Gaming  bool `yaml:"gaming"`
	} `yaml:"settings"`
}

type Tokens struct {
	Auth    string `json:"auth"`
	Refresh string `json:"refresh"`
}

func LoadConfig() (*Config, error) {
	var cfg Config
	viper.AddConfigPath(filepath.Join("config"))
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	if envWebAppInit := os.Getenv("WEB_APP_INIT_DATA"); envWebAppInit != "" {
		cfg.Auth.WebAppInit = envWebAppInit
	}
	if daemonMode := os.Getenv("DAEMON"); daemonMode != "" {
		cfg.Settings.Daemon = isTruthy(daemonMode)
	}

	return &cfg, nil
}

func (c *Config) LoadTokens() (*Tokens, error) {
	tokens := Tokens{}

	tokenPath := filepath.Join(c.Auth.Folder, c.Auth.TokenFile)
	if _, err := os.Stat(tokenPath); os.IsNotExist(err) {
		return &tokens, errors.New("tokens.json does not exist")
	}

	file, err := os.ReadFile(tokenPath)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(file, &tokens); err != nil {
		return nil, errors.New("failed unmarshalling tokens.json")
	}

	return &tokens, nil
}

func (c *Config) SaveTokens(tokens *Tokens) error {
	tokensJson, err := json.Marshal(&tokens)
	if err != nil {
		return err
	}

	tokenPath := filepath.Join(c.Auth.Folder, c.Auth.TokenFile)

	if err := os.WriteFile(tokenPath, tokensJson, 0644); err != nil {
		return err
	}

	return nil
}

func isTruthy(value string) bool {
	switch strings.ToLower(value) {
	case "true", "1", "t":
		return true
	default:
		return false
	}
}
