package render

import (
	"fmt"
	"math"

	"github.com/pterm/pterm"

	datetime "github.com/fubarhouse/personal-performance-metrics/internal/time"
	"github.com/fubarhouse/personal-performance-metrics/internal/types"
)

// PrintTable will print a table showing all the metrics which are going to be pushed.
func PrintTable(data types.PerformanceData, config types.Config) error {
	alternateStyle := pterm.NewStyle(pterm.BgDarkGray)
	tableData := pterm.TableData{
		{"Metric name", "Value", "Dimensions", "Machine Name", "Timestamp (unix)"},
	}

	for key, val := range data {
		var dimensions string
		for _, v := range config.MetricMappings[key].Dimensions {
			dimensions += fmt.Sprintf("%s=%s ", v.Name, v.Value)
		}
		datetimeCalculated, _ := datetime.ParseTimeString(config.MetricMappings[key].Timestamp)
		tableData = append(tableData, []string{config.MetricMappings[key].Name, fmt.Sprint(math.Round(val*100) / 100), dimensions, key, datetimeCalculated.String()})
	}

	//fmt.Println("Metrics to be published:")
	return pterm.DefaultTable.WithHasHeader().WithBoxed().WithData(tableData).WithStyle(alternateStyle).Render()
}
