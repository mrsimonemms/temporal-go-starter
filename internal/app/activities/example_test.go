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

package activities_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/mrsimonemms/temporal-go-starter/internal/app/activities"
	"github.com/stretchr/testify/assert"
	"go.temporal.io/sdk/testsuite"
)

// TestExampleActivity exercises the Example activity through the Temporal
// activity test environment. This is the recommended way to test activities
// that depend on the activity context (for logging, heartbeats, and so on).
//
// The activity sleeps for several seconds, so the test is skipped when
// `go test -short` is used.
func TestExampleActivity(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping slow activity test in short mode")
	}

	tests := []struct {
		name  string
		input *activities.ExampleInput
	}{
		{
			name:  "returns a response for a simple input",
			input: &activities.ExampleInput{ID: "abc-123"},
		},
		{
			name:  "returns a response for an empty input ID",
			input: &activities.ExampleInput{ID: ""},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ts := &testsuite.WorkflowTestSuite{}
			env := ts.NewTestActivityEnvironment()

			a, err := activities.NewActivities()
			assert.NoError(t, err)
			env.RegisterActivity(a)

			val, err := env.ExecuteActivity(a.Example, tc.input)
			assert.NoError(t, err)

			var result activities.ExampleResult
			assert.NoError(t, val.Get(&result))

			assert.NotEmpty(t, result.Response, "expected a generated response value")

			_, parseErr := uuid.Parse(result.Response)
			assert.NoError(t, parseErr, "expected response to be a valid UUID")
		})
	}
}

// TestNewActivities verifies that the activities constructor returns a usable
// value. The test exists so that customers extending NewActivities have a
// natural place to assert any new wiring (for example, that a database client
// is not nil).
func TestNewActivities(t *testing.T) {
	a, err := activities.NewActivities()

	assert.NoError(t, err)
	assert.NotNil(t, a)
}
