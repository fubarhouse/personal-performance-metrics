package types

// Config provides global configuration for metrics.
type Config struct {
	Interactive     bool                     `yaml:"interactive"`
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
