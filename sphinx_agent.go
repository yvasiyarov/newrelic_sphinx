package main

import (
	"flag"
	"github.com/yvasiyarov/newrelic_platform_go"
	"log"
)

var sphinxPort = flag.Int("sphinx-port", 9312, "Sphinx port")
var sphinxHost = flag.String("sphinx-host", "127.0.0.1", "Sphinx host")

var newrelicName = flag.String("newrelic-name", "Sphinx", "Name in Sphinx")

var newrelicLicense = flag.String("newrelic-license", "", "Newrelic license")

var verbose = flag.Bool("verbose", false, "Verbose mode")

const (
	MIN_PAUSE_TIME            = 30 //do not query sphinx often than once in 30 seconds
	SPHINX_CONNECTION_TIMEOUT = 0  //no timeout
	NEWRELIC_POLL_INTERVAL    = 60 //Send data to newrelic every 60 seconds

	AGENT_GUID    = "com.github.yvasiyarov.Sphinx"
	AGENT_VERSION = "0.0.2"
)

func addMetrcsToComponent(component newrelic_platform_go.IComponent, metrics []newrelic_platform_go.IMetrica) {
	for _, m := range metrics {
		component.AddMetrica(m)
	}
}

func plainMetricsBuilder(metrics []*Metrica, dataSource *MetricsDataSource) []newrelic_platform_go.IMetrica {
	result := make([]newrelic_platform_go.IMetrica, len(metrics))
	for i, m := range metrics {
		m.DataSource = dataSource
		result[i] = m
	}
	return result
}
func incrementalMetricsBuilder(metrics []*Metrica, dataSource *MetricsDataSource) []newrelic_platform_go.IMetrica {
	incMetrics := make([]newrelic_platform_go.IMetrica, len(metrics))
	for i, m := range metrics {
		m.DataSource = dataSource
		incMetrics[i] = &IncrementalMetrica{*m}
	}
	return incMetrics
}

func main() {
	flag.Parse()
	if *newrelicLicense == "" {
		log.Fatalf("Please, pass a valid newrelic license key.\n Use --help to get more information about available options\n")
	}

	plugin := newrelic_platform_go.NewNewrelicPlugin(AGENT_VERSION, *newrelicLicense, NEWRELIC_POLL_INTERVAL)
	component := newrelic_platform_go.NewPluginComponent(*newrelicName, AGENT_GUID, *verbose)
	plugin.AddComponent(component)

	ds := NewMetricsDataSource(*sphinxHost, *sphinxPort, SPHINX_CONNECTION_TIMEOUT)
	addMetrcsToComponent(component, plainMetricsBuilder(plainMetrics, ds))
	addMetrcsToComponent(component, incrementalMetricsBuilder(incrementalMetrics, ds))

	plugin.Verbose = *verbose
	plugin.Run()
}
