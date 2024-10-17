package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
)

// Config provides global configuration
type Config struct {
	Region          string            `yaml:"region"`
	Profile         string            `yaml:"profile"`
	SkipPublish     bool              `yaml:"skipPublish"`
	MetricNamespace string            `yaml:"metricNamespace"`
	MetricMappings  map[string]string `yaml:"metricMappings"`
}

// PerformanceData is the data being captured and sent to AWS.
type PerformanceData map[string]int

// run will execute the main logic component for error handling.
func run() error {

	configInput, err := loadConfig("config.yml")
	if err != nil {
		return err
	}

	dataInput, err := loadData("data.yml")
	if err != nil {
		return err
	}

	// Prepare AWS configuration options
	var opts []func(*config.LoadOptions) error

	// Add profile option if provided
	if configInput.Profile != "" {
		opts = append(opts, config.WithSharedConfigProfile(configInput.Profile))
	}

	// Add region option if provided
	if configInput.Region != "" {
		opts = append(opts, config.WithRegion(configInput.Region))
	}

	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(), opts...)
	if err != nil {
		fmt.Println("Error creating AWS config:", err)
		return err
	}

	// Create CloudWatch client
	client := cloudwatch.NewFromConfig(cfg)

	// Publish metrics
	err = publishMetrics(client, dataInput, configInput)
	if err != nil {
		return err
	}

	fmt.Println("Metrics published successfully!")
	return nil
}

// loadConfig will load the configuration file.
func loadConfig(filename string) (Config, error) {
	var config Config
	file, err := os.ReadFile(filename)
	if err != nil {
		return config, err
	}
	err = yaml.Unmarshal(file, &config)
	return config, err
}

// lodaData will load the data file.
func loadData(filename string) (PerformanceData, error) {
	var data PerformanceData
	file, err := os.ReadFile(filename)
	if err != nil {
		return data, err
	}
	err = yaml.Unmarshal(file, &data)
	return data, err
}

// publishMetrics will publish the metrics to the nominated AWS account.
func publishMetrics(client *cloudwatch.Client, data PerformanceData, config Config) error {
	fmt.Printf("Metrics: %v\n", data)

	// Do not publish until we're ready.
	if config.SkipPublish {
		fmt.Println("Skipping.")
		return nil
	}

	var metricData []types.MetricDatum

	for key, value := range data {
		metricName, ok := config.MetricMappings[key]
		if !ok {
			continue
		}

		metricData = append(metricData, types.MetricDatum{
			MetricName: aws.String(metricName),
			Value:      aws.Float64(float64(value)),
			Timestamp:  aws.Time(time.Now()),
			Unit:       types.StandardUnitCount,
		})
	}

	input := &cloudwatch.PutMetricDataInput{
		Namespace:  aws.String(config.MetricNamespace),
		MetricData: metricData,
	}

	_, err := client.PutMetricData(context.TODO(), input)
	return err
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
