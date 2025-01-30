package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"

	"github.com/fubarhouse/personal-performance-metrics/internal/types"
)

func TestLoadConfig(t *testing.T) {
	input := types.Config{
		Interactive:     true,
		Region:          "ap-southeast-2",
		Profile:         "custom",
		SkipPublish:     true,
		MetricNamespace: "AWS/RDS",
		MetricMappings:  nil,
	}

	inputBytes, err := yaml.Marshal(input)
	assert.NoError(t, err)

	cfg, err := LoadConfig(inputBytes)
	assert.NoError(t, err)

	assert.Equal(t, cfg.Interactive, true)
	assert.Equal(t, cfg.Region, "ap-southeast-2")
	assert.Equal(t, cfg.Profile, "custom")
	assert.Equal(t, cfg.SkipPublish, true)
	assert.Equal(t, cfg.MetricNamespace, "AWS/RDS")
	assert.Empty(t, cfg.MetricMappings)
}

func TestLoadData(t *testing.T) {
	input := types.PerformanceData{
		"0": 0,
		"1": 1,
		"2": 2,
		"3": 3,
		"4": 4,
	}

	inputBytes, err := yaml.Marshal(input)
	assert.NoError(t, err)

	data, err := LoadData(inputBytes)
	assert.NoError(t, err)

	assert.Len(t, data, 5)
	assert.Equal(t, data["0"], float64(0))
	assert.Equal(t, data["1"], float64(1))
	assert.Equal(t, data["2"], float64(2))
	assert.Equal(t, data["3"], float64(3))
	assert.Equal(t, data["4"], float64(4))
}
