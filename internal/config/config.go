package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type (
	REST struct {
		Host    string `yaml:"host"`
		Port    string `yaml:"port"`
		WebPath string `yaml:"web_path"`
	}

	Storage struct {
		Host   string `yaml:"host"`
		Port   string `yaml:"port"`
		DBName string `yaml:"db_name"`
		DBUser string `yaml:"db_user"`
		DBPass string `yaml:"db_pass"`
	}

	GameVariables struct {
		CardsPath              string  `yaml:"cards_path"`
		EarnedCoinsForInvestor uint64  `yaml:"earned_coins_for_investor"`
		PercentsForInvestor    float64 `yaml:"percents_for_investor"`
	}

	Config struct {
		REST          REST          `yaml:"rest"`
		Storage       Storage       `yaml:"storage"`
		GameVariables GameVariables `yaml:"game_variables"`
	}
)

// New creates a new Config instance
func New() *Config {
	return &Config{}
}

// Read loads the configuration from the given path
func (c *Config) Read(path string) (err error) {
	var configBytes []byte

	if configBytes, err = os.ReadFile(path); err != nil {
		return err
	}

	if err = yaml.Unmarshal(configBytes, c); err != nil {
		return err
	}

	return nil
}
