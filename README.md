# grpc

client streaming grpc
- client sends a stream of messages to the server instead of a single message.


server streaming grpc
- server returns a stream of messages in response to a client's request. 


proto build
- protoc --go_out=build/api/. --go_opt=paths=source_relative --go-grpc_out=build/api/. --go-grpc_opt=paths=source_relative proto/api.proto


issue. 
1) protoc-gen-go: unable to determine Go import path for "proto/person.proto"
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/person.proto
> option go_package = "grpc/proto"; 추가