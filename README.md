# Personal performance tracking

Measure your own personally defined metrics in AWS CloudWatch.

## The idea

This is a simple app that takes arbitrary inputs which are defined in `config.yml`, with the data defined in `data.yml`
and pushes the metrics to AWS. This will enable you to track them over time in graphs of your choosing without being
restricted to keeping all the data in a single spot for manual analysis.

* Track metrics you define on a regular cadence
* Build yourself a CloudWatch dashboard (not provided here)
* Measure your success over time to achieve your goals.

## Why

I am using this utility to push AWS CloudWatch Metrics to my account in order to achieve the following:

- Tracking time spent on specific volunteering responsibilities.
- To facilitate an evidence-based approach to self-improvement.
- To compare the amount of time spent between very specific responsibilities to optimize time and personal expectations.
- To uphold myself to specific KPI's.

## Roadmap

* Break the structure up a bit more
* Unit tests
* Command-line parameters for configuration
* Environment variables for configuration

## The project

Still rudimentary and designed for my own use at this point, it is not designed to encompass more than what is defined
above. It is still missing some niceties such as CLI parameters, but essentially here's how you should use it.

### Set up the configuration file

Config contains the definition of the fields being tracked which marry up to the data file. This includes the information
required to connect to AWS, such as the profile being used and the region. This is a great opportunity to define the
metrics you want to capture.

```yaml
region: ap-southeast-2
profile: my-aws-profile
metricNamespace: Personal/Performance
skipPublish: false
metricMappings:
  your-metric-here:
    name: MyCustomMetricName
    dimensions:
      - name: Goal
        value: Fitness
```

### Set up the data file

Matching up the data to the defined data set above, you can now identify the metric values you wish to send to your
AWS account before pushing them.

```yaml
your-metric-here: 100
```

### Pushing your metrics

Everything is now set up, so all that is left is for you to push the data.

```
go run main.go
```

## License

MIT, use at your own risk.

Note that CloudWatch metrics do have a cost associated to their usage. Please read the code and understand the costs
of what is being created in your account before you decide to use this.