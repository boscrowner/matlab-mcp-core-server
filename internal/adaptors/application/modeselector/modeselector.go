// Copyright 2025-2026 The MathWorks, Inc.

package modeselector

import (
	"context"
	"fmt"
	"io"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/config"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/telemetry"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
)

type ConfigFactory interface {
	Config() (config.Config, messages.Error)
}

type Parser interface {
	Usage() (string, messages.Error)
}

type TelemetryFactory interface {
	Telemetry() (telemetry.Telemetry, messages.Error)
}

type WatchdogProcess interface { //nolint:iface // Intentional interface for deps injection
	StartAndWaitForCompletion(ctx context.Context) error
}

type Orchestrator interface { //nolint:iface // Intentional interface for deps injection
	StartAndWaitForCompletion(ctx context.Context) error
}

type OSLayer interface {
	Stdout() io.Writer
}

type LifecycleSignaler interface {
	RequestShutdown()
	WaitForShutdownToComplete() error
}

type LoggerFactory interface {
	GetGlobalLogger() (entities.Logger, messages.Error)
}

type ModeSelector struct {
	configFactory     ConfigFactory
	telemetryFactory  TelemetryFactory
	watchdogProcess   WatchdogProcess
	orchestrator      Orchestrator
	osLayer           OSLayer
	parser            Parser
	lifecycleSignaler LifecycleSignaler
	loggerFactory     LoggerFactory
}

func New(
	configFactory ConfigFactory,
	parser Parser,
	telemetryFactory TelemetryFactory,
	watchdogProcess WatchdogProcess,
	orchestrator Orchestrator,
	osLayer OSLayer,
	lifecycleSignaler LifecycleSignaler,
	loggerFactory LoggerFactory,
) *ModeSelector {
	return &ModeSelector{
		configFactory:     configFactory,
		parser:            parser,
		telemetryFactory:  telemetryFactory,
		watchdogProcess:   watchdogProcess,
		orchestrator:      orchestrator,
		osLayer:           osLayer,
		lifecycleSignaler: lifecycleSignaler,
		loggerFactory:     loggerFactory,
	}
}

func (m *ModeSelector) StartAndWaitForCompletion(ctx context.Context) error {
	config, err := m.configFactory.Config()
	if err != nil {
		return err
	}

	logger, err := m.loggerFactory.GetGlobalLogger()
	if err != nil {
		return err
	}

	telemetryInstance, err := m.telemetryFactory.Telemetry()
	if err != nil {
		return err
	}

	telemetryInstance.RecordServerStart(ctx)

	switch {
	case config.HelpMode():
		usage, messagesErr := m.parser.Usage()
		if messagesErr != nil {
			return m.shutdownAndReturn(logger, messagesErr)
		}
		_, err := fmt.Fprintf(m.osLayer.Stdout(), "%s\n", usage)
		return m.shutdownAndReturn(logger, err)
	case config.VersionMode():
		_, err := fmt.Fprintf(m.osLayer.Stdout(), "%s\n", config.Version())
		return m.shutdownAndReturn(logger, err)
	case config.WatchdogMode():
		return m.watchdogProcess.StartAndWaitForCompletion(ctx)
	default:
		return m.orchestrator.StartAndWaitForCompletion(ctx)
	}
}

func (m *ModeSelector) shutdownAndReturn(logger entities.Logger, err error) error {
	m.lifecycleSignaler.RequestShutdown()
	if shutdownErr := m.lifecycleSignaler.WaitForShutdownToComplete(); shutdownErr != nil {
		logger.WithError(shutdownErr).Warn("Shutdown failed")
	}
	return err
}
