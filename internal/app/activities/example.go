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

package activities

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mrsimonemms/temporal-go-starter/internal/observability"
	"go.temporal.io/sdk/activity"
)

type ExampleInput struct {
	ID string `json:"id"`
}

type ExampleResult struct {
	Response string `json:"response"`
}

func (a *Activities) Example(ctx context.Context, input *ExampleInput) (*ExampleResult, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Starting Example activity", "input", input)

	startTime := time.Now()

	// Increment the Prometheus counter for started example activities
	observability.ExampleStartedTotal.Inc()

	logger.Info("Sleeping to indicate activity")
	time.Sleep(5 * time.Second)

	// Record the duration of the activity execution in the Prometheus histogram
	observability.ExampleDurationSeconds.Observe(time.Since(startTime).Seconds())

	return &ExampleResult{
		Response: uuid.NewString(),
	}, nil
}
