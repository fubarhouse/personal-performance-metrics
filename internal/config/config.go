package config

import (
	"os"

	"gopkg.in/yaml.v3"

	"github.com/fubarhouse/personal-performance-metrics/internal/types"
)

// LoadConfig will load the configuration file.
func LoadConfig(inputData []byte) (types.Config, error) {
	var cfg types.Config
	if inputData == nil {
		file, err := os.ReadFile("config.yml")
		if err != nil {
			return cfg, err
		}
		inputData = file
	}
	err := yaml.Unmarshal(inputData, &cfg)
	return cfg, err
}

// LoadData will load the data file.
func LoadData(inputData []byte) (types.PerformanceData, error) {
	var data types.PerformanceData
	if inputData == nil {
		file, err := os.ReadFile("data.yml")
		if err != nil {
			return data, err
		}
		inputData = file
	}
	err := yaml.Unmarshal(inputData, &data)
	return data, err
}
