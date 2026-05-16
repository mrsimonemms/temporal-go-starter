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

package cmd

import (
	"encoding/json"
	"fmt"

	gh "github.com/mrsimonemms/golang-helpers"
	"github.com/mrsimonemms/golang-helpers/temporal"
	"github.com/mrsimonemms/temporal-go-starter/internal/app"
	"github.com/mrsimonemms/temporal-go-starter/internal/app/workflows"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"go.temporal.io/sdk/client"
)

func newStartCmd() *cobra.Command {
	opts := struct {
		temporal *temporal.TemporalOpts
	}{
		temporal: &temporal.TemporalOpts{},
	}

	cmd := &cobra.Command{
		Use:   "start",
		Short: "Trigger the workflow",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := temporal.NewConnection(temporal.ParseCobraOpts(
				opts.temporal,
				temporal.WithZerolog(&log.Logger),
			)...)
			if err != nil {
				return gh.FatalError{
					Cause: err,
					Msg:   "Unable to create client",
				}
			}
			defer func() {
				log.Trace().Msg("Closing Temporal connection")
				c.Close()
				log.Trace().Msg("Temporal connection closed")
			}()

			input := &workflows.ExampleInput{
				ID: "123",
			}

			workflowOptions := client.StartWorkflowOptions{
				TaskQueue: app.TaskQueue,
			}

			we, err := c.ExecuteWorkflow(cmd.Context(), workflowOptions, workflows.ExampleWorkflow, input)
			if err != nil {
				return gh.FatalError{
					Cause: err,
					Msg:   "Error executing workflow",
				}
			}

			log.Info().Str("workflowId", we.GetID()).Str("runId", we.GetRunID()).Msg("Started workflow")

			var result workflows.ExampleResult
			if err := we.Get(cmd.Context(), &result); err != nil {
				return gh.FatalError{
					Cause: err,
					Msg:   "Error getting response",
				}
			}

			log.Info().Interface("result", result).Msg("Workflow completed")

			f, err := json.MarshalIndent(result, "", "  ")
			if err != nil {
				return gh.FatalError{
					Cause: err,
					Msg:   "Error marshalling result",
				}
			}
			fmt.Println("===")
			fmt.Println(string(f))
			fmt.Println("===")

			return nil
		},
	}

	temporal.NewCobraOpts(cmd, opts.temporal)

	return cmd
}
