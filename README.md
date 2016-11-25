# GPATs

These are the Garden Performance Acceptance Tests.

## About

These tests run as part of the [garden-runc pipeline](https://garden.ci.cf-app.com/) and will gate us from shipping a new release if they fail.
With this in mind, the goal is to keep this test suite as minimal and quick as possible.

Metrics generated by the test suite are currently sent to DataDog. This is achieved by wiring up the custom DataDog Reporter into Ginkgo.
In the future, should we wish to send metrics to other services, we can extend the suite by writing and wiring up a new custom reporter for that service.

Note that the reporter package used here is based on this [reporter package](https://github.com/cloudfoundry/benchmarkbbs/tree/master/reporter) from the [cloudfoundry/benchmarkbbs](https://github.com/cloudfoundry/benchmarkbbs) repo.

## Usage

```
DATADOG_API_KEY="x" DATADOG_APP_KEY="x" DATADOG_METRIC_PREFIX="x" GARDEN_ADDRESS="127.0.0.1:7777" ginkgo
```

**NB**: This test suite will destroy ALL containers on the Garden server as part of the run. You have been warned.
