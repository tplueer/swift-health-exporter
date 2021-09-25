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

package main

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sapcc/go-bits/assert"

	"github.com/sapcc/swift-health-exporter/internal/collector"
	"github.com/sapcc/swift-health-exporter/internal/collector/dispersion"
	"github.com/sapcc/swift-health-exporter/internal/collector/recon"
)

func TestCollector(t *testing.T) {
	testCollector(t,
		"build/mock-swift-dispersion-report",
		"build/mock-swift-recon",
		"test/fixtures/successful_collect.prom")
}

func TestCollectorWithErrors(t *testing.T) {
	testCollector(t,
		"build/mock-swift-dispersion-report-with-errors",
		"build/mock-swift-recon-with-errors",
		"test/fixtures/failed_collect.prom")
}

func testCollector(t *testing.T, dispersionReportPath, reconPath, fixturesPath string) {
	recon.IsTest = true

	dispersionReportAbsPath, err := filepath.Abs(dispersionReportPath)
	if err != nil {
		t.Error(err)
	}
	reconAbsPath, err := filepath.Abs(reconPath)
	if err != nil {
		t.Error(err)
	}

	registry := prometheus.NewPedanticRegistry()
	c := collector.New()
	s := collector.NewScraper(0)

	dispersionExitCode := dispersion.GetTaskExitCodeGaugeVec(registry)
	addTask(true, c, s, dispersion.NewReportTask(dispersionReportAbsPath, 20*time.Second), dispersionExitCode)

	reconExitCode := recon.GetTaskExitCodeGaugeVec(registry)
	opts := &recon.TaskOpts{
		PathToExecutable: reconAbsPath,
		HostTimeout:      1,
		CtxTimeout:       4 * time.Second,
	}
	addTask(true, c, s, recon.NewDiskUsageTask(opts), reconExitCode)
	addTask(true, c, s, recon.NewDriveAuditTask(opts), reconExitCode)
	addTask(true, c, s, recon.NewMD5Task(opts), reconExitCode)
	addTask(true, c, s, recon.NewQuarantinedTask(opts), reconExitCode)
	addTask(true, c, s, recon.NewReplicationTask(opts), reconExitCode)
	addTask(true, c, s, recon.NewUnmountedTask(opts), reconExitCode)
	addTask(true, c, s, recon.NewUpdaterSweepTask(opts), reconExitCode)

	registry.MustRegister(c)

	s.UpdateAllMetrics()
	assert.HTTPRequest{
		Method:       "GET",
		Path:         "/metrics",
		ExpectStatus: 200,
		ExpectBody:   assert.FixtureFile(fixturesPath),
	}.Check(t, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
}
