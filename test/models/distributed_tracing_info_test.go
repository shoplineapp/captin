package models_test

import (
	"context"
	"encoding/json"
	"testing"

	. "github.com/shoplineapp/captin/v2/models"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
)

func TestDistributedTracingInfoNormalFlow(t *testing.T) {
	spanCtxConfig := trace.SpanContextConfig{
		TraceID:    trace.TraceID([16]byte{0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11}),
		SpanID:     trace.SpanID([8]byte{0x22, 0x22, 0x22, 0x22, 0x22, 0x22, 0x22, 0x22}),
		TraceFlags: trace.TraceFlags(0x01),
	}
	var err error
	spanCtxConfig.TraceState, err = trace.ParseTraceState("a=b, c=d")
	require.NoError(t, err)

	spanContext := trace.NewSpanContext(spanCtxConfig)
	require.True(t, spanContext.IsValid())
	d := NewDistributedTracingInfoFromContext(trace.ContextWithSpanContext(context.Background(), spanContext))

	// serialize to JSON
	bz, err := json.Marshal(d)
	require.NoError(t, err)

	jsonStructure := make(map[string]string)
	err = json.Unmarshal(bz, &jsonStructure)
	require.NoError(t, err)
	require.Equal(t, map[string]string{
		"traceparent": "00-11111111111111111111111111111111-2222222222222222-01",
		"tracestate":  "a=b,c=d",
	}, jsonStructure)

	// recover from JSON
	var d2 DistributedTracingInfo
	err = json.Unmarshal(bz, &d2)
	require.NoError(t, err)
	require.Equal(t, d, d2)

	// recover the trace context
	ctx := d2.PropagateIntoContext(context.Background())
	spanContext2 := trace.SpanContextFromContext(ctx)
	require.Equal(t, spanContext.TraceID(), spanContext2.TraceID())
	require.Equal(t, spanContext.SpanID(), spanContext2.SpanID())
	require.Equal(t, spanContext.TraceFlags(), spanContext2.TraceFlags())
	require.Equal(t, spanContext.TraceState().String(), spanContext2.TraceState().String())
	// since the recovered span context is from propagator, it will be a remote context
	require.True(t, spanContext2.IsRemote())
}

func TestDistributedTracingInfoNilFlow(t *testing.T) {
	// default DistributedTracingInfo should work like a no-op propagator without nil pointer dereference
	d := DistributedTracingInfo{}

	// serialize to JSON should work without nil pointer dereference
	bz, err := json.Marshal(d)
	require.NoError(t, err)
	require.Equal(t, "{}", string(bz))

	// deserialize from null should work
	var d2 DistributedTracingInfo
	err = json.Unmarshal([]byte("null"), &d2)
	require.NoError(t, err)
	require.Equal(t, d, d2)

	// PropagateIntoContext should work without nil pointer dereference
	ctx := d.PropagateIntoContext(context.Background())
	spanContext2 := trace.SpanContextFromContext(ctx)
	require.False(t, spanContext2.IsValid())
}
