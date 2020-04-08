package tracergo

import (
	"context"
	"testing"

	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uber/jaeger-client-go"
)

func TestGetSpanFromRestfulContextWithoutSpan(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)

	closer := InitGlobalTracer("", "", "test", "")
	defer closer.Close()

	span := GetSpanFromRestfulContext(context.Background())
	require.NotNil(t, span)

	require.NotNil(t, span.Context())
	require.IsType(t, jaeger.SpanContext{}, span.Context())
	require.NotNil(t, span.Context().(jaeger.SpanContext))
	require.NotEmpty(t, span.Context().(jaeger.SpanContext).TraceID())
	assert.NotEmpty(t, span.Context().(jaeger.SpanContext).TraceID().String())
}

func TestGetSpanFromRestfulContextWithSpan(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)

	closer := InitGlobalTracer("", "", "test", "")
	defer closer.Close()

	expectedSpan, _ := StartSpanFromContext(context.Background(), "test")
	ctx := context.WithValue(context.Background(), SpanContextKey, expectedSpan)

	span := GetSpanFromRestfulContext(ctx)
	require.NotNil(t, span)

	require.NotNil(t, span.Context())
	require.IsType(t, jaeger.SpanContext{}, span.Context())
	require.NotNil(t, span.Context().(jaeger.SpanContext))
	require.NotEmpty(t, span.Context().(jaeger.SpanContext).TraceID())

	assert.Equal(t,
		expectedSpan.Context().(jaeger.SpanContext).TraceID().String(),
		span.Context().(jaeger.SpanContext).TraceID().String(),
	)
}

func TestChildSpanFromRemoteSpan(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)

	closer := InitGlobalTracer("", "", "test", "")
	defer closer.Close()

	expectedSpan, _ := opentracing.StartSpanFromContext(context.Background(), "test")

	spanContextStr := expectedSpan.Context().(jaeger.SpanContext).String()

	span, _ := ChildSpanFromRemoteSpan(context.Background(), "test", spanContextStr)

	assert.Equal(t,
		expectedSpan.Context().(jaeger.SpanContext).TraceID().String(),
		span.Context().(jaeger.SpanContext).TraceID().String(),
	)

	assert.Equal(t,
		expectedSpan.Context().(jaeger.SpanContext).SpanID().String(),
		span.Context().(jaeger.SpanContext).ParentID().String(),
	)
}

func TestChildSpanFromRemoteSpan_EmptySpanContextString(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)

	closer := InitGlobalTracer("", "", "test", "")
	defer closer.Close()

	scope, _ := ChildSpanFromRemoteSpan(context.Background(), "test", "")

	assert.NotEmpty(t,
		scope.Context().(jaeger.SpanContext).TraceID().String(),
	)

	assert.NotEmpty(t,
		scope.Context().(jaeger.SpanContext).ParentID().String(),
	)
}

func TestGetSpanContextString_NotEmptySpanContext(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)

	closer := InitGlobalTracer("", "", "test", "")
	defer closer.Close()

	span := GetSpanFromRestfulContext(context.Background())
	require.NotNil(t, span)

	require.NotNil(t, span.Context())
	require.IsType(t, jaeger.SpanContext{}, span.Context())
	require.NotNil(t, span.Context().(jaeger.SpanContext))
	require.NotEmpty(t, span.Context().(jaeger.SpanContext).TraceID())
	assert.NotEmpty(t, span.Context().(jaeger.SpanContext).TraceID().String())

	spanContextString := GetSpanContextString(span)
	assert.NotEmpty(t, spanContextString)
}

func TestGetSpanContextString_EmptySpanContext(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)

	closer := InitGlobalTracer("", "", "test", "")
	defer closer.Close()

	span := opentracing.Span(nil)
	require.Nil(t, span)

	spanContextString := GetSpanContextString(span)
	assert.Empty(t, spanContextString)
}
