package render

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/fubarhouse/personal-performance-metrics/internal/types"
)

func TestPrintTable(t *testing.T) {
	// Prepare test data
	testData := types.PerformanceData{
		"Machine1": 95.5,
		"Machine2": 87.3,
	}

	testConfig := types.Config{
		MetricMappings: map[string]types.MetricMapping{
			"Machine1": {
				Name: "CPU Usage",
				Dimensions: []types.MetricMappingDimensions{
					{Name: "Environment", Value: "Production"},
				},
			},
			"Machine2": {
				Name: "Memory Usage",
				Dimensions: []types.MetricMappingDimensions{
					{Name: "Environment", Value: "Staging"},
					{Name: "Region", Value: "US-West"},
				},
			},
		},
	}

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Call the function
	err := PrintTable(testData, testConfig)

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout

	// Read captured output
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Check for errors
	if err != nil {
		t.Errorf("PrintTable returned an error: %v", err)
	}

	// Check if the output contains expected content
	expectedContents := []string{
		"Metric name", "Value", "Dimensions", "Machine Name",
		"CPU Usage", "95.5", "Environment=Production", "Machine1",
		"Memory Usage", "87.3", "Environment=Staging Region=US-West", "Machine2",
	}

	assert.Equal(t, expectedContents[0], strings.Split(output, "\n"))

	// TODO assertions.
}
