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
	gh "github.com/mrsimonemms/golang-helpers"
	"github.com/mrsimonemms/golang-helpers/temporal"
	"github.com/mrsimonemms/temporal-go-starter/internal/app"
	"github.com/mrsimonemms/temporal-go-starter/internal/app/activities"
	"github.com/mrsimonemms/temporal-go-starter/internal/app/workflows"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"go.temporal.io/sdk/contrib/sysinfo"
	"go.temporal.io/sdk/worker"
)

const workerCmdName = "worker"

func newWorkerCmd() *cobra.Command {
	opts := struct {
		temporal *temporal.TemporalOpts
	}{
		temporal: &temporal.TemporalOpts{},
	}

	cmd := &cobra.Command{
		Use:   workerCmdName,
		Short: "Run the Temporal worker",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := temporal.NewConnection(temporal.ParseCobraOpts(
				opts.temporal,
				temporal.WithZerolog(&log.Logger),
				temporal.WithPrometheusMetrics(opts.temporal.MetricsListenAddress, opts.temporal.MetricsPrefix, nil),
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

			w := worker.New(c, app.TaskQueue, worker.Options{
				SysInfoProvider: sysinfo.SysInfoProvider(),
			})

			log.Debug().Msg("Registering workflows")
			w.RegisterWorkflow(workflows.ExampleWorkflow)

			acts, err := activities.NewActivities()
			if err != nil {
				return gh.FatalError{
					Cause: err,
					Msg:   "Unable to create activities",
				}
			}

			log.Debug().Msg("Registering activities")
			w.RegisterActivity(acts)

			log.Debug().Msg("Starting health check server")
			temporal.NewHealthCheck(cmd.Context(), []string{app.TaskQueue}, opts.temporal.HealthListenAddress, c)

			if err := w.Run(worker.InterruptCh()); err != nil {
				return gh.FatalError{
					Cause: err,
					Msg:   "Workflow stopped",
				}
			}

			return nil
		},
	}

	temporal.NewCobraOpts(cmd, opts.temporal)

	return cmd
}
