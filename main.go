package main

import (
	"context"
	"fmt"
	"log"

	"github.com/alecthomas/kingpin/v2"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"

	"github.com/fubarhouse/personal-performance-metrics/internal/config"
	"github.com/fubarhouse/personal-performance-metrics/internal/metrics"
)

var (
	cliRegion      = kingpin.Flag("region", "AWS Region to push metrics").Envar("AWS_REGION").String()
	cliProfile     = kingpin.Flag("profile", "Configured AWS profile to use").Envar("AWS_PROFILE").String()
	cliSkipPublish = kingpin.Flag("skip-publish", "Skip publishing metrics").Default("false").Bool()
	cliInteractive = kingpin.Flag("interactive", "Perform work using interactions").Default("false").Bool()
)

// run will execute the main logic component for error handling.
func run() error {

	ctx := context.Background()

	configInput, err := config.LoadConfig(nil)
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

	if *cliInteractive {
		configInput.Interactive = true
	}

	if *cliSkipPublish {
		configInput.SkipPublish = true
	}

	dataInput, err := config.LoadData(nil)
	if err != nil {
		return err
	}

	// Prepare AWS configuration options
	var opts []func(*awsConfig.LoadOptions) error

	// Add profile option if provided
	if configInput.Profile != "" {
		opts = append(opts, awsConfig.WithSharedConfigProfile(configInput.Profile))
	}

	// Add region option if provided
	if configInput.Region != "" {
		opts = append(opts, awsConfig.WithRegion(configInput.Region))
	}

	// Load AWS configuration
	cfg, err := awsConfig.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		fmt.Println("Error creating AWS config:", err)
		return err
	}

	// Create CloudWatch client
	client := cloudwatch.NewFromConfig(cfg)

	// Filter unwanted metrics
	filteredData := metrics.Filter(configInput, dataInput)

	// Process metrics
	processedData, err := metrics.Process(configInput, filteredData)
	if err != nil {
		fmt.Println("Error processing metrics:", err)
	}

	// Print metrics
	err = metrics.Print(configInput, filteredData)
	if err != nil {
		fmt.Println("Error printing metrics:", err)
	}

	// Publish metrics
	err = metrics.Publish(client, processedData, configInput)
	if err != nil {
		fmt.Println("Error publishing metrics:", err)
	}

	return nil
}

func main() {
	kingpin.Parse()
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
