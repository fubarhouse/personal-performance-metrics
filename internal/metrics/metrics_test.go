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
			"invalid": {
				Name:       "",
				Dimensions: nil,
			},
		},
	}

	data := types.PerformanceData{
		"invalid": 0,
		"valid":   1,
	}

	return config, data
}

func TestFilter(t *testing.T) {

	// Generate test data
	config, data := testSetup()

	// Filter the data.
	out := Filter(config, data)

	// Compare the old data-set vs the filtered data-set.
	assert.Len(t, data, 2)
	assert.Len(t, out, 1)
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
	assert.Len(t, processed, 2)

	assert.Equal(t, *processed[0].MetricName, "")
	assert.Equal(t, *processed[0].Value, float64(0))
	assert.Equal(t, *processed[1].MetricName, "valid")
	assert.Equal(t, *processed[1].Value, float64(1))
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
