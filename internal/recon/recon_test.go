// Copyright 2019 SAP SE
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

package recon

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sapcc/go-bits/assert"
	"github.com/sapcc/swift-health-exporter/internal/collector"
)

func TestReconCollector(t *testing.T) {
	isTest = true

	pathToExecutable, err := filepath.Abs("../../build/mock-swift-recon")
	if err != nil {
		t.Error(err)
	}

	registry := prometheus.NewPedanticRegistry()
	c := collector.New(0)
	exitCode := GetTaskExitCodeTypedDesc(registry)
	opts := &TaskOpts{
		PathToExecutable: pathToExecutable,
		HostTimeout:      1,
		CtxTimeout:       4 * time.Second,
	}
	c.AddTask(true, NewDiskUsageTask(opts), exitCode)
	c.AddTask(true, NewDriveAuditTask(opts), exitCode)
	c.AddTask(true, NewMD5Task(opts), exitCode)
	c.AddTask(true, NewQuarantinedTask(opts), exitCode)
	c.AddTask(true, NewReplicationTask(opts), exitCode)
	c.AddTask(true, NewUnmountedTask(opts), exitCode)
	c.AddTask(true, NewUpdaterSweepTask(opts), exitCode)
	registry.MustRegister(c)

	assert.HTTPRequest{
		Method:       "GET",
		Path:         "/metrics",
		ExpectStatus: 200,
		ExpectBody:   assert.FixtureFile("fixtures/recon_successful_collect.prom"),
	}.Check(t, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
}

func TestReconCollectorWithErrors(t *testing.T) {
	isTest = true

	pathToExecutable, err := filepath.Abs("../../build/mock-swift-recon-with-errors")
	if err != nil {
		t.Error(err)
	}

	registry := prometheus.NewPedanticRegistry()
	c := collector.New(0)
	exitCode := GetTaskExitCodeTypedDesc(registry)
	opts := &TaskOpts{
		PathToExecutable: pathToExecutable,
		HostTimeout:      1,
		CtxTimeout:       4 * time.Second,
	}
	c.AddTask(true, NewDiskUsageTask(opts), exitCode)
	c.AddTask(true, NewDriveAuditTask(opts), exitCode)
	c.AddTask(true, NewMD5Task(opts), exitCode)
	c.AddTask(true, NewQuarantinedTask(opts), exitCode)
	c.AddTask(true, NewReplicationTask(opts), exitCode)
	c.AddTask(true, NewUnmountedTask(opts), exitCode)
	c.AddTask(true, NewUpdaterSweepTask(opts), exitCode)
	registry.MustRegister(c)

	// For first attempt, we'll get metric results and exit code will be 0 because error
	// is reported only after max failure count has been exceeded.
	assert.HTTPRequest{
		Method:       "GET",
		Path:         "/metrics",
		ExpectStatus: 200,
		ExpectBody:   assert.FixtureFile("fixtures/recon_failed_collect.prom"),
	}.Check(t, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
}
