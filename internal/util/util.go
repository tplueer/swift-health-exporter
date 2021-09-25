// Copyright 2021 SAP SE
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package util

import (
	"context"
	"os/exec"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/sapcc/swift-health-exporter/internal/collector"
)

// RunCommandWithTimeout runs a command with the provided timeout duration and returns its
// combined output.
func RunCommandWithTimeout(timeout time.Duration, name string, args ...string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return exec.CommandContext(ctx, name, args...).CombinedOutput()
}

// CmdArgsToStr returns a space separated string for cmdArgs.
func CmdArgsToStr(cmdArgs []string) string {
	return strings.Join(cmdArgs, " ")
}

// AddTask adds a Task to the given Collector and the Scraper along
// with its corresponding exit code GaugeVec.
func AddTask(
	shouldAdd bool,
	c *collector.Collector,
	s *collector.Scraper,
	t collector.Task,
	exitCode *prometheus.GaugeVec) {

	if shouldAdd {
		name := t.Name()
		c.Tasks[name] = t
		s.Tasks[name] = t
		s.ExitCodeGaugeVec[name] = exitCode
	}
}
