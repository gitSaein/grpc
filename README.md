
## Installation
> go get -u go.opencensus.io
>
> go get -u contrib.go.opencensus.io/exporter/stackdriver

## Prerequisites
OpenCensus Go libraries require Go 1.8 or later.

## Command

protoc -I rpc rpc/defs.proto --go_out=plugins=grpc:rpc

protoc -I proto proto/service.proto --go_out=plugins=grpc:proto

## Refs
https://opencensus.io/guides/grpc/go/#0


https://pkg.go.dev/go.opencensus.io#section-readme
