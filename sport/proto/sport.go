package proto

//go:generate protoc --go_out=. --go-grpc_out=require_unimplemented_servers=false:. sport/sport.proto
