package garden_performance_acceptance_tests_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"os"
	"testing"

	"code.cloudfoundry.org/garden"
	"code.cloudfoundry.org/garden-performance-acceptance-tests/reporter"
	"code.cloudfoundry.org/garden/client"
	"code.cloudfoundry.org/garden/client/connection"
	"code.cloudfoundry.org/lager"
	datadog "github.com/zorkian/go-datadog-api"
)

var (
	gardenClient           garden.Client
	ignorePerfExpectations bool // allows us to report metrics even when an expectation fails
)

var _ = BeforeSuite(func() {
	gardenHost := os.Getenv("GARDEN_ADDRESS")
	if gardenHost == "" {
		gardenHost = "127.0.0.1"
	}
	gardenPort := os.Getenv("GARDEN_PORT")
	if gardenPort == "" {
		gardenPort = "7777"
	}
	gardenClient = client.New(connection.New("tcp", fmt.Sprintf("%s:%s", gardenHost, gardenPort)))

	// ensure a 'clean' starting state
	cleanupContainers()
})

func TestGardenPerformanceAcceptanceTests(t *testing.T) {
	dataDogAPIKey := os.Getenv("DATADOG_API_KEY")
	dataDogAppKey := os.Getenv("DATADOG_APP_KEY")
	metricPrefix := os.Getenv("DATADOG_METRIC_PREFIX")
	if metricPrefix == "" {
		metricPrefix = "gpats"
	}

	if os.Getenv("IGNORE_PERF_EXPECTATIONS") != "" {
		ignorePerfExpectations = true
	}

	logger := lager.NewLogger("garden-performance-acceptance-tests")
	logger.RegisterSink(lager.NewWriterSink(GinkgoWriter, lager.INFO))

	customReporters := []Reporter{}
	if dataDogAPIKey != "" && dataDogAppKey != "" {
		dataDogClient := datadog.NewClient(dataDogAPIKey, dataDogAppKey)
		reporter := reporter.NewDataDogReporter(logger, metricPrefix, dataDogClient)
		customReporters = append(customReporters, &reporter)
	}

	RegisterFailHandler(Fail)
	RunSpecsWithDefaultAndCustomReporters(t, "GardenPerformanceAcceptanceTests Suite", customReporters)
}

func cleanupContainers() {
	containers, err := gardenClient.Containers(garden.Properties{})
	Expect(err).NotTo(HaveOccurred())

	for _, container := range containers {
		Expect(gardenClient.Destroy(container.Handle())).To(Succeed())
	}
}

func Conditionally(expectation func(), condition bool) {
	if condition {
		expectation()
	}
}
