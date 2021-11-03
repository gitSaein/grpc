# grpc

client streaming grpc
- client sends a stream of messages to the server instead of a single message.
- stub local object, wrapping protocol buffer message


server streaming grpc 
- server returns a stream of messages in response to a client's request. 
- service, function, decoding -> service -> encoding -> client

metadata
- rpc call key-value form

channel
- host, port grpc server connection


proto build
- protoc --go_out=build/api/. --go_opt=paths=source_relative --go-grpc_out=build/api/. --go-grpc_opt=paths=source_relative proto/api.proto


issue. 
1) protoc-gen-go: unable to determine Go import path for "proto/person.proto"

> protoc --go_out=build/. --go_opt=paths=source_relative --go-grpc_out=build/. --go-grpc_opt=paths=source_relative proto/api/api.proto 
> build source) go_package = "grpc/proto/api"; 추가