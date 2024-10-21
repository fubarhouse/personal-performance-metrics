package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/pterm/pterm"
	"gopkg.in/yaml.v3"
)

// Config provides global configuration
type Config struct {
	Region          string                   `yaml:"region"`
	Profile         string                   `yaml:"profile"`
	SkipPublish     bool                     `yaml:"skipPublish"`
	MetricNamespace string                   `yaml:"metricNamespace"`
	MetricMappings  map[string]MetricMapping `yaml:"metricMappings"`
}

// MetricMapping is the configuration data for the metrics.
type MetricMapping struct {
	Name       string                    `yaml:"name"`
	Dimensions []MetricMappingDimensions `yaml:"dimensions"`
}

// MetricMappingDimensions is the definition for the dimensions associated to the metric.
type MetricMappingDimensions struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

// PerformanceData is the data being captured and sent to AWS.
type PerformanceData map[string]float64

var (
	cliRegion         = kingpin.Flag("region", "AWS Region to push metrics").Envar("AWS_REGION").String()
	cliProfile        = kingpin.Flag("profile", "Configured AWS profile to use").Envar("AWS_PROFILE").String()
	cliSkipPublish    = kingpin.Flag("skip-publish", "Skip publishing metrics").Default("false").Bool()
	cliNoninteractive = kingpin.Flag("non-interactive", "Perform work without interactions").Default("false").Bool()
)

// run will execute the main logic component for error handling.
func run() error {

	configInput, err := loadConfig()
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

	dataInput, err := loadData()
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
func loadConfig() (Config, error) {
	var cfg Config
	file, err := os.ReadFile("config.yml")
	if err != nil {
		return cfg, err
	}
	err = yaml.Unmarshal(file, &cfg)
	return cfg, err
}

// lodaData will load the data file.
func loadData() (PerformanceData, error) {
	var data PerformanceData
	file, err := os.ReadFile("data.yml")
	if err != nil {
		return data, err
	}
	err = yaml.Unmarshal(file, &data)
	return data, err
}

// printTable will print a table showing all the metrics which are going to be pushed.
func printTable(data PerformanceData, config Config) error {
	alternateStyle := pterm.NewStyle(pterm.BgDarkGray)
	tableData := pterm.TableData{
		{"Metric name", "Value", "Dimensions"},
	}
	fmt.Sprint(alternateStyle, tableData)

	for key, val := range data {
		var dimensions string
		for _, v := range config.MetricMappings[key].Dimensions {
			dimensions += fmt.Sprintf("%s=%s ", v.Name, v.Value)
		}
		tableData = append(tableData, []string{config.MetricMappings[key].Name, fmt.Sprint(math.Round(val*100) / 100), dimensions})
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
		metric, ok := config.MetricMappings[key]
		if !ok {
			continue
		}

		metricValue := math.Round(value*100) / 100
		metricDatum := types.MetricDatum{
			MetricName: aws.String(metric.Name),
			Value:      aws.Float64(metricValue),
			Timestamp:  aws.Time(time.Now()),
			Unit:       types.StandardUnitCount,
		}

		for _, dimension := range metric.Dimensions {
			metricDatum.Dimensions = append(metricDatum.Dimensions, types.Dimension{
				Name:  &dimension.Name,
				Value: &dimension.Value,
			})
		}

		metricData = append(metricData, metricDatum)
	}

	input := &cloudwatch.PutMetricDataInput{
		Namespace:  aws.String(config.MetricNamespace),
		MetricData: metricData,
	}

	if *cliNoninteractive || confirm("Do you want to proceed?") {
		_, err = client.PutMetricData(context.TODO(), input)
		if err != nil {
			return err
		}
		fmt.Println("Metrics published successfully!")
	} else {
		fmt.Println("Operation cancelled.")
	}

	return nil
}

// confirm will accept input for a prompt.
func confirm(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/n]: ", prompt)

		response, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			return false
		}

		response = strings.ToLower(strings.TrimSpace(response))

		switch response {
		case "y", "yes":
			return true
		case "n", "no":
			return false
		default:
			fmt.Println("Invalid input. Please enter 'y' or 'n'.")
		}
	}
}

func main() {
	kingpin.Parse()
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
