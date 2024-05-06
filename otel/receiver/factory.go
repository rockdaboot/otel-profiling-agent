package otelreceiver

import (
	"context"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
)

var (
	typeStr = component.MustNewType("otelreceiver")
)

const (
	defaultInterval  = 5 * time.Second
	defaultProjectID = "1"
)

// NewFactory creates a factory for otelreceiver receiver.
func NewFactory() receiver.Factory {
	return receiver.NewFactory(
		typeStr,
		createDefaultConfig,
		receiver.WithTraces(createTracesReceiver, component.StabilityLevelAlpha))
}

func createTracesReceiver(_ context.Context, params receiver.CreateSettings,
	baseCfg component.Config, consumer consumer.Traces) (receiver.Traces, error) {
	logger := params.Logger
	tailtracerCfg := baseCfg.(*Config)

	traceRcvr := &otelReceiver{
		logger:       logger,
		nextConsumer: consumer,
		config:       tailtracerCfg,
	}

	return traceRcvr, nil
}
