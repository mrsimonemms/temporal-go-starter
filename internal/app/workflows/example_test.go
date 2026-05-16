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

package workflows_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/mrsimonemms/temporal-go-starter/internal/app/activities"
	"github.com/mrsimonemms/temporal-go-starter/internal/app/workflows"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.temporal.io/sdk/testsuite"
)

// TestExampleWorkflow_Success runs the workflow against a mocked activity so
// that the workflow logic can be exercised in milliseconds. Timers created by
// `workflow.Sleep` are skipped automatically by the test environment.
func TestExampleWorkflow_Success(t *testing.T) {
	tests := []struct {
		name             string
		input            *workflows.ExampleInput
		activityResponse string
		wantResult       *workflows.ExampleResult
	}{
		{
			name:             "passes the input ID through and uses the activity response as the greeting",
			input:            &workflows.ExampleInput{ID: "123"},
			activityResponse: "hello",
			wantResult: &workflows.ExampleResult{
				ID:       "123",
				Greeting: "hello",
			},
		},
		{
			name:             "handles an empty input ID",
			input:            &workflows.ExampleInput{ID: ""},
			activityResponse: "anything",
			wantResult: &workflows.ExampleResult{
				ID:       "",
				Greeting: "anything",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ts := &testsuite.WorkflowTestSuite{}
			env := ts.NewTestWorkflowEnvironment()

			var a *activities.Activities
			env.RegisterActivity(a)

			env.OnActivity(a.Example, mock.Anything, &activities.ExampleInput{ID: tc.input.ID}).
				Return(&activities.ExampleResult{Response: tc.activityResponse}, nil)

			env.ExecuteWorkflow(workflows.ExampleWorkflow, tc.input)

			assert.True(t, env.IsWorkflowCompleted())
			assert.NoError(t, env.GetWorkflowError())

			var result workflows.ExampleResult
			assert.NoError(t, env.GetWorkflowResult(&result))
			assert.Equal(t, *tc.wantResult, result)
		})
	}
}

// TestExampleWorkflow_ActivityError verifies that a failing activity surfaces
// as a workflow error rather than being silently swallowed.
func TestExampleWorkflow_ActivityError(t *testing.T) {
	ts := &testsuite.WorkflowTestSuite{}
	env := ts.NewTestWorkflowEnvironment()

	var a *activities.Activities
	env.RegisterActivity(a)

	env.OnActivity(a.Example, mock.Anything, &activities.ExampleInput{ID: "boom"}).
		Return(nil, errors.New("activity failed"))

	env.ExecuteWorkflow(workflows.ExampleWorkflow, &workflows.ExampleInput{ID: "boom"})

	assert.True(t, env.IsWorkflowCompleted())

	err := env.GetWorkflowError()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error executing example activity")
}

// TestExampleInput_JSON locks in the JSON encoding of ExampleInput. The
// encoding is part of the public contract of any client that starts the
// workflow, so accidental field renames or tag changes should fail the test.
func TestExampleInput_JSON(t *testing.T) {
	tests := []struct {
		name  string
		input workflows.ExampleInput
		want  string
	}{
		{
			name:  "encodes the ID using its json tag",
			input: workflows.ExampleInput{ID: "abc"},
			want:  `{"id":"abc"}`,
		},
		{
			name:  "encodes an empty ID",
			input: workflows.ExampleInput{},
			want:  `{"id":""}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := json.Marshal(tc.input)
			assert.NoError(t, err)
			assert.JSONEq(t, tc.want, string(got))
		})
	}
}

// TestExampleResult_JSON locks in the JSON encoding of ExampleResult.
func TestExampleResult_JSON(t *testing.T) {
	tests := []struct {
		name   string
		result workflows.ExampleResult
		want   string
	}{
		{
			name:   "encodes both fields using their json tags",
			result: workflows.ExampleResult{ID: "abc", Greeting: "hi"},
			want:   `{"id":"abc","greeting":"hi"}`,
		},
		{
			name:   "encodes empty fields as empty strings",
			result: workflows.ExampleResult{},
			want:   `{"id":"","greeting":""}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := json.Marshal(tc.result)
			assert.NoError(t, err)
			assert.JSONEq(t, tc.want, string(got))
		})
	}
}
