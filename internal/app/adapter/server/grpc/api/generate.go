package api

//go:generate kratos proto client --proto_path=../api --proto_path=../../../../../proto v1/greet.proto
//go:generate protoc-go-inject-tag -input=./v1/greet.pb.go

//go:generate kratos proto client --proto_path=../api --proto_path=../../../../../proto v1/user.proto
//go:generate protoc-go-inject-tag -input=./v1/user.pb.go

//go:generate kratos proto client --proto_path=../api --proto_path=../../../../../proto v1/role.proto
//go:generate protoc-go-inject-tag -input=./v1/role.pb.go

//go:generate kratos proto client --proto_path=../api --proto_path=../../../../../proto v1/permission.proto
//go:generate protoc-go-inject-tag -input=./v1/permission.pb.go

//go:generate kratos proto client --proto_path=../api --proto_path=../../../../../proto v1/product.proto
//go:generate protoc-go-inject-tag -input=./v1/product.pb.go
