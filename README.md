# go

## setting
- GOPATH : go lang work folder   
- GOROOT : go lang setting files installed on this path ( /usr/local/go )


## env
- GO111MODULE: 

# protocol buffer

## compiler
--go_out : Go compiler ouput path ( .proto -> .pb.go )
--proto_path : import directory
--go-grpc_out :  *_grpc.pb.go 로 빌드
--go_out :            *.pb.go 로 빌드
 > go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
 > protoc --go-grpc_out=build/. ./proto/api/api.proto
 > 


# grpca

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
- protoc --go_out=build/. --go_opt=paths=source_relative --go-grpc_out=build/. --go-grpc_opt=paths=source_relative proto/api.proto


issue. 
1) protoc-gen-go: unable to determine Go import path for "proto/person.proto"

> protoc --go_out=build/. --go_opt=paths=source_relative --go-grpc_out=build/. --go-grpc_opt=paths=source_relative proto/api/api.proto 
> build source) go_package = "grpc/proto/api"; 추가

2) package [name] is not in GOROOT 에러

> go env -w GO111MODULE="off"
>  project path ==> $GOPATH/src/{go project path}

3) could not import

> go mod init    
> go mod tiny

