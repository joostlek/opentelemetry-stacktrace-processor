package opentelemetry_stacktrace_processor

import (
	"context"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/processor/processortest"
	"testing"
	"time"
)

func TestBatchProcessorSpansDeliveredEnforceBatchSize(t *testing.T) {
	sink := new(consumertest.TracesSink)
	cfg := createDefaultConfig().(*Config)
	creationSet := processortest.NewNopCreateSettings()
	processor, err := newStackTraceProcessor(creationSet, sink, cfg)
	require.NoError(t, err)
	require.NoError(t, processor.Start(context.Background(), componenttest.NewNopHost()))

	td := GenerateTraces()
	assert.NoError(t, processor.ConsumeTraces(context.Background(), td))
	for {
		if sink.SpanCount() == 1 {
			break
		}
	}
	traces := sink.AllTraces()
	event := traces[0].ResourceSpans().At(0).ScopeSpans().At(0).Spans().At(0).Events().At(0)
	stackTrace, valid := event.Attributes().Get("exception.stacktrace")
	assert.True(t, valid)
	assert.Equal(t, "padStart@lib/lineSlicer.js:11:53\npadStart@http://localhost:4203/lineSlicers.min.js:1:228\n", stackTrace.Str())

	require.NoError(t, processor.Shutdown(context.Background()))
}

var (
	spanStartTimestamp = pcommon.NewTimestampFromTime(time.Date(2020, 2, 11, 20, 26, 12, 321, time.UTC))
	spanEventTimestamp = pcommon.NewTimestampFromTime(time.Date(2020, 2, 11, 20, 26, 13, 123, time.UTC))
	spanEndTimestamp   = pcommon.NewTimestampFromTime(time.Date(2020, 2, 11, 20, 26, 13, 789, time.UTC))
)

func GenerateTraces() ptrace.Traces {
	td := ptrace.NewTraces()
	td.ResourceSpans().AppendEmpty().Resource().Attributes().PutStr("telemetry.sdk.language", "webjs")
	ss := td.ResourceSpans().At(0).ScopeSpans().AppendEmpty().Spans()
	ss.EnsureCapacity(1)
	fillExceptionSpan(ss.AppendEmpty())
	return td
}

func fillExceptionSpan(span ptrace.Span) {
	span.SetName("@pi/error-handler")
	span.SetStartTimestamp(spanStartTimestamp)
	span.SetEndTimestamp(spanEndTimestamp)
	span.SetDroppedAttributesCount(0)
	span.SetTraceID([16]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10})
	span.SetSpanID([8]byte{0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18})
	evs := span.Events()
	ev0 := evs.AppendEmpty()
	ev0.SetTimestamp(spanEventTimestamp)
	ev0.SetName("exception")
	ev0.Attributes().PutStr("exception.type", "TypeError")
	ev0.Attributes().PutStr("exception.message", "symbol.viewBox.baseVal is null")
	ev0.Attributes().PutStr("exception.stacktrace", "padStart@http://localhost:4203/lineSlicer.min.js:1:228\npadStart@http://localhost:4203/lineSlicers.min.js:1:228\n")
	span.Status().SetCode(2)
}
