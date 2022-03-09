package v1

//go:generate kratos proto client --proto_path=../../../../../proto greet/greet.proto
//go:generate protoc-go-inject-tag -input=./greet/greet.pb.go

//go:generate kratos proto client --proto_path=../../../../../proto user/user.proto
//go:generate protoc-go-inject-tag -input=./user/user.pb.go
