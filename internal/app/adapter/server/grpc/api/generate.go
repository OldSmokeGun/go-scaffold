package api

//go:generate kratos proto client --proto_path=../api --proto_path=../../../../../proto v1/greet.proto
//go:generate protoc-go-inject-tag -input=./v1/greet.pb.go

//go:generate kratos proto client --proto_path=../api --proto_path=../../../../../proto v1/system_user.proto
//go:generate protoc-go-inject-tag -input=./v1/system_user.pb.go

//go:generate kratos proto client --proto_path=../api --proto_path=../../../../../proto v1/system_role.proto
//go:generate protoc-go-inject-tag -input=./v1/system_role.pb.go

//go:generate kratos proto client --proto_path=../api --proto_path=../../../../../proto v1/system_permission.proto
//go:generate protoc-go-inject-tag -input=./v1/system_permission.pb.go

//go:generate kratos proto client --proto_path=../api --proto_path=../../../../../proto v1/product.proto
//go:generate protoc-go-inject-tag -input=./v1/product.pb.go
