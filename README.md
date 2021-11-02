# grpc

issue. 
1) protoc-gen-go: unable to determine Go import path for "proto/person.proto"
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/person.proto
> option go_package = "grpc/proto"; 추가