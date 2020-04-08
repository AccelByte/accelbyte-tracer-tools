[![Build Status](https://travis-ci.com/AccelByte/tracer-go.svg?branch=master)](https://travis-ci.com/AccelByte/tracer-go)

# tracer-go
Distributed tracing tools used by Accelbyte
Implemented to use jaeger tracing system and zipkin headers to pass spans

# usage
```go
    jaegerAgentHost := "" // jaeger:6831
    jaegerAgentEndpoint := "http://jaeger:14268/api/traces"
    serviceName := "test"
    realm := "node1"
    closer := InitGlobalTracer(jaegerAgentHost, jaegerAgentEndpoint, serviceName, realm)
    defer closer.Close()
```