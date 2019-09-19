package log

import (
	"fmt"

	"github.com/Microsoft/ApplicationInsights-Go/appinsights"
	"github.com/Microsoft/ApplicationInsights-Go/appinsights/contracts"
)

type AppInsightsLogger struct {
	c     appinsights.TelemetryClient
	props map[string]string
}

func (l *AppInsightsLogger) Logf(level Level, format string, v ...interface{}) {
	ailevel := contracts.Verbose
	switch level {
	case LevelInfo:
		ailevel = contracts.Information
	case LevelWarning:
		ailevel = contracts.Warning
	case LevelError:
		ailevel = contracts.Error
	case LevelCritical:
		ailevel = contracts.Critical
	}

	trace := appinsights.NewTraceTelemetry(fmt.Sprintf(format, v...), ailevel)
	if l.props != nil {
		for k, v := range l.props {
			trace.Properties[k] = v
		}
	}

	l.c.Track(trace)
}

func (l *AppInsightsLogger) Push() error { return nil }

type AppInsightsConfig struct {
	TelemetryClient appinsights.TelemetryClient
	Properties      map[string]string
}

func NewAppInsightsLogger(config *AppInsightsConfig) *AppInsightsLogger {
	return &AppInsightsLogger{
		c:     config.TelemetryClient,
		props: config.Properties,
	}
}
