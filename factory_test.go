package opentelemetry_stacktrace_processor // import "github.com/joostlek/opentelemetry-stacktrace-processor"

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/processor/processortest"
)

func TestCreateDefaultConfig(t *testing.T) {
	factory := NewFactory()

	cfg := factory.CreateDefaultConfig()
	assert.NotNil(t, cfg, "failed to create default config")
	assert.NoError(t, componenttest.CheckConfigStruct(cfg))
}

func TestCreateProcessor(t *testing.T) {
	factory := NewFactory()

	cfg := factory.CreateDefaultConfig()
	creationSet := processortest.NewNopCreateSettings()
	tp, err := factory.CreateTracesProcessor(context.Background(), creationSet, cfg, nil)
	assert.NotNil(t, tp)
	assert.NoError(t, err, "cannot create trace processor")
}
