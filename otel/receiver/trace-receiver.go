package otelreceiver

import (
	"context"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.uber.org/zap"
)

type otelReceiver struct {
	host         component.Host
	cancel       context.CancelFunc
	logger       *zap.Logger
	nextConsumer consumer.Traces
	config       *Config
}

func (otelRcvr *otelReceiver) Start(ctx context.Context, host component.Host) error {
	otelRcvr.host = host
	ctx, otelRcvr.cancel = context.WithCancel(context.Background())

	interval, _ := time.ParseDuration(otelRcvr.config.Interval)
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				otelRcvr.logger.Info("I should start processing traces now!")
			case <-ctx.Done():
				return
			}
		}
	}()
	return nil
}

func (otelRcvr *otelReceiver) Shutdown(ctx context.Context) error {
	otelRcvr.cancel()
	return nil
}
