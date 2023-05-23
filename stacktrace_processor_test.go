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
	ev0.Attributes().PutStr("exception.stacktrace", "loadSvg@http://localhost:4203/main.js:1514:21\nngOnInit/this.svgDefinition$<@http://localhost:4203/main.js:1397:23\nmap/</<@http://localhost:4203/vendor.js:47989:31\nOperatorSubscriber/this._next<@http://localhost:4203/vendor.js:47268:15\nnext@http://localhost:4203/vendor.js:46165:12\n_next@http://localhost:4203/vendor.js:46192:22\nnext@http://localhost:4203/vendor.js:46165:12\nmap/</<@http://localhost:4203/vendor.js:47989:18\nOperatorSubscriber/this._next<@http://localhost:4203/vendor.js:47268:15\nnext@http://localhost:4203/vendor.js:46165:12\nfilter/</<@http://localhost:4203/vendor.js:47771:175\nOperatorSubscriber/this._next<@http://localhost:4203/vendor.js:47268:15\nnext@http://localhost:4203/vendor.js:46165:12\ndoInnerSub/<@http://localhost:4203/vendor.js:48106:20\nOperatorSubscriber/this._next<@http://localhost:4203/vendor.js:47268:15\nnext@http://localhost:4203/vendor.js:46165:12\nswitchMap/</</innerSubscriber<@http://localhost:4203/vendor.js:48637:243\nOperatorSubscriber/this._next<@http://localhost:4203/vendor.js:47268:15\nnext@http://localhost:4203/vendor.js:46165:12\nonLoad@http://localhost:4203/vendor.js:83794:22\ninvokeTask@http://localhost:4203/polyfills.js:8205:171\nonInvokeTask@http://localhost:4203/vendor.js:110067:22\ninvokeTask@http://localhost:4203/polyfills.js:8205:54\nonInvokeTask@http://localhost:4203/vendor.js:110367:25\ninvokeTask@http://localhost:4203/polyfills.js:8205:54\nrunTask@http://localhost:4203/polyfills.js:8007:37\npatchRunTask@http://localhost:4203/vendor.js:3777:27\ninvokeTask@http://localhost:4203/polyfills.js:8282:26\ninvokeTask@http://localhost:4203/polyfills.js:9406:12\nglobalCallback@http://localhost:4203/polyfills.js:9447:33\nglobalZoneAwareCallback@http://localhost:4203/polyfills.js:9467:12\nEventListener.handleEvent*customScheduleGlobal@http://localhost:4203/polyfills.js:9555:37\nscheduleTask@http://localhost:4203/polyfills.js:8195:16\nonScheduleTask@http://localhost:4203/vendor.js:110060:21\nscheduleTask@http://localhost:4203/polyfills.js:8190:43\nonScheduleTask@http://localhost:4203/polyfills.js:8107:61\nscheduleTask@http://localhost:4203/polyfills.js:8190:43\nscheduleTask@http://localhost:4203/polyfills.js:8046:35\npatchScheduleTask@http://localhost:4203/vendor.js:3740:25\nscheduleEventTask@http://localhost:4203/polyfills.js:8071:19\nmakeAddListener/<@http://localhost:4203/polyfills.js:9706:27\nhandle/</<@http://localhost:4203/vendor.js:83881:13\n_trySubscribe@http://localhost:4203/vendor.js:45790:19\nsubscribe/<@http://localhost:4203/vendor.js:45784:113\nerrorContext@http://localhost:4203/vendor.js:50241:5\nsubscribe@http://localhost:4203/vendor.js:45779:69\nswitchMap/</<@http://localhost:4203/vendor.js:48637:100\nOperatorSubscriber/this._next<@http://localhost:4203/vendor.js:47268:15\nnext@http://localhost:4203/vendor.js:46165:12\nfromArrayLike/<@http://localhost:4203/vendor.js:46991:18\n_trySubscribe@http://localhost:4203/vendor.js:45790:19\nsubscribe/<@http://localhost:4203/vendor.js:45784:113\nerrorContext@http://localhost:4203/vendor.js:50241:5\nsubscribe@http://localhost:4203/vendor.js:45779:69\nswitchMap/<@http://localhost:4203/vendor.js:48633:12\noperate/</<@http://localhost:4203/vendor.js:50530:18\nsubscribe/<@http://localhost:4203/vendor.js:45784:42\nerrorContext@http://localhost:4203/vendor.js:50241:5\nsubscribe@http://localhost:4203/vendor.js:45779:69\ndoInnerSub@http://localhost:4203/vendor.js:48101:95\nouterNext@http://localhost:4203/vendor.js:48096:52\nOperatorSubscriber/this._next<@http://localhost:4203/vendor.js:47268:15\nnext@http://localhost:4203/vendor.js:46165:12\nfromArrayLike/<@http://localhost:4203/vendor.js:46991:18\n_trySubscribe@http://localhost:4203/vendor.js:45790:19\nsubscribe/<@http://localhost:4203/vendor.js:45784:113\nerrorContext@http://localhost:4203/vendor.js:50241:5\nsubscribe@http://localhost:4203/vendor.js:45779:69\nmergeInternals@http://localhost:4203/vendor.js:48129:10\nmergeMap/<@http://localhost:4203/vendor.js:48167:149\noperate/</<@http://localhost:4203/vendor.js:50530:18\nsubscribe/<@http://localhost:4203/vendor.js:45784:42\nerrorContext@http://localhost:4203/vendor.js:50241:5\nsubscribe@http://localhost:4203/vendor.js:45779:69\nfilter/<@http://localhost:4203/vendor.js:47771:12\noperate/</<@http://localhost:4203/vendor.js:50530:18\nsubscribe/<@http://localhost:4203/vendor.js:45784:42\nerrorContext@http://localhost:4203/vendor.js:50241:5\nsubscribe@http://localhost:4203/vendor.js:45779:69\nmap/<@http://localhost:4203/vendor.js:47988:12\noperate/</<@http://localhost:4203/vendor.js:50530:18\nsubscribe/<@http://localhost:4203/vendor.js:45784:42\nerrorContext@http://localhost:4203/vendor.js:50241:5\nsubscribe@http://localhost:4203/vendor.js:45779:69\ncatchError/<@http://localhost:4203/vendor.js:47404:23\noperate/</<@http://localhost:4203/vendor.js:50530:18\nsubscribe/<@http://localhost:4203/vendor.js:45784:42\nerrorContext@http://localhost:4203/vendor.js:50241:5\nsubscribe@http://localhost:4203/vendor.js:45779:69\ncatchError/</innerSub<@http://localhost:4203/vendor.js:47409:23\nOperatorSubscriber/this._error<@http://localhost:4203/vendor.js:47275:16\nerror@http://localhost:4203/vendor.js:46173:12\n_error@http://localhost:4203/vendor.js:46196:24\nerror@http://localhost:4203/vendor.js:46173:12\n_error@http://localhost:4203/vendor.js:46196:24\nerror@http://localhost:4203/vendor.js:46173:12\n_error@http://localhost:4203/vendor.js:46196:24\nerror@http://localhost:4203/vendor.js:46173:12\n_error@http://localhost:4203/vendor.js:46196:24\nerror@http://localhost:4203/vendor.js:46173:12\ninit@http://localhost:4203/vendor.js:47189:41\n_trySubscribe@http://localhost:4203/vendor.js:45790:19\nsubscribe/<@http://localhost:4203/vendor.js:45784:113\nerrorContext@http://localhost:4203/vendor.js:50241:5\nsubscribe@http://localhost:4203/vendor.js:45779:69\ncatchError/</innerSub<@http://localhost:4203/vendor.js:47409:23\nOperatorSubscriber/this._error<@http://localhost:4203/vendor.js:47275:16\nerror@http://localhost:4203/vendor.js:46173:12\ninit@http://localhost:4203/vendor.js:47189:41\n_trySubscribe@http://localhost:4203/vendor.js:45790:19\nsubscribe/<@http://localhost:4203/vendor.js:45784:113\nerrorContext@http://localhost:4203/vendor.js:50241:5\nsubscribe@http://localhost:4203/vendor.js:45779:69\ncatchError/</innerSub<@http://localhost:4203/vendor.js:47409:23\nOperatorSubscriber/this._error<@http://localhost:4203/vendor.js:47275:16\nerror@http://localhost:4203/vendor.js:46173:12\n_error@http://localhost:4203/vendor.js:46196:24\nerror@http://localhost:4203/vendor.js:46173:12\nonLoad@http://localhost:4203/vendor.js:83806:22\ninvokeTask@http://localhost:4203/polyfills.js:8205:171\nonInvokeTask@http://localhost:4203/vendor.js:110067:22\ninvokeTask@http://localhost:4203/polyfills.js:8205:54\nonInvokeTask@http://localhost:4203/vendor.js:110367:25\ninvokeTask@http://localhost:4203/polyfills.js:8205:54\nrunTask@http://localhost:4203/polyfills.js:8007:37\npatchRunTask@http://localhost:4203/vendor.js:3777:27\ninvokeTask@http://localhost:4203/polyfills.js:8282:26\ninvokeTask@http://localhost:4203/polyfills.js:9406:12\nglobalCallback@http://localhost:4203/polyfills.js:9447:33\n")
	span.Status().SetCode(2)
}
