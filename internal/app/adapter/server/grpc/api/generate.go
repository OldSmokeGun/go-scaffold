package api

//go:generate kratos proto client --proto_path=../../../../../proto v1/greet.proto
//go:generate protoc-go-inject-tag -input=./v1/greet.pb.go

//go:generate kratos proto client --proto_path=../../../../../proto v1/user.proto
//go:generate protoc-go-inject-tag -input=./v1/user.pb.go
