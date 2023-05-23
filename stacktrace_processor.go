package opentelemetry_stacktrace_processor

import (
	"context"
	"fmt"
	"github.com/go-sourcemap/sourcemap"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.opentelemetry.io/collector/processor"
	"io"
	"os"
	"strconv"
	"strings"
)

type stackTraceProcessor struct {
	nextConsumer consumer.Traces
}

func (s *stackTraceProcessor) Start(ctx context.Context, host component.Host) error {
	return nil
}

func (s *stackTraceProcessor) Shutdown(ctx context.Context) error {
	return nil
}

func (s *stackTraceProcessor) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: true}
}

func (s *stackTraceProcessor) ConsumeTraces(ctx context.Context, td ptrace.Traces) error {
	for resourceSpanId := 0; resourceSpanId < td.ResourceSpans().Len(); resourceSpanId++ {
		resourceSpan := td.ResourceSpans().At(resourceSpanId)
		sdkLanguage, valid := resourceSpan.Resource().Attributes().Get("telemetry.sdk.language")
		if valid && sdkLanguage.Str() == "webjs" {
			ConsumeScopeSpans(resourceSpan.ScopeSpans())
		}
	}
	return s.nextConsumer.ConsumeTraces(ctx, td)
}

func ConsumeScopeSpans(scopeSpans ptrace.ScopeSpansSlice) {
	for scopeSpanId := 0; scopeSpanId < scopeSpans.Len(); scopeSpanId++ {
		spans := scopeSpans.At(scopeSpanId).Spans()
		ConsumeSpans(spans)
	}
}

func ConsumeSpans(spans ptrace.SpanSlice) {
	for spanId := 0; spanId < spans.Len(); spanId++ {
		span := spans.At(spanId)
		ConsumeSpan(span)
	}
}

func ConsumeSpan(span ptrace.Span) {
	for eventId := 0; eventId < span.Events().Len(); eventId++ {
		event := span.Events().At(eventId)
		if event.Name() == "exception" {
			ConsumeException(event)
		}
	}
}

func ConsumeException(event ptrace.SpanEvent) {
	stacktrace, valid := event.Attributes().Get("exception.stacktrace")
	if valid != true {
		return
	}
	var res []string
	lines := strings.Split(stacktrace.Str(), "\n")
	for n := 0; n < len(lines); n++ {
		line := lines[n]
		pieces := strings.Split(line, "@")
		if len(pieces) == 1 {
			res = append(res, line)
			continue
		}
		trace := pieces[0]
		source := pieces[1]
		sourceFilePieces := strings.Split(source, "/")[len(strings.Split(source, "/"))-1]
		mapPieces := strings.Split(sourceFilePieces, ":")
		sourceFile := mapPieces[0]
		mapLine, err := strconv.Atoi(mapPieces[1])
		if err != nil {
			res = append(res, line)
			continue
		}
		mapColumn, err := strconv.Atoi(mapPieces[2])
		if err != nil {
			res = append(res, line)
			continue
		}
		file, err := os.Open(fmt.Sprintf("%s.map", sourceFile))
		if err != nil {
			res = append(res, line)
			continue
		}
		b, err := io.ReadAll(file)
		if err != nil {
			res = append(res, line)
			continue
		}

		smap, err := sourcemap.Parse(sourceFile, b)
		if err != nil {
			panic(err)
		}

		finalFile, fn, sourceLine, sourceColumn, ok := smap.Source(mapLine, mapColumn)
		fmt.Println(finalFile, fn, sourceLine, sourceColumn, ok)
		res = append(res, fmt.Sprintf("%s@%s:%d:%d", trace, finalFile, sourceLine, sourceColumn))
	}
	fmt.Println(res)
	event.Attributes().PutStr("exception.stacktrace", strings.Join(res, "\n"))
}

var _ consumer.Traces = (*stackTraceProcessor)(nil)

func newStackTraceProcessor(set processor.CreateSettings, next consumer.Traces, config *Config) (processor.Traces, error) {
	return &stackTraceProcessor{nextConsumer: next}, nil
}
