package main

import (
	"context"
	"fmt"
	"github.com/pterm/pterm"
	"log"
	"math"
	"os"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"gopkg.in/yaml.v2"
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
type PerformanceData map[string]float64

var (
	cliRegion      = kingpin.Flag("region", "AWS Region to push metrics").Envar("AWS_REGION").String()
	cliProfile     = kingpin.Flag("profile", "Configured AWS profile to use").Envar("AWS_PROFILE").String()
	cliSkipPublish = kingpin.Flag("skip-publish", "Skip publishing metrics").Default("false").Bool()
)

// run will execute the main logic component for error handling.
func run() error {

	configInput, err := loadConfig("config.yml")
	if err != nil {
		return err
	}

	if configInput.Region == "" {
		configInput.Region = *cliRegion
		if configInput.Region == "" {
			return fmt.Errorf("AWS_REGION environment variable not set")
		}
	}

	if configInput.Profile == "" {
		configInput.Profile = *cliProfile
		if configInput.Profile == "" {
			return fmt.Errorf("AWS_PROFILE environment variable not set")
		}
	}

	if *cliSkipPublish {
		configInput.SkipPublish = true
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

// printTable will print a table showing all of the metrics which are going to be pushed.
func printTable(data PerformanceData, config Config) error {
	alternateStyle := pterm.NewStyle(pterm.BgDarkGray)
	tableData := pterm.TableData{
		{"Metric name", "Value"},
	}

	for key, val := range data {
		tableData = append(tableData, []string{config.MetricMappings[key], fmt.Sprint(math.Round(val*10) / 10)})
	}

	fmt.Println("Metrics to be published:")
	return pterm.DefaultTable.WithHasHeader().WithBoxed().WithData(tableData).WithStyle(alternateStyle).Render()
}

// publishMetrics will publish the metrics to the nominated AWS account.
func publishMetrics(client *cloudwatch.Client, data PerformanceData, config Config) error {
	err := printTable(data, config)
	if err != nil {
		return err
	}

	// Do not publish until we're ready.
	if config.SkipPublish {
		fmt.Println("You have elected to not publish these metrics, exiting...")
		return nil
	}

	var metricData []types.MetricDatum

	for key, value := range data {
		metricName, ok := config.MetricMappings[key]
		if !ok {
			continue
		}

		metricValue := math.Round(value*10) / 10

		metricData = append(metricData, types.MetricDatum{
			MetricName: aws.String(metricName),
			Value:      aws.Float64(metricValue),
			Timestamp:  aws.Time(time.Now()),
			Unit:       types.StandardUnitCount,
		})

		input := &cloudwatch.PutMetricDataInput{
			Namespace:  aws.String(config.MetricNamespace),
			MetricData: metricData,
		}

		_, err := client.PutMetricData(context.TODO(), input)
		if err != nil {
			return err
		}
	}

	fmt.Println("Metrics published successfully!")
	return nil
}

func main() {
	kingpin.Parse()
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
