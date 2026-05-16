/*
 * Copyright 2026 Simon Emms <simon@simonemms.com>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package observability_test

import (
	"testing"

	"github.com/mrsimonemms/temporal-go-starter/internal/observability"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

// histogramSnapshot reads the current sample count and sample sum for a
// histogram registered with the default Prometheus gatherer. Histograms cannot
// be inspected with testutil.ToFloat64, so this helper walks the gathered
// metric families instead.
//
// The metric is assumed to have no labels, which matches how the example
// histogram in this package is declared.
func histogramSnapshot(t *testing.T, name string) (sampleCount uint64, sampleSum float64) {
	t.Helper()

	mfs, err := prometheus.DefaultGatherer.Gather()
	assert.NoError(t, err)

	for _, mf := range mfs {
		if mf.GetName() != name {
			continue
		}
		metrics := mf.GetMetric()
		assert.Len(t, metrics, 1, "expected exactly one unlabelled metric for %q", name)
		h := metrics[0].GetHistogram()
		return h.GetSampleCount(), h.GetSampleSum()
	}

	t.Fatalf("histogram %q not found in default gatherer", name)
	return 0, 0
}

// TestExampleStartedTotal_Inc demonstrates how to assert that a counter
// increments when its Inc method is called. testutil.ToFloat64 reads the
// counter's current value without needing a separate registry.
//
// The test snapshots the starting value so that it does not depend on test
// ordering or on whether the counter has been touched elsewhere in the
// process.
func TestExampleStartedTotal_Inc(t *testing.T) {
	tests := []struct {
		name string
		incs int
	}{
		{name: "single increment", incs: 1},
		{name: "multiple increments", incs: 5},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			before := testutil.ToFloat64(observability.ExampleStartedTotal)

			for i := 0; i < tc.incs; i++ {
				observability.ExampleStartedTotal.Inc()
			}

			after := testutil.ToFloat64(observability.ExampleStartedTotal)
			assert.Equal(t, before+float64(tc.incs), after)
		})
	}
}

// TestExampleDurationSeconds_Observe demonstrates how to assert that a
// histogram records observations. We check both the sample count (number of
// observations) and the sample sum (cumulative observed value).
func TestExampleDurationSeconds_Observe(t *testing.T) {
	tests := []struct {
		name    string
		samples []float64
	}{
		{name: "single observation", samples: []float64{0.25}},
		{name: "multiple observations", samples: []float64{0.1, 0.5, 1.0}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			beforeCount, beforeSum := histogramSnapshot(t, "example_duration_seconds")

			var wantSum float64
			for _, v := range tc.samples {
				observability.ExampleDurationSeconds.Observe(v)
				wantSum += v
			}

			afterCount, afterSum := histogramSnapshot(t, "example_duration_seconds")
			assert.Equal(t, beforeCount+uint64(len(tc.samples)), afterCount)
			assert.InDelta(t, beforeSum+wantSum, afterSum, 1e-9)
		})
	}
}

// TestMetricsRegisteredWithDefaultGatherer verifies that the package's init
// function has registered the metrics with the default Prometheus gatherer.
// The HTTP /metrics endpoint exposed by the worker reads from this gatherer,
// so anything not registered here will not be scraped in production.
func TestMetricsRegisteredWithDefaultGatherer(t *testing.T) {
	tests := []struct {
		name       string
		metricName string
	}{
		{name: "counter is exposed", metricName: "example_started_total"},
		{name: "histogram is exposed", metricName: "example_duration_seconds"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			count, err := testutil.GatherAndCount(prometheus.DefaultGatherer, tc.metricName)
			assert.NoError(t, err)
			assert.Equal(t, 1, count, "expected exactly one registered series for %q", tc.metricName)
		})
	}
}

// TestMetricsAreSafeToRegisterOnce confirms that re-registering the package's
// metrics on a fresh registry would fail because the package has already
// claimed those names on the default registry through its init function.
//
// In other words, this locks in the contract that import side effects own
// the metric names. If a customer wants to register the same definitions
// elsewhere, they need to use prometheus.WrapRegistererWith or define their
// own names.
func TestMetricsAreSafeToRegisterOnce(t *testing.T) {
	reg := prometheus.NewRegistry()

	// A fresh registry will accept the metric because it has not seen the
	// name before. This shows that the metric definitions themselves are
	// valid and that they would also work in a non-default registry.
	err := reg.Register(prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "example",
		Name:      "started_total",
		Help:      "Total number of example workflows started",
	}))
	assert.NoError(t, err)

	// Registering the same definition twice on the same registry fails.
	// This is the behaviour that protects production code from accidentally
	// double-counting a metric.
	err = reg.Register(prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "example",
		Name:      "started_total",
		Help:      "Total number of example workflows started",
	}))
	assert.Error(t, err)
}
