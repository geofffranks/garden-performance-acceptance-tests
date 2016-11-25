package reporter

import (
	"errors"
	"fmt"
	"time"

	"code.cloudfoundry.org/lager"
	"github.com/onsi/ginkgo/config"
	"github.com/onsi/ginkgo/types"
	"github.com/zorkian/go-datadog-api"
)

type DataDogReporter struct {
	logger        lager.Logger
	metricPrefix  string
	dataDogClient *datadog.Client
}

func NewDataDogReporter(
	logger lager.Logger,
	metricPrefix string,
	dataDogClient *datadog.Client,
) DataDogReporter {
	return DataDogReporter{
		logger:        logger.Session("datadog-reporter"),
		metricPrefix:  metricPrefix,
		dataDogClient: dataDogClient,
	}
}

func (r *DataDogReporter) SpecSuiteWillBegin(config config.GinkgoConfigType, summary *types.SuiteSummary) {
}

func (r *DataDogReporter) BeforeSuiteDidRun(setupSummary *types.SetupSummary) {
}

func (r *DataDogReporter) AfterSuiteDidRun(setupSummary *types.SetupSummary) {
}

func (r *DataDogReporter) SpecWillRun(specSummary *types.SpecSummary) {
}

func (r *DataDogReporter) SpecDidComplete(specSummary *types.SpecSummary) {
	if specSummary.Failed() {
		eventTitle := fmt.Sprintf("%s-test-failure", r.metricPrefix)
		eventText := fmt.Sprintf("%s - %s", specSummary.ComponentTexts[1], specSummary.ComponentTexts[2])
		eventTags := []string{r.metricPrefix}

		failEvent := &datadog.Event{
			Title: eventTitle,
			Text:  eventText,
			Tags:  eventTags,
		}

		_, err := r.dataDogClient.PostEvent(failEvent)
		if err != nil {
			r.logger.Error("failed-sending-events-to-datadog", err, lager.Data{"metric": "failevent", "prefix": r.metricPrefix})
		}
	}

	if specSummary.Passed() && specSummary.IsMeasurement {
		for _, measurement := range specSummary.Measurements {
			if measurement.Info == nil {
				panic(fmt.Sprintf("%#v", specSummary))
			}

			info, ok := measurement.Info.(ReporterInfo)
			if !ok {
				r.logger.Error("failed-type-assertion-on-measurement-info", errors.New("type-assertion-failed"))
				continue
			}

			if info.MetricName == "" {
				r.logger.Error("failed-blank-metric-name", errors.New("blank-metric-name"))
				continue
			}

			timestamp := float64(time.Now().Unix())
			r.logger.Info("sending-metrics-to-datadog", lager.Data{"metric": info.MetricName, "prefix": r.metricPrefix})
			err := r.dataDogClient.PostMetrics([]datadog.Metric{
				{
					Metric: fmt.Sprintf("%s.%s-slowest", r.metricPrefix, info.MetricName),
					Points: []datadog.DataPoint{
						{timestamp, measurement.Largest},
					},
				},
				{
					Metric: fmt.Sprintf("%s.%s-fastest", r.metricPrefix, info.MetricName),
					Points: []datadog.DataPoint{
						{timestamp, measurement.Smallest},
					},
				},
				{
					Metric: fmt.Sprintf("%s.%s-average", r.metricPrefix, info.MetricName),
					Points: []datadog.DataPoint{
						{timestamp, measurement.Average},
					},
				},
			})
			if err != nil {
				r.logger.Error("failed-sending-metrics-to-datadog", err, lager.Data{"metric": info.MetricName, "prefix": r.metricPrefix})
				continue
			}
			r.logger.Info("sending-metrics-to-datadog-complete", lager.Data{"metric": info.MetricName, "prefix": r.metricPrefix})
		}
	}
}

func (r *DataDogReporter) SpecSuiteDidEnd(summary *types.SuiteSummary) {
}
