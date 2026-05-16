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

package workflows

import (
	"fmt"
	"time"

	"github.com/mrsimonemms/temporal-go-starter/internal/app/activities"
	"go.temporal.io/sdk/workflow"
)

type ExampleInput struct {
	ID string `json:"id"`
}

type ExampleResult struct {
	ID       string `json:"id"`
	Greeting string `json:"greeting"`
}

func ExampleWorkflow(ctx workflow.Context, input *ExampleInput) (*ExampleResult, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Starting ExampleWorkflow", "input", input)

	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Second,
	})

	var a *activities.Activities

	var result *activities.ExampleResult
	if err := workflow.ExecuteActivity(ctx, a.Example, &activities.ExampleInput{ID: input.ID}).Get(ctx, &result); err != nil {
		logger.Error("Error executing Example activity", "error", err)
		return nil, fmt.Errorf("error executing example activity: %w", err)
	}

	if err := workflow.Sleep(ctx, 5*time.Second); err != nil {
		logger.Error("Error sleeping", "error", err)
		return nil, fmt.Errorf("error sleeping: %w", err)
	}

	return &ExampleResult{
		ID:       input.ID,
		Greeting: result.Response,
	}, nil
}
