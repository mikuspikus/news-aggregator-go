package tracer

import (
	"fmt"
	"io"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	jaegerconfig "github.com/uber/jaeger-client-go/config"
)

// NewTracer returns new Trace, Closer
func NewTracer(serviceName, host string) (opentracing.Tracer, io.Closer, error) {
	jcfg := jaegerconfig.Configuration{
		Sampler: &jaegerconfig.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &jaegerconfig.ReporterConfig{
			LogSpans:            false,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  host,
		},
		ServiceName: serviceName,
	}

	tracer, closer, err := jcfg.NewTracer()
	if err != nil {
		return nil, nil, fmt.Errorf("new tracer error: %v", err)
	}

	return tracer, closer, nil
}
