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

package observability

import "github.com/prometheus/client_golang/prometheus"

const (
	metricNamespace = "example"
)

var (
	ExampleStartedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: metricNamespace,
			Name:      "started_total",
			Help:      "Total number of example workflows started",
		},
	)

	ExampleDurationSeconds = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: metricNamespace,
			Name:      "duration_seconds",
			Help:      "Duration of example application workflow execution",
			Buckets:   prometheus.DefBuckets,
		},
	)
)

func init() {
	prometheus.MustRegister(
		ExampleStartedTotal,
		ExampleDurationSeconds,
	)
}
