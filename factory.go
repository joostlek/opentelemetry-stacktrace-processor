package opentelemetry_stacktrace_processor // import "github.com/joostlek/opentelemetry-stacktrace-processor"
import (
	"context"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/processor"
)

const typeStr = "stacktrace"

func NewFactory() processor.Factory {
	return processor.NewFactory(
		typeStr,
		createDefaultConfig,
		processor.WithTraces(createTraces, component.StabilityLevelStable),
	)
}

func createDefaultConfig() component.Config {
	return &Config{
		SourceMapDirs: make([]string, 0),
	}
}

func createTraces(
	_ context.Context,
	set processor.CreateSettings,
	cfg component.Config,
	nextConsumer consumer.Traces,
) (processor.Traces, error) {
	return newStackTraceProcessor(set, nextConsumer, cfg.(*Config))
}
