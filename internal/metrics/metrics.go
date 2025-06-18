package metrics

import (
	"context"
	"fmt"
	datetime "github.com/fubarhouse/personal-performance-metrics/internal/time"
	"math"
	"slices"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	cloudwatchtypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"

	"github.com/fubarhouse/personal-performance-metrics/internal/prompt"
	"github.com/fubarhouse/personal-performance-metrics/internal/render"
	"github.com/fubarhouse/personal-performance-metrics/internal/types"
)

// Filter will clean the data from empty entries.
func Filter(config types.Config, data types.PerformanceData) types.PerformanceData {
	filteredData := types.PerformanceData{}
	for k, _ := range data {
		if config.MetricMappings[k].Name != "" {
			filteredData[k] = data[k]
		}
	}

	return filteredData
}

func Print(config types.Config, data types.PerformanceData) error {
	return render.PrintTable(data, config)
}

func Process(config types.Config, data types.PerformanceData) ([]cloudwatchtypes.MetricDatum, error) {

	var metricData []cloudwatchtypes.MetricDatum

	for key, value := range data {
		metric, ok := config.MetricMappings[key]
		if !ok {
			continue
		}

		metricValue := math.Round(value*100) / 100
		metricDatum := cloudwatchtypes.MetricDatum{
			MetricName: aws.String(metric.Name),
			Value:      aws.Float64(metricValue),
			Timestamp:  aws.Time(time.Now()),
			Unit:       cloudwatchtypes.StandardUnitCount,
		}

		// Handle timestamps if provided.
		if metric.Timestamp != "" {
			t, check, _ := datetime.ParseTimeString(config.MetricMappings[key].Timestamp)
			if check {
				metricDatum.Timestamp = aws.Time(t)
			}
		}

		for _, dimension := range metric.Dimensions {
			metricDatum.Dimensions = append(metricDatum.Dimensions, cloudwatchtypes.Dimension{
				Name:  &dimension.Name,
				Value: &dimension.Value,
			})
		}

		metricData = append(metricData, metricDatum)
	}

	return metricData, nil
}

// Publish will publish the metrics to the nominated AWS account.
func Publish(client *cloudwatch.Client, data []cloudwatchtypes.MetricDatum, config types.Config) error {

	// Do not publish until we're ready.
	if config.SkipPublish {
		fmt.Println("You have elected to not publish these metrics, exiting...")
		return nil
	}

	input := &cloudwatch.PutMetricDataInput{
		Namespace:  aws.String(config.MetricNamespace),
		MetricData: data,
	}

	if config.Interactive || prompt.Confirm("Do you want to proceed?") {
		_, err := client.PutMetricData(context.TODO(), input)
		if err != nil {
			return err
		}
		fmt.Println("Metrics published successfully!")
	} else {
		fmt.Println("Operation cancelled.")
	}

	return nil
}

func Sort(config types.Config) []string {
	// Get our keys, so that we can sort the data.
	keys := make([]string, 0, len(config.MetricMappings))
	for k := range config.MetricMappings {
		keys = append(keys, k)
	}

	// Sort the data.
	slices.Sort(keys)

	// Return the sorted array
	return keys
}
