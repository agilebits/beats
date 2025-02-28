// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package beats

import (
	"crypto/sha256"
	"fmt"
	"path/filepath"

	"github.com/elastic/beats/v7/x-pack/elastic-agent/pkg/agent/application/paths"
	"github.com/elastic/beats/v7/x-pack/elastic-agent/pkg/agent/program"
)

const (
	// args: data path, pipeline name, application name
	logFileFormat = "%s/logs/%s/%s-json.log"
	// args: data path, install path, pipeline name, application name
	logFileFormatWin = "%s\\logs\\%s\\%s-json.log"

	// args: pipeline name, application name
	mbEndpointFileFormatWin = `npipe:///%s-%s`

	// args: pipeline name, application name
	agentMbEndpointFileFormatWin = `npipe:///elastic-agent`
)

func getMonitoringEndpoint(spec program.Spec, operatingSystem, pipelineID string) string {
	if endpoint, ok := spec.MetricEndpoints[operatingSystem]; ok {
		return endpoint
	}
	if operatingSystem == "windows" {
		return fmt.Sprintf(mbEndpointFileFormatWin, pipelineID, spec.Cmd)
	}
	// unix socket path cannot be longer than 107 characters
	path := fmt.Sprintf("unix://%s.sock", filepath.Join(paths.TempDir(), pipelineID, spec.Cmd, spec.Cmd))
	if len(path) <= 107 {
		return path
	}
	// place in global /tmp to ensure that its small enough to fit; current path is way to long
	// for it to be used, but needs to be unique per Agent (in the case that multiple are running)
	return fmt.Sprintf(`unix:///tmp/elastic-agent-%x.sock`, sha256.Sum256([]byte(path)))
}

func getLoggingFile(spec program.Spec, operatingSystem, installPath, pipelineID string) string {
	if path, ok := spec.LogPaths[operatingSystem]; ok {
		return path
	}
	if operatingSystem == "windows" {
		return fmt.Sprintf(logFileFormatWin, paths.Home(), pipelineID, spec.Cmd)
	}
	return fmt.Sprintf(logFileFormat, paths.Home(), pipelineID, spec.Cmd)
}

// AgentMonitoringEndpoint returns endpoint with exposed metrics for agent.
func AgentMonitoringEndpoint(operatingSystem string) string {
	if operatingSystem == "windows" {
		return agentMbEndpointFileFormatWin
	}
	// unix socket path cannot be longer than 107 characters
	path := fmt.Sprintf("unix://%s.sock", filepath.Join(paths.TempDir(), "elastic-agent"))
	if len(path) <= 107 {
		return path
	}
	// place in global /tmp to ensure that its small enough to fit; current path is way to long
	// for it to be used, but needs to be unique per Agent (in the case that multiple are running)
	return fmt.Sprintf(`unix:///tmp/elastic-agent-%x.sock`, sha256.Sum256([]byte(path)))
}

// AgentPrefixedMonitoringEndpoint returns endpoint with exposed metrics for agent.
func AgentPrefixedMonitoringEndpoint(operatingSystem string) string {
	return httpPlusPrefix + AgentMonitoringEndpoint(operatingSystem)
}
