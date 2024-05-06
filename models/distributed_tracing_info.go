package models

import (
	"context"
	"encoding/json"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

// DistributedTracingInfo is basically the same as propagation.MapCarrier
// But we reimplement it for easier interop (e.g. serialization and deserialization without caring map being nil)
type DistributedTracingInfo struct {
	Carrier propagation.MapCarrier
}

var _ propagation.TextMapCarrier = (*DistributedTracingInfo)(nil)
var _ json.Marshaler = (*DistributedTracingInfo)(nil)
var _ json.Unmarshaler = (*DistributedTracingInfo)(nil)

// To prevent nil pointer dereference, we need to initialize the MapCarrier before any action
func (d *DistributedTracingInfo) setup() {
	if d.Carrier == nil {
		d.Carrier = make(propagation.MapCarrier)
	}
}

func (d *DistributedTracingInfo) Get(s string) string {
	d.setup()
	return d.Carrier.Get(s)
}

func (d *DistributedTracingInfo) Set(key, value string) {
	d.setup()
	d.Carrier.Set(key, value)
}

func (d *DistributedTracingInfo) Keys() []string {
	d.setup()
	return d.Carrier.Keys()
}

func (d DistributedTracingInfo) MarshalJSON() ([]byte, error) {
	d.setup()
	return json.Marshal(d.Carrier)
}

func (d *DistributedTracingInfo) UnmarshalJSON(bz []byte) error {
	return json.Unmarshal(bz, &d.Carrier)
}

// If it's not set, then the default value of otel.GetTextMapPropagator() will be a no-op propagator.
// But we want to make sure that it works, so we prepend our propagator by a composite propagator.
// If someone use cases don't want propagation, ClearContext can be used.
func getPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, otel.GetTextMapPropagator())
}

// PropagateIntoContext takes a context, extract the tracing info from the DistributedTracingInfo,
// then return a new context with the tracing info injected based on the original context
func (d *DistributedTracingInfo) PropagateIntoContext(ctx context.Context) context.Context {
	// extract the tracing info from the carrier d into ctx
	return getPropagator().Extract(ctx, d)
}

func NewDistributedTracingInfoFromContext(ctx context.Context) DistributedTracingInfo {
	var d DistributedTracingInfo
	// inject the tracing info from ctx into the carrier d
	getPropagator().Inject(ctx, &d)
	return d
}

func (d *DistributedTracingInfo) InjectContext(ctx context.Context) *DistributedTracingInfo {
	getPropagator().Inject(ctx, d)
	return d
}

func (d *DistributedTracingInfo) ClearContext() *DistributedTracingInfo {
	d.Carrier = nil
	d.setup()
	return d
}

// GetTraceParent returns the traceparent header value, which provides a convenient way to get a representation of the current trace context
func (d *DistributedTracingInfo) GetTraceParent() string {
	return d.Get("traceparent")
}
