package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/fubarhouse/personal-performance-metrics/internal/types"
)

func testSetup() (types.Config, types.PerformanceData) {
	config := types.Config{
		MetricMappings: map[string]types.MetricMapping{
			"valid": {
				Name:       "valid",
				Dimensions: nil,
			},
			"valid-with-timestamp-one": {
				Name:       "valid-with-timestamp-one",
				Dimensions: nil,
				Timestamp:  "1750246865",
			},
			"valid-with-timestamp-two": {
				Name:       "valid-with-timestamp-two",
				Dimensions: nil,
				Timestamp:  "30/05/2025",
			},
			"valid-with-timestamp-three": {
				Name:       "valid-with-timestamp-three",
				Dimensions: nil,
				Timestamp:  "invalid timestamp",
			},
			"invalid": {
				Name:       "",
				Dimensions: nil,
			},
		},
	}

	// Specify valid data inputs we want to process.
	data := types.PerformanceData{
		"valid":                      1,
		"valid-with-timestamp-one":   2,
		"valid-with-timestamp-two":   3,
		"valid-with-timestamp-three": 4,
	}

	return config, data
}

func TestFilter(t *testing.T) {

	// Generate test data
	config, data := testSetup()

	// Filter the data.
	out := Filter(config, data)

	// Compare the old data-set vs the filtered data-set.
	assert.Len(t, config.MetricMappings, 5)
	assert.Len(t, out, 4)
}

func TestPrint(t *testing.T) {

	// Generate test data
	config, data := testSetup()

	// Print the unfiltered data.
	err := Print(config, data)
	assert.NoError(t, err)
}

func TestProcess(t *testing.T) {

	// Generate test data
	config, data := testSetup()

	processed, err := Process(config, data)

	assert.NoError(t, err)
	assert.Len(t, processed, 4)
}

func TestPublish(t *testing.T) {
	// TODO.
}

func TestSort(t *testing.T) {

	before := types.Config{
		MetricMappings: map[string]types.MetricMapping{
			"zero": {
				Name:       "unknown-0",
				Dimensions: nil,
			},
			"one": {
				Name:       "unknown-1",
				Dimensions: nil,
			},
		},
	}

	expected := []string{
		"one",
		"zero",
	}

	after := Sort(before)

	assert.Len(t, after, len(expected))
	assert.Equal(t, expected, after)
}
