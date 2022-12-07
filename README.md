# GPATs

These are the Garden Performance Acceptance Tests.

## About

These tests run as part of the
[garden-runcpipeline](https://garden.ci.cf-app.com/) and will gate us from
shipping a new release if they fail. With this in mind, the goal is to keep
this test suite as minimal and quick as possible.

Metrics generated by the test suite are currently sent to Wavefront. This is
achieved by wiring up the custom [Wavefront
Reporter](https://github.com/cloudfoundry/garden-performance-acceptance-tests/blob/main/reporter/wavefront.go)
into Ginkgo. In the future, should we wish to send metrics to other services,
we can extend the suite by writing and wiring up a new custom reporter for that
service.

Note that the reporter package used here is based on this [reporter
package](https://github.com/cloudfoundry/benchmarkbbs/tree/main/reporter)
from the
[cloudfoundry/benchmarkbbs](https://github.com/cloudfoundry/benchmarkbbs) repo.

## Usage

In order to run the performance test suite you need to [deploy
garden](https://github.com/cloudfoundry/garden-runc-release/wiki/Creating-sandbox-environments-for-debugging#eden-deployments)
and run the fololowing command (you will need the garden server ip as shown in bosh -d <deployment> vms):

```
cd $HOME/workspace/garden-ci/directors/eden
fly -t garden-ci login -c https://garden.ci.cf-app.com/
GARDEN_ADDRESS="<garden-server-ip>" fly -t garden-ci execute \
  -c $HOME/workspace/garden-performance-acceptance-tests/ci/performance-tests.yml \
  -i gpats-main=$HOME/workspace/garden-performance-acceptance-tests
```

**NB**: This test suite will destroy ALL containers on the Garden server as
part of the run. You have been warned.

## Conditionally Expect Metrics

This suite is used both to gate releases via expectations on performance but
also to provide metrics when thresholds are succeeded. To achieve this, metrics
related expectations are wrapped in `Conditionally()` functions, to ensure they
only fail tests when required. To turn off metrics related expectations, set
the `IGNORE_PERF_EXPECTATIONS` environment variable.
